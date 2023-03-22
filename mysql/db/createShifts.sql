use seeft_db;

CREATE TABLE IF NOT EXISTS shifts (
    id int(10) unsigned not null auto_increment,
    user_id int(10) not null,
    task_id int(10) not null,
    year_id int(10) not null,
    date_id int(10) not null,
    time_id int(10) not null,
    weather_id int(10) not null,
    is_attendance boolean default false,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);
