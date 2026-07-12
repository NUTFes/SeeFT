"""解説スライドの機械的検証スクリプト（AI 呼び出しなし）

pandoc -t plain で両側（元 HTML / 生成 HTML）をテキスト化し、文字列比較で差分を
検出する。`verify_slide.py` (AI 検証版) のペアスクリプト。

設計方針:
- AI を呼ばないので決定的。同じ入力に対して常に同じ結果を返す
- 章順の組み替えは許容（substring 検索で順不同にマッチさせる）
- テキスト断片の追加・欠落・改変はすべて検出
- 「欠落」と「改変」の区別は difflib で最近傍候補を探し、類似度で判定

過検出対策（2026-05 改善）:
比較の前に「内容ではないノイズ」を両側から取り除き、純粋な本文差分だけを残す。
- (A) ナビゲーション chrome 除去: 生成 HTML 側の目次 (<nav>)・FAB ボタン・閉じる
      ボタンは、プロンプトが生成を義務付けた UI 部品であって元 Doc には存在しない。
      pandoc に渡す前に DOM ごと削除する。これで「目次見出しの重複」「<ol> の連番
      (1. 2. 3.)」が追加扱いされる誤検出が消える。
- (B) 表記ゆれの分離: 〜/~、＆/&、全角半角、括弧、CJK 間スペースだけの差は
      「改変」ではなく「表記ゆれ（許容差）」として別カテゴリに切り出す。デフォルト
      では VERDICT を NG にしない。`--strict-symbols` で従来どおり NG にできる。
- (C) 表の罫線ノイズ除去: pandoc が吐く表の罫線 (+---+, -----, |) は内容ではない。
      src/gen 双方で除去し、表組みの記法差 (グリッド表 vs シンプル表) の誤検出を防ぐ。

使い方:
  uv run --project scripts/claude-slide python scripts/claude-slide/verify_slide_mechanical.py docs/manuals/44th_幼稚園WARSコラボブース当日マニュアル

入力:
  - manual_dir 内の元 .html (slide* / verify* で始まらないもの)
  - manual_dir/slide_claude.card-strict.html (デフォルト)

出力:
  - 標準出力に検証レポート
  - manual_dir/verify_mechanical.card-strict.txt にも保存
  - 終了コード: VERDICT: OK なら 0、NG なら 1
"""

import argparse
import difflib
import os
import re
import subprocess
import sys
import unicodedata
from collections import Counter


SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))

DEFAULT_HTML_NAME = "slide_claude.card-strict.html"
DEFAULT_OUTPUT_NAME = "verify_mechanical.card-strict.txt"

# 改変判定のしきい値。difflib.SequenceMatcher.ratio() がこの値以上なら「改変」、
# 未満なら「欠落」と分類する。0.6 は経験則（短い句で誤分類しないが、軽微な書き換え
# は改変として拾える境界）。
SIMILARITY_THRESHOLD = 0.6

# 検証対象の最小文字数。これ未満の断片は誤検出が多いのでスキップする
# （例: "OK"、"完了" 等の超短文を厳密追跡すると誤検出が増える）。
DEFAULT_MIN_CHARS = 4


def resolve_manual_dir(arg: str) -> str:
    arg = arg.rstrip("/")
    if os.path.isabs(arg):
        return arg
    if os.path.isdir(arg):
        return os.path.abspath(arg)
    return os.path.join(PROJECT_ROOT, arg)


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


# 生成 HTML から除去するナビゲーション chrome のパターン。
# いずれもプロンプトが生成を義務付けた UI 部品で、元 Doc には対応物がない。
# <nav> は入れ子にならないので非貪欲マッチで安全に 1 要素ずつ消せる
# （トップ目次 id="toc"、サブ目次 class="pill-nav"、オーバーレイ目次 class="toc"
#  の 3 つの <nav> をまとめて落とす）。
_NAV_CHROME_PATTERNS = [
    re.compile(r"<nav\b[^>]*>.*?</nav>", re.DOTALL | re.IGNORECASE),
    # FAB / 閉じるボタン（≡ / ×）。class に fab-toc か toc-close を含む button。
    re.compile(
        r'<button\b[^>]*class="[^"]*(?:fab-toc|toc-close)[^"]*"[^>]*>.*?</button>',
        re.DOTALL | re.IGNORECASE,
    ),
    # 折りたたみ <details> の <summary> ラベル（「○○を表示」「表示／非表示」等）。
    # プロンプトが付ける開閉ボタンの文言で元 Doc には無い。中身(details 本体)は残す。
    re.compile(r"<summary\b[^>]*>.*?</summary>", re.DOTALL | re.IGNORECASE),
]

