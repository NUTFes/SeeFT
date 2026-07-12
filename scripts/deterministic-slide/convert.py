"""Google Docs HTML → SeeFT card HTML 決定的変換 (生成 AI なし)

pandoc -t json の AST を読み込み、ノードごとに規則写像で SeeFT 規格の HTML を組み立てる。
同じ入力に対して常に同じ出力を返す。Claude / API キー一切不要。

使い方:
  uv run --project scripts/deterministic-slide python scripts/deterministic-slide/convert.py docs/manuals/44th_幼稚園WARSコラボブース当日マニュアル

入力:
  manual_dir/ の中の元 .html (slide*/verify* で始まらないもの)
  manual_dir/images/ 配下の画像 (base64 埋め込み)

出力:
  manual_dir/slide_det.html
"""

import argparse
import base64
import html
import json
import mimetypes
import os
import re
import subprocess
import sys
from urllib.parse import parse_qs, urlparse


SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))

OUTPUT_NAME = "slide_det.html"

# 電話番号として認識するパターン (国内向け、3-4桁 - 2-4桁 - 3-4桁)
PHONE_RE = re.compile(r"(?<![0-9])(\d{2,4}-\d{2,4}-\d{3,4})(?![0-9])")


# ---------------------------------------------------------------------------
# AST 取得 + 平坦化
# ---------------------------------------------------------------------------

def find_source_html(manual_dir: str) -> str:
    """元 HTML (slide_* / verify_* 以外の .html) を 1 つ返す。"""
    candidates = sorted(
        f for f in os.listdir(manual_dir)
        if f.endswith(".html")
        and not f.startswith("slide")
        and not f.startswith("verify")
    )
    if not candidates:
        raise FileNotFoundError(f"No source HTML found in {manual_dir}")
    if "source.html" in candidates:
        return os.path.join(manual_dir, "source.html")
    if len(candidates) == 1:
        return os.path.join(manual_dir, candidates[0])
    raise RuntimeError(
        f"Ambiguous source HTML in {manual_dir}: {candidates}. "
        "Keep exactly one source HTML (or name it source.html)."
    )


def load_ast(html_path: str) -> dict:
    """pandoc -t json で HTML を AST 化。"""
    result = subprocess.run(
        ["pandoc", "-f", "html", "-t", "json", html_path],
        capture_output=True, text=True, check=True,
    )
    return json.loads(result.stdout)


def flatten_blocks(blocks: list) -> list:
    """Google Docs HTML 由来のラッパー (Div, 単一 Header を含む OrderedList) を剥がす。

    Google Docs は番号付き見出しを <ol><li><h2>title</h2></li></ol> として出力するため、
    AST レベルでは OrderedList の中身が Header になっている。これをほどく。"""
    out: list = []
    for b in blocks:
        t = b.get("t")
        c = b.get("c")
        if t == "Div":
            # Div は [attrs, [blocks]]、属性無視で中身を平坦化
            out.extend(flatten_blocks(c[1]))
        elif t == "OrderedList":
            # OrderedList の各 item が「Header から始まる」なら、Google Docs の
            # 番号付き見出しラッパー。Header 群として展開する。
            items = c[1]
            all_lead_with_header = items and all(
                len(item) > 0 and item[0].get("t") == "Header"
                for item in items
            )
            if all_lead_with_header:
                for item in items:
                    out.extend(flatten_blocks(item))
            else:
                out.append(b)
        else:
            out.append(b)
    return out


# ---------------------------------------------------------------------------
# インライン (Str / Space / Strong / Link / Image / Span 等) の HTML 化
# ---------------------------------------------------------------------------

def unwrap_google_redirect(url: str) -> str:
    """Google Docs の `https://www.google.com/url?q=...&sa=...` から実 URL を抽出。"""
    parsed = urlparse(url)
    if parsed.netloc == "www.google.com" and parsed.path == "/url":
        qs = parse_qs(parsed.query)
        if "q" in qs:
            return qs["q"][0]
    return url


