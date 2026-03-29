#!/usr/bin/env python3
"""Google Doc HTML → スライドHTML変換スクリプト（SKILL.md準拠）

Usage:
    python3 scripts/generate_slide_html.py docs/manuals/{マニュアル名}/
    python3 scripts/generate_slide_html.py --all   # 全フォルダ一括変換

入力: docs/manuals/{名前}/{名前}.html + images/
出力: docs/manuals/{名前}/slide.html
"""

import argparse
import base64
import glob
import mimetypes
import os
import re
import sys
from html.parser import HTMLParser


# ── 除外パターン ─────────────────────────────
PHONE_RE = re.compile(r'0\d{1,4}[-‐ー]?\d{1,4}[-‐ー]?\d{2,4}')
EMAIL_RE = re.compile(r'[\w.+-]+@[\w-]+\.[\w.]+')
TIME_RANGE_RE = re.compile(r'\d{1,2}:\d{2}\s*[〜~～ー−\-]\s*\d{1,2}:\d{2}')
KV_RE = re.compile(r'^(.+?)[：:]\s*(.+)$')

EXCLUDE_SECTIONS = ['代表連絡先', '連絡先一覧', '緊急連絡先']
ALERT_KEYWORDS = ['必ず', 'すること', '禁止', '注意', '厳禁', '絶対']


def strip_private(text):
    text = PHONE_RE.sub('', text)
    text = EMAIL_RE.sub('', text)
    return text.strip()


# ── HTML パーサー ─────────────────────────────

class GoogleDocParser(HTMLParser):
    """Google Doc HTMLを構造化データに変換"""

    def __init__(self):
        super().__init__()
        self.elements = []  # [(type, data), ...]
        self.tag_stack = []
        self.current_text = ''
        self.in_table = False
        self.table_rows = []
        self.current_row = []
        self.current_cell_text = ''
        self.skip = False  # script/style内はスキップ

    def handle_starttag(self, tag, attrs):
        attrs_dict = dict(attrs)

        if tag in ('script', 'style'):
            self.skip = True
            return

        if tag == 'table':
            self._flush_text()
            self.in_table = True
            self.table_rows = []
        elif tag == 'tr':
            self.current_row = []
        elif tag in ('td', 'th'):
            self.current_cell_text = ''
        elif tag == 'img':
            self._flush_text()
            src = attrs_dict.get('src', '')
            if src:
                self.elements.append(('img', src))
        elif tag == 'br':
            self.current_text += '\n'
        elif tag in ('h1', 'h2', 'h3', 'h4', 'h5', 'h6'):
            self._flush_text()
            self.tag_stack.append(tag)
        elif tag == 'li':
            self._flush_text()
            self.tag_stack.append('li')
        elif tag == 'hr':
            self._flush_text()
            self.elements.append(('hr', ''))

    def handle_endtag(self, tag):
        if tag in ('script', 'style'):
            self.skip = False
            return

        if tag == 'table':
            if self.table_rows:
                self.elements.append(('table', self.table_rows))
            self.in_table = False
            self.table_rows = []
        elif tag == 'tr':
            if self.current_row:
                self.table_rows.append(self.current_row)
        elif tag in ('td', 'th'):
            self.current_row.append(self.current_cell_text.strip())
            self.current_cell_text = ''
        elif tag in ('h1', 'h2', 'h3', 'h4', 'h5', 'h6'):
            text = self.current_text.strip()
            if text:
                self.elements.append((tag, text))
            self.current_text = ''
            if self.tag_stack and self.tag_stack[-1] == tag:
                self.tag_stack.pop()
        elif tag == 'li':
            text = self.current_text.strip()
            if text:
                self.elements.append(('li', text))
            self.current_text = ''
            if self.tag_stack and self.tag_stack[-1] == 'li':
                self.tag_stack.pop()
        elif tag == 'p':
            self._flush_text()

    def handle_data(self, data):
        if self.skip:
            return
        if self.in_table and not self.tag_stack:
            self.current_cell_text += data
        else:
            self.current_text += data

    def _flush_text(self):
        text = self.current_text.strip()
        if text:
            tag = self.tag_stack[-1] if self.tag_stack else 'p'
            self.elements.append((tag, text))
        self.current_text = ''

    def finish(self):
        self._flush_text()


