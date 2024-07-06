-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE roles
ALTER COLUMN create_user_id TYPE TEXT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
