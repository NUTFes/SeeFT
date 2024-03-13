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

INSERT INTO dates
  (date)
VALUES
  ('準備日'),
  ('1日目'),
  ('2日目'),
  ('片付け日');

INSERT INTO bureaus
  (bureau)
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
  (id, year)
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
  ('量子・原子力統合工学分野/原子力システム安全工学専攻'),
  ('技術科学イノベーション');

INSERT INTO users
  (name,mail,grade_id,department_id,bureau_id,role_id,student_number,tel,password)
VALUES
  ('test','test1@testmail.com','1','1','1','1','12345678','09012345678','123456');
-- 以下テスト用のデータなので本番環境で起こらないようにする
-- INSERT INTO users
--   (`name`, `mail`, `grade_id`, `department_id`, `bureau_id`, `role_id`, 'student_number', `tel`, `password`)
-- VALUES
--   ('Admin', 'nutfes@gmail.com', 1, 1, 1, 1, 11111111, '00000000000', "gidaifes");

-- INSERT INTO permissions
--   (`user_id`, `allow_shift`, `allow_task`, `allow_user`, `allow_property`)
-- VALUES
--   (1, True, True, True, True);

-- INSERT INTO tasks
--   (`task`, `place`, `url`, `superviser`, `notes`, `year_id`)
-- VALUES
--   ('テスト1', '24', 'https://example.com', 'Admin', 'a', 42),
--   ('テスト2', '体育館', 'https://example.com', 'Admin', 'b', 42),
--   ('テスト3', 'D講', 'https://nutfes.net', 'Admin', 'c', 42);

-- INSERT INTO shifts
--   (user_id, task_id, year_id, date_id, time_id, weather_id)
-- VALUES
--   (1, 1, 41, 1, 1, 1),
--   (1, 1, 41, 1, 2, 1),
--   (1, 1, 41, 1, 3, 1);