# ── セクション構造化 ──────────────────────────

def build_sections(elements):
    """パース結果をセクション（見出し区切り）に構造化"""
    sections = []
    current = None
    title_text = None
    subtitle_parts = []

    # タイトル検出: 最初のh1より前のテキストからタイトルを推定
    first_heading_idx = None
    for i, (tag, data) in enumerate(elements):
        if tag in ('h1', 'h2'):
            first_heading_idx = i
            break

    # タイトルセクション: 最初のheadingより前のp要素
    pre_heading = []
    start_idx = 0
    if first_heading_idx is not None:
        for i in range(first_heading_idx):
            tag, data = elements[i]
            if tag == 'p' and data:
                pre_heading.append(data)
        start_idx = first_heading_idx

    # pre_heading からタイトルを推定
    # 一番長い or 「マニュアル」を含むものをタイトルとする
    for text in pre_heading:
        if 'マニュアル' in text and not title_text:
            title_text = text
        elif not title_text and len(text) > 5:
            title_text = text
        else:
            subtitle_parts.append(text)

    if title_text:
        sections.append({
            'type': 'title',
            'title': strip_private(title_text),
            'subtitle': ' / '.join(subtitle_parts[:2]) if subtitle_parts else ''
        })

    # メインコンテンツをheading区切りでセクション化
    excluded = False
    for i in range(start_idx, len(elements)):
        tag, data = elements[i]

        if tag in ('h1', 'h2'):
            text = strip_private(data)
            if any(kw in text for kw in EXCLUDE_SECTIONS):
                excluded = True
                continue
            excluded = False
            current = {'heading': text, 'heading_level': tag, 'items': []}
            sections.append(current)
            continue

        if excluded:
            continue

        if current is None:
            # heading前のコンテンツ → タイトルに含まれなかったもの
            if not sections:
                sections.append({'type': 'title', 'title': 'マニュアル', 'subtitle': ''})
            continue

        if tag == 'table':
            current['items'].append(('table', data))
        elif tag == 'img':
            current['items'].append(('img', data))
        elif tag == 'li':
            current['items'].append(('li', strip_private(data)))
        elif tag in ('h3', 'h4', 'h5', 'h6'):
            current['items'].append(('subheading', strip_private(data)))
        elif tag == 'p' and data.strip():
            current['items'].append(('p', strip_private(data)))

    return sections


def classify_section(section):
    """セクションのレイアウトパターンを判定"""
    if section.get('type') == 'title':
        return 'A'

    items = section.get('items', [])
    if not items:
        return None

    texts = [data for tag, data in items if tag == 'p']
    tables = [data for tag, data in items if tag == 'table']
    images = [data for tag, data in items if tag == 'img']
    list_items = [data for tag, data in items if tag == 'li']

    # テーブルがある → C
    if tables:
        return 'C'

    # 画像がある → F
    if images:
        return 'F'

    # key-value形式が2つ以上 → B
    kv_count = sum(1 for t in texts if KV_RE.match(t))
    if 2 <= kv_count <= 6:
        return 'B'

    # タイムライン
    timeline_count = sum(1 for t in texts if TIME_RANGE_RE.search(t))
    if timeline_count >= 2:
        return 'D'

    # アラート
    alert_count = sum(1 for t in texts if any(kw in t for kw in ALERT_KEYWORDS))
    if alert_count >= 2:
        return 'E'

    # 箇条書き
    if list_items:
        return 'G'

    # デフォルト
    return 'G'


# ── 画像 Base64 変換 ──────────────────────────

def image_to_base64(img_src, manual_dir):
    """画像パスをBase64 data URIに変換"""
    # img_srcは "images/image1.png" のような相対パス
    img_path = os.path.join(manual_dir, img_src)
    if not os.path.exists(img_path):
        # images/ プレフィックスなしでも試す
        img_path = os.path.join(manual_dir, 'images', os.path.basename(img_src))
    if not os.path.exists(img_path):
        return None

    mime, _ = mimetypes.guess_type(img_path)
    if not mime:
        mime = 'image/png'
    with open(img_path, 'rb') as f:
        b64 = base64.b64encode(f.read()).decode()
    return f'data:{mime};base64,{b64}'


# ── HTML生成 ──────────────────────────────

