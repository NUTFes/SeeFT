use seeft_db;

CREATE TABLE permission (
    id int(10) unsigned not null auto_increment,
    allowShift boolean,
    allowTask boolean,
    allowUser boolean,
    allowProperty boolean,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);