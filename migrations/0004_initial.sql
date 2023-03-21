-- +migrate Up
create table if not exists users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);
-- +migrate Down
drop table users;