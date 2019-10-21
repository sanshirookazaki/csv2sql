CREATE DATABASE todo;
use todo;

CREATE TABLE user ( id int NOT NULL AUTO_INCREMENT,
                    name varchar(100) NOT NULL,
                    email varchar(100) NOT NULL,
                    primary key(id));

CREATE TABLE user_task ( id int NOT NULL AUTO_INCREMENT,
                         user_id int NOT NULL,
                         task varchar(100) NOT NULL,
                         primary key(id));

CREATE TABLE user_detail ( user_id int NOT NULL,
                           tel varchar(100) NOT NULL,
                           address varchar(100) NOT NULL
                           );
