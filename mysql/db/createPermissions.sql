use seeft_db;

CREATE TABLE permissions (
    user_id int(10) unsigned,
    allow_shift boolean not null default false,
    allow_task boolean not null default false,
    allow_user boolean not null default false,
    allow_property boolean not null default false,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (user_id)
);