# <nav> 除去後にオーバーレイパネルへ残る目次見出し・装飾トークン。
# これらは単独で本文に現れることはない（目次ラベル・FAB 記号）ので最後に落とす。
_LEFTOVER_CHROME_TOKENS = re.compile(r"(?<![^\s>])(?:目次|≡|×)(?![^\s<])")

# 折りたたみカードのトグルラベル（「表示／非表示」「表示/非表示」）。UI 部品で本文ではない。
_TOGGLE_LABEL_RE = re.compile(r"表示\s*[／/]\s*非表示")


def preprocess_generated_html(html_path: str) -> str:
    """生成 HTML を pandoc に渡す前に整形したテキストを返す（元ファイルは不変）。
    - base64 画像を [image] に置換（巨大 data URI で pandoc が死ぬのを防ぐ）
    - ナビゲーション chrome（目次 <nav>・FAB・閉じるボタン）を DOM ごと削除
    """
    with open(html_path, "r", encoding="utf-8") as f:
        html = f.read()
    html = re.sub(r'src="data:[^"]*"', 'src="[image]"', html)
    html = re.sub(r"src='data:[^']*'", "src='[image]'", html)
    for pat in _NAV_CHROME_PATTERNS:
        html = pat.sub(" ", html)
    return html


def strip_table_art(text: str) -> str:
    """pandoc plain が吐く表の罫線・セル区切りを除去する（内容は残す）。
    src 側のグリッド表 (+---+) と gen 側のシンプル表 (-----) の記法差を吸収する。
    改行を潰す前（行指向のまま）に呼ぶこと。"""
    # 罫線だけの行を消す（- + | : = と box-drawing と空白のみで構成される行）
    text = re.sub(r"(?m)^[\s|+:=\-─━┼┃╋╂│┌┐└┘├┤┬┴]{2,}$", "", text)
    # 行頭のキャプション記号 ": 表 ..." の先頭コロン、行頭セル区切り "| "
    text = re.sub(r"(?m)^\s*[:|]\s*", "", text)
    # インラインに残る罫線ラン（3 連以上のダッシュ、box-drawing 連、+--+ グリッド）
    text = re.sub(r"[─━]{2,}", " ", text)
    text = re.sub(r"-{3,}", " ", text)
    text = re.sub(r"\+[-+]{2,}\+?", " ", text)
    # セル区切りパイプを空白化（電話番号等の単独ダッシュは -{3,} に当たらず温存される）
    text = text.replace("|", " ")
    return text


def fold_cosmetic(s: str) -> str:
    """表記ゆれ判定用の正規化。これで一致するペアは「改変」でなく「表記ゆれ」扱い。
    - NFKC で全角→半角（＆→&、（→(、：→:、全角数字→半角 等）
    - 波ダッシュ系を ~ に統一（NFKC で畳まれない 〜 U+301C を明示的に処理）
    - すべての空白を除去（CJK 間の余分なスペース差・改行差を無視）
    内容（可読文字）の追加・削除・置換は畳まれないので、真の改変は検出され続ける。"""
    s = unicodedata.normalize("NFKC", s)
    s = s.replace("〜", "~").replace("～", "~").replace("∼", "~")
    # 箇条書き記号の差（元 Doc の「・」↔ 生成の「-」「‐」「•」等）は内容ではなく表記。
    # src/gen 双方に対称適用するので、ダッシュだけ違う文以外を誤って一致させることはない。
    # カタカナ長音符「ー」は実単語の一部なので除外する。
    s = re.sub(r"[・•·‐―−\-]", "", s)
    s = re.sub(r"\s+", "", s)
    return s


def html_to_plain(html_path: str, preprocess: bool = False) -> str:
    """pandoc -t plain で HTML をプレーンテキスト化する。
    preprocess=True なら base64 とナビ chrome を剥がしてから渡す（生成 HTML 用）。"""
    if preprocess:
        cleaned = preprocess_generated_html(html_path)
        result = subprocess.run(
            ["pandoc", "-f", "html", "-t", "plain", "--wrap=none"],
            input=cleaned,
            capture_output=True,
            text=True,
            check=True,
        )
    else:
        result = subprocess.run(
            ["pandoc", "-f", "html", "-t", "plain", "--wrap=none", html_path],
            capture_output=True,
            text=True,
            check=True,
        )
    out = result.stdout
    # 表罫線ノイズは両側で除去（行指向のまま、改行を潰す normalize の前段）
    out = strip_table_art(out)
    if preprocess:
        # <nav> 除去後にパネルへ残った目次見出し等の単独トークンを掃除
        out = _LEFTOVER_CHROME_TOKENS.sub(" ", out)
        # 折りたたみカードのトグルラベル「表示／非表示」はプロンプトが付ける UI 部品で
        # 元 Doc には無い。本文比較の前に落として誤検出（改変・追加）を防ぐ。
        out = _TOGGLE_LABEL_RE.sub(" ", out)
    return out


