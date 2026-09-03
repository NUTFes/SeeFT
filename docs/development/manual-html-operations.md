# 解説マニュアルHTML 運用手順

Googleドキュメントで書かれたマニュアルを、スマホ最適化された解説HTMLに変換し、SeeFTアプリのシフトカードから開けるようにするまでの手順書である。2026-08-31時点のライブのApps Script（`clasp` で取得）と develop のコードに基づく。

この文書が対象にするのは「マニュアル1本を委員が読める状態にするまで」であり、マニュアル本文そのものの作成（部門長が書く工程）は範囲外である。

## 要約

工程は5つある。従来ここが1人に集中していたが、必要な権限は工程ごとに違っており、分担できる。

```
① Googleドキュメントを取得   → ドキュメントの閲覧権限があれば誰でも
② 解説HTMLに変換             → Claudeの契約とMacのセットアップが要る
③ 生成結果を確認             → 誰でも（判断基準は本書の「確認項目」）
④ サーバーへアップロード      → アップロードトークンがあれば誰でも（curl 1発）
⑤ タスクへ紐付け             → シフトスプシの編集権限があれば誰でも
```

②だけが環境構築を伴う。①③④⑤は権限を渡せばその日から任せられる。

紐付け（⑤）はスプレッドシートで完結する。DBを直接触る必要はなく、**触ってはいけない**（後述の「やってはいけないこと」）。

## 全体像

```
[Googleドキュメント]
      │ ①「ウェブページ(.html, zip)」でダウンロード → prepare_manual.py
      ▼
[docs/manuals/{マニュアル名}/]  source.html + images/
      │ ② generate_slide.py（Claudeが変換）
      ▼
[slide_claude.card-strict.html]  自己完結HTML（CSS・JS・画像を1ファイルに内包）
      │ ③ 目視 + 機械検証
      │ ④ upload_manual.py → PUT /manuals/{id}
      ▼
[https://seeft-api.nutfes.net/manuals/{id}]  ← manual_url と、対応表へ貼る1行が返る
      │ ⑤ シフトスプシ「マニュアルURL」シートへ貼る → タスク送信
      ▼
[tasks.url / tasks.manual_url]  → mobileのシフトカードにボタンが出る
```

## 必要な権限とアカウント

| 工程 | 必要なもの | 現状の保有者 |
| --- | --- | --- |
| ① 取得 | 対象Googleドキュメントの閲覧権限 | 部門長・執行部・SeeFT |
| ② 変換 | Claudeの契約（`claude login` 済み）、pandoc、uv、リポジトリのclone | 上林のMaxプランを共用（要見直し） |
| ③ 確認 | なし（ブラウザでHTMLを開くだけ） | 誰でも |
| ④ アップロード | `MANUAL_UPLOAD_TOKEN` の値 | 上林（共有すれば誰でも） |
| ⑤ 紐付け | シフトスプシ `45th_シフト_ver0` の編集権限 | シフト担当・SeeFT |
| 閲覧 | nutfesのGoogleアカウント（`〜.nutfes@gmail.com`） | 全委員 |

②のClaude契約が個人サブスクの共用になっている点は、9月以降のPM不在期間に向けて体制を決め直す必要がある。他の4工程は権限の受け渡しだけで分担できる。

## まとめて実行する（①〜④）

①〜④を1人で担当するなら、通しで実行するスクリプトがある。工程ごとに実行したい場合は以降の①〜⑤を参照する。

```bash
scripts/automation/run_pipeline.sh docs/manuals/_zips/45th_企画マニュアル_縁日.zip --id en-nichi --doc-url "https://docs.google.com/document/d/xxxx/edit"
```

判断が必要な2箇所で確認を挟む。①の後に「マニュアル名がタスク一覧M列と一致しているか」、③の後に「内容を見て問題ないか」。どちらも間違えたまま進むと後で静かに失敗する。