def autolink_phone(text: str) -> str:
    """電話番号らしき文字列を <a href="tel:..."> でラップする。
    HTML 化済みのテキストを受け取る前提なので、既にタグの中にある番号は触らない (簡易判定)。"""
    # 「" 又は >」の直後にある番号だけ対象にする雑な区切り。電話番号がプレーンに段落中で
    # 出現する場合をターゲットにし、既に href="tel:..." の中にあるものはスキップする。
    def repl(m: re.Match) -> str:
        num = m.group(1)
        return f'<a href="tel:{num}">{num}</a>'
    return PHONE_RE.sub(repl, text)


# RawInline (fmt=="html") で許可する生 HTML。pandoc が Google Docs HTML から実際に
# 出力する既知パターンのみ (改行 <br> と 上付き/下付き)。それ以外はエスケープする。
_RAW_INLINE_ALLOWLIST_RE = re.compile(r"</?(br|sub|sup)\s*/?>", re.IGNORECASE)


def render_inlines(inlines: list, images_b64: dict, autolink: bool = True) -> str:
    """Pandoc の inline ノード列を HTML 文字列に変換する。

    autolink: True の場合のみ Str ノードのプレーンテキストを電話番号リンク化する。
    既存の <a> の内側 (Link) や HTML 属性値 (Image alt) では False を渡し、
    タグ構造やページを破壊するリンクの二重ラップを防ぐ。
    """
    parts: list = []
    for x in inlines:
        if not isinstance(x, dict):
            continue
        t = x.get("t")
        c = x.get("c")
        if t == "Str":
            text = html.escape(c)
            parts.append(autolink_phone(text) if autolink else text)
        elif t in ("Space", "SoftBreak"):
            parts.append(" ")
        elif t == "LineBreak":
            parts.append("<br>")
        elif t == "Strong":
            parts.append(f"<strong>{render_inlines(c, images_b64, autolink)}</strong>")
        elif t == "Emph":
            parts.append(f"<em>{render_inlines(c, images_b64, autolink)}</em>")
        elif t == "Underline":
            parts.append(f"<u>{render_inlines(c, images_b64, autolink)}</u>")
        elif t == "Strikeout":
            parts.append(f"<s>{render_inlines(c, images_b64, autolink)}</s>")
        elif t == "Code":
            # Code は [attrs, str]
            parts.append(f"<code>{html.escape(c[1])}</code>")
        elif t == "Link":
            attrs, inner, target = c
            url, _title = target
            url = unwrap_google_redirect(url)
            # リンクの内側テキストは既に <a> の中なので二重にラップしない
            inner_html = render_inlines(inner, images_b64, autolink=False)
            parts.append(f'<a href="{html.escape(url)}">{inner_html}</a>')
        elif t == "Image":
            attrs, alt_inlines, target = c
            url, _title = target
            fname = os.path.basename(url)
            data_uri = images_b64.get(fname, url)
            # alt は HTML 属性値なのでタグを差し込めない
            alt = render_inlines(alt_inlines, images_b64, autolink=False)
            parts.append(
                f'<img src="{html.escape(data_uri)}" alt="{alt}" '
                f'onclick="openLightbox(this)" loading="lazy">'
            )
        elif t == "Span":
            # Google Docs CSS クラス用ラッパー。中身だけ展開
            _attrs, inner = c
            parts.append(render_inlines(inner, images_b64, autolink))
        elif t == "Quoted":
            # ["DoubleQuote" or "SingleQuote", inlines]
            quote_type, inner = c
            q = '"' if quote_type.get("t") == "DoubleQuote" else "'"
            parts.append(f"{q}{render_inlines(inner, images_b64, autolink)}{q}")
        elif t == "RawInline":
            # ["html", "<raw>"] みたいなやつ。フォーマット指定が html でも
            # allowlist (br/sub/sup) 以外は実行可能な形で通さずエスケープする
            fmt, raw = c
            if fmt == "html":
                safe = raw.strip()
                if _RAW_INLINE_ALLOWLIST_RE.fullmatch(safe):
                    parts.append(safe)
                else:
                    parts.append(html.escape(raw))
        elif t == "Note":
            # 脚注。MVP では無視
            pass
    return "".join(parts)


