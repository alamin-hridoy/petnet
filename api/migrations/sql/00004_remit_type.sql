-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE IF EXISTS remittances
    ADD COLUMN remit_type text NOT NULL DEFAULT 'Send';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
