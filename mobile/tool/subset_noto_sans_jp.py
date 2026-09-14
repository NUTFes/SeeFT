"""mobile に同梱する Noto Sans JP を、使う字だけに削って書き出す（#513）

背景
  Flutter 3.27.3 の Web（CanvasKit）は、Regular しかない日本語フォントに FontWeight.bold を
  指定すると疑似ボールドで描き、「閉」「開」「関」などの一部が輪郭線だけになる。
  これを避けるため Regular と Bold を同梱しているが、元のままでは 2 ファイルで 10.9MB あり、
  アプリの起動時に読み込むので遅くなる。そこで、Windows の Shift_JIS（CP932）で表せる字
  （JIS 第1・第2水準の漢字、かな、記号、英数、半角カナ、丸数字、髙・﨑などの IBM 拡張文字）だけに削る。

元フォント
  Google Fonts の Noto Sans JP（https://fonts.google.com/noto/specimen/Noto+Sans+JP）で
  「Get font」→「Download all」を押して得られる zip の、static/ にある 2 ファイルを使う。
  可変フォント（NotoSansJP-VariableFont_wght.ttf）は FontWeight で太さが変わらないので使わない。
  2026-09-13 に使ったファイルの sha256:
    NotoSansJP-Regular.ttf  d930d5d52d15231c283089760f84584272ad5e37e14607ba0d19c798e7a9caec
    NotoSansJP-Bold.ttf     c5b7b9d6a6eb682b0d4e6bbb38509575fd2759a28f147daa74714d1359a7909e

使い方（mobile/ で実行）
  python3 -m venv .venv && .venv/bin/pip install fonttools
  .venv/bin/python tool/subset_noto_sans_jp.py <zip を展開した static/ のパス>

範囲外の字（絵文字など）は、これまでどおり Flutter のフォールバックフォントで描かれる。
"""

import argparse
import hashlib
import unicodedata
from pathlib import Path

from fontTools import subset
from fontTools.ttLib import TTFont

OUT_DIR = Path(__file__).resolve().parent.parent / "lib" / "assets" / "fonts"
WEIGHTS = {"Regular": 400, "Bold": 700}


def cp932_chars() -> set[str]:
    chars = {chr(c) for c in range(0x20, 0x7F)}  # ASCII
    chars |= {bytes([b]).decode("cp932") for b in range(0xA1, 0xE0)}  # 半角カナ
    for lead in [*range(0x81, 0xA0), *range(0xE0, 0xFD)]:
        for trail in [*range(0x40, 0x7F), *range(0x80, 0xFD)]:
            try:
                ch = bytes([lead, trail]).decode("cp932")
            except UnicodeDecodeError:
                continue
            # 外字領域（私用領域）はフォントに字形が無いので除く
            if unicodedata.category(ch) != "Co":
                chars.add(ch)
    return chars


def jis_x_0208_chars() -> set[str]:
    # EUC-JP 経由の対応表は「〜」を U+301C、CP932 は U+FF5E に割り当てるなど揺れがあるので、両方を含める
    chars = set()
    for row in range(0xA1, 0xFF):
        for cell in range(0xA1, 0xFF):
            try:
                chars.add(bytes([row, cell]).decode("euc_jp"))
            except UnicodeDecodeError:
                continue
    return chars


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument("static_dir", type=Path, help="Google Fonts の zip を展開した static/ のパス")
    args = parser.parse_args()

    text = "".join(sorted(cp932_chars() | jis_x_0208_chars()))
    print(f"対象の字: {len(text)} 字")

    options = subset.Options()
    options.layout_features = ["*"]  # palt や縦書きなどの OpenType 機能は残す
    options.name_IDs = ["*"]  # 著作権表示などの name テーブルは残す（OFL の条件）
    options.name_languages = ["*"]
    options.notdef_outline = True

    OUT_DIR.mkdir(parents=True, exist_ok=True)
    for name, weight in WEIGHTS.items():
        src = args.static_dir / f"NotoSansJP-{name}.ttf"
        dst = OUT_DIR / f"NotoSansJP-{name}.ttf"
        source = TTFont(src)
        if "fvar" in source:
            raise SystemExit(f"{src} は可変フォントです。static/ の {src.name} を指定してください")
        if source["OS/2"].usWeightClass != weight:
            raise SystemExit(f"{src} のウェイトが {source['OS/2'].usWeightClass} です（期待値 {weight}）")

        font = subset.load_font(str(src), options)
        subsetter = subset.Subsetter(options)
        subsetter.populate(text=text)
        subsetter.subset(font)
        subset.save_font(font, str(dst), options)
        print(f"{name}: {src.stat().st_size:,} → {dst.stat().st_size:,} bytes")
        print(f"  入力 sha256 {sha256(src)}")
        print(f"  出力 sha256 {sha256(dst)}")


if __name__ == "__main__":
    main()
