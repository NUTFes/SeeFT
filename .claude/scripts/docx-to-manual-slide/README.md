# docx → マニュアルスライド変換 リファレンス

SKILL.mdの生成ルールを補足する詳細リファレンスと検証ガイド。

---

## docx読み取り手順

### 1. テキスト・構造の抽出

```python
import docx

doc = docx.Document("対象.docx")
for p in doc.paragraphs:
    # p.style.name → "Title", "Heading 1", "Heading 2", "normal" 等
    # p.text → テキスト内容
```

- `Title` → タイトルスライド（パターンA）の見出し
- `Heading 1` → スライドグループの区切り
- `Heading 2`〜 → スライド内のセクション区切り
- `normal` → 本文。前後の文脈でパターンB〜Gを判定

### 2. テーブルの抽出

```python
for table in doc.tables:
    for row in table.rows:
        cells = [cell.text for cell in row.cells]
```

### 3. 画像の抽出

```python
import zipfile, os

IMAGES_DIR = "scripts/manuals/01_配線/images"  # マニュアルごとのimagesディレクトリ
os.makedirs(IMAGES_DIR, exist_ok=True)

with zipfile.ZipFile("対象.docx", 'r') as z:
    for name in z.namelist():
        if name.startswith('word/media/'):
            fname = os.path.basename(name)
            with open(os.path.join(IMAGES_DIR, fname), 'wb') as f:
                f.write(z.read(name))
```

### 4. 画像とパラグラフの対応付け

```python
from lxml import etree

nsmap = {
    'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main',
    'r': 'http://schemas.openxmlformats.org/officeDocument/2006/relationships',
    'a': 'http://schemas.openxmlformats.org/drawingml/2006/main',
}

for i, para in enumerate(doc.paragraphs):
    drawings = para._element.findall('.//w:drawing', nsmap)
    for d in drawings:
        blips = d.findall('.//a:blip', nsmap)
        for blip in blips:
            rId = blip.get('{http://schemas.openxmlformats.org/officeDocument/2006/relationships}embed')
            if rId and rId in doc.part.rels:
                print(f"Para[{i}] -> {doc.part.rels[rId].target_ref}")
```

---

## python-pptx コードテンプレート

### 共通定数・初期化（SKILL.md準拠）

```python
from pptx import Presentation
from pptx.util import Inches, Pt
from pptx.dml.color import RGBColor
from pptx.enum.text import PP_ALIGN, MSO_ANCHOR
from pptx.enum.shapes import MSO_SHAPE

# カラー定数（SKILL.md準拠）
DARK_BLUE = RGBColor(0x1B, 0x3A, 0x5C)
ACCENT_BLUE = RGBColor(0x2E, 0x75, 0xB6)
LIGHT_BG = RGBColor(0xF2, 0xF2, 0xF2)
WHITE = RGBColor(0xFF, 0xFF, 0xFF)
ORANGE = RGBColor(0xED, 0x7D, 0x31)
TABLE_ALT_BG = RGBColor(0xE8, 0xF0, 0xF8)
BLACK = RGBColor(0x33, 0x33, 0x33)

MARGIN = 0.4  # inch
CONTENT_W = 7.5 - MARGIN * 2  # 6.7 inch

def create_presentation():
    prs = Presentation()
    prs.slide_width = Inches(7.5)
    prs.slide_height = Inches(13.333)
    return prs
```

---

## レイアウトパターン詳細

### A: タイトルスライド

- タイトルバー: なし（唯一の例外）
- 背景: DARK_BLUE全面
- タイトル: 48pt 白太字 中央寄せ（Y=4.0inch付近）
- サブタイトル: 24pt 薄青(`#A0C4E8`) 中央寄せ
- 補足: 18pt 薄青 中央寄せ

### B: 情報カード

- 角丸矩形（WHITE背景、ACCENT_BLUE枠線2pt）
- 2列配置、1行あたり2カード
- ラベル: 12pt ACCENT_BLUE太字 中央寄せ
- 値: 20pt DARK_BLUE太字 中央寄せ

### C: テーブル

- ヘッダ: ACCENT_BLUE背景、白12pt太字、中央寄せ
- 偶数行: TABLE_ALT_BG背景
- セル: 11pt BLACK、中央寄せ
- 幅: コンテンツ幅全体（6.7inch）を均等分割
- 最大8行/スライド

### D: タイムライン

- 時間バー: 角丸矩形（ACCENT_BLUE背景、白18pt太字、中央寄せ）、幅=コンテンツ幅
- 説明: その下に16pt テキスト
- 間隔: 各項目1.6inch

