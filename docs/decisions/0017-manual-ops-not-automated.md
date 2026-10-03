# 0017: マニュアル運用のスプレッドシートの操作は自動化しない

- 状態：見送り
- 決めた日：2026-06-18（状態の監視。手順書 `docs/proposals/manual-slide-operations.md` を足したコミット 0f4c921 の日付）、2026-08-31（対応表への記入。PR #487 をマージした日）
- 決めた人：45th の PM
- 決定の信頼度：記録なし
  - 根拠：`docs/proposals/manual-slide-operations.md`・`docs/proposals/manual-proposal-v4-slides/automation-design.md`・#486・PR #487 には、決めたことと理由、見直す条件（状態の管理を自動化したくなったとき）は書かれているが、どれだけ確かだと思っていたかは書かれていない
- 出典：`docs/proposals/manual-slide-operations.md`、`docs/proposals/manual-proposal-v4-slides/automation-design.md`、#486、PR #487

## 背景

解説マニュアルを作って配るまでには、生成・確認・部門長のチェック・アップロード・タスクへの紐付けと、いくつもの手作業がある。

2026 年 5 月に、この流れをスプレッドシートで管理する自動化を設計し、コードも書いた。スプレッドシートの行を常に監視し、状態が変わったら次の作業を進める仕組みである（`scripts/automation/` の watcher・sheets_client・drive_client）。

8 月には、手順を文書にして誰でも作業できるようにしたが、マニュアル1本を配るのに手作業が 10 手かかるのは変わっていなかった（#486）。

## 候補

| 候補 | 良い点 | 悪い点 |
| --- | --- | --- |
| スプレッドシートを中心に、状態の監視から対応表への記入まで自動化する | 手作業が減る | 部門長にスプレッドシートを開いてもらう必要がある。常に監視するには、watcher をどこかのマシンで常駐させるか cron で回し続ける必要がある。Google の API を使うための設定（GCP のプロジェクトと認証情報）が、作業する人ごとに要る |
| 状態の管理は Slack のスレッドで行い、手作業を減らせるところだけコマンドにする | 部門長は Slack だけ見ればよい。新しく入った人も、設定なしで作業を始められる | 対応表への記入とタスクの送信は、手作業のまま残る |

## 決定

マニュアル運用のスプレッドシートの操作は、自動化しない。

- 2026-06-18：本番は、Slack のスレッドを中心にした手作業の流れで運用する。状態の監視の自動化は、コードを残したまま保留にする（`docs/proposals/manual-slide-operations.md`）
- 2026-08-31：手作業を 10 手から 7 手に減らした（PR #487）。残る7手のうち、対応表への記入とタスクの送信の2手は自動化しない

## 理由

状態の監視を保留にしたのは、部門長にスプレッドシートを開かせずに済ませたかったからである。状態の管理は SeeFT の仕事にして、部門長は Slack のスレッドだけを見れば済む形にした。常に監視するには watcher をどこかのマシンで動かし続ける必要があり、ほぼ1人の運用ではそこまで手が回らなかった。

対応表への記入を自動化しなかったのは、Google の API を使うと、作業する人ごとに GCP のプロジェクトや認証情報の設定が要るからである。PR #487 では、「新しく入った人が最初の1本を出すまでの時間」で比べると、自動化すると逆に遅くなると判断した。

## 前提

- 保留にした自動化のコード：`scripts/automation/watcher.py`、`scripts/automation/sheets_client.py`、`scripts/automation/drive_client.py`。今の運用では使っていない
- 設計：`docs/proposals/manual-proposal-v4-slides/automation-design.md`
- 今の運用の手順：`docs/proposals/manual-slide-operations.md`、`docs/development/manual-html-operations.md`
- ほぼ1人で運用していること。担当する人が増えて、状態の管理が Slack で追いきれなくなったら、自動化を考え直す

## 結果

45th は、Slack のスレッドを中心にした手作業の流れで、マニュアルを配った。

## 追記