- `--prompt` / `--model` で②の設定を変えられる（既定は `card-strict` / `claude-opus-4-7`）
- ④のトークンは `MANUAL_UPLOAD_TOKEN` を設定しておけば対話入力なしで進む。**`--yes` は確認を飛ばすだけでトークン入力は省略しないため、無人実行するならこの環境変数の設定が前提になる**
- ⑤（シフトスプシへの紐付け）は対象外。最後に④の出力（対応表に貼る行）が表示されるので、それを「⑤ タスクに紐づける」の手順1で貼る

### 再生成する（①をやり直さない）

元のGoogleドキュメントを直したがzipを取り直したくない、`instructions.md` に修正指示を足して作り直したい、というときは `--manual-dir` で①をスキップできる。

```bash
scripts/automation/run_pipeline.sh --manual-dir docs/manuals/45th_企画マニュアル_縁日 --id en-nichi --doc-url "https://docs.google.com/document/d/xxxx/edit"
```

②〜④だけが回り、出力HTMLは無条件に上書きされる。

## ① Googleドキュメントを取得する

対象のドキュメントを開き、`ファイル > ダウンロード > ウェブページ（.html、zip形式）` を選ぶ。PDFやWordではなくHTMLである点が重要で、これ以外の形式では画像が取り出せない。

ダウンロードしたzipを、そのまま次のコマンドに渡す。展開・リネーム・配置をまとめて行う。

```bash
python3 scripts/automation/prepare_manual.py docs/manuals/_zips/45th_企画マニュアル_縁日.zip
```

`docs/manuals/{マニュアル名}/` に次の形で配置される。マニュアル名はzipのファイル名から取る（`--name` で明示もできる）。

```
docs/manuals/45th_企画マニュアル_縁日/
├── source.html        ← zip内のHTMLをリネームしたもの
└── images/            ← zip内の画像
    ├── image1.png
    └── image2.jpg
```

手作業でやらないのは、次の3点を毎回間違えるためである。

zip内のHTMLは名前が切り詰められている。実例として `45th_企画マニュアル_ホールインワン.zip` の中身は `45th_.html` だった。zip名はタイトルを保つが、中のファイル名は保たない。

`source.html` という名前が必要なのは、生成スクリプトが「`slide` と `verify` で始まらない `.html`」を入力として探すためで、複数あると曖昧だとして停止する（`scripts/claude-slide/generate_slide.py` の `load_source()`）。

`images/` の位置がずれると画像が埋め込まれない。

実行後、表示されたマニュアル名が**タスク一覧M列の値と完全に一致しているか確認する**。ここがずれていると⑤で紐付けが静かに失敗する。

## ② 解説HTMLに変換する

初回のみ、Macに環境を作る。

```bash
brew install pandoc uv
claude login
uv sync --project scripts/claude-slide
```

変換を実行する。`card-strict` は本文を一字一句変えないプロンプトで、元Docが執行部と部門長の間で調整済みであることを前提にしている。

```bash
uv run --project scripts/claude-slide python scripts/claude-slide/generate_slide.py --prompt card-strict --model claude-opus-4-7 docs/manuals/45th_企画マニュアル_縁日
```

同じディレクトリに `slide_claude.card-strict.html` が出る。1本あたり3〜6分かかる。

生の応答は `slide_claude.card-strict.html.raw.md` にも保存される。抽出や画像埋め込みに失敗した場合、10分かけて再生成せずここから復旧できる（LLMの呼び出しは非決定的なので、再生成すると別のHTMLになる）。

## ③ 生成結果を確認する

### 確認項目

ファイルサイズを最初に見る。画像のあるマニュアルなら通常は数MBになる。数十KBしかない場合は画像が埋め込まれていないので、再生成せずに埋め直せる。

```bash
uv run --project scripts/claude-slide python scripts/claude-slide/generate_slide.py --prompt card-strict --embed-only docs/manuals/45th_企画マニュアル_縁日
```

