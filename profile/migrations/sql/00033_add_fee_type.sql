-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE fee_commission ADD COLUMN fee_commision_type smallint DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE fee_commission DROP COLUMN fee_commision_type;
