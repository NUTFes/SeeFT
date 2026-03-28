---
name: html-manual-slide
description: "技大祭マニュアルのGoogle DocumentをHTML縦スライドに変換する。マニュアルHTML生成時に必ず参照すること。"
---

# Google Document → HTMLマニュアルスライド 生成ルール

## 概要

Google DocumentをHTML形式でダウンロードし、そのHTML+画像をClaudeが読み取り、モバイル向け縦スライド形式の単一HTMLファイルに再構成する。

## 入力

Google DocumentをHTML形式（.html.zip）でダウンロードして展開した以下のファイル群:

```
docs/manuals/{マニュアル名}/
├── {マニュアル名}.html   ← メインHTML（構造・テキスト）
└── images/               ← 画像ファイル群
    ├── image1.png
    ├── image2.jpg
    ...
```

## 出力

単一の自己完結型HTMLファイル（CSS埋め込み、画像はBase64 data URI）:

```
docs/manuals/{マニュアル名}/slide.html
```

## スライド基本設定

- **方向**: 縦長（モバイルファースト、幅100vw）
- **構造**: セクション区切りの縦スクロール形式（1セクション = 1画面分）
- **各セクション高さ**: `min-height: 100svh`（スマホのビューポート高さ1画面分を目安）
- **余白**: 左右 5vw
- **フォント**: system-ui（デバイスネイティブ）

## カラースキーム

| 変数名 | HEX | 用途 |
|--------|-----|------|
| --gold-dark | `#C89932` | タイトルバー背景、見出しアクセント |
| --gold-mid | `#D8B660` | テーブルヘッダ背景、タイムバー |
| --gold-light | `#FFD966` | アラートアイコン背景、ホバー |
| --warm-beige | `#E0DACF` | スライド背景、カード背景 |
| --brown-red | `#96514D` | 注意アイコン、強調テキスト |
| --text-main | `#000000` | メインテキスト |
| --text-sub | `#595959` | サブテキスト、キャプション |
| --white | `#FFFFFF` | カード背景、テーブルセル |

## レイアウトパターン（7種）と選択基準

```
1. ドキュメント先頭 → A（タイトル）
2. 見出し配下を分類:
   - 「〇〇：△△」形式が2〜6個 → B（情報カード）
   - テーブル → C（テーブル）
   - HH:MM〜HH:MM を含む → D（タイムライン）
   - 「必ず」「〜すること」「禁止」「注意」を含む → E（アラート）
   - 画像 + 手順説明 → F（写真ステップ）
   - 上記以外のテキスト・箇条書き → G（箇条書き）
```

### A: タイトルセクション
- 背景: `--gold-dark` グラデーション
- タイトル: 白文字、太字、2rem
- サブタイトル: 白文字 半透明、1.2rem
- 中央寄せ、`height: 100svh` で1画面全体を使う

### B: 情報カード
- CSS Grid 2列配置
- 各カード: `--white` 背景、`--gold-dark` 左ボーダー 4px
- ラベル: `--text-sub` 0.8rem
- 値: `--text-main` 太字 1.2rem

### C: テーブル
- ヘッダ: `--gold-mid` 背景、白文字
- 偶数行: `--warm-beige` 背景
- セル: 0.85rem、`--text-main`
- 横スクロール対応（`overflow-x: auto`）
- 最大8行で次セクションに分割

### D: タイムライン
- 時間バー: `--gold-mid` 背景、白文字太字、角丸
- 説明テキスト: その下に `--text-main`
- 左に縦線（`--gold-dark` 2px）でタイムライン接続

### E: アラート
- 左: `--brown-red` 背景の「!」アイコン（角丸正方形）
- 右: `--text-main` テキスト
- 各項目間にセパレータ

### F: 写真ステップ
- ステップラベル: `--gold-dark` テキスト太字
- 画像: 幅100%、角丸、アスペクト比維持
- 縦積み（1セクション最大2枚、超えたら次セクションへ）

### G: 箇条書き
- ビュレット: `--gold-dark` の丸（カスタムリストスタイル）
- テキスト: `--text-main` 1rem
- 行間: 1.8

## コンテンツルール

- **セクション分割**: 見出し（H1/H2）ごとに1セクション。コンテンツが多い場合は分割可
- **密度**: 1セクションに100svh分のコンテンツを **しっかり詰める**。スカスカ禁止。ただしはみ出しも禁止（はみ出る場合はセクション分割）
- **テキスト**: フォント下限 0.85rem（約14px）
- **除外**: 電話番号（0X0-XXXX-XXXX等）、メールアドレス、「代表連絡先」セクション全体
- **重複禁止**: 同一マニュアル内でほぼ同じ内容のセクションを複数生成しない。見出し違い・表現違いでも内容が重複していれば1つにまとめる
- **画像**: Base64 data URIに変換して埋め込み。アスペクト比維持。`max-width: 100%`
- **テーブル**: 最大8行で分割。列数が多い場合は横スクロール

## HTML構造テンプレート

```html
<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{マニュアル名}</title>
  <style>
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
    body {
      font-family: system-ui, sans-serif;
      color: var(--text-main);
      background: var(--warm-beige);
      -webkit-text-size-adjust: 100%;
    }
    section {
      min-height: 100svh;
      padding: 5vw;
      display: flex;
      flex-direction: column;
      justify-content: center;
    }
    /* ... パターンごとのCSS ... */
  </style>
</head>
<body>
  <section class="title"> ... </section>
  <section class="info-cards"> ... </section>
  <section class="table"> ... </section>
  ...
</body>
</html>
```

## 生成手順

1. `docs/manuals/{マニュアル名}/` 内のHTMLファイルと画像ファイルをReadツールで読み取る
2. HTMLの構造（見出し・段落・テーブル・画像参照）を解析し、レイアウトパターンを選択
3. 画像をBase64に変換し、data URIとして埋め込む
4. 上記テンプレート・カラースキーム・レイアウトパターンに従い、単一HTMLファイル `slide.html` を生成
5. 生成した `slide.html` を `docs/manuals/{マニュアル名}/slide.html` に出力

## 検証チェックリスト

- [ ] 全セクションが縦長・モバイル幅に収まっているか
- [ ] セクションがスカスカでないか（100svhを有効活用しているか）
- [ ] セクションからコンテンツがはみ出していないか
- [ ] 電話番号・メールアドレスが含まれていないか
- [ ] 「代表連絡先」セクションが除外されているか
- [ ] フォントサイズが0.85rem未満の箇所がないか
- [ ] テーブルが8行を超えていないか
- [ ] 画像のアスペクト比が崩れていないか
- [ ] 画像がBase64で埋め込まれているか（外部参照なし）
- [ ] スマホ実機でスクロール表示・ピンチズームが動作するか
