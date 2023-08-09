use seeft_db;

CREATE TABLE IF NOT EXISTS users (
  id int(10) unsigned not null auto_increment,
  name varchar(255) not null,
  mail varchar(255) not null unique,
  grade_id int(10) unsigned not null,
  department_id int(10) unsigned not null,
  bureau_id int(10) unsigned not null,
  role_id int(10) unsigned not null,
  student_number int(10) unsigned not null,
  tel varchar(15) not null unique,
  password varchar(255) not null,
  created_at datetime not null default current_timestamp,
  updated_at datetime not null default current_timestamp on update current_timestamp,
  PRIMARY KEY (id),
  FOREIGN KEY (grade_id) REFERENCES grades (id) ON UPDATE CASCADE,
  FOREIGN KEY (department_id) REFERENCES departments (id) ON UPDATE CASCADE,
  FOREIGN KEY (bureau_id) REFERENCES bureaus (id) ON UPDATE CASCADE,
  FOREIGN KEY (role_id) REFERENCES roles (id) ON UPDATE CASCADE
);
