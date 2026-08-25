-- 45thの年度・日程マスタを投入する（issue #453）
-- shift_usecase.goの日付名→dateID変換は固定値（準備日=1/1日目=2/2日目=3/片付け日=4/準々備日=5）
-- のため、45thの日程もその並びに合わせる。準々備日は既存の並びを壊さないよう5番に置く。

-- 既存のdates 1〜4を書き換えるため、想定外のデータを壊さないよう前提を検証する。
-- 期待する状態: id 1〜4が43rd(2024年)のseed値のまま、かつそれらを参照するシフトが無いこと。
DO $$
DECLARE
  seed_rows INTEGER;
  ref_shifts INTEGER;
BEGIN
  SELECT count(*) INTO seed_rows
  FROM dates
  WHERE (id, year_id, name,       date)
     IN ((1, 43, '準備日',   '2024/09/13'),
         (2, 43, '1日目',    '2024/09/14'),
         (3, 43, '2日目',    '2024/09/15'),
         (4, 43, '片付け日', '2024/09/16'));

  -- 既に45thへ移行済み（このmigrationを適用済み）の場合は再実行を許可する
  IF seed_rows <> 4 THEN
    SELECT count(*) INTO seed_rows
    FROM dates
    WHERE id IN (1, 2, 3, 4) AND year_id = 45;

    IF seed_rows <> 4 THEN
      RAISE EXCEPTION 'dates の id 1〜4 が想定(43rdのseed値または45th適用済み)と異なります。手動で確認してください';
    END IF;
  END IF;

  SELECT count(*) INTO ref_shifts
  FROM shifts
  WHERE date_id IN (1, 2, 3, 4) AND year_id <> 45;

  IF ref_shifts > 0 THEN
    RAISE EXCEPTION '45th以外のシフトが dates.id 1〜4 を参照しています(%件)。日程の上書きは行えません', ref_shifts;
  END IF;
END $$;

INSERT INTO years (id, year)
VALUES (45, 2026)
ON CONFLICT (id) DO NOTHING;

-- 既存の1〜4は43rd(2024年)の日程。45thの日程で上書きする
UPDATE dates SET year_id = 45, name = '準備日',   date = '2026/09/18' WHERE id = 1;
UPDATE dates SET year_id = 45, name = '1日目',    date = '2026/09/19' WHERE id = 2;
UPDATE dates SET year_id = 45, name = '2日目',    date = '2026/09/20' WHERE id = 3;
UPDATE dates SET year_id = 45, name = '片付け日', date = '2026/09/21' WHERE id = 4;

INSERT INTO dates (id, year_id, name, date)
VALUES (5, 45, '準々備日', '2026/09/17')
ON CONFLICT (id) DO UPDATE
  SET year_id = EXCLUDED.year_id,
      name    = EXCLUDED.name,
      date    = EXCLUDED.date;

-- idを明示指定したため、SERIALの採番カウンタを最大idに追従させる
-- （これを忘れると次のINSERTで主キー重複になる）
SELECT setval('dates_id_seq', (SELECT MAX(id) FROM dates));
