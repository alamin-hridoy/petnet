-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remittances ADD COLUMN dsa_order_id text NOT NULL DEFAULT '';
ALTER TABLE remittances ALTER COLUMN dsa_order_id DROP DEFAULT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
