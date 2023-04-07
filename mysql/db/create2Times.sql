use seeft_db;

CREATE TABLE IF NOT EXISTS times (
    id int(10) unsigned not null auto_increment,
    time varchar(255) not null,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);