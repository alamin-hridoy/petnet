-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remit_history ADD COLUMN transaction_type text default null;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE remit_history DROP COLUMN transaction_type;
