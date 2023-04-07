use seeft_db;

CREATE TABLE IF NOT EXISTS tasks (
    id int(10) unsigned not null auto_increment,
    name varchar(255) not null,
    place varchar(255) not null,
    url varchar(255) not null,
    superviser varchar(255) not null,
    color varchar(255) not null default "fffafa",
    notes varchar(255),
    year_id int(10) unsigned not null,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);

