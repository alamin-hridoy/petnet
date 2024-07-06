-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE mfa_event
    ADD COLUMN IF NOT EXISTS expired boolean NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS validation boolean NOT NULL DEFAULT FALSE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
