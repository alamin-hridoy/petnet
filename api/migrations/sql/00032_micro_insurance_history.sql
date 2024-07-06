-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS micro_insurance_history (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    dsa_id text NOT NULL,
    coy text NOT NULL,
    location_id text NOT NULL,
    user_code text NOT NULL,
    trx_date timestamptz NOT NULL,
    promo_amount text NOT NULL,
    promo_code text NOT NULL,
    amount text NOT NULL,
    coverage_count text NOT NULL,
    product_code text NOT NULL,
    processing_branch text NOT NULL,
    processed_by text NOT NULL,
    user_email text NOT NULL,
    last_name text NOT NULL,
    first_name text NOT NULL,
    middle_name text NOT NULL,
    gender text NOT NULL,
    birthdate timestamptz NOT NULL,
    mobile_number text NOT NULL,
    province_code text NOT NULL,
    city_code text NOT NULL,
    address text NOT NULL,
    marital_status text NOT NULL,
    occupation text NOT NULL,
    card_number text NOT NULL DEFAULT '',
    number_units text NOT NULL,
    beneficiaries jsonb NOT NULL,
    dependents jsonb NOT NULL,
    trx_status text NOT NULL DEFAULT '',
    trace_number text NULL,
    insurance_details jsonb NOT NULL,
    error_code TEXT NOT NULL DEFAULT '',
    error_message TEXT NOT NULL DEFAULT '',
    error_type TEXT NOT NULL DEFAULT '',
    error_time timestamptz NULL,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now()
);

-- staging issue about existing trigger
DROP TRIGGER IF EXISTS micro_insurance_history_updated ON micro_insurance_history;

CREATE TRIGGER micro_insurance_history_updated
    BEFORE UPDATE ON micro_insurance_history
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TRIGGER IF EXISTS micro_insurance_history_updated ON micro_insurance_history;
DROP TABLE IF EXISTS micro_insurance_history;
