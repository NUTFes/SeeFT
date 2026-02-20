-- シフト変更通知のテスト用: 既存の action_logs を「未送信の UPDATE」に更新する
-- 実行後、最大5分以内にワーカーが処理し、Slack に「〇〇 → △△」形式の変更通知が届く

UPDATE action_logs
SET
  action_type = 'UPDATE',
  diff_payload = '{"changes":[{"field":"task_name","old":"テストタスク","new":"変更後タスク"}]}'::jsonb,
  is_sent = false
WHERE id = 1;