def inlines_to_plain(inlines: list) -> str:
    """目次用にインライン群を装飾なしテキスト化する。"""
    parts: list = []
    for x in inlines:
        if not isinstance(x, dict):
            continue
        t = x.get("t")
        c = x.get("c")
        if t == "Str":
            parts.append(c)
        elif t in ("Space", "SoftBreak", "LineBreak"):
            parts.append(" ")
        elif t in ("Strong", "Emph", "Underline", "Strikeout"):
            parts.append(inlines_to_plain(c))
        elif t == "Link":
            parts.append(inlines_to_plain(c[1]))
        elif t == "Span":
            parts.append(inlines_to_plain(c[1]))
    return "".join(parts).strip()


# ---------------------------------------------------------------------------
# ブロック (Header / Para / Lists / Table / Image) の HTML 化
# ---------------------------------------------------------------------------

def render_block(block: dict, images_b64: dict) -> str:
    """ブロック単位で HTML を組み立てる (Header 以外)。Header は呼び出し側のセクション制御で扱う。"""
    t = block.get("t")
    c = block.get("c")
    if t == "Para" or t == "Plain":
        text = render_inlines(c, images_b64)
        if not text.strip():
            return ""
        return f"<p>{text}</p>"
    if t == "BulletList":
        items = [
            f"<li>{render_blocks_inline(item, images_b64)}</li>"
            for item in c
        ]
        return f"<ul>{''.join(items)}</ul>"
    if t == "OrderedList":
        items = [
            f"<li>{render_blocks_inline(item, images_b64)}</li>"
            for item in c[1]
        ]
        return f"<ol>{''.join(items)}</ol>"
    if t == "Table":
        return render_table(c, images_b64)
    if t == "CodeBlock":
        # [attrs, code]
        return f"<pre><code>{html.escape(c[1])}</code></pre>"
    if t == "BlockQuote":
        return f"<blockquote>{render_blocks_inline(c, images_b64)}</blockquote>"
    if t == "HorizontalRule":
        return "<hr>"
    if t == "Div":
        # 残った Div は中身を展開して連結
        return render_blocks_inline(c[1], images_b64)
    if t == "Header":
        # render_blocks_inline 経由でリスト内見出しが来た場合は普通の h タグ
        level, _attrs, inlines = c
        return f"<h{level}>{render_inlines(inlines, images_b64)}</h{level}>"
    return ""


def render_blocks_inline(blocks: list, images_b64: dict) -> str:
    """リスト項目内・引用内などの「セクション分割を伴わないブロック列」を連結する。"""
    return "".join(render_block(b, images_b64) for b in flatten_blocks(blocks))


