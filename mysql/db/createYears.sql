use seeft_db;

CREATE TABLE years (
    id int(10) unsigned unique not null,
    year int(10) unsigned unique not null,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);
