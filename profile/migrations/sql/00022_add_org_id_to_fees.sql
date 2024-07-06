-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE fee_commission ADD COLUMN org_id uuid DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