def render_table(table_c: list, images_b64: dict) -> str:
    """Pandoc Table の AST を簡易な <table> に落とす。

    Pandoc Table の構造 (pandoc-api-version 1.23):
      [attr, caption, colspecs, head, [bodies], foot]
    各 row は (attr, [cells])、各 cell は (attr, alignment, rowspan, colspan, [blocks])
    """
    # 雑にだけど構造は読み取る
    try:
        _attr, caption, _colspecs, head, bodies, foot = table_c
    except (ValueError, TypeError):
        return "<!-- table parse failed -->"

    def render_cell(cell):
        try:
            _a, _align, _rs, _cs, blocks = cell
        except (ValueError, TypeError):
            return ""
        return render_blocks_inline(blocks, images_b64)

    def render_rows(rows, cell_tag="td"):
        out = []
        for row in rows:
            try:
                _a, cells = row
            except (ValueError, TypeError):
                continue
            tds = "".join(f"<{cell_tag}>{render_cell(c)}</{cell_tag}>" for c in cells)
            out.append(f"<tr>{tds}</tr>")
        return "".join(out)

    head_rows = head[1] if isinstance(head, list) and len(head) >= 2 else []
    thead_html = ""
    if head_rows:
        thead_html = f"<thead>{render_rows(head_rows, 'th')}</thead>"

    body_html_parts = []
    for body in bodies:
        try:
            _a, _rh, _bh, body_rows = body
        except (ValueError, TypeError):
            continue
        body_html_parts.append(render_rows(body_rows))
    tbody_html = f"<tbody>{''.join(body_html_parts)}</tbody>"

    caption_inlines = caption[1] if isinstance(caption, list) and len(caption) >= 2 else []
    caption_text = ""
    if caption_inlines:
        # caption は blocks のリスト
        caption_text = f"<caption>{render_blocks_inline(caption_inlines, images_b64)}</caption>"

    return f"<table>{caption_text}{thead_html}{tbody_html}</table>"


# ---------------------------------------------------------------------------
# セクション分割 (H2 ベースで <details> カード化)
# ---------------------------------------------------------------------------

def slugify(text: str, fallback_idx: int) -> str:
    """ID 用のスラグ生成。日本語はそのまま使うと URL エンコードが面倒なので簡易番号で逃げる。"""
    safe = re.sub(r"[^\w\-]+", "-", text, flags=re.UNICODE).strip("-")
    return safe if safe else f"section-{fallback_idx}"


def split_into_sections(blocks: list, images_b64: dict) -> tuple[dict, list]:
    """flat な blocks 列から、表紙情報 + セクションリストを取り出す。

    戻り値:
      cover: {"title": str, "subtitle": str | None, "meta": list[str]}
      sections: [{"id": str, "title": str, "html": str}]

    分割ルール:
      最初の H1 (見つかれば) → cover.title
      その後の H2 (または H1 が複数) → 新セクション開始
      H3+ → 現セクション内見出し
    """
    cover = {"title": "", "subtitle": None, "meta": []}
    sections: list = []
    current: dict | None = None
    h1_seen_for_cover = False

    def push_to_current(html_fragment: str):
        nonlocal current
        if not html_fragment.strip():
            return
        if current is None:
            # H2 が来る前のリード文 → 仮想セクション "概要" を作るか、cover.meta に流す
            cover["meta"].append(html_fragment)
        else:
            current["html_parts"].append(html_fragment)

    for b in blocks:
        t = b.get("t")
        if t == "Header":
            level, _attrs, inlines = b["c"]
            title_text = inlines_to_plain(inlines)
            title_html = render_inlines(inlines, images_b64)
            if level == 1 and not h1_seen_for_cover:
                cover["title"] = title_text
                h1_seen_for_cover = True
                continue
            if level <= 2:
                # 新セクション
                if current is not None:
                    sections.append(current)
                current = {
                    "id": slugify(title_text, len(sections) + 1),
                    "title": title_text,
                    "title_html": title_html,
                    "html_parts": [],
                }
                continue
            # H3+ はセクション内見出しとして emit
            push_to_current(f"<h{level}>{title_html}</h{level}>")
            continue
        # 非 Header ブロックは現セクションへ
        push_to_current(render_block(b, images_b64))

    if current is not None:
        sections.append(current)

    # html_parts を連結して html フィールドに
    for s in sections:
        s["html"] = "\n".join(s["html_parts"])
        del s["html_parts"]

    return cover, sections


# ---------------------------------------------------------------------------
# 画像読み込み (base64)
# ---------------------------------------------------------------------------