CSS = """
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
:root {
  --gold-dark: #C89932;
  --gold-mid: #D8B660;
  --gold-light: #FFD966;
  --warm-beige: #E0DACF;
  --brown-red: #96514D;
  --text-main: #000000;
  --text-sub: #595959;
  --white: #FFFFFF;
}
html { scroll-snap-type: y proximity; }
body {
  font-family: system-ui, -apple-system, sans-serif;
  color: var(--text-main);
  background: var(--warm-beige);
  -webkit-text-size-adjust: 100%;
  line-height: 1.6;
}
section {
  min-height: 100svh;
  padding: 6vw 5vw;
  display: flex;
  flex-direction: column;
  justify-content: center;
  scroll-snap-align: start;
  border-bottom: 1px solid rgba(0,0,0,0.08);
}

/* A: タイトル */
.slide-title {
  background: linear-gradient(135deg, var(--gold-dark), #A07828);
  color: var(--white);
  text-align: center;
  justify-content: center;
  align-items: center;
}
.slide-title h1 {
  font-size: clamp(1.6rem, 6vw, 2.4rem);
  font-weight: 800;
  margin-bottom: 0.5em;
  line-height: 1.3;
}
.slide-title .subtitle {
  font-size: clamp(0.9rem, 3vw, 1.2rem);
  opacity: 0.85;
}

/* セクション見出し */
.section-header {
  background: var(--gold-dark);
  color: var(--white);
  padding: 0.6em 1em;
  border-radius: 6px;
  font-size: clamp(1rem, 4vw, 1.3rem);
  font-weight: 700;
  margin-bottom: 1.2em;
}

/* B: 情報カード */
.info-cards {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.8em;
}
.info-card {
  background: var(--white);
  border-left: 4px solid var(--gold-dark);
  border-radius: 6px;
  padding: 0.8em;
  box-shadow: 0 1px 3px rgba(0,0,0,0.08);
}
.info-card .label {
  font-size: 0.8rem;
  color: var(--text-sub);
  margin-bottom: 0.2em;
}
.info-card .value {
  font-size: clamp(1rem, 3.5vw, 1.2rem);
  font-weight: 700;
  color: var(--text-main);
}

/* C: テーブル */
.table-wrap {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  margin: 0.8em 0;
}
.table-wrap table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}
.table-wrap th {
  background: var(--gold-mid);
  color: var(--white);
  font-weight: 700;
  padding: 0.6em 0.5em;
  text-align: center;
  white-space: nowrap;
}
.table-wrap td {
  padding: 0.5em;
  text-align: center;
  border-bottom: 1px solid rgba(0,0,0,0.08);
}
.table-wrap tr:nth-child(even) td {
  background: rgba(224,218,207,0.4);
}

/* D: タイムライン */
.timeline {
  position: relative;
  padding-left: 1.2em;
  border-left: 3px solid var(--gold-dark);
}
.timeline-item {
  margin-bottom: 1.5em;
  position: relative;
}
.timeline-item::before {
  content: '';
  position: absolute;
  left: -1.55em;
  top: 0.3em;
  width: 12px;
  height: 12px;
  background: var(--gold-dark);
  border-radius: 50%;
}
.timeline-time {
  display: inline-block;
  background: var(--gold-mid);
  color: var(--white);
  font-weight: 700;
  padding: 0.2em 0.8em;
  border-radius: 4px;
  font-size: 0.95rem;
  margin-bottom: 0.3em;
}
.timeline-desc {
  color: var(--text-main);
  font-size: 0.95rem;
  padding-left: 0.2em;
}

/* E: アラート */
.alert-item {
  display: flex;
  align-items: flex-start;
  gap: 0.8em;
  margin-bottom: 1em;
  padding: 0.8em;
  background: var(--white);
  border-radius: 6px;
  border-left: 4px solid var(--brown-red);
}
.alert-icon {
  flex-shrink: 0;
  width: 2em;
  height: 2em;
  background: var(--brown-red);
  color: var(--white);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 1.1rem;
}
.alert-text {
  font-size: 0.95rem;
  line-height: 1.6;
}

/* F: 写真ステップ */
.photo-step {
  margin-bottom: 1.5em;
}
.photo-label {
  color: var(--gold-dark);
  font-weight: 700;
  font-size: 0.95rem;
  margin-bottom: 0.4em;
}
.photo-step img {
  width: 100%;
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  display: block;
}

/* G: 箇条書き */
.bullet-list {
  list-style: none;
  padding: 0;
}
.bullet-list li {
  position: relative;
  padding-left: 1.3em;
  margin-bottom: 0.8em;
  font-size: 1rem;
  line-height: 1.8;
}
.bullet-list li::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0.55em;
  width: 8px;
  height: 8px;
  background: var(--gold-dark);
  border-radius: 50%;
}

/* 汎用テキスト */
.text-block {
  font-size: 1rem;
  line-height: 1.8;
  margin-bottom: 0.8em;
}
.text-block.sub {
  color: var(--text-sub);
  font-size: 0.9rem;
}
.subheading {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--gold-dark);
  margin: 1em 0 0.4em;
  padding-bottom: 0.2em;
  border-bottom: 2px solid var(--gold-light);
}

/* 画像（セクション内汎用） */
section img {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  display: block;
  margin: 0.8em 0;
}
"""


