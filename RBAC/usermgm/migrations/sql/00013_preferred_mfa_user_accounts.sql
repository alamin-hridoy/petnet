-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE user_account
    ADD COLUMN IF NOT EXISTS preferred_mfa text DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
