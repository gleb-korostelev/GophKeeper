-- +goose Up
create table if not exists auth.cards
(
    id              bigint generated always as identity primary key,
    user_id         bigint not null references auth.users(id) on delete cascade,
    card_number     text not null unique,
    card_holder     text not null,
    expiration_date date not null,
    cvv             text not null,
    metadata        text not null,
    created_at      timestamp default (now() at time zone 'utc'),
    updated_at      timestamp default (now() at time zone 'utc')
);


-- +goose Down

DROP TABLE IF EXISTS auth.cards;