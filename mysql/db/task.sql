use seeft_db;

CREATE TABLE task (
    id int(10) unsigned not null auto_increment,
    name varchar(255) not null,
    url varchar(255),
    palce varchar(255) not null,
    president varchar(255) not null,
    tel varchar(255),
    color varchar(255),
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);