def load_images_base64(manual_dir: str) -> dict:
    images: dict = {}
    img_dir = os.path.join(manual_dir, "images")
    if not os.path.isdir(img_dir):
        return images
    for f in sorted(os.listdir(img_dir)):
        path = os.path.join(img_dir, f)
        mime = mimetypes.guess_type(path)[0] or "image/png"
        with open(path, "rb") as fh:
            b64 = base64.b64encode(fh.read()).decode()
        images[f] = f"data:{mime};base64,{b64}"
    return images


# ---------------------------------------------------------------------------
# SeeFT テンプレート
# ---------------------------------------------------------------------------

TEMPLATE = """<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{title}</title>
<style>
:root {{
  --teal: #009688;
  --teal-medium: #4DB6AC;
  --teal-light: #CCE8E2;
  --teal-pale: #F4FBF8;
  --link: #1264A3;
  --warn: #F3AE56;
  --gold: #C9A227;
  --gray: #D9D9D9;
  --gray-light: #F5F5F5;
  --text: #222;
  --text-secondary: #666;
}}
* {{ box-sizing: border-box; }}
html {{ scroll-behavior: smooth; }}
body {{
  margin: 0; padding: 1rem 0.8rem 5rem;
  font-family: -apple-system, BlinkMacSystemFont, "Hiragino Sans", "Yu Gothic", system-ui, sans-serif;
  font-size: clamp(15px, 4vw, 16px);
  line-height: 1.6;
  color: var(--text);
  background: var(--teal-pale);
}}
h1, h2, h3, h4 {{ margin: 0.5em 0; color: var(--teal); }}
h1 {{ font-size: 18px; }}
h2 {{ font-size: 16px; }}
h3 {{ font-size: 15px; }}
h4 {{ font-size: 14px; }}
p {{ margin: 0.5em 0; }}
a {{ color: var(--link); }}

.cover {{
  background: #fff;
  border: 1px solid var(--gold);
  border-radius: 8px;
  padding: 0.85rem 1.1rem;
  margin-bottom: 1rem;
}}
.cover h1 {{ font-size: 17px; margin: 0; }}
.cover .meta {{ color: var(--text-secondary); font-size: 13px; margin-top: 0.3rem; }}

#toc {{
  background: #fff;
  border: 1px solid var(--gold);
  border-radius: 8px;
  padding: 0.8rem 1rem;
  margin-bottom: 1rem;
  counter-reset: toc;
}}
#toc h2 {{ font-size: 15px; margin: 0 0 0.6rem; padding-bottom: 0.3rem; border-bottom: 1px solid var(--teal-medium); }}
#toc ol {{ list-style: none; padding: 0; margin: 0; }}
#toc li {{ counter-increment: toc; }}
#toc li a {{
  display: flex; align-items: center;
  gap: 0.6rem;
  padding: 0.5rem 0.7rem;
  margin-bottom: 0.35rem;
  background: #fff;
  border: 1px solid var(--teal-medium);
  border-radius: 6px;
  color: var(--teal);
  font-weight: bold;
  text-decoration: none;
  transition: background 0.15s, transform 0.05s;
}}
#toc li a::before {{
  content: counter(toc);
  display: inline-flex; align-items: center; justify-content: center;
  width: 22px; height: 22px;
  border-radius: 50%;
  background: var(--teal);
  color: #fff;
  font-size: 12px;
  flex-shrink: 0;
}}
#toc li a::after {{
  content: "›";
  margin-left: auto;
  color: var(--teal);
}}
#toc li a:hover {{ background: var(--teal-light); }}
#toc li a:active {{ background: var(--teal); color: #fff; transform: scale(0.98); }}

details.section {{
  background: #fff;
  border: 1px solid var(--gold);
  border-radius: 8px;
  padding: 0.4rem 1rem;
  margin-bottom: 1rem;
}}
details.section > summary {{
  font-size: 15px; font-weight: bold;
  color: var(--teal);
  padding: 0.4rem 0;
  cursor: pointer;
  list-style: none;
}}
details.section > summary::-webkit-details-marker {{ display: none; }}
details.section > summary::after {{
  content: "▾"; float: right; color: var(--teal-medium);
}}
details.section[open] > summary::after {{ content: "▴"; }}

table {{ width: 100%; border-collapse: collapse; margin: 0.6rem 0; font-size: 13px; }}
th, td {{ padding: 0.35rem 0.5rem; border: 1px solid var(--gray); text-align: left; vertical-align: top; }}
th {{ background: var(--teal-light); color: var(--teal); }}
caption {{ caption-side: bottom; font-size: 12px; color: var(--text-secondary); padding-top: 0.3rem; }}

img {{ max-width: 100%; height: auto; border-radius: 4px; cursor: zoom-in; }}

ul, ol {{ margin: 0.5em 0; padding-left: 1.4em; }}
li {{ margin: 0.2em 0; }}

/* Lightbox */
.lightbox {{
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.9);
  display: none;
  align-items: center; justify-content: center;
  z-index: 1000;
  cursor: zoom-out;
}}
.lightbox.show {{ display: flex; }}
.lightbox img {{ max-width: 95vw; max-height: 95vh; cursor: zoom-out; }}

/* FAB + Overlay */
button.fab-toc {{
  position: fixed; bottom: 1.25rem; right: 1.25rem;
  width: 52px; height: 52px;
  border-radius: 50%;
  background: var(--teal); color: #fff;
  border: 2px solid var(--gold);
  font-size: 22px;
  cursor: pointer;
  padding: 0; font-family: inherit; outline: none;
  z-index: 100;
}}
.toc-overlay {{
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.5);
  display: none;
  align-items: flex-start; justify-content: center;
  z-index: 99;
  padding: 2rem 1rem;
  overflow-y: auto;
}}
.toc-overlay.show {{ display: flex; }}
.toc-panel {{
  background: #fff;
  border: 2px solid var(--gold);
  border-radius: 12px;
  padding: 1.25rem 1.5rem;
  max-width: 480px; width: 100%;
  position: relative;
  max-height: calc(100vh - 4rem);
  overflow-y: auto;
}}
.toc-panel h3 {{
  margin: 0 0 0.75rem;
  color: var(--teal);
  font-size: 18px;
  border-bottom: 2px solid var(--teal);
  padding-bottom: 0.4rem;
}}
.toc-close {{
  position: absolute;
  top: 0.5rem; right: 0.6rem;
  background: var(--gray-light);
  border: 1px solid var(--gray);
  border-radius: 50%;
  width: 32px; height: 32px;
  font-size: 18px;
  cursor: pointer;
}}
.toc-overlay ol {{ list-style: none; padding: 0; margin: 0; counter-reset: tocov; }}
.toc-overlay li {{ counter-increment: tocov; }}
.toc-overlay li a {{
  display: flex; align-items: center; gap: 0.6rem;
  padding: 0.5rem 0.7rem; margin-bottom: 0.35rem;
  background: #fff; border: 1px solid var(--teal-medium); border-radius: 6px;
  color: var(--teal); font-weight: bold; text-decoration: none;
}}
.toc-overlay li a::before {{
  content: counter(tocov);
  display: inline-flex; align-items: center; justify-content: center;
  width: 22px; height: 22px; border-radius: 50%;
  background: var(--teal); color: #fff; font-size: 12px; flex-shrink: 0;
}}
</style>
</head>
<body>
<div class="cover">
  <h1>{cover_title}</h1>
  {cover_meta}
</div>

<nav id="toc">
  <h2>目次</h2>
  <ol>
    {toc_items}
  </ol>
</nav>

{sections_html}

<button class="fab-toc" onclick="toggleTocOverlay()" aria-label="目次を開く" title="目次を開く">≡</button>
<div id="toc-overlay" class="toc-overlay" onclick="closeTocOverlayBackground(event)">
  <div class="toc-panel">
    <button class="toc-close" onclick="closeTocOverlay()" aria-label="閉じる">×</button>
    <h3>目次</h3>
    <ol>
      {toc_items_overlay}
    </ol>
  </div>
</div>

<div id="lightbox" class="lightbox" onclick="closeLightbox()">
  <img id="lightbox-img" alt="">
</div>

<script>
function openLightbox(el) {{
  document.getElementById('lightbox-img').src = el.src;
  document.getElementById('lightbox').classList.add('show');
}}
function closeLightbox() {{
  document.getElementById('lightbox').classList.remove('show');
}}
function toggleTocOverlay() {{
  document.getElementById('toc-overlay').classList.toggle('show');
}}
function closeTocOverlay() {{
  document.getElementById('toc-overlay').classList.remove('show');
}}
function closeTocOverlayBackground(e) {{
  if (e.target.id === 'toc-overlay') closeTocOverlay();
}}
</script>
</body>
</html>
"""


