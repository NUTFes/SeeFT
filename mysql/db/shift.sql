use seeft-db;

CREATE TABLE shift (
    id int(10) unsigned not null auto_increment,
    user_id int(10) not null,
    date varchar(255) not null,
    time varchar(255) not null,
    work_id int(10) not null,
    weather varchar(255) not null,
    attendance boolean,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);

INSERT into shift (user_id, date, time, work_id, weather, attendance) values (1, "2023-03-16", "10:00", 1, "椅子", false);