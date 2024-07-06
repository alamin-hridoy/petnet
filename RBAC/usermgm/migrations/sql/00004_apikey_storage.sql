-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE service_account
    ADD COLUMN IF NOT EXISTS challenge text NOT NULL DEFAULT '',
    ALTER COLUMN disable_user_id SET DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