def render_title_section(section):
    title = section.get('title', 'マニュアル')
    subtitle = section.get('subtitle', '')
    subtitle_html = f'<p class="subtitle">{subtitle}</p>' if subtitle else ''
    return f'''<section class="slide-title">
  <h1>{title}</h1>
  {subtitle_html}
  <p class="subtitle" style="margin-top:1em;font-size:0.85rem;opacity:0.7">第44回 技大祭</p>
</section>'''


def render_info_cards(section, manual_dir):
    heading = section['heading']
    items = section['items']
    cards_html = ''
    extra_html = ''

    for tag, data in items:
        if tag == 'p':
            m = KV_RE.match(data)
            if m:
                label, value = m.group(1).strip(), m.group(2).strip()
                cards_html += f'''<div class="info-card">
  <div class="label">{label}</div>
  <div class="value">{value}</div>
</div>\n'''
            else:
                extra_html += f'<p class="text-block">{data}</p>\n'
        elif tag == 'img':
            b64 = image_to_base64(data, manual_dir)
            if b64:
                extra_html += f'<img src="{b64}" alt="">\n'
        elif tag == 'subheading':
            extra_html += f'<div class="subheading">{data}</div>\n'

    return f'''<section>
  <div class="section-header">{heading}</div>
  <div class="info-cards">{cards_html}</div>
  {extra_html}
</section>'''


def render_table_section(section, manual_dir):
    heading = section['heading']
    parts = []

    for tag, data in section['items']:
        if tag == 'table':
            rows = data
            if not rows:
                continue
            # 8行ずつ分割
            headers = rows[0]
            body_rows = rows[1:]
            for chunk_start in range(0, max(len(body_rows), 1), 8):
                chunk = body_rows[chunk_start:chunk_start + 8]
                th_html = ''.join(f'<th>{h}</th>' for h in headers)
                tr_html = ''
                for row in chunk:
                    td_html = ''.join(f'<td>{strip_private(c)}</td>' for c in row)
                    tr_html += f'<tr>{td_html}</tr>\n'
                parts.append(f'''<div class="table-wrap">
  <table><tr>{th_html}</tr>{tr_html}</table>
</div>''')
        elif tag == 'p':
            parts.append(f'<p class="text-block">{data}</p>')
        elif tag == 'img':
            b64 = image_to_base64(data, manual_dir)
            if b64:
                parts.append(f'<img src="{b64}" alt="">')
        elif tag == 'subheading':
            parts.append(f'<div class="subheading">{data}</div>')
        elif tag == 'li':
            parts.append(f'<p class="text-block">・{data}</p>')

    content = '\n'.join(parts)
    return f'''<section>
  <div class="section-header">{heading}</div>
  {content}
</section>'''


