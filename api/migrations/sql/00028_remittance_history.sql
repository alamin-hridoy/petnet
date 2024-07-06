-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS remittance_history (
    remittance_history_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    dsa_id text DEFAULT '',
    user_id text DEFAULT '',
    phrn text DEFAULT '',
    send_validate_reference_number text DEFAULT '',
    cancel_send_reference_number text DEFAULT '',
    payout_validate_reference_number text DEFAULT '',
    txn_status text DEFAULT '',
    error_code text DEFAULT '',
    error_message text DEFAULT '',
    error_time text DEFAULT '',
    error_type text DEFAULT '',
    details jsonb DEFAULT '{}',
    remarks text DEFAULT '',
    txn_created_time timestamptz DEFAULT NOW(),
    txn_updated_time timestamptz DEFAULT NULL,
    txn_confirm_time timestamptz DEFAULT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS payment_history;
DROP TABLE IF EXISTS remittance_history;