def normalize(s: str) -> str:
    """空白の正規化のみ。文字種・全半角の変換はしない（不変ポリシーのため）。
    改行も空白に統一する。これにより、生成側で文の途中に改行が挿入されても
    文単位の比較がズレない（src/gen 双方を同じ手続きで文リスト化するため、
    比較粒度が対称になる）。"""
    s = s.replace(" ", " ")  # NBSP → ASCII space
    s = s.replace("　", " ")  # 全角スペース → ASCII space（pandoc plain で残ることがある）
    s = re.sub(r"[ \t\r\n]+", " ", s)
    return s.strip()


def split_sentences(text: str) -> list[str]:
    """テキストを「文」に分割する。文末記号のみで区切る。
    normalize() が改行を空白化しているので、ここでは改行は出てこない前提。
    終端記号を持たない箇条書き等は 1 つのまとまりとして扱われる。"""
    parts: list[str] = []
    for chunk in re.split(r"(?<=[。．\.！\!？\?])", text):
        c = chunk.strip()
        if c:
            parts.append(c)
    return parts


def filter_meaningful(sentences: list[str], min_chars: int) -> list[str]:
    """超短い断片や箇条書きの装飾だけの行を除外する。"""
    out = []
    for s in sentences:
        if len(s) < min_chars:
            continue
        # 記号と空白だけの行を除外
        if re.fullmatch(r"[\s\-\*\+\・\〇\○\●\▼\◆\◇\■\□\【\】\「\」\:：\(\)\（\）]+", s):
            continue
        out.append(s)
    return out


def _nospace_len(s: str) -> int:
    """空白を除いた可読文字数。一致度は「文字数」で測るので空白は数えない。"""
    return len(re.sub(r"\s+", "", s))


def _matched_chars(needle: str, haystack: str, min_block: int = 2) -> int:
    """needle の可読文字のうち、haystack 内に（順不同の位置でも）連続で現れる文字数。
    min_block 以上の連続一致ブロックだけ数えて、'の' 'を' 等の単発偶然一致の水増しを避ける。"""
    sm = difflib.SequenceMatcher(None, needle, haystack, autojunk=False)
    return sum(b.size for b in sm.get_matching_blocks() if b.size >= min_block)


def _unmatched_spans(
    needle: str, haystack: str, min_match: int = 2, min_span: int = 2
) -> list[str]:
    """needle のうち haystack に現れなかった連続部分（span）を返す。
    『AI が足した実テキスト』『元Doc から消えた実テキスト』を取り出すのに使う。
    min_match 未満の偶然一致は「一致」と見なさず span をつなげる。"""
    sm = difflib.SequenceMatcher(None, needle, haystack, autojunk=False)
    blocks = [b for b in sm.get_matching_blocks() if b.size >= min_match]
    spans: list[str] = []
    i = 0
    for b in blocks:
        if b.a > i and len(needle[i:b.a]) >= min_span:
            spans.append(needle[i:b.a])
        i = max(i, b.a + b.size)
    if len(needle) - i >= min_span:
        spans.append(needle[i:])
    return spans


# 表示前に span の前後から削るノイズ記号（画像プレースホルダ [] や箇条書き記号など）
_SPAN_TRIM = " 　.・|[]()（）【】●○◯-—ー　"
# 「意味のある可読語」を含むか: 漢字・かな・カナ、または 2 文字以上の英字
_MEANINGFUL_RE = re.compile(r"[一-龥぀-ゟ゠-ヿ々〆ヶ]|[A-Za-z]{2,}")


def _summarize_spans(spans: list[str], max_items: int = 20, max_len: int = 24) -> str:
    """span のリストを読みやすい 1 行に。
    前後のブラケット等を削り、意味のある語を含む span だけに絞り、重複除去・件数打ち切り。
    （'1][' '[][][' のような記号断片は除外する。文字数カウントには含むが内訳には出さない）"""
    seen: dict[str, None] = {}
    for s in spans:
        cleaned = s.strip(_SPAN_TRIM)
        if cleaned and _MEANINGFUL_RE.search(cleaned):
            seen.setdefault(cleaned, None)
    uniq = list(seen.keys())
    shown = [s if len(s) <= max_len else s[:max_len] + "…" for s in uniq[:max_items]]
    body = " / ".join(shown) if shown else "（記号・画像枠のみ）"
    if len(uniq) > max_items:
        body += f" …他 {len(uniq) - max_items} 件"
    return body