### E: アラート

- 左: 角丸矩形（ORANGE背景、"!" 白22pt太字）
- 右: テキスト（14pt BLACK）
- 間隔: 各項目1.5inch

### F: 写真ステップ

- ラベル: "Step N: 説明" 13pt ACCENT_BLUE太字 中央寄せ
- 画像: ラベルの下、コンテンツ幅いっぱい。アスペクト比維持
- 縦に積む（横並びにしない）
- 間隔: 各ステップ2.5inch
- 1スライド最大2枚

### G: 箇条書き

- "● " プレフィックス
- 20pt BLACK
- space_after: 8pt

---

## 検証チェックリスト

### レイアウト

- [ ] 全スライドが縦長（9:16）になっているか
- [ ] タイトルバーがスライド上部に表示されているか
- [ ] テキストがスライド幅からはみ出していないか
- [ ] 画像のアスペクト比が崩れていないか
- [ ] 1スライドに画像が3枚以上詰め込まれていないか

### コンテンツ

- [ ] 電話番号・メールアドレスが含まれていないか
- [ ] 「代表連絡先」セクションが除外されているか
- [ ] 1スライドの本文が200文字を大幅に超えていないか
- [ ] フォントサイズが14pt未満の箇所がないか
- [ ] テーブルが8行を超えていないか

### PDF変換後

- [ ] LibreOfficeでのPDF変換でレイアウトが崩れていないか
- [ ] pdf.js viewerでスクロール表示できるか
- [ ] ピンチズーム/ズームボタンが動作するか

---

## PDF変換・配信の詳細手順

### LibreOfficeでPDF変換

```bash
/Applications/LibreOffice.app/Contents/MacOS/soffice \
  --headless --convert-to pdf \
  --outdir docs/manuals/ \
  docs/manuals/対象.pptx
```

### build/webへの配置

```bash
# ASCII名にリネームして配置（日本語ファイル名はURLエンコード問題を起こす）
# manuals/ ディレクトリに複数PDF を格納
mkdir -p mobile/build/web/manuals/
cp docs/manuals/対象.pdf mobile/build/web/manuals/haisen.pdf
```

### pdf.jsビューアの配置

```
mobile/build/web/pdfjs/
├── viewer.html        ← ピンチズーム・ズームボタン付きビューア
├── pdf.min.mjs        ← pdf.js本体（v4.10.38）
└── pdf.worker.min.mjs ← pdf.js Worker
```

ソースは `scripts/pdfjs/` に保管。`flutter build web` 後に再配置が必要。

### ワークディレクトリ構造

各マニュアルは `scripts/manuals/{短縮名}/` に専用ディレクトリを持つ。

```
scripts/manuals/
├── 01_配線/
│   ├── create_pptx.py    ← このマニュアル専用の生成スクリプト
│   └── images/           ← このマニュアルのdocxから抽出した画像
├── 02_駐車場設営/
│   ├── create_pptx.py
│   └── images/
...
```

画像やスクリプトがマニュアル間で混ざるのを防ぐ。入力docxと出力pptx/PDFは `docs/manuals/` にフラット格納。

---

## 正解例: 配線マニュアルの変換対応表

`01_44th_配線マニュアル.docx` の変換結果を正解例として参照する。

| docxの構造 | パターン | スライドタイトル |
|---|---|---|
| ドキュメント冒頭 | A: タイトル | 配線マニュアル 第44回技大祭 |
| 集合場所 + key-value | B: 情報カード | 該当シフト・集合情報 |
| 表 晴天時チーム分け + テーブル | C: テーブル | チーム分け（晴天時） |
| 表 雨天時チーム分け + テーブル | C: テーブル | チーム分け（雨天時） |
| 必要物品 + 箇条書き | G: 箇条書き | 必要物品 |
| タイムスケジュール + 時間付き | D: タイムライン | タイムスケジュール |
| 配線箇所 + テーブル | C: テーブル | 配線箇所（晴天時/雨天時） |
| 配線時の注意点 + 箇条書き | E: アラート | 配線時の注意点 |
| 電源付近の処理 + 画像 | F: 写真ステップ | 電源付近の処理 |
| コードの固定処理 + Step付き画像群 | F: 写真ステップ | コードの固定処理 |
| コードの固定場所 + 画像群 | F: 写真ステップ | コードの固定場所 |
| 配線計画図（画像6枚） | F: 写真 | 配線計画図 |
| 代表連絡先 + 電話番号 | **除外** | — |
