use seeft_db;

CREATE TABLE users (
  id int(10) unsigned not null auto_increment,
  name varchar(255) not null,
  mail varchar(255) not null unique,
  grade_id int(10) not null,
  department_id int(10) not null,
  bureau_id int(10) not null,
  role_id int(10) not null,
  created_at datetime not null default current_timestamp,
  updated_at datetime not null default current_timestamp on update current_timestamp,
  PRIMARY KEY (id)
);
