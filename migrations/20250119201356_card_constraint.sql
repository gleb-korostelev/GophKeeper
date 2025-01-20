-- +goose Up
ALTER TABLE auth.cards
ADD CONSTRAINT unique_user_card UNIQUE (user_id, card_number);

-- +goose Down
ALTER TABLE auth.cards
DROP CONSTRAINT IF EXISTS unique_user_card;
