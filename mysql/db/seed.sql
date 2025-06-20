INSERT INTO weathers
  (weather)
VALUES
  ('晴れ'),
  ('雨'),
  ('none');

INSERT INTO times
  (time)
VALUES
  ('0:00'),
  ('0:15'),
  ('0:30'),
  ('0:45'),
  ('1:00'),
  ('1:15'),
  ('1:30'),
  ('1:45'),
  ('2:00'),
  ('2:15'),
  ('2:30'),
  ('2:45'),
  ('3:00'),
  ('3:15'),
  ('3:30'),
  ('3:45'),
  ('4:00'),
  ('4:15'),
  ('4:30'),
  ('4:45'),
  ('5:00'),
  ('5:15'),
  ('5:30'),
  ('5:45'),
  ('6:00'),
  ('6:15'),
  ('6:30'),
  ('6:45'),
  ('7:00'),
  ('7:15'),
  ('7:30'),
  ('7:45'),
  ('8:00'),
  ('8:15'),
  ('8:30'),
  ('8:45'),
  ('9:00'),
  ('9:15'),
  ('9:30'),
  ('9:45'),
  ('10:00'),
  ('10:15'),
  ('10:30'),
  ('10:45'),
  ('11:00'),
  ('11:15'),
  ('11:30'),
  ('11:45'),
  ('12:00'),
  ('12:15'),
  ('12:30'),
  ('12:45'),
  ('13:00'),
  ('13:15'),
  ('13:30'),
  ('13:45'),
  ('14:00'),
  ('14:15'),
  ('14:30'),
  ('14:45'),
  ('15:00'),
  ('15:15'),
  ('15:30'),
  ('15:45'),
  ('16:00'),
  ('16:15'),
  ('16:30'),
  ('16:45'),
  ('17:00'),
  ('17:15'),
  ('17:30'),
  ('17:45'),
  ('18:00'),
  ('18:15'),
  ('18:30'),
  ('18:45'),
  ('19:00'),
  ('19:15'),
  ('19:30'),
  ('19:45'),
  ('20:00'),
  ('20:15'),
  ('20:30'),
  ('20:45'),
  ('21:00'),
  ('21:15'),
  ('21:30'),
  ('21:45'),
  ('22:00'),
  ('22:15'),
  ('22:30'),
  ('22:45'),
  ('23:00'),
  ('23:15'),
  ('23:30'),
  ('23:45');
  
INSERT INTO grades
  (grade)
VALUES
  ('B1'),
  ('B2'),
  ('B3'),
  ('B4'),
  ('M1'),
  ('M2'),
  ('D1'),
  ('D2'),
  ('D3'),
  ('OB');

INSERT INTO bureaus
  (bureau)
VALUES
  ('執行部'),
  ('執行部補佐'),
  ('総務局'),
  ('企画局'),
  ('渉外局'),
  ('財務局'),
  ('制作局'),
  ('情報局');

INSERT INTO years
  (id, year)
VALUES
  (40, 2021),
  (41, 2022),
  (42, 2023),
  (43, 2024);

INSERT INTO dates
  (year_id, name, date)
VALUES
  (43, '準備日', '2024/09/13'),
  (43, '1日目', '2024/09/14'),
  (43, '2日目', '2024/09/15'),
  (43, '片付け日', '2024/09/16');

INSERT INTO roles 
    (role) 
VALUES
    ('user'),
    ('admin'),
    ('SeeFT Director'),
    ('SeeFT Staff');

INSERT INTO departments 
  (department) 
VALUES 
  ('未所属'),
  ('機械工学分野/機械創造工学課程・機械創造工学専攻'),
  ('電気電子情報工学分野/電気電子情報工学課程/電気電子情報工学専攻'),
  ('情報・経営システム工学分野/情報・経営システム工学課程/情報・経営システム工学専攻'),
  ('物質生物工学分野/物質材料工学課程/生物機能工学課程/物質材料工学専攻/生物機能工学専攻'),
  ('環境社会基盤工学分野/環境社会基盤工学課程/環境社会基盤工学専攻'),
  ('量子・原子力統合工学分野/原子力システム安全工学専攻'),
  ('技術科学イノベーション');

-- placesテーブルの初期データ
INSERT INTO places
  (id, place, remark)
VALUES
  -- (1, '未定', '本部に指示を聞いてください'),
  (1, '本部(電気棟1F)', ''),
  (2, '体育館', '体育館での作業'),
  (3, 'D講', 'D講義室'),
  (4, '屋外', '屋外エリア');

