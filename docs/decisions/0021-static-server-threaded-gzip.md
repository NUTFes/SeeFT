# 0021: mobileの静的配信はserver.pyのまま、並行処理とgzipの圧縮にする

- 状態：採用
- 決めた日：2026-09-14（#518の開発の開始日。PR #519は同じ日にマージ）
- 決めた人：45thのPM
- 決定の信頼度：記録なし
  - 根拠：#518・PR #519の本文と`docs/development/load-test-plan.md`のissue 6には、決めたことと計測の結果は書かれているが、どれだけ確かだと思っていたかは書かれていない
- 出典：#518、PR #519、`docs/development/load-test-plan.md`のissue 6

## 背景

本番のmobileは、Flutter Webのビルド成果物を`mobile/python/server.py`が配っている。当時の`server.py`は`socketserver.TCPServer`で、1本のスレッドで1件ずつ返し、圧縮もしていなかった。

負荷試験の計画（`docs/development/load-test-plan.md`）のissue 6は、朝に大勢が一斉に開くと、APIより先にこの静的配信が詰まるおそれがあると挙げていた。

そこに、日本語フォントの同梱（[0020](0020-bundle-subset-japanese-font.md)）で、初めて開く人1人あたりに`server.py`が送る量が2.51MBから8.03MB（約3倍）に増えることになった。本番では、途中の経路でファイルがキャッシュされず、`server.py`が送った量がそのまま流れていた（#518。2026-09-13に本番の応答で確かめた）。

## 候補

| 候補 | 良い点 | 悪い点 |
| --- | --- | --- |
| `server.py`を並行処理にし、gzipで圧縮して返す | 数行の変更で済む。起動の仕方（compose）を変えなくてよい | Pythonの標準のサーバーのままで、同時の接続が増えるとスレッドが増える |
| nginxなどに置き換える | 静的配信に向いたサーバーを使える | イメージとcomposeを作り直す必要がある |
| Cloudflareでキャッシュする | `server.py`に届く量そのものが減る | Cloudflareの設定はSeeFTの外の管理者に相談が要る |

nginxとCloudflareの案は、load-test-plan.mdと#518に候補として挙がっているが、比べて退けた記録は無い。

## 決定

`server.py`の`socketserver.TCPServer`を`http.server.ThreadingHTTPServer`に置き換え、リクエストごとにスレッドで処理する。

- html・js・json・css・wasm・フォント・svg・txtをgzipで圧縮して返す。圧縮するのは、`Accept-Encoding`のgzipに付いたqの値が0より大きく1以下のときだけ（qが無ければ1とみなす）。qが0のときや、数として読めないときは圧縮しない。圧縮はファイルごとに1回だけ行い、結果をメモリに持つ
- 接続の受け付け待ちの列（`request_queue_size`）を、標準の5から1024に広げる。スレッドの数に上限（スレッドプール）は付けない
- フォントの同梱（PR #515）だけが先に本番に出ないよう、PR #519を先にマージする

## 理由

フォントの同梱で、初回の負担が約3倍になるのを避けたかった（PMの作業メモ）。この変更で、初回に`server.py`が送る量は8.03MBから4.03MBに減る（PR #519）。

受け付け待ちの列を広げたのは、負荷試験の結果による。遅い接続を4,200本（名簿の約350人 × 初回に取るファイル約12）張ると、標準の5では1,273〜1,737本が列から溢れて切られた。1024にすると0本だった（PR #519の本文）。

スレッドの数に上限を付けなかったのは、4,200スレッドが同時に動いてもメモリが最大879MiB（4GiBの約22%）で足りたからである（PR #519の本文）。上限を付けると待ちが戻る、という判断はPMの作業メモにある。

nginxなどに置き換えなかった理由は、記録に無い。数行の変更で技大祭（9/19）の前に間に合うから、というのは推測である。

## 前提

- 配信は`mobile/python/server.py#Server`、圧縮は`mobile/python/server.py#gzipped_body`と`mobile/python/server.py#MyHandler.accepts_gzip`、圧縮する拡張子は`mobile/python/server.py#COMPRESSIBLE_EXTENSIONS`にある
- 本番のサーバーに、この負荷を受けられるだけのスレッドの数とメモリの余裕があること（2026-09-14に本番で確かめた。PR #519の本文）。スレッドの数やメモリに上限を付けるときや、もっと小さいサーバーに移すときは、このADRを見直す
- `request_queue_size`の1024が、OSの受け付け待ちの列の上限（`somaxconn`）に収まっていること
- フォントの同梱（[0020](0020-bundle-subset-japanese-font.md)）を続けていること。PR #519だけを元に戻すと、初回に送る量が8.03MBに戻る

## 結果

PR #519を2026-09-14にマージし、その直後にPR #515をマージした。

## 追記
