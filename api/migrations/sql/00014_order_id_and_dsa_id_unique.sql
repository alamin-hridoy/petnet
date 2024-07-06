-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remit_history ADD UNIQUE (dsa_id,dsa_order_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