`--embed-only` はLLMを呼ばず、既存HTMLの `<img>` をローカル画像で埋め直すだけなので、何度実行しても同じ結果になる。

次に、元Docと本文を機械的に比較する。デザインの良し悪しは判定せず、文章の「消えた・書き換わった・増えた」だけを出す。

```bash
uv run --project scripts/claude-slide python scripts/claude-slide/verify_slide_mechanical.py docs/manuals/45th_企画マニュアル_縁日
```

`文章チェック結果.card-strict.md` が部門長に見せられる日本語サマリ、`verify_mechanical.card-strict.txt` が開発用の詳細である。

最後にブラウザで開き、スマホ幅（開発者ツールの375px程度）で目次・折りたたみ・画像の表示を見る。

### 既知の不具合

透過PNGを含むマニュアルは、背景が白く焼き込まれた状態で埋め込む必要がある。透過のまま埋めるとカード背景と干渉して図が読めなくなる。

応答が途中で切れた場合、生成スクリプトは `</html>` で閉じているかを検査して停止する（PR #478 で追加）。停止したら `.raw.md` を開いて、どこで切れたかを確認する。

## ④ サーバーへアップロードする

### エンドポイント

```
PUT https://seeft-api.nutfes.net/manuals/{id}
```

| 項目 | 内容 |
| --- | --- |
| 認証 | `Authorization: Bearer {MANUAL_UPLOAD_TOKEN}` |
| ボディ | HTMLファイルそのもの（バイナリ送信） |
| `{id}` の規則 | `a-z` `0-9` `_` `-` のみ、1〜64文字。日本語・大文字・スラッシュは不可 |
| サイズ上限 | 20MB |
| 反映 | 即時。APIの再起動は不要 |

`{id}` はそのまま公開URLになるので、後から見て何のマニュアルか分かる英数字を付ける（`en-nichi`、`stamp-rally` など）。

### コマンド

`--doc-url` にはGoogleドキュメントの共有URLを渡す。⑤で対応表のB列に入れる値になる。

```bash
python3 scripts/automation/upload_manual.py --id en-nichi --doc-url "https://docs.google.com/document/d/xxxx/edit" docs/manuals/45th_企画マニュアル_縁日
```

トークンは訊かれるので貼り付ける。入力は画面にもシェルの履歴にも残らない。環境変数 `MANUAL_UPLOAD_TOKEN` を設定してあればそちらが使われる。

送信前にHTMLが `</html>` で閉じているかとサイズ上限を検査するので、壊れたファイルや大きすぎるファイルを配信してしまうことはない。

成功すると公開URLと、⑤で対応表にそのまま貼れるタブ区切りの1行が表示される。

```
成功: https://seeft-api.nutfes.net/manuals/en-nichi

「マニュアルURL」シートに貼る行（タブ区切り）:

45th_企画マニュアル_縁日	https://docs.google.com/document/d/xxxx/edit	https://seeft-api.nutfes.net/manuals/en-nichi
```

### curl で送る場合

スクリプトを使わずに送ることもできる。トークンは `export TOKEN='値'` と直接書かない。コマンドがシェルの履歴ファイルに残り、トークンが平文で保存されるためである。

```bash
printf 'アップロードトークンを貼り付けてEnter: '; read -rs MANUAL_UPLOAD_TOKEN; echo; export MANUAL_UPLOAD_TOKEN
```

```bash
curl -X PUT "https://seeft-api.nutfes.net/manuals/en-nichi" -H "Authorization: Bearer $MANUAL_UPLOAD_TOKEN" -H "Content-Type: text/html" --data-binary @docs/manuals/45th_企画マニュアル_縁日/slide_claude.card-strict.html
```

作業が終わったらターミナルを閉じる。環境変数はそのウィンドウにしか残らないため、閉じれば消える。

### レスポンス

成功すると `201 Created` で、次のJSONが返る。

```json
{
  "id": "en-nichi",
  "manual_url": "https://seeft-api.nutfes.net/manuals/en-nichi"
}
```

`manual_url` の値をそのまま⑤で使う。キー名がDBの `tasks.manual_url` と揃えてあるのは、貼り先を取り違えないためである。

このキー名は issue #479 で `url` から改名した。**本番APIへデプロイされるまでは旧キー `url` で返る**（2026-08-31 時点で本番は改名前のビルド）。`upload_manual.py` は両方を受け付けるため、デプロイの前後どちらでも動く。curl で直接叩く場合は、返ってきたキーがどちらかを見て読み替えること。

失敗時は次のいずれかが返る。

| ステータス | ボディ | 原因 |
| --- | --- | --- |
| 401 | `{"error":"unauthorized"}` | トークンが違う、`Bearer ` が抜けている、またはサーバー側でトークン未設定 |
| 400 | `{"error":"invalid manual id"}` | `{id}` に日本語・大文字・記号が入っている |
| 413 | `{"error":"payload too large"}` | 20MB超。画像の多いマニュアルで起きうる |
| 500 | `{"error":"failed to save"}` | サーバー側の保存失敗。ログを確認する |

401が返る場合、トークンの誤りとサーバー側の未設定を区別できない。これは意図的な設計で、設定漏れのままアップロードが素通りする事故を防ぐため、`MANUAL_UPLOAD_TOKEN` が未設定なら全リクエストを拒否する。

### アップロードできたか確認する

ブラウザで `https://seeft-api.nutfes.net/manuals/en-nichi` を開く。nutfesのGoogleアカウントでのログインを求められ、ログイン後にマニュアルが表示されれば成功である。個人のGmailアカウントでは403になる。

## ⑤ タスクへ紐付ける

### 仕組み

紐付けはシフトスプシ `45th_シフト_ver0` で完結する。人が触るのは「マニュアルURL」シート1枚だけである。

```
[マニュアルURL シート]  A=マニュアル名 / B=ドキュメントURL / C=スライドURL
      │ VLOOKUP（タスク一覧のM列をキーにする）
      ▼
[タスク一覧シート]  M列=マニュアル名 → R列=ドキュメントURL / S列=スライドURL（自動）
      │ SeeFTメニュー > タスク送信
      ▼
[POST /api/update_tasks_and_places] → tasks.url / tasks.manual_url
```

この構造の利点は、**対応表に1行足すとそのマニュアル名を持つ全タスクへ自動的に広がる**ことである。1つのマニュアルに対してタスクは5〜9件あるのが普通で（縁日は7件）、その分の作業が1行で済む。

### 手順

「マニュアルURL」シートに行を追加する。④の `upload_manual.py` が出力したタブ区切りの1行を、空き行のA列に貼れば3列が一度に埋まる。

- A列: マニュアル名
- B列: Googleドキュメントの共有URL（ドキュメント版ボタン用）
- C列: ④のレスポンスの `manual_url`（HTML版ボタン用）

**貼った後、A列がタスク一覧M列と完全に一致しているか必ず確認する。** スクリプトが組み立てるA列はディレクトリ名の写しであって、M列の値である保証はない。一致していなければM列のセルをコピーして貼り直す。手で打ち直さない。

B列とC列は別のDBカラムに入る。名前が似ていて紛らわしいので、確認するときは対応を間違えないこと。

| 対応表 | タスク一覧 | DBカラム | シフトカードのボタン |
| --- | --- | --- | --- |
| B列 ドキュメントURL | R列 | `tasks.url` | ドキュメント版 |
| C列 スライドURL | S列 | `tasks.manual_url` | HTML版 |

B列だけ埋めた状態は正常な運用形態である。HTMLを生成していないマニュアルでも、ドキュメント版のボタンは出せる。このときC列は空で、`tasks.manual_url` も空のままになる。

