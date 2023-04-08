-- +migrate Up
create table if not exists templates
(
    id text primary key,
    name text,
    subject text,
    description text,
    assessment text,
    recommendation text
);
-- +migrate Down
drop table templates;