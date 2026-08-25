-- 45thの日程マスタ投入を巻き戻し、43rd(2024年)のseed相当の状態へ戻す（issue #453）

-- shifts.date_id は dates.id を ON DELETE 指定なしで参照しているため、
-- 45thのシフトが投入済みの状態では巻き戻せない。データを黙って壊さないよう明示的に失敗させる。
DO $$
DECLARE
  ref_shifts INTEGER;
BEGIN
  SELECT count(*) INTO ref_shifts FROM shifts WHERE date_id IN (1, 2, 3, 4, 5);

  IF ref_shifts > 0 THEN
    RAISE EXCEPTION 'dates を参照するシフトが%件あるため巻き戻せません。先にシフトを削除するか、移行方針を決めてください', ref_shifts;
  END IF;
END $$;

DELETE FROM dates WHERE id = 5;

UPDATE dates SET year_id = 43, name = '準備日',   date = '2024/09/13' WHERE id = 1;
UPDATE dates SET year_id = 43, name = '1日目',    date = '2024/09/14' WHERE id = 2;
UPDATE dates SET year_id = 43, name = '2日目',    date = '2024/09/15' WHERE id = 3;
UPDATE dates SET year_id = 43, name = '片付け日', date = '2024/09/16' WHERE id = 4;

SELECT setval('dates_id_seq', (SELECT MAX(id) FROM dates));

-- yearsの45は他のレコードから参照されている可能性があるため削除しない
-- （tasks.year_id等が45を指している場合、削除すると外部キー違反になる）
