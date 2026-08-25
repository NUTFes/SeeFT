-- 45thの年度・日程マスタを投入する（issue #453）
-- shift_usecase.goの日付名→dateID変換は固定値（準備日=1/1日目=2/2日目=3/片付け日=4/準々備日=5）
-- のため、45thの日程もその並びに合わせる。準々備日は既存の並びを壊さないよう5番に置く。

-- 固定IDの既存行を書き換えるため、想定外のデータを黙って壊さないよう前提を検証する。
-- 許可する状態は次の3つだけで、いずれにも当てはまらない場合はmigrationを失敗させる。
--   (a) 初期状態: datesが空（初期化直後で、seedをまだ流していないDB）
--   (b) 未適用:   id 1〜4が43rd(2024年)のseed値のままで、id 5が存在しない
--   (c) 適用済み: id 1〜5が本migrationの投入値と完全一致（seed適用後の再実行を含む）
DO $$
DECLARE
  total_rows   INTEGER;
  seed_rows    INTEGER;
  applied_rows INTEGER;
  extra_rows   INTEGER;
  year45       INTEGER;
  ref_shifts   INTEGER;
BEGIN
  -- years.id = 45 が既にある場合、yearが2026でなければ想定外の状態
  SELECT year INTO year45 FROM years WHERE id = 45;
  IF year45 IS NOT NULL AND year45 <> 2026 THEN
    RAISE EXCEPTION 'years.id = 45 の year が想定(2026)と異なります: %', year45;
  END IF;

  -- (a) 初期状態: 1行も無ければこのmigrationがそのまま45thの日程を作る
  SELECT count(*) INTO total_rows FROM dates;

  -- (b) 未適用: 43rdのseed値と完全一致し、id 5が未作成
  SELECT count(*) INTO seed_rows
  FROM dates
  WHERE (id, year_id, name,       date)
     IN ((1, 43, '準備日',   '2024/09/13'),
         (2, 43, '1日目',    '2024/09/14'),
         (3, 43, '2日目',    '2024/09/15'),
         (4, 43, '片付け日', '2024/09/16'));

  SELECT count(*) INTO extra_rows FROM dates WHERE id = 5;

  -- (c) 適用済み: 本migrationの投入値と完全一致（seed.sqlが投入した状態もこれに一致する）
  SELECT count(*) INTO applied_rows
  FROM dates
  WHERE (id, year_id, name,       date)
     IN ((1, 45, '準備日',   '2026/09/18'),
         (2, 45, '1日目',    '2026/09/19'),
         (3, 45, '2日目',    '2026/09/20'),
         (4, 45, '片付け日', '2026/09/21'),
         (5, 45, '準々備日', '2026/09/17'));

  IF NOT (total_rows = 0 OR (seed_rows = 4 AND extra_rows = 0) OR applied_rows = 5) THEN
    RAISE EXCEPTION 'dates の id 1〜5 が想定と異なります（空、43rdのseed値のまま、本migration適用済みのいずれかである必要があります）。手動で確認してください';
  END IF;

  -- 45th以外のシフトが対象の日程を参照している場合、上書きすると別日程を指す状態になる
  SELECT count(*) INTO ref_shifts
  FROM shifts
  WHERE date_id IN (1, 2, 3, 4, 5) AND year_id <> 45;

  IF ref_shifts > 0 THEN
    RAISE EXCEPTION '45th以外のシフトが dates.id 1〜5 を参照しています(%件)。日程の上書きは行えません', ref_shifts;
  END IF;
END $$;

INSERT INTO years (id, year)
VALUES (45, 2026)
ON CONFLICT (id) DO NOTHING;

-- 既存行があれば45thの日程で上書きし、無ければ作成する。
-- 初期化直後のDB（datesが空）でもUPDATEが0行にならないよう、全idをINSERT ... ON CONFLICTで揃える
INSERT INTO dates (id, year_id, name, date)
VALUES (1, 45, '準備日',   '2026/09/18'),
       (2, 45, '1日目',    '2026/09/19'),
       (3, 45, '2日目',    '2026/09/20'),
       (4, 45, '片付け日', '2026/09/21'),
       (5, 45, '準々備日', '2026/09/17')
ON CONFLICT (id) DO UPDATE
  SET year_id = EXCLUDED.year_id,
      name    = EXCLUDED.name,
      date    = EXCLUDED.date;

-- idを明示指定したため、SERIALの採番カウンタを最大idに追従させる
-- （これを忘れると次のINSERTで主キー重複になる）
SELECT setval('dates_id_seq', (SELECT MAX(id) FROM dates));
