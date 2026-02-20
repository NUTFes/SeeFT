-- shifts の id=1 の一部カラムを更新する
-- task_id: 1 → 4（タスク「テスト1」に変更）
-- weather_id: 1 → 2（晴れ → 雨）

UPDATE shifts
SET
  task_id = 4,
  weather_id = 2,
  updated_at = CURRENT_TIMESTAMP
WHERE id = 1;