```sql
-- ドキュメント版が入ったか（B列の結果）
SELECT id, task, COALESCE(url,'') AS doc_url FROM tasks WHERE task LIKE '%カジノ%' ORDER BY id;

-- 両方の埋まり具合をまとめて見る
SELECT COALESCE(url,'') <> '' AS has_doc, COALESCE(manual_url,'') <> '' AS has_slide, count(*) FROM tasks GROUP BY 1,2 ORDER BY 3 DESC;
```

タスク一覧のR/S列にVLOOKUPの数式が入っていない場合は、`SeeFTメニュー` から `fillManualUrlFormulas` を実行する。4行目から最終行+50行まで数式を貼るので、行が増えても追従する。

貼ったら `SeeFTメニュー > マスタをSeeFTに送信する > 【事前確認】マニュアルURLの紐付けを点検する` を実行する。対応表とM列のずれを送信前にすべて列挙する（issue #495）。

```
不足    M列にあるのに対応表に無い → まだ作っていないマニュアルの一覧（バックログ）
惜しい  空白を除けば一致 → 末尾スペース等の揺れ。直す
余り    対応表がどのタスクからも参照されていない → 行を見直す
重複    対応表A列に同じ値が2行 → 2行目は引かれない。消す
先勝ち  同名タスクの最初の行だけM列が空 → 局ファイルで全行に入れる
URL欠け 使われている行のB/C列が空 → ボタンが出ない
```

「問題N件」には不足（バックログ）も含まれる。紐づけようとしているマニュアルに関する指摘がなければ送信してよい。

問題がなければ `SeeFTメニュー > タスク送信` を実行する。これでDBの `tasks.manual_url` が更新され、mobileのシフトカードに「HTML版」ボタンが出る。

### A列を手で打ってはいけない理由

`VLOOKUP(..., FALSE)` は完全一致で、外れた場合は `IFERROR(..., "")` が空に潰す。**エラーは一切出ず、R/S列が空になるだけ**である。次のような差異で静かに外れる。

```
45th_企画マニュアル_縁日                                  ← M列と同じ。引ける
45th_企画マニュアル_縁日                                  ← 末尾に半角スペース。外れる
45th＿企画マニュアル＿縁日                                ← アンダースコアが全角。外れる
45ｔh_企画マニュアル_縁日                                 ← t が全角。外れる
```

やむを得ず手入力した場合は、送信後にタスク件数が想定と一致するかで検証する。

### 同名タスクは最初の行だけが送られる

タスク一覧には同じタスクが日程分だけ行として並ぶ。送信時はタスク名で先勝ちの重複排除がかかるため、**採用されるのは最初に現れた行だけ**である。

```js
// 名簿タスク送信.js:225-228
if (seen[taskName]) return; // タスク一覧には同名が複数日程分並ぶため先勝ちで1件にまとめる
```

最初の行のM列が空で、2行目以降だけにマニュアル名が入っていると、URLは送られない。M列は同名タスクの全行に入れておく。

### ドキュメント名を変えると紐付けが静かに切れる

各局のタスクファイルのM列は、Googleドキュメントへの**スマートチップ**になっていることがある。チップが表示するのは**参照先ドキュメントの現在のファイル名**である。タスク一覧はそれをIMPORTRANGEで取り込むが、IMPORTRANGEは値しか運ばないため、タスク一覧のM列に届くのは表示名の**文字列**である（タスク一覧側にチップは存在しない）。

つまりドキュメントをリネームすると、誰もスプレッドシートを触っていなくてもタスク一覧のM列の値が変わる。

```
ドキュメントをリネーム
  → 局ファイルのチップの表示が変わる（自動）
  → IMPORTRANGE経由でタスク一覧M列の文字列が変わる
  → 対応表A列に古い名前が残っていると VLOOKUP が外れる
  → S列が空になる
  → 次のタスク送信で manual_url が空で上書きされ、シフトカードのボタンが消える
```

この過程でエラーは出ない。ある日突然ボタンが消えていたら、まずドキュメントがリネームされていないかを疑う。