def compute_fidelity(
    src_sentences: list[str],
    gen_sentences: list[str],
    additions: list[dict],
    omissions: list[dict],
    modifications: list[dict],
    cosmetic: list[dict],
) -> dict:
    """「実際にどれだけ文章が変わったか」を文字ベースで集計する。

    件数（改変N件…）はノイズで水増しされるので、第三者に渡す指標としては
    『元Doc の可読文字のうち何 % が生成 HTML に残っているか（保持率）』と
    『AI が新たに足した可読文字数』の方が実態を表す。

    保持率は文ペアリングに依存させない（依存させると、表ブロックがペアに失敗した
    だけで「消失」と誤算定され、保持率が不当に下がる）。代わりに、元Doc の各文を
    生成 HTML 全体に対して照合し、どこかに（章順を入れ替えても）現れる文字を保持と
    数える。これで章の組み替えに強く、ペアリングの粗さに引きずられない。
    """
    src_join = [re.sub(r"\s+", "", s) for s in src_sentences]
    gen_join = [re.sub(r"\s+", "", s) for s in gen_sentences]
    src_all = "".join(src_join)
    gen_all = "".join(gen_join)
    total_src = len(src_all)
    total_gen = len(gen_all)

    # 元Doc の各文 → 生成 HTML 全体に現れる文字数（保持）
    matched_src = sum(_matched_chars(s, gen_all) for s in src_join)
    matched_src = min(matched_src, total_src)
    # 生成 HTML の各文 → 元Doc 全体に現れる文字数（残りが「追加」）
    matched_gen = sum(_matched_chars(g, src_all) for g in gen_join)
    matched_gen = min(matched_gen, total_gen)

    coverage = (matched_src / total_src) if total_src else 1.0
    added_chars = total_gen - matched_gen
    lost_chars = total_src - matched_src

    # 実テキストの内訳: 何が足され／何が消えたか（数字だけでは判断できないため）
    added_spans: list[str] = []
    for g in gen_join:
        added_spans.extend(_unmatched_spans(g, src_all))
    dropped_spans: list[str] = []
    for s in src_join:
        dropped_spans.extend(_unmatched_spans(s, gen_all))

    total_sent = len(src_sentences)
    exact_sent = total_sent - (len(omissions) + len(modifications) + len(cosmetic))

    return {
        "total_src": total_src,
        "total_gen": total_gen,
        "preserved": matched_src,
        "coverage": coverage,
        "added_chars": added_chars,
        "lost_chars": lost_chars,
        "added_spans": added_spans,
        "dropped_spans": dropped_spans,
        "exact_sent": exact_sent,
        "total_sent": total_sent,
    }


def find_closest(query: str, candidates: list[str]) -> tuple[str | None, float]:
    """difflib で最も近い候補と類似度を返す。"""
    if not candidates:
        return None, 0.0
    best_match = None
    best_ratio = 0.0
    for c in candidates:
        ratio = difflib.SequenceMatcher(None, query, c).ratio()
        if ratio > best_ratio:
            best_ratio = ratio
            best_match = c
    return best_match, best_ratio


