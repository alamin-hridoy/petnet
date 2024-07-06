-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remit_history ADD COLUMN dsa_order_id text NOT NULL DEFAULT '';
ALTER TABLE remit_history ALTER COLUMN dsa_order_id DROP DEFAULT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
