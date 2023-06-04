CREATE TABLE IF NOT EXISTS users (
    id varchar(64) primary key not null,
    registration_time timestamp without time zone default now() not null,
    password varchar(512) not null,
    username varchar(64) not null unique,
    avatar varchar(512) not null,
    sex varchar(32) not null,
    email varchar(64) not null
);

CREATE TABLE IF NOT EXISTS tasks (
    id varchar(64) primary key not null,
    creation_time timestamp without time zone default now() not null,
    user_id varchar(64) not null references users(id),
    result bytea
);