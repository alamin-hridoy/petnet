-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE fee_commission_rate
DROP COLUMN txn_fee,
DROP COLUMN commission_rate,
DROP COLUMN commission_fee;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
