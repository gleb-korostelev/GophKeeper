-- +goose Up
create schema if not exists auth;

create table if not exists auth.users
(
    id              bigint generated always as identity primary key,
    username        text  not null unique,
    secret          bytea not null,
    account_type    int       default 1,
    created_at      timestamp default (now() at time zone 'utc'),
    updated_at      timestamp default (now() at time zone 'utc')
);


-- +goose Down

DROP TABLE IF EXISTS auth.users;