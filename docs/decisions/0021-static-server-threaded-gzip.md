# 0021: mobile の静的配信は server.py のまま、並行処理と gzip の圧縮にする

- 状態：採用
- 決めた日：2026-09-14（#518 の開発の開始日。PR #519 は同じ日にマージ）
- 決めた人：45th の PM
- 決めたときの自信：記録なし
- 出典：#518、PR #519、`docs/development/load-test-plan.md` の issue 6

## 背景

本番の mobile は、Flutter Web のビルド成果物を `mobile/python/server.py` が配っている。当時の `server.py` は `socketserver.TCPServer` で、1本のスレッドで1件ずつ返し、圧縮もしていなかった。

負荷試験の計画（`docs/development/load-test-plan.md`）の issue 6 は、朝に大勢が一斉に開くと、API より先にこの静的配信が詰まるおそれがあると挙げていた。

そこに、日本語フォントの同梱（[0020](0020-bundle-subset-japanese-font.md)）で、初めて開く人1人あたりに `server.py` が送る量が 2.51MB から 8.03MB（約3倍）に増えることになった。本番の Cloudflare は、どのファイルもキャッシュしていなかった（`cf-cache-status: DYNAMIC`）。`server.py` から Cloudflare までは、圧縮しないまま流れる（#518。2026-09-13 に本番の応答ヘッダーで確かめた）。

## 候補

| 候補 | 良い点 | 悪い点 |
| --- | --- | --- |
| `server.py` を並行処理にし、gzip で圧縮して返す | 数行の変更で済む。起動の仕方（compose）を変えなくてよい | Python の標準のサーバーのままで、同時の接続が増えるとスレッドが増える |
| nginx などに置き換える | 静的配信に向いたサーバーを使える | イメージと compose を作り直す必要がある |
| Cloudflare でキャッシュする | `server.py` に届く量そのものが減る | Cloudflare の設定は SeeFT の外の管理者に相談が要る |

nginx と Cloudflare の案は、load-test-plan.md と #518 に候補として挙がっているが、比べて退けた記録は無い。

## 決定

`server.py` の `socketserver.TCPServer` を `http.server.ThreadingHTTPServer` に置き換え、リクエストごとにスレッドで処理する。

- `Accept-Encoding` に gzip があれば、js・json・html・wasm・フォントを gzip で圧縮して返す。圧縮はファイルごとに1回だけ行い、結果をメモリに持つ
- 接続の受け付け待ちの列（`request_queue_size`）を、標準の 5 から 1024 に広げる。スレッドの数に上限（スレッドプール）は付けない
- フォントの同梱（PR #515）だけが先に本番に出ないよう、PR #519 を先にマージする

## 理由

フォントの同梱で、初回の負担が約3倍になるのを避けたかった（PM の作業メモ）。この変更で、初回に `server.py` が送る量は 8.03MB から 4.03MB に減る（PR #519）。

受け付け待ちの列を広げたのは、負荷試験の結果による。遅い接続を 4,200 本（名簿の約 350 人 × 初回に取るファイル約 12）張ると、標準の 5 では 1,273〜1,737 本が列から溢れて切られた。1024 にすると 0 本だった（PR #519 の本文）。

スレッドの数に上限を付けなかったのは、4,200 スレッドが同時に動いてもメモリが最大 879MiB（4GiB の約 22%）で足りたからである（PR #519 の本文）。上限を付けると待ちが戻る、という判断は PM の作業メモにある。

nginx などに置き換えなかった理由は、記録に無い。数行の変更で技大祭（9/19）の前に間に合うから、というのは推測である。

## 前提

- 配信：`mobile/python/server.py#Server`、圧縮：`mobile/python/server.py#gzipped_body`、`mobile/python/server.py#MyHandler.accepts_gzip`
- 本番の mobile のコンテナとホストに、スレッドの数やメモリを制限する上限が無いこと（2026-09-14 に本番で確かめた。PR #519 の本文）。上限を付けるなら、この ADR を見直す
- `request_queue_size` の 1024 が、OS の `somaxconn`（本番は 4096）に収まっていること
- フォントの同梱（[0020](0020-bundle-subset-japanese-font.md)）を続けていること。PR #519 だけを元に戻すと、初回に送る量が 8.03MB に戻る

## 結果

PR #519 を 2026-09-14 にマージし、その直後に PR #515 をマージした。

## 追記
