-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS bill_payment (
    bill_payment_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    bill_id text NOT NULL,
    biller_tag text NOT NULL,
    location_id text NOT NULL,
    user_id text NOT NULL,
    sender_member_id text NOT NULL,
    currency_id text NOT NULL,
    account_number text NOT NULL,
    amount text NOT NULL,
    identifier text NOT NULL,
    coy text NOT NULL DEFAULT '',
    service_charge text NOT NULL,
    total_amount text NOT NULL,
    bill_payment_status TEXT NOT NULL DEFAULT '',
    error_code TEXT NOT NULL DEFAULT '',
    error_message TEXT NOT NULL DEFAULT '',
    error_type TEXT NOT NULL DEFAULT '',
    partner_id text NOT NULL,
    biller_name text NOT NULL,
    trx_date timestamptz NOT NULL,
    remote_user_id text NOT NULL,
    customer_id text NOT NULL,
    remote_location_id text NOT NULL,
    location_name text NOT NULL,
    form_type text NOT NULL,
    form_number text NOT NULL,
    payment_method text NOT NULL,
    other_info jsonb NOT NULL,
    bills jsonb NOT NULL,
    bill_payment_date timestamptz NOT NULL,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER bill_payment_updated
    BEFORE UPDATE ON bill_payment
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS bill_payment;