def categorize(
    src_sentences: list[str],
    gen_sentences: list[str],
) -> tuple[list[dict], list[dict], list[dict], list[dict]]:
    """source / generated を文リスト同士の多重集合 (Counter) で比較し、
    追加・欠落・改変・表記ゆれ の 4 リストを返す。

    比較粒度を双方とも「split_sentences を通った文」に統一することで、
    片側だけ改行が入った文の誤検知 (欠落として報告されるケース) を抑える。
    Counter による多重集合差を取るので、同一文が複数回出る場合の片方欠落も検出できる。

    ペアリングは 2 段階:
      pass1) fold_cosmetic() で一致するペア → 表記ゆれ（許容差）として確保
      pass2) 残りを類似度でペア → 改変、付かなければ欠落／追加
    pass1 を先に通すことで、〜/~・全半角・空白だけの差が「改変」に混ざらない。
    """
    src_counter = Counter(src_sentences)
    gen_counter = Counter(gen_sentences)

    # 多重集合差: 要素ごとに max(0, count_src - count_gen)
    missing_in_gen = list((src_counter - gen_counter).elements())
    extra_in_gen = list((gen_counter - src_counter).elements())

    additions: list[dict] = []
    omissions: list[dict] = []
    modifications: list[dict] = []
    cosmetic: list[dict] = []

    remaining_missing = list(missing_in_gen)
    remaining_extras = list(extra_in_gen)

    # pass1: 表記ゆれペア（fold_cosmetic 一致）を先に回収する。
    # extras を畳んだキーで引けるよう索引を作る（同一キー複数コピーに対応）。
    extras_by_fold: dict[str, list[str]] = {}
    for e in remaining_extras:
        extras_by_fold.setdefault(fold_cosmetic(e), []).append(e)

    still_missing: list[str] = []
    for m in remaining_missing:
        bucket = extras_by_fold.get(fold_cosmetic(m))
        if bucket:
            gen_match = bucket.pop(0)
            remaining_extras.remove(gen_match)
            cosmetic.append({"src": m, "gen": gen_match})
        else:
            still_missing.append(m)

    # pass2: 残りを類似度でペア → 改変。付かなければ欠落。
    for m in still_missing:
        candidate, ratio = find_closest(m, remaining_extras)
        if candidate is not None and ratio >= SIMILARITY_THRESHOLD:
            modifications.append({
                "src": m,
                "gen": candidate,
                "ratio": ratio,
            })
            remaining_extras.remove(candidate)  # 1 件だけ消費
        else:
            omissions.append({"src": m, "closest": candidate, "ratio": ratio})

    # ペアが付かなかった extras は純粋な追加
    for e in remaining_extras:
        additions.append({"gen": e})

    return additions, omissions, modifications, cosmetic


def _inline_diff(src: str, gen: str) -> str:
    """文字レベルの差分を [+追加] [-削除] 表記で 1 行にする。"""
    d = list(difflib.ndiff(src, gen))
    return "".join(c[2:] if c.startswith("  ") else f"[{c[0]}{c[2:]}]" for c in d)


# 部門長向け分類で使う「自動整形トークン」のパターン。
_SECTION_NUM_RE = re.compile(r"^\s*\d+\s+")  # 見出し先頭に付いた章番号「2 手順」
_FIGURE_ONLY_RE = re.compile(r"^[\s\[\]]*(?:図|表|画像)\s*\d*\s*[\.．]?\s*$")  # 「図1.」「[] 図3.」
# 文末に残るリスト連番・図表番号・目次のページ番号。pandoc plain が順序付きリスト・
# キャプション・目次（「見出し 3 2.」= 見出し+ページ+連番）の数字を文末に落とす副作用で、
# 本文ではなく採番／ナビの差。連続する数字トークン列をまとめて畳む。
_TRAILING_NUM_RE = re.compile(r"(?:\s*(?:図|表|画像)?\s*\d+\s*[\.．]?)+\s*$")

# 日本語（ひらがな・カタカナ・漢字）。マニュアル本文は日本語なので、これを一切含まない
# 断片（Windows の画像パス `C:\Users\...\INetCache` 等、元Doc由来のゴミ）は本文ではない。
_CJK_RE = re.compile(r"[぀-ヿ㐀-鿿]")
# キャプション標識「図 」「表 」「画像 」（番号なしの空白付き）。図のラベルを示すだけで本文ではない。
# 「図4」のように数字直結のインライン参照はマッチしない（本文として温存する）。
_CAPTION_MARK_RE = re.compile(r"(?:図|表|画像)\s+")
# 元Docから漏れた画像ファイル名の残骸「jpg]」「JPG]」「png]」等。
_IMG_EXT_JUNK_RE = re.compile(r"(?:jpe?g|png|gif)\s*\]?", re.IGNORECASE)


def _has_cjk(s: str) -> bool:
    return bool(_CJK_RE.search(s))


def _strip_structure(s: str) -> str:
    """自動整形で入る構造トークン（章番号・角括弧・トグル語・末尾連番・キャプション標識・
    画像ファイル名残骸）を畳んだ正規化キー。これで src と gen が一致／一方が他方に包含される
    なら、差は「本文」ではなく「整形」だけ。fold_cosmetic で NFKC・波ダッシュ・空白も吸収する。"""
    s = _TOGGLE_LABEL_RE.sub("", s)
    s = _SECTION_NUM_RE.sub("", s)
    s = _CAPTION_MARK_RE.sub("", s)
    s = _IMG_EXT_JUNK_RE.sub("", s)
    s = _TRAILING_NUM_RE.sub("", s)
    s = s.replace("[", "").replace("]", "")
    return fold_cosmetic(s)


