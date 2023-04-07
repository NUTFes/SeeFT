use seeft_db;

CREATE TABLE IF NOT EXISTS shifts (
    id int(10) unsigned not null auto_increment,
    task_id int(10) unsigned not null,
    user_id int(10) unsigned not null,
    year_id int(10) unsigned not null,
    date_id int(10) unsigned not null,
    time_id int(10) unsigned not null,
    weather_id int(10) unsigned not null,
    is_attendance boolean default false,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id), 
    FOREIGN KEY (task_id) REFERENCES tasks (id) ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE,
    FOREIGN KEY (year_id) REFERENCES years (id) ON UPDATE CASCADE,
    FOREIGN KEY (date_id) REFERENCES dates (id) ON UPDATE CASCADE,
    FOREIGN KEY (time_id) REFERENCES times (id) ON UPDATE CASCADE,
    FOREIGN KEY (weather_id) REFERENCES weathers (id) ON UPDATE CASCADE
);
