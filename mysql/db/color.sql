use seeft_db;

CREATE TABLE color (
    id int(10) unsigned not null auto_increment,
    value int(10) not null,
    PRIMARY KEY (id)
);

INSERT into color (value) values (1);