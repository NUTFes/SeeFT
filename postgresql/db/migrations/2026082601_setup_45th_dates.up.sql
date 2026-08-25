-- 45thの年度・日程マスタを投入する（issue #453）
-- shift_usecase.goの日付名→dateID変換は固定値（準備日=1/1日目=2/2日目=3/片付け日=4）
-- のため、45thの日程もその並びに合わせる。準々備日は既存の並びを壊さないよう5番に置く。

INSERT INTO years (id, year)
VALUES (45, 2026)
ON CONFLICT (id) DO NOTHING;

-- 既存の1〜4は43rd(2024年)の日程。45thの日程で上書きする
-- （本番は投入時点でseedのサンプルのみでシフト実データが無いことを確認済み）
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
