-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS fee_commission_rate (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    fee_commission_id uuid NOT NULL,
    min_volume text NOT NULL DEFAULT '',
    max_volume text NOT NULL DEFAULT '',
    txn_rate text NOT NULL DEFAULT '',
    txn_fee jsonb NOT NULL,
    commission_rate text NOT NULL DEFAULT '',
    commission_fee jsonb NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS fee_commission_rate;