-- tasksテーブルの初期データ
INSERT INTO tasks
  (task, place_id, url, bureau_id, max_member, color, remark, year_id)
VALUES
  ('', 1, '', 1, 1, 'ffffff', '', 43),
  ('NG', 1, '', 1, 1, '949593', '', 43),
  ('テスト1', 2, 'https://example.com', 1, 10, 'ff0000', 'テストタスク1', 43),
  ('テスト2', 3, 'https://example.com', 2, 8, '00ff00', 'テストタスク2', 43),
  ('テスト3', 4, 'https://nutfes.net', 3, 5, '0000ff', 'テストタスク3', 43);

-- INSERT INTO places
--   (place, remark)
-- VALUES
--   ('未定', '本部に指示を聞いてください'),
--   ('体育館', ''),
--   ('D講', ''),
--   ('24', '');

INSERT INTO users
  (name,mail,grade_id,department_id,bureau_id,role_id,student_number,tel,password)
VALUES
  ('root','test1@example.com',1,1,1,1,'12345678','','shiftroot'),
  ('test','test1@testmail.com',1,1,1,1,'12345678','09012345678','123456'),
  ('Admin', 'nutfes@gmail.com', 1, 1, 1, 1, '11111111', '00000000000', 'gidaifes'),
  ('田中太郎','tanaka@example.com',2,2,2,1,'22222222','09011111111','password1'),
  ('佐藤花子','sato@example.com',3,3,3,1,'33333333','09022222222','password2'),
  ('山田次郎','yamada@example.com',4,4,4,1,'44444444','09033333333','password3');

INSERT INTO permissions
  (user_id, allow_shift, allow_task, allow_user, allow_property)
VALUES
  (1, TRUE, TRUE, TRUE, TRUE),
  (2, TRUE, FALSE, FALSE, FALSE),
  (3, TRUE, TRUE, TRUE, TRUE),
  (4, TRUE, FALSE, FALSE, FALSE),
  (5, TRUE, FALSE, FALSE, FALSE),
  (6, TRUE, FALSE, FALSE, FALSE);

-- レスキューテスト用データ
INSERT INTO trouble_rescues
  (user_id, task_id, place, detail, status, response)
VALUES
  (1, 1, '案内所', '機器が故障しています', 'todo', ''),
  (2, 2, '講義棟', 'プリンターが動きません', 'inProgress', '確認中です');

INSERT INTO question_rescues
  (user_id, question, status, response)
VALUES
  (1, '開催時間について教えてください', 'done', '10:00-17:00です'),
  (3, '駐車場の場所はどこですか', 'todo', '');

INSERT INTO shorthanded_rescues
  (user_id, task_id, missing_number, place, status, response)
VALUES
  (2, 1, 2, '案内所', 'todo', ''),
  (1, 3, 1, '会場入口', 'inProgress', '応援を手配中です');

-- INSERT INTO tasks
--   (task, place_id, url, bureau_id, max_member, color, remark, year_id)
-- VALUES
--   ('テスト1', 2, 'https://example.com', 1, 3, 'ff0000', 'テスト用タスク1', 43),
--   ('テスト2', 3, 'https://example.com', 1, 2, '00ff00', 'テスト用タスク2', 43),
--   ('テスト3', 4, 'https://nutfes.net', 1, 4, '0000ff', 'テスト用タスク3', 43),
--   ('NG', 1, '', 1, 1, '949593', '', 43);

-- INSERT INTO shifts
--   (user_id, task_id, year_id, date_id, time_id, weather_id)
-- VALUES
--   -- before_membersのためのデータ（11:45の時間帯）
--   (2, 1, 43, 1, 48, 1),
--   (4, 1, 43, 1, 48, 1),
--   -- メインのシフトデータ（12:00-12:45）
--   (1, 1, 43, 1, 49, 1),
--   (1, 1, 43, 1, 50, 1),
--   (1, 1, 43, 1, 51, 1),
--   (4, 1, 43, 1, 49, 1),
--   (5, 1, 43, 1, 49, 1),
--   (4, 1, 43, 1, 50, 1),
--   (6, 1, 43, 1, 50, 1),
--   (5, 1, 43, 1, 51, 1),
--   (6, 1, 43, 1, 51, 1),
--   -- after_membersのためのデータ（12:45-13:00の時間帯）
--   (3, 1, 43, 1, 52, 1),
--   (5, 1, 43, 1, 52, 1),
--   (6, 1, 43, 1, 52, 1),
--   -- その他のテストデータ
--   (2, 2, 43, 1, 49, 1),
--   (2, 2, 43, 1, 50, 1),
--   (3, 3, 43, 1, 49, 1);
