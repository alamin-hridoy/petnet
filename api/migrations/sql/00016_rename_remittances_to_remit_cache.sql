-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remittances RENAME TO remit_cache;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
