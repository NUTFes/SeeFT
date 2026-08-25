-- 45thの日程マスタ投入を巻き戻し、43rd(2024年)のseed相当の状態へ戻す（issue #453）

DELETE FROM dates WHERE id = 5;

UPDATE dates SET year_id = 43, name = '準備日',   date = '2024/09/13' WHERE id = 1;
UPDATE dates SET year_id = 43, name = '1日目',    date = '2024/09/14' WHERE id = 2;
UPDATE dates SET year_id = 43, name = '2日目',    date = '2024/09/15' WHERE id = 3;
UPDATE dates SET year_id = 43, name = '片付け日', date = '2024/09/16' WHERE id = 4;

SELECT setval('dates_id_seq', (SELECT MAX(id) FROM dates));

-- yearsの45は他のレコードから参照されている可能性があるため削除しない
-- （tasks.year_id等が45を指している場合、削除すると外部キー違反になる）
