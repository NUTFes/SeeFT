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
      │ ①「ウェブページ(.html, zip)」でダウンロード
      ▼
[docs/manuals/{マニュアル名}/ に展開]  source.html + images/
      │ ② generate_slide.py（Claudeが変換）
      ▼
[slide_claude.card-strict.html]  自己完結HTML（CSS・JS・画像を1ファイルに内包）
      │ ③ 目視 + 機械検証
      │ ④ curl で PUT /manuals/{id}
      ▼
[https://seeft-api.nutfes.net/manuals/{id}]  ← レスポンスで manual_url が返る
      │ ⑤ シフトスプシ「マニュアルURL」シートC列へ貼る → タスク送信
      ▼
[tasks.manual_url]  → mobileのシフトカードに「HTML版」ボタンが出る
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

## ① Googleドキュメントを取得する

対象のドキュメントを開き、`ファイル > ダウンロード > ウェブページ（.html、zip形式）` を選ぶ。PDFやWordではなくHTMLである点が重要で、これ以外の形式では画像が取り出せない。

ダウンロードしたzipを展開し、リポジトリの `docs/manuals/` の下にマニュアル名のディレクトリを作って配置する。

```
docs/manuals/45th_企画マニュアル_縁日/
├── source.html        ← zip内のHTMLをこの名前にリネームする
└── images/            ← zip内の画像フォルダをそのまま置く
    ├── image1.png
    └── image2.jpg
```

`source.html` という名前にする理由は、生成スクリプトが「`slide` と `verify` で始まらない `.html`」を入力として探すためで、複数あると曖昧だとして停止するからである（`scripts/claude-slide/generate_slide.py` の `load_source()`）。

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

`verify_report.card-strict.md` が部門長に見せられる日本語サマリ、`verify_mechanical.card-strict.txt` が開発用の詳細である。

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

```bash
export MANUAL_UPLOAD_TOKEN='（共有されたトークン）'
curl -X PUT "https://seeft-api.nutfes.net/manuals/en-nichi" -H "Authorization: Bearer $MANUAL_UPLOAD_TOKEN" -H "Content-Type: text/html" --data-binary @docs/manuals/45th_企画マニュアル_縁日/slide_claude.card-strict.html
```

### レスポンス

成功すると `201 Created` で、次のJSONが返る。

```json
{
  "id": "en-nichi",
  "manual_url": "https://seeft-api.nutfes.net/manuals/en-nichi"
}
```

`manual_url` の値をそのまま⑤で使う。キー名がDBの `tasks.manual_url` と揃えてあるのは、貼り先を取り違えないためである。

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

「マニュアルURL」シートに行を追加する。

- A列: タスク一覧のM列から**コピーして貼る**。手で打たない
- B列: Googleドキュメントの共有URL（ドキュメント版ボタン用）
- C列: ④のレスポンスの `manual_url`（HTML版ボタン用）

タスク一覧のR/S列にVLOOKUPの数式が入っていない場合は、`SeeFTメニュー` から `fillManualUrlFormulas` を実行する。4行目から最終行+50行まで数式を貼るので、行が増えても追従する。

引けているかを `inspectManualUrlLookup` で確認する。対応表の中身と、タスク一覧でM列が埋まっている行のR/S列が表示される。

確認できたら `SeeFTメニュー > タスク送信` を実行する。これでDBの `tasks.manual_url` が更新され、mobileのシフトカードに「HTML版」ボタンが出る。

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

タスク一覧のM列は、多くの場合Googleドキュメントへの**スマートチップ**である。セルに文字列が保存されているのではなく、ドキュメントへの参照が入っており、表示されるのは**参照先ドキュメントの現在のファイル名**である。`getValues()` やIMPORTRANGEはその表示名を文字列として読む。

つまりドキュメントをリネームすると、誰もスプレッドシートを触っていなくてもM列の値が変わる。

```
ドキュメントをリネーム
  → チップの表示が変わる（自動）
  → タスク一覧M列の値が変わる
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

2026-08-29時点で10件以上、マニュアル名ではなくURLがそのまま入っている行がある。これはデータ品質の問題で、そのままではVLOOKUPのキーにならない。マニュアル名に置き換えるか、対応表のA列にそのURL文字列を入れて引かせるかのどちらかで処理する。

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
gas/node_modules/.bin/clasp clone 10tkh3yiMzbEmYjs0I6nmUFfu7_GTnozZnWtCxzL6G43rNLHd4IMCXxJC
```

`.clasp.json` はスクリプトIDを含むため `gas/.gitignore` で除外されている。作業はスクラッチパッドで行う。

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
  scripts/claude-slide/generate_slide.py          変換スクリプト
  scripts/claude-slide/verify_slide_mechanical.py 機械検証
  .claude/manual-prompt-card-strict.md            変換プロンプト（文章不変版）
  api/lib/internals/controller/manual_controller.go  アップロード・配信エンドポイント
  api/lib/usecase/manual_usecase.go               認証・保存ロジック

ライブのApps Script（スクリプトID 10tkh3yiMzbEmYjs0I6nmUFfu7_GTnozZnWtCxzL6G43rNLHd4IMCXxJC）
  名簿タスク送信.js      タスク送信本体。TASK_COL_MANUAL_URL = 19（S列）
  調査_マニュアルURL.js  fillManualUrlFormulas / inspectManualUrlLookup

スプレッドシート
  45th_シフト_ver0        1b5FhiuT7M6kcAM_BkFRu1UVLj-Ssbt-VoBjGmEAN3-I
  マニュアル割り当て      1a2pvM1M8NWQNLNaYnsqTzpB_oGE129-ibpbN3Be1Z1Q

Slack
  #081_執行部マニュアル窓口  C0B65H5FBQ8（マニュアル1本 = スレッド1本）
```

`docs/proposals/manual-slide-operations.md` は2026-06-18時点の手順書で、GitHub Pagesへの配置とSlackスレッドでの部門長レビューを前提にしている。配信経路が自前APIに移った現在、④以降は本書が正となる。①〜③の生成と検証、および部門長レビューの回し方は同文書がなお詳しい。
