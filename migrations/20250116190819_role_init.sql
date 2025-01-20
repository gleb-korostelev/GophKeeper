-- +goose Up
ALTER TABLE auth.users ADD COLUMN role_changed_at timestamp default (now() at time zone 'utc');

-- +goose Down
ALTER TABLE auth.users DROP COLUMN role_changed_at;
