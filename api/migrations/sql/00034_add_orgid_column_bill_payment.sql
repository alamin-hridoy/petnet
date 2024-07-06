-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE bill_payment
    ADD COLUMN IF NOT EXISTS org_id text NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE bill_payment
DROP COLUMN org_id;
