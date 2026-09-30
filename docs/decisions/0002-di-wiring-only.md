# 0002: di.go は配線だけにし、定期実行などは externals の下に切る

- 状態：採用
- 決めた日：2026-04-29
- 決めた人：45th の PM
- 確信度：記録なし
- 出典：コミット `a44d742`（`di.go` から `slackService` を外した）、`7cc59c8`（`scheduler` を作った）

## 背景

`api/lib/di/di.go` は、Repository → UseCase → Controller → Router の順に部品を組み立てる場所である。ここに `slackService` に関する処理が紛れ込んでいた。

その後、通知を5分ごとに送る処理を API のプロセスの中で動かすことになり、その繰り返しをどこに書くかを決める必要があった。

## 候補

| 候補 | 良い点 | 悪い点 |
| --- | --- | --- |
| `di.go` に処理（for ループや goroutine）を書く | 部品がすべてそろっている場所なので、手早く書ける | 組み立てと処理の責務が混ざる。処理だけを取り出してテストできない |
| 処理は `api/lib/externals/` の下に専用のパッケージを切り、`di.go` は組み立てて渡すだけにする | 層の分け方が保たれる。処理を単体で試せる | パッケージが1つ増える |

## 決定

`di.go` は部品の組み立て（配線）だけを担当する。定期実行・バックグラウンドの処理などは、`api/lib/externals/` の下に専用のパッケージを切る。`di.go` では、それを組み立てて `Start()` を1行呼ぶだけにする。

## 理由

配線の層に処理が混ざると、層に分けた構成（controller / usecase / repository と、外の世界を扱う externals）の意味がなくなる。処理が `di.go` にあると、その処理だけを取り出してテストすることもできない。`slackService` が紛れ込んでいたのを外したのも同じ理由である。

外の実行環境を抽象化するパッケージ（`api/lib/externals/server/` など）が、すでに externals の下にあった。定期実行も「外から時間で起こされる」ものなので、同じ位置に置くのが自然だった。

## 前提

- 配線の本体は `api/lib/di/di.go#InitializeServer`
- 定期実行は `api/lib/externals/scheduler/scheduler.go#Scheduler.Start`

## 結果

通知の定期実行は `api/lib/externals/scheduler/` に置き、`di.go` では `scheduler.New(...).Start(ctx)` の1行で起動している。

## 追記