def render_timeline_section(section, manual_dir):
    heading = section['heading']
    items_html = ''
    extra_html = ''

    for tag, data in section['items']:
        if tag == 'p':
            m = TIME_RANGE_RE.search(data)
            if m:
                time_str = m.group(0)
                rest = data[m.end():].strip().lstrip('：:： ').strip()
                if not rest:
                    parts = re.split(r'[：:]', data, 1)
                    rest = parts[-1].strip() if len(parts) > 1 else ''
                items_html += f'''<div class="timeline-item">
  <span class="timeline-time">{time_str}</span>
  <p class="timeline-desc">{rest}</p>
</div>\n'''
            else:
                extra_html += f'<p class="text-block">{data}</p>\n'
        elif tag == 'img':
            b64 = image_to_base64(data, manual_dir)
            if b64:
                extra_html += f'<img src="{b64}" alt="">\n'
        elif tag == 'subheading':
            extra_html += f'<div class="subheading">{data}</div>\n'
        elif tag == 'table':
            extra_html += render_table_inline(data)
        elif tag == 'li':
            extra_html += f'<p class="text-block">・{data}</p>\n'

    return f'''<section>
  <div class="section-header">{heading}</div>
  <div class="timeline">{items_html}</div>
  {extra_html}
</section>'''


def render_alert_section(section, manual_dir):
    heading = section['heading']
    alerts_html = ''
    extra_html = ''

    for tag, data in section['items']:
        if tag == 'p':
            alerts_html += f'''<div class="alert-item">
  <div class="alert-icon">!</div>
  <div class="alert-text">{data}</div>
</div>\n'''
        elif tag == 'li':
            alerts_html += f'''<div class="alert-item">
  <div class="alert-icon">!</div>
  <div class="alert-text">{data}</div>
</div>\n'''
        elif tag == 'img':
            b64 = image_to_base64(data, manual_dir)
            if b64:
                extra_html += f'<img src="{b64}" alt="">\n'
        elif tag == 'subheading':
            extra_html += f'<div class="subheading">{data}</div>\n'

    return f'''<section>
  <div class="section-header">{heading}</div>
  {alerts_html}
  {extra_html}
</section>'''


def render_photo_section(section, manual_dir):
    heading = section['heading']
    parts = []
    img_count = 0
    current_parts = []

    for tag, data in section['items']:
        if tag == 'img':
            b64 = image_to_base64(data, manual_dir)
            if b64:
                img_count += 1
                step_label = f'Step {img_count}'
                current_parts.append(f'''<div class="photo-step">
  <div class="photo-label">{step_label}</div>
  <img src="{b64}" alt="">
</div>''')
                # 2枚ごとにセクション分割
                if img_count % 2 == 0:
                    parts.append('\n'.join(current_parts))
                    current_parts = []
        elif tag == 'p':
            current_parts.append(f'<p class="text-block">{data}</p>')
        elif tag == 'subheading':
            current_parts.append(f'<div class="subheading">{data}</div>')
        elif tag == 'li':
            current_parts.append(f'<p class="text-block">・{data}</p>')
        elif tag == 'table':
            current_parts.append(render_table_inline(data))

    if current_parts:
        parts.append('\n'.join(current_parts))

    sections_html = ''
    for i, content in enumerate(parts):
        h = heading if i == 0 else f'{heading}（続き）'
        sections_html += f'''<section>
  <div class="section-header">{h}</div>
  {content}
</section>\n'''
    return sections_html


def render_bullet_section(section, manual_dir):
    heading = section['heading']
    parts = []
    li_items = []

    for tag, data in section['items']:
        if tag == 'li':
            li_items.append(f'<li>{data}</li>')
        elif tag == 'p':
            if li_items:
                parts.append(f'<ul class="bullet-list">{"".join(li_items)}</ul>')
                li_items = []
            parts.append(f'<p class="text-block">{data}</p>')
        elif tag == 'img':
            if li_items:
                parts.append(f'<ul class="bullet-list">{"".join(li_items)}</ul>')
                li_items = []
            b64 = image_to_base64(data, manual_dir)
            if b64:
                parts.append(f'<img src="{b64}" alt="">')
        elif tag == 'subheading':
            if li_items:
                parts.append(f'<ul class="bullet-list">{"".join(li_items)}</ul>')
                li_items = []
            parts.append(f'<div class="subheading">{data}</div>')
        elif tag == 'table':
            if li_items:
                parts.append(f'<ul class="bullet-list">{"".join(li_items)}</ul>')
                li_items = []
            parts.append(render_table_inline(data))

    if li_items:
        parts.append(f'<ul class="bullet-list">{"".join(li_items)}</ul>')

    content = '\n'.join(parts)
    return f'''<section>
  <div class="section-header">{heading}</div>
  {content}
</section>'''


