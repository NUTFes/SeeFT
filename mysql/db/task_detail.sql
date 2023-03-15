use seeft-db;

CREATE TABLE task_detail (
    id int(10) unsigned not null auto_increment,
    task varchar(255) not null,
    palce varchar(255) not null,
    url varchar(255),
    superviser varchar(255) not null,
    notes varchar(255),
    yaer_id int(10) not null,
    users varchar(255) not null,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);