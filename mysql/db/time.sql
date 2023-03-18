use seeft_db;

CREATE TABLE time (
    id int(10) unsigned not null auto_increment,
    time varchar(255) not null,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);

INSERT into time (time) values ("8:00");
INSERT into time (time) values ("8:30");
INSERT into time (time) values ("9:00");
INSERT into time (time) values ("9:30");
INSERT into time (time) values ("10:00");
INSERT into time (time) values ("10:30");
INSERT into time (time) values ("11:00");
INSERT into time (time) values ("11:30");
INSERT into time (time) values ("12:00");
INSERT into time (time) values ("12:30");
INSERT into time (time) values ("13:00");
INSERT into time (time) values ("13:30");
INSERT into time (time) values ("14:00");
INSERT into time (time) values ("14:30");
INSERT into time (time) values ("15:00");
INSERT into time (time) values ("15:30");
INSERT into time (time) values ("16:00");
INSERT into time (time) values ("16:30");
INSERT into time (time) values ("17:00");
INSERT into time (time) values ("17:30");