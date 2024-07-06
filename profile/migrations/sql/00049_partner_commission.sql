-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS partner_commission_config (
	id uuid NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
	remit_type text NOT NULL DEFAULT '', -- The remit type is the different services we are offering for example REMITTANCE, BILLSPAYMENT
	bound_type text NOT NULL DEFAULT '', -- The bound type is two type of transaction one is inbound, other one is outbound
	partner text NOT NULL DEFAULT '',
	transaction_type text NOT NULL DEFAULT '', -- There is transaction type one is digital, other one is otc
	tier_type text NOT NULL DEFAULT '', -- There is tier type is fixed, percentage, fixed_tier, percentage_tier
	amount text NOT NULL DEFAULT '',
	start_date timestamptz DEFAULT NULL,
	end_date timestamptz DEFAULT NULL,
	created_by text NOT NULL DEFAULT '',
	updated_by text NOT NULL DEFAULT '',
	updated timestamptz NOT NULL DEFAULT NOW(),
	created timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS partner_commission_tier (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	partner_commission_config_id uuid NOT NULL REFERENCES partner_commission_config (id),
	min_value text NOT NULL DEFAULT '',
	max_value text NOT NULL DEFAULT '',
	amount text NOT NULL DEFAULT ''
);

comment on column partner_commission_config.remit_type is 'The remit type is the different services we are offering for example REMITTANCE, BILLSPAYMENT';

comment on column partner_commission_config.bound_type is 'The bound type is two type of transaction one is inbound, other one is outbound';

comment on column partner_commission_config.transaction_type is 'There is transaction type one is digital, other one is otc';

comment on column partner_commission_config.tier_type is 'There is tier type is fixed, percentage, fixed_tier, percentage_tier';
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS partner_commission_config;
DROP TABLE IF EXISTS partner_commission_tier;