def classify_for_buncho(
    additions: list[dict],
    omissions: list[dict],
    modifications: list[dict],
    src_sentences: list[str],
    gen_sentences: list[str],
) -> dict:
    """追加・欠落・改変を「本文差（要確認）」と「整形差／図表差（確認不要）」に振り分ける。

    判定の核は 2 つ:
      1) 構造トークンを畳んだキー (_strip_structure) で src と gen が一致 → 整形差
      2) 整形キーが反対側の本文全体に substring で存在 → 並べ替え（実在する）= 整形差
    どちらにも当たらない欠落・追加・改変だけが「本文差」として残る。
    """
    src_blob = _strip_structure(" ".join(src_sentences))
    gen_blob = _strip_structure(" ".join(gen_sentences))

    content_omissions: list[dict] = []
    content_additions: list[dict] = []
    content_modifications: list[dict] = []
    format_notes: list[str] = []
    figure_notes: list[str] = []

    for o in omissions:
        src = o["src"]
        if _FIGURE_ONLY_RE.match(src):
            figure_notes.append(src)  # 「図1.」等のキャプションのみ
            continue
        if not _has_cjk(src):
            format_notes.append(src)  # 元Doc由来のゴミ（画像パス等）。本文ではない
            continue
        key = _strip_structure(src)
        if key and key in gen_blob:
            format_notes.append(src)  # 整形・並べ替えで実は生成HTMLに存在
            continue
        content_omissions.append(o)

    for a in additions:
        gen = a["gen"]
        if _FIGURE_ONLY_RE.match(gen):
            figure_notes.append(gen)
            continue
        if not _has_cjk(gen):
            format_notes.append(gen)
            continue
        key = _strip_structure(gen)
        if key and key in src_blob:
            format_notes.append(gen)  # 元Docに存在する文の並べ替え（捏造ではない）
            continue
        content_additions.append(a)

    for m in modifications:
        src_key = _strip_structure(m["src"])
        # 整形キーが一致、または元の語句が生成HTML全体に実在（並べ替え）なら整形差。
        if src_key == _strip_structure(m["gen"]) or (src_key and src_key in gen_blob):
            format_notes.append(m["gen"])  # 章番号付与・角括弧化・キャプション並べ替え等
            continue
        content_modifications.append(m)

    return {
        "content_omissions": content_omissions,
        "content_additions": content_additions,
        "content_modifications": content_modifications,
        "format_notes": format_notes,
        "figure_notes": figure_notes,
    }


def render_buncho_report(
    buckets: dict,
    fidelity: dict | None,
) -> str:
    """部門長（非エンジニア）向けの文章チェック結果。Slack スレッドにそのまま貼れる散文。
    デザインは対象外、文章の欠落・改変・追加だけを日本語で提示する。"""
    omissions = buckets["content_omissions"]
    additions = buckets["content_additions"]
    modifications = buckets["content_modifications"]
    n_real = len(omissions) + len(additions) + len(modifications)
    n_format = len(buckets["format_notes"])
    n_figure = len(buckets["figure_notes"])

    lines: list[str] = []
    lines.append("解説マニュアル 文章チェック結果")
    lines.append(
        "（元の Google ドキュメントの「文章」が解説HTMLに正しく引き継がれているかの"
        "自動チェックです。デザイン・レイアウトは対象外です。）"
    )
    lines.append("")

    lines.append("■ 全体")
    if fidelity is not None:
        lines.append(
            f"元の文章は約 {fidelity['coverage']*100:.0f}% がそのまま引き継がれています。"
        )
    if n_real == 0:
        lines.append("文章の欠落・書き換え・追加は見つかりませんでした。安心して公開できます。")
    else:
        lines.append(
            f"確認してほしい本文差は {n_real} 件です"
            f"（消えたかも {len(omissions)} / 書き換わったかも {len(modifications)} / 増えたかも {len(additions)}）。"
        )
    lines.append(
        f"このほか、見出し番号やレイアウト上の差が {n_format} 件、"
        f"図番号の整理が {n_figure} 件ありますが、内容は変わっていないため確認は不要です。"
    )
    lines.append("")

    if n_real > 0:
        lines.append("■ 確認してほしい本文差")
        if omissions:
            lines.append("▼ 消えたかもしれない文（元にあって解説HTMLに見当たらない）")
            for o in omissions:
                lines.append(f"・元の文: 「{o['src']}」")
            lines.append("")
        if modifications:
            lines.append("▼ 書き換わったかもしれない文")
            for m in modifications:
                lines.append(f"・元　: 「{m['src']}」")
                lines.append(f"  生成: 「{m['gen']}」")
            lines.append("")
        if additions:
            lines.append("▼ 増えたかもしれない文（元になかった文）")
            for a in additions:
                lines.append(f"・生成: 「{a['gen']}」")
            lines.append("")
    else:
        lines.append("■ 確認してほしい本文差")
        lines.append("なし")
        lines.append("")

    lines.append("■ 確認不要の差（参考）")
    lines.append(
        "見出しに番号が付いた／キャプションが […] で囲まれた／折りたたみの「表示／非表示」ボタンが"
        "付いた、などレイアウト上の差です。本文の内容は変わっていません。"
    )
    lines.append(
        "「図1」「図3」などの図番号ラベルは整理されていますが、画像そのものは解説HTMLに含まれています。"
    )

    return "\n".join(lines)


def render_report(
    additions: list[dict],
    omissions: list[dict],
    modifications: list[dict],
    cosmetic: list[dict],
    fidelity: dict | None = None,
    strict_symbols: bool = False,
) -> str:
    # 内容差分（追加・欠落・改変）があれば常に NG。
    # 表記ゆれは strict_symbols のときだけ NG に算入する。
    is_ng = bool(additions or omissions or modifications)
    if strict_symbols and cosmetic:
        is_ng = True

    lines: list[str] = []
    lines.append("VERDICT: NG" if is_ng else "VERDICT: OK")
    lines.append("")

    # 冒頭サマリ: 「件数」ではなく「文章がどれだけ変わったか」を文字ベースで示す。
    # 第三者・AI に結果を渡すとき、忠実度がこの 1 ブロックで伝わる。
    if fidelity is not None:
        lines.append("## サマリ（文章の一致度・文字ベース）")
        lines.append(
            f"- 本文保持率: {fidelity['coverage']*100:.1f}% "
            f"（元Docの可読文字 {fidelity['total_src']:,} 字のうち "
            f"{fidelity['preserved']:,} 字が生成HTMLに残存）"
        )
        lines.append(
            f"- AIが追加した可読文字: {fidelity['added_chars']:,} 字"
            f"（生成HTML 全体 {fidelity['total_gen']:,} 字）"
        )
        lines.append(f"  内訳: {_summarize_spans(fidelity['added_spans'])}")
        lines.append(
            f"- AIが落とした/書き換えた可読文字: {fidelity['lost_chars']:,} 字"
        )
        lines.append(f"  内訳: {_summarize_spans(fidelity['dropped_spans'])}")
        lines.append(
            f"- 完全一致した文: {fidelity['exact_sent']}/{fidelity['total_sent']} 文"
        )
        lines.append(
            "# 注: 件数（下記）は画像alt移動・表組み替え等で水増しされる。"
            "文章改変の大きさはこの保持率で判断する。"
        )
        lines.append("")

    lines.append("## 追加された情報（HTML にあって元 Markdown にない）")
    if not additions:
        lines.append("- なし")
    else:
        for a in additions:
            lines.append(f"- HTML: {a['gen']!r}")
    lines.append("")

    lines.append("## 欠落した情報（元 Markdown にあって HTML にない）")
    if not omissions:
        lines.append("- なし")
    else:
        for o in omissions:
            lines.append(f"- 元 Markdown: {o['src']!r}")
            if o["closest"] is not None:
                lines.append(
                    f"  (最も近い HTML 文: {o['closest']!r}, 類似度 {o['ratio']:.2f})"
                )
    lines.append("")

    lines.append("## 改変された情報（数値・固有名詞・内容の差異）")
    if not modifications:
        lines.append("- なし")
    else:
        for m in modifications:
            lines.append(f"- 類似度 {m['ratio']:.2f}")
            lines.append(f"  元 Markdown: {m['src']!r}")
            lines.append(f"  生成 HTML : {m['gen']!r}")
            lines.append(f"  差分    : {_inline_diff(m['src'], m['gen'])}")
    lines.append("")

    note = "（VERDICT に算入）" if strict_symbols else "（許容差・VERDICT に非算入）"
    lines.append(f"## 表記ゆれ {note}")
    lines.append("# 〜/~・全角半角・括弧・空白だけの差。可読文字の内容は同一。")
    if not cosmetic:
        lines.append("- なし")
    else:
        for c in cosmetic:
            lines.append(f"- 元 Markdown: {c['src']!r}")
            lines.append(f"  生成 HTML : {c['gen']!r}")
            lines.append(f"  差分    : {_inline_diff(c['src'], c['gen'])}")
    lines.append("")

    lines.append("## 統計")
    lines.append(f"- 追加: {len(additions)} 件")
    lines.append(f"- 欠落: {len(omissions)} 件")
    lines.append(f"- 改変: {len(modifications)} 件")
    lines.append(f"- 表記ゆれ: {len(cosmetic)} 件")

    return "\n".join(lines)


