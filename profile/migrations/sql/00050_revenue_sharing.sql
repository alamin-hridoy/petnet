-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS revenue_sharing (
   	id uuid NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
   	org_id uuid NOT NULL,
   	user_id uuid NOT NULL,
	remit_type text NOT NULL DEFAULT '', -- The remit type is the different services we are offering for example REMITTANCE, BILLSPAYMENT
	bound_type text NOT NULL DEFAULT '', -- The bound type is two type of transaction one is inbound, other one is outbound
	partner text NOT NULL DEFAULT '',
	transaction_type text NOT NULL DEFAULT '', -- There is transaction type one is digital, other one is otc
	tier_type text NOT NULL DEFAULT '', -- There is tier type is fixed, percentage, fixed_tier, percentage_tier
	amount text NOT NULL DEFAULT '',
	start_date timestamptz DEFAULT NULL,
	created_by text NOT NULL DEFAULT '',
	updated_by text NOT NULL DEFAULT '',
	updated timestamptz NOT NULL DEFAULT NOW(),
	created timestamptz NOT NULL DEFAULT NOW(),
	CONSTRAINT revenue_sharing_org_id_remit_type_bound_type_partner_transaction_type_tier_type UNIQUE (org_id, remit_type, bound_type, partner, transaction_type, tier_type)
);
CREATE TABLE IF NOT EXISTS revenue_sharing_tier (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	revenue_sharing_id uuid NOT NULL REFERENCES revenue_sharing (id),
	min_value text NOT NULL DEFAULT '',
	max_value text NOT NULL DEFAULT '',
	amount text NOT NULL DEFAULT ''
);

comment on column revenue_sharing.remit_type is 'The remit type is the different services we are offering for example REMITTANCE, BILLSPAYMENT';

comment on column revenue_sharing.bound_type is 'The bound type is two type of transaction one is inbound, other one is outbound';

comment on column revenue_sharing.transaction_type is 'There is transaction type one is digital, other one is otc';

comment on column revenue_sharing.tier_type is 'There is tier type is fixed, percentage, fixed_tier, percentage_tier';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS revenue_sharing;
DROP TABLE IF EXISTS revenue_sharing_tier;