def render_table_inline(table_data):
    """テーブルデータをHTMLに変換（セクション外で使う用）"""
    if not table_data:
        return ''
    headers = table_data[0]
    rows = table_data[1:]
    th_html = ''.join(f'<th>{h}</th>' for h in headers)
    tr_html = ''
    for row in rows[:8]:
        td_html = ''.join(f'<td>{strip_private(c)}</td>' for c in row)
        tr_html += f'<tr>{td_html}</tr>\n'
    return f'<div class="table-wrap"><table><tr>{th_html}</tr>{tr_html}</table></div>'


RENDER_MAP = {
    'A': lambda s, d: render_title_section(s),
    'B': render_info_cards,
    'C': render_table_section,
    'D': render_timeline_section,
    'E': render_alert_section,
    'F': render_photo_section,
    'G': render_bullet_section,
}


def deduplicate_sections(sections):
    """内容が重複するセクションを除去"""
    seen = set()
    unique = []
    for s in sections:
        if s.get('type') == 'title':
            unique.append(s)
            continue
        # headingとitemsの内容でハッシュ
        heading = s.get('heading', '')
        items_key = str([(t, str(d)[:50]) for t, d in s.get('items', [])])
        key = f"{heading}|{items_key}"
        if key not in seen:
            seen.add(key)
            unique.append(s)
    return unique


# ── メイン処理 ────────────────────────────────

def convert_manual(manual_dir):
    """1つのマニュアルフォルダを変換"""
    # HTMLファイルを探す
    html_files = glob.glob(os.path.join(manual_dir, '*.html'))
    # slide.html自体は除外
    html_files = [f for f in html_files if os.path.basename(f) != 'slide.html']
    if not html_files:
        print(f"  SKIP: No HTML file found in {manual_dir}")
        return None

    html_file = html_files[0]
    manual_name = os.path.basename(manual_dir)
    output_path = os.path.join(manual_dir, 'slide.html')

    print(f"Converting: {manual_name}")
    print(f"  Source: {os.path.basename(html_file)}")

    # HTML読み取り
    with open(html_file, encoding='utf-8') as f:
        raw_html = f.read()

    # パース
    parser = GoogleDocParser()
    parser.feed(raw_html)
    parser.finish()

    img_count = sum(1 for t, _ in parser.elements if t == 'img')
    print(f"  Parsed {len(parser.elements)} elements, {img_count} images")

    # セクション構造化
    sections = build_sections(parser.elements)
    sections = deduplicate_sections(sections)
    print(f"  {len(sections)} sections")

    # パターン分類・HTML生成
    body_html = ''
    for section in sections:
        pattern = classify_section(section)
        if pattern is None:
            continue
        renderer = RENDER_MAP.get(pattern, RENDER_MAP['G'])
        body_html += renderer(section, manual_dir) + '\n'

    # タイトル取得
    page_title = manual_name
    for s in sections:
        if s.get('type') == 'title':
            page_title = s.get('title', manual_name)
            break

    # 最終HTML
    full_html = f'''<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{page_title}</title>
<style>
{CSS}
</style>
</head>
<body>
{body_html}
</body>
</html>'''

    with open(output_path, 'w', encoding='utf-8') as f:
        f.write(full_html)

    size_kb = os.path.getsize(output_path) / 1024
    print(f"  Output: slide.html ({size_kb:.0f} KB)")
    return output_path


def main():
    parser = argparse.ArgumentParser(description='Google Doc HTML → スライドHTML変換')
    parser.add_argument('path', nargs='?', help='マニュアルフォルダパス or --all')
    parser.add_argument('--all', action='store_true', help='docs/manuals/ 内の全フォルダを変換')
    args = parser.parse_args()

    if args.all or args.path == '--all':
        project_root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
        manuals_dir = os.path.join(project_root, 'docs', 'manuals')
        folders = sorted(glob.glob(os.path.join(manuals_dir, '*', '')))
        print(f"Found {len(folders)} manual folders\n")
        for folder in folders:
            folder = folder.rstrip('/')
            convert_manual(folder)
            print()
    elif args.path:
        convert_manual(args.path.rstrip('/'))
    else:
        parser.print_help()
        sys.exit(1)


if __name__ == '__main__':
    main()