**マニュアルのドキュメント名を変更したら、対応表A列も同じ値に更新すること。** 逆に、対応表A列の文字列がM列と合わなくなったときは、手打ちミスだけでなくリネームの可能性も考える。

なお、この性質は利点にもなる。M列のセルが編集できない形式（チップ）になっている場合、末尾スペースなどの表記の汚れは**ドキュメント名を直せば全セルまとめて直る**。スプレッドシート側のセルを編集する必要はなく、IMPORTRANGEにも表にも触らずに済む。

### 空のS列は既存の値を消す

`UpdateWithManualURL` は `manual_url` を無条件にSETするため、S列が空の行は空で上書きされる。タスク送信は差分更新ではなく全面同期であり、スプレッドシートに無い値はDBから消える。

これは「スプレッドシートが唯一の入力先」という設計を強制するための挙動であって、不具合ではない。裏を返すと、M列が未設定のタスクにURLを持たせたい場合、対応表とM列を埋める以外の方法は無い。

## どのタスクがどのマニュアルに対応するかの決め方

### 基本: タスク一覧のM列を見る

タスク一覧シートのM列は見出しが「リンク」だが、中身はハイパーリンクではなく**マニュアル名の文字列**である（例: `45th_企画マニュアル_ビンゴ大会`）。つまりM列が既にタスク→マニュアルの対応表になっている。

```
A シフト名 / B 天気 / C 開始時間 / D 終了時間 / E 拘束時間 / F 管轄局 / G レベル
H 集合場所 / I 作業場所 / J 最低人数 / K 適性人数 / L 最大人数 / M リンク
N タスク内容 / O 備考 / P 色 / Q 日付 / R ドキュメントURL / S スライドURL
```

データは4行目から始まる（3行目がヘッダー）。2026-08-29時点でM列のユニーク値は59種ある。

同じマニュアル名が複数行に入っているのが正常な状態で、それがそのまま「複数タスクで1マニュアルを共有する」という意味になる。

### M列にGoogleドキュメントのURLが直接入っている行

マニュアル名ではなくURLがそのまま入っている行がある（2026-09-01時点で66行）。これらのURLは**44th（去年）のマニュアルを指しており、リンク先として価値がない**。紐づけるときは、担当局のタスクファイルのM列を45thのマニュアル名に書き換える（キーはマニュアル名で統一する。2026-09-01決定）。

URL文字列をそのまま対応表A列のキーとして使うことも技術的には可能だが、どのみち捨てる値に依存することになるため使わない。

### M列が空のタスク

マニュアル本文の側から拾える。マニュアルには次のような書き方でタスク名が埋まっていることが多い。

```
当日、【駐車場管理】のシフトに該当する者は、8:00に本部前に集合すること。
```

この角括弧の中がタスク一覧のシフト名（A列）に対応する。本文から角括弧を抽出してM列へ書き戻すのが逆方向の解になる。1つのマニュアルが複数タスクに対応するため、M列側には同じマニュアル名が複数行入る形になる。

抽出できなかったものは一覧に吐き出して、人が見る範囲を絞る。

### 割り当てスプシは対応表に使わない

マニュアルの割り当てスプシ（`1a2pvM1M8NWQNLNaYnsqTzpB_oGE129-ibpbN3Be1Z1Q`）にもマニュアル名の列があるが、こちらは局ごとに命名規則が違い、`配線マニュアル` `物品移動計画書` のような短い形になっている。タスク一覧のM列は `45th_企画マニュアル_縁日` の完全形なので、両者は機械的に突き合わせられない。対応表のキーはタスク一覧M列に揃える。

## やってはいけないこと

### SQLで tasks.manual_url を直接UPDATEする

タスク送信は `UpdateWithManualURL` で `url` と `manual_url` の両方をSETする。SQLで直接書いた値は、**次のタスク送信で上書きされて消える**。入力先は「マニュアルURL」シートの1枚に統一する。

