-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE mfa_event
    ADD COLUMN IF NOT EXISTS attempt int NOT NULL DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE mfa_event DROP COLUMN IF EXISTS attempt;

