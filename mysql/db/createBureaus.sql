use seeft_db;

CREATE TABLE IF NOT EXISTS bureaus (
    id int(10) unsigned not null auto_increment,
    name varchar(255) not null,
    color varchar(255) not null default "fffafa",
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);
