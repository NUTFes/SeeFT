use seeft_db;

CREATE TABLE grade (
    id int(10) unsigned not null auto_increment,
    grade varchar(255) not null,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);

INSERT into grade (grade) values ('B1');
INSERT into grade (grade) values ('B2');
INSERT into grade (grade) values ('B3');
INSERT into grade (grade) values ('B4');
INSERT into grade (grade) values ('M1');
INSERT into grade (grade) values ('M2');
INSERT into grade (grade) values ('D1');
INSERT into grade (grade) values ('D2');
INSERT into grade (grade) values ('D3');