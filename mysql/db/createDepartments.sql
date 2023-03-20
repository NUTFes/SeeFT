use seeft_db;

CREATE TABLE departments (
  id int(10) unsigned not null auto_increment,
  name varchar(255),
  created_at datetime not null default current_timestamp,
  updated_at datetime not null default current_timestamp on update current_timestamp,
  PRIMARY KEY (id)
);