```go
// api/lib/internals/repository/task_repository.go:93
func (b *taskRepository) UpdateWithManualURL(c context.Context, id string, name string, placeID string, url string, manualURL string, ...) error {
	SET task = $1, place_id = $2, url = $3, manual_url = $4, bureau_id = $5,
```

### 管理画面からマニュアルURLを直そうとする

Web管理画面のタスク編集は `manual_url` を更新対象から意図的に外している。編集しても値は変わらない。

```go
// api/lib/usecase/task_usecase.go:223-225
// 管理画面からの編集ではマニュアル（スライド版）を扱わないため、
// Update は manual_url をSET句に含めない（既存値がそのまま残る）
```

### リポジトリの gas/ を現行コードとして読む

`gas/` はライブのApps Scriptと乖離することがある。GASを触る作業は必ず `clasp` での取得から始める。

```bash
cd gas && npm install
node_modules/.bin/clasp clone <スクリプトID>
```

スクリプトIDは、対象スプレッドシートの `拡張機能 > Apps Script` を開いたときのURLから取る。`.clasp.json` はスクリプトIDを含むため `gas/.gitignore` で除外されており、本リポジトリは public のため文書にも書かない。作業はスクラッチパッドで行い、`gas/` を汚さない。

詳しくは `gas/README.md` を参照。

## トラブル対応

| 症状 | 確認すること |
| --- | --- |
| シフトカードにHTML版ボタンが出ない | タスク一覧S列が空でないか。空なら対応表A列の文字列がM列と完全一致しているか |
| S列は埋まっているがDBに入らない | タスク送信を実行したか。`httpUrlOrEmpty_()` が `http(s)` 以外を空に落とすため、メモ書きが入っていないか |
| マニュアルを開くと403 | 個人のGmailでログインしていないか。`〜.nutfes@gmail.com` に切り替える |
| マニュアルを開くと404 | アップロード時の `{id}` とURLの `{id}` が一致しているか |
| アップロードが401 | トークンの値、`Bearer ` の有無。サーバー側で `MANUAL_UPLOAD_TOKEN` が未設定でも401になる |

## 関連する場所

```
リポジトリ
  scripts/automation/prepare_manual.py            zipの展開・配置（①）
  scripts/claude-slide/generate_slide.py          変換スクリプト（②）
  scripts/claude-slide/verify_slide_mechanical.py 機械検証（③）
  scripts/claude-slide/embed_images.py            画像埋め込み・完全性検査（③、依存ゼロ）
  scripts/automation/upload_manual.py             アップロードと貼付行の出力（④）
  .claude/manual-prompt-card-strict.md            変換プロンプト（文章不変版）
  api/lib/internals/controller/manual_controller.go  アップロード・配信エンドポイント
  api/lib/usecase/manual_usecase.go               認証・保存ロジック

ライブのApps Script（45th_シフト_ver0 の 拡張機能 > Apps Script から開く）
  名簿タスク送信.js      タスク送信本体。TASK_COL_MANUAL_URL = 19（S列）
  調査_マニュアルURL.js  checkManualUrlMapping（紐付けの点検） / fillManualUrlFormulas / inspectManualUrlLookup

スプレッドシート
  45th_シフト_ver0        1b5FhiuT7M6kcAM_BkFRu1UVLj-Ssbt-VoBjGmEAN3-I
  マニュアル割り当て      1a2pvM1M8NWQNLNaYnsqTzpB_oGE129-ibpbN3Be1Z1Q

Slack
  #081_執行部マニュアル窓口  C0B65H5FBQ8（マニュアル1本 = スレッド1本）
```

`docs/proposals/manual-slide-operations.md` は2026-06-18時点の手順書で、GitHub Pagesへの配置とSlackスレッドでの部門長レビューを前提にしている。配信経路が自前APIに移った現在、④以降は本書が正となる。①〜③の生成と検証、および部門長レビューの回し方は同文書がなお詳しい。
