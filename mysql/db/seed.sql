use seeft_db;

INSERT INTO weathers
  (`weather`)
VALUES
  ('晴れ'),
  ('雨'),
  ('none');

INSERT INTO times
  (`time`)
VALUES
  ('6:00'),
  ('6:30'),
  ('7:00'),
  ('7:30'),
  ('8:00'),
  ('8:30'),
  ('9:00'),
  ('9:30'),
  ('10:00'),
  ('10:30'),
  ('11:00'),
  ('11:30'),
  ('12:00'),
  ('12:30'),
  ('13:00'),
  ('13:30'),
  ('14:00'),
  ('14:30'),
  ('15:00'),
  ('15:30'),
  ('16:00'),
  ('16:30'),
  ('17:00'),
  ('17:30'),
  ('18:00'),
  ('18:30'),
  ('19:00'),
  ('19:30'),
  ('20:00'),
  ('20:30'),
  ('21:00'),
  ('21:30'),
  ('22:00'),
  ('22:30');
  
INSERT INTO grades
  (`grade`)
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

INSERT INTO dates
  (`date`)
VALUES
  ('準備日'),
  ('1日目'),
  ('2日目'),
  ('片付け日');

INSERT INTO bureaus
  (`bureau`)
VALUES
  ('委員長'),
  ('副委員長'),
  ('総務局'),
  ('企画局'),
  ('渉外局'),
  ('財務局'),
  ('制作局'),
  ('情報局');

INSERT INTO years
  (id, `year`)
VALUES
  (40, 2021),
  (41, 2022),
  (42, 2023);

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
  ('量子・原子力統合工学分野/原子力システム安全工学専攻');

-- 以下テスト用のデータなので本番環境で起こらないようにする
INSERT INTO users
  (`name`, `mail`, `grade_id`, `department_id`, `bureau_id`, `role_id`, `tel`)
VALUES
  ('Admin', 'nutfes@gmail.com', 1, 1, 1, 1, '00000000000');

INSERT INTO permissions
  (`user_id`, `allow_shift`, `allow_task`, `allow_user`, `allow_property`)
VALUES
  (1, True, True, True, True);

INSERT INTO tasks
  (`task`, `place`, `url`, `superviser`, `notes`, `year_id`)
VALUES
  ('テスト1', '24', 'https://example.com', 'Admin', 'a', 42),
  ('テスト2', '体育館', 'https://example.com', 'Admin', 'b', 42),
  ('テスト3', 'D講', 'https://nutfes.net', 'Admin', 'c', 42);

INSERT INTO shifts
  (user_id, task_id, year_id, date_id, time_id, weather_id)
VALUES
  (1, 1, 41, 1, 1, 1),
  (1, 1, 41, 1, 2, 1),
  (1, 1, 41, 1, 3, 1);