def render_full_document(cover: dict, sections: list) -> str:
    cover_meta = ""
    if cover["meta"]:
        cover_meta = '<div class="meta">' + "\n".join(cover["meta"]) + "</div>"

    toc_items = "\n".join(
        f'<li><a href="#{s["id"]}">{html.escape(s["title"])}</a></li>'
        for s in sections
    )
    toc_items_overlay = "\n".join(
        f'<li><a href="#{s["id"]}" onclick="closeTocOverlay()">{html.escape(s["title"])}</a></li>'
        for s in sections
    )
    sections_html = "\n".join(
        f'<details open class="section" id="{s["id"]}">\n'
        f'  <summary>{s["title_html"]}</summary>\n'
        f'  <div class="section-body">{s["html"]}</div>\n'
        f"</details>"
        for s in sections
    )
    return TEMPLATE.format(
        title=html.escape(cover["title"] or "マニュアル"),
        cover_title=html.escape(cover["title"] or "マニュアル"),
        cover_meta=cover_meta,
        toc_items=toc_items,
        toc_items_overlay=toc_items_overlay,
        sections_html=sections_html,
    )


# ---------------------------------------------------------------------------
# main
# ---------------------------------------------------------------------------

def main() -> int:
    parser = argparse.ArgumentParser(
        description="決定的 Google Docs HTML → SeeFT card HTML 変換器 (AI なし)",
    )
    parser.add_argument("manual_dir", help="マニュアルディレクトリ")
    parser.add_argument(
        "--output-name", default=OUTPUT_NAME,
        help=f"出力ファイル名 (default: {OUTPUT_NAME})",
    )
    args = parser.parse_args()

    manual_dir = os.path.abspath(args.manual_dir)
    output_path = os.path.join(manual_dir, args.output_name)

    print(f"=== 決定的 SeeFT HTML 生成 (AI なし) ===")
    print(f"  Source: {manual_dir}")
    print(f"  Output: {output_path}")

    src_html = find_source_html(manual_dir)
    print(f"  Source HTML: {src_html}")

    ast = load_ast(src_html)
    print(f"  AST blocks: {len(ast['blocks'])}")

    flat = flatten_blocks(ast["blocks"])
    print(f"  Flattened blocks: {len(flat)}")

    images = load_images_base64(manual_dir)
    print(f"  Images: {len(images)} files")

    cover, sections = split_into_sections(flat, images)
    print(f"  Cover title: {cover['title']!r}")
    print(f"  Sections: {len(sections)}")
    for s in sections:
        print(f"    - {s['id']}: {s['title']}")

    output_html = render_full_document(cover, sections)
    with open(output_path, "w", encoding="utf-8") as f:
        f.write(output_html)

    print(f"  Output: {output_path} ({os.path.getsize(output_path)//1024}KB)")
    print("=== 完了 ===")
    return 0


if __name__ == "__main__":
    sys.exit(main())
