---
name: docx-to-manual-slide
description: "技大祭マニュアルのdocxを縦長pptxスライドに変換する。docxからpptx変換、マニュアルスライド生成時に参照すること。"
---

# docx → マニュアルスライド 生成ルール

## スライド基本設定

- サイズ: 縦長 9:16（Inches(7.5) × Inches(13.333)）
- 余白: 左右 0.4 inch、コンテンツ幅 6.7 inch
- 全スライド上部にタイトルバー（DARK_BLUE背景, 白文字, 太字26pt）。タイトルスライドのみ例外

## カラースキーム

| 色名 | HEX | 用途 |
|------|-----|------|
| DARK_BLUE | `#1B3A5C` | タイトルバー背景 |
| ACCENT_BLUE | `#2E75B6` | テーブルヘッダ、タイムバー、ラベル |
| LIGHT_BG | `#F2F2F2` | スライド背景 |
| WHITE | `#FFFFFF` | カード・セル背景 |
| ORANGE | `#ED7D31` | 注意アイコン |
| TABLE_ALT_BG | `#E8F0F8` | テーブル偶数行 |
| BLACK | `#333333` | 本文テキスト |

## レイアウトパターン（7種）と選択基準

```
1. ドキュメント先頭 → A（タイトル）
2. 見出し配下を分類:
   - 「〇〇：△△」形式が2〜6個 → B（情報カード）
   - 「表　〇〇」+ テーブル → C（テーブル）
   - HH:MM〜HH:MM を含む → D（タイムライン）
   - 「必ず」「〜すること」「禁止」を含む → E（アラート）
   - 画像参照 + 手順説明 → F（写真ステップ）
   - 上記以外の箇条書き → G（箇条書き）
```

## コンテンツルール

- **スライド分割**: H2ごとに1〜2枚。晴天/雨天等の条件分岐は別スライド
- **テキスト**: 1スライド200文字以内。フォント下限14pt
- **除外**: 電話番号、メールアドレス、「代表連絡先」セクション全体、ローカルファイルパス
- **画像**: キャプション（「図　〇〇」）から用途判断。アスペクト比維持。1スライド最大2枚
- **テーブル**: 最大8行で分割。ヘッダはACCENT_BLUE背景白文字

## 入出力パス（マニュアルごとに分離）

各マニュアルは `scripts/manuals/{短縮名}/` に専用ワークディレクトリを持つ。

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

| 種別 | パス | 備考 |
|------|------|------|
| 入力docx | `docs/manuals/` | 全docxをフラットに格納 |
| 抽出画像 | `scripts/manuals/{短縮名}/images/` | マニュアルごとに分離 |
| 生成スクリプト | `scripts/manuals/{短縮名}/create_pptx.py` | マニュアルごとに分離 |
| pptx/PDF出力 | `docs/manuals/` | 入力docxと同階層 |
| PDF配信先（build時） | `mobile/build/web/manuals/{ascii名}.pdf` | ASCII名で配置 |

**短縮名の命名規則**: docxファイル名から `01_44th_` 等のプレフィックスと `マニュアル` サフィックスを除いた部分（例: `01_44th_配線マニュアル.docx` → `01_配線`）

## 生成手順（概要）

```bash
MANUAL_DIR=scripts/manuals/01_配線  # 対象マニュアルのワークディレクトリ

# 1. docxから画像抽出（python-docx + zipfile）→ ${MANUAL_DIR}/images/
# 2. docxの構造を読み取り、パターン選択基準に従いスライド設計
# 3. python-pptxでpptx生成
python3 ${MANUAL_DIR}/create_pptx.py
# 4. LibreOfficeでPDF変換
/Applications/LibreOffice.app/Contents/MacOS/soffice --headless --convert-to pdf --outdir docs/manuals/ docs/manuals/対象.pptx
# 5. 配信先に配置（ASCII名で）
cp docs/manuals/対象.pdf mobile/build/web/manuals/haisen.pdf
```
