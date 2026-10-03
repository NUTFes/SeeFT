# 0017: マニュアル運用のスプレッドシートの操作は自動化しない

- 状態：見送り
- 決めた日：2026-06-18（状態の監視。手順書`docs/proposals/manual-slide-operations.md`を足したコミット0f4c921の日付）、2026-08-31（対応表への記入。PR #487をマージした日）
- 決めた人：45thのPM
- 決定の信頼度：記録なし
  - 根拠：`docs/proposals/manual-slide-operations.md`・`docs/proposals/manual-proposal-v4-slides/automation-design.md`・#486・PR #487には、決めたことと理由、見直す条件（状態の管理を自動化したくなったとき）は書かれているが、どれだけ確かだと思っていたかは書かれていない
- 出典：`docs/proposals/manual-slide-operations.md`、`docs/proposals/manual-proposal-v4-slides/automation-design.md`、#486、PR #487

## 背景

解説マニュアルを作って配るまでには、生成・確認・部門長のチェック・アップロード・タスクへの紐付けと、いくつもの手作業がある。

2026年5月に、この流れをスプレッドシートで管理する自動化を設計し、コードも書いた。スプレッドシートの行を常に監視し、状態が変わったら次の作業を進める仕組みである（`scripts/automation/`のwatcher・sheets_client・drive_client）。

8月には、手順を文書にして誰でも作業できるようにしたが、マニュアル1本を配るのに手作業が10手かかるのは変わっていなかった（#486）。

## 候補

| 候補 | 良い点 | 悪い点 |
| --- | --- | --- |
| スプレッドシートを中心に、状態の監視から対応表への記入まで自動化する | 手作業が減る | 部門長にスプレッドシートを開いてもらう必要がある。常に監視する仕組みは重い。GoogleのAPIを使うための設定（GCPのプロジェクトと認証情報）が、作業する人ごとに要る |
| 状態の管理はSlackのスレッドで行い、手作業を減らせるところだけコマンドにする | 部門長はSlackだけ見ればよい。新しく入った人も、設定なしで作業を始められる | 対応表への記入とタスクの送信は、手作業のまま残る |

## 決定

マニュアル運用のスプレッドシートの操作は、自動化しない。

- 2026-06-18：本番は、Slackのスレッドを中心にした手作業の流れで運用する。状態の監視の自動化は、コードを残したまま保留にする（`docs/proposals/manual-slide-operations.md`）
- 2026-08-31：手作業を10手から7手に減らした（PR #487）。残る7手のうち、対応表への記入とタスクの送信の2手は自動化しない

## 理由

状態の監視を保留にしたのは、部門長にスプレッドシートを開かせずに済ませたかったからである。状態の管理はSeeFTの仕事にして、部門長はSlackのスレッドだけを見れば済む形にした。常に監視する仕組みは、ほぼ1人の運用には重すぎた。

対応表への記入を自動化しなかったのは、GoogleのAPIを使うと、作業する人ごとにGCPのプロジェクトや認証情報の設定が要るからである。PR #487では、「新しく入った人が最初の1本を出すまでの時間」で比べると、自動化すると逆に遅くなると判断した。

## 前提

- 保留にした自動化のコードは`scripts/automation/watcher.py`、`scripts/automation/sheets_client.py`、`scripts/automation/drive_client.py`で、今の運用では使っていない
- 設計は`docs/proposals/manual-proposal-v4-slides/automation-design.md`にある
- 今の運用の手順は`docs/proposals/manual-slide-operations.md`と`docs/development/manual-html-operations.md`にある
- ほぼ1人で運用していること。担当する人が増えて、状態の管理がSlackで追いきれなくなったら、自動化を考え直す

## 結果

45thは、Slackのスレッドを中心にした手作業の流れで、マニュアルを配った。

## 追記
