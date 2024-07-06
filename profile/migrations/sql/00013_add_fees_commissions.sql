-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS fee_commission (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    org_profile_id uuid NOT NULL,
    fee_amount text NOT NULL DEFAULT '',
    commission_amount text NOT NULL DEFAULT '',
    start_date timestamptz DEFAULT NULL,
	 end_date timestamptz DEFAULT NULL,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now(),
    deleted timestamptz DEFAULT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS fee_commission;