def is_ok(report: str) -> bool:
    stripped = report.strip()
    if not stripped:
        return False
    first_line = stripped.splitlines()[0].strip()
    return first_line.startswith("VERDICT: OK")


def main() -> int:
    parser = argparse.ArgumentParser(
        description="解説スライド 機械的検証（AI 呼び出しなし）",
    )
    parser.add_argument(
        "manual_dir",
        help="マニュアルディレクトリ",
    )
    parser.add_argument(
        "--html-name",
        default=DEFAULT_HTML_NAME,
        help=f"検証対象の生成 HTML ファイル名 (default: {DEFAULT_HTML_NAME})",
    )
    parser.add_argument(
        "--output-name",
        default=DEFAULT_OUTPUT_NAME,
        help=f"検証レポートの出力ファイル名 (default: {DEFAULT_OUTPUT_NAME})",
    )
    parser.add_argument(
        "--min-chars",
        type=int,
        default=DEFAULT_MIN_CHARS,
        help="検出対象とする文の最小文字数（短い断片は誤検出が多いので除外）",
    )
    parser.add_argument(
        "--strict-symbols",
        action="store_true",
        help="表記ゆれ（〜/~・全半角・括弧・空白差）も VERDICT を NG にする（厳格）",
    )
    parser.add_argument(
        "--report-name",
        default="verify_report.card-strict.md",
        help="部門長向け文章チェック結果の出力ファイル名（Slack 貼り付け用）",
    )
    args = parser.parse_args()

    manual_dir = resolve_manual_dir(args.manual_dir)
    output_path = os.path.join(manual_dir, args.output_name)
    gen_path = os.path.join(manual_dir, args.html_name)

    print("=== 機械検証（決定的、AI 呼び出しなし）===")
    print(f"  Source dir: {manual_dir}")
    print(f"  Target HTML: {args.html_name}")
    print(f"  Output: {output_path}")
    print(f"  Min chars: {args.min_chars}")

    if not os.path.isfile(gen_path):
        print(f"ERROR: Generated HTML not found: {gen_path}", file=sys.stderr)
        return 2

    src_html_path = find_source_html(manual_dir)
    print(f"  Source HTML: {src_html_path}")

    src_text = normalize(html_to_plain(src_html_path, preprocess=False))
    gen_text = normalize(html_to_plain(gen_path, preprocess=True))

    print(f"  Source text: {len(src_text)} chars")
    print(f"  Generated text: {len(gen_text)} chars")

    src_sentences = filter_meaningful(split_sentences(src_text), args.min_chars)
    gen_sentences = filter_meaningful(split_sentences(gen_text), args.min_chars)

    print(f"  Sentences: src={len(src_sentences)}, gen={len(gen_sentences)}")

    additions, omissions, modifications, cosmetic = categorize(
        src_sentences, gen_sentences
    )

    fidelity = compute_fidelity(
        src_sentences, gen_sentences, additions, omissions, modifications, cosmetic
    )

    report = render_report(
        additions, omissions, modifications, cosmetic, fidelity, args.strict_symbols
    )
    with open(output_path, "w", encoding="utf-8") as f:
        f.write(report + "\n")

    # 部門長向け（Slack 貼り付け用）: ノイズを畳んで「本文差」だけを日本語で提示する。
    buckets = classify_for_buncho(
        additions, omissions, modifications, src_sentences, gen_sentences
    )
    buncho_report = render_buncho_report(buckets, fidelity)
    report_path = os.path.join(manual_dir, args.report_name)
    with open(report_path, "w", encoding="utf-8") as f:
        f.write(buncho_report + "\n")

    print()
    print("--- 検証レポート（開発用・詳細）---")
    print(report)
    print("--- ここまで ---")
    print()
    print("--- 部門長向けレポート（Slack 用）---")
    print(buncho_report)
    print("--- ここまで ---")
    print()
    print(f"  詳細: {output_path} ({os.path.getsize(output_path)//1024}KB)")
    print(f"  部門長向け: {report_path}")

    if is_ok(report):
        print("=== VERDICT: OK (元 Doc と HTML は機械的に一致) ===")
        return 0
    print("=== VERDICT: NG (差分あり、レポート参照) ===")
    return 1


if __name__ == "__main__":
    sys.exit(main())
