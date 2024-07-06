-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS cico_history (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    org_id uuid NOT NULL,
    svc_provider text DEFAULT '', -- There is service type is gcash cashin, gcash cashout, dragonpay cashin, dragonpay cashout, paymaya cashin, coins cashin, etc
    partner_code text DEFAULT '',
    trx_provider text DEFAULT '',
    trx_type text DEFAULT '',
    reference_number text DEFAULT '',
    petnet_trackingno text DEFAULT '',
    provider_trackingno text DEFAULT '',
    principal_amount text DEFAULT '0',
    charges text DEFAULT '0',
    total_amount text DEFAULT '0',
    trx_date timestamptz NOT NULL,
    txn_status text DEFAULT '',
    error_code text DEFAULT '',
    error_message text DEFAULT '',
    error_time text DEFAULT '',
    error_type text DEFAULT '',
    details jsonb DEFAULT '{}',
    created_by text NOT NULL DEFAULT '',
	updated_by text NOT NULL DEFAULT '',
	updated timestamptz NOT NULL DEFAULT NOW(),
	created timestamptz NOT NULL DEFAULT NOW()
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS cico_history;
