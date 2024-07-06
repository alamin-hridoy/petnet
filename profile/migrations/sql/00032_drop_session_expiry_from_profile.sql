-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE user_profile 
DROP COLUMN session_expiry;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
