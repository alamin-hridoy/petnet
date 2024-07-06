-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS remittances (
    transaction_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    dsa_id uuid NOT NULL,
    user_id text NOT NULL,
    remco_id text NOT NULL,
    remco_member_id text NOT NULL DEFAULT '',
    remco_reference text NOT NULL,
    remco_alternate_reference text NOT NULL DEFAULT '',
    status text NOT NULL,
    remittance jsonb NOT NULL,
    updated timestamptz NOT NULL DEFAULT now(),
    created timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS remittance_history (
    transaction_id uuid NOT NULL,
    dsa_id uuid NOT NULL,
    user_id text NOT NULL,
    remco_id text NOT NULL,
    remco_member_id text NOT NULL DEFAULT '',
    remco_reference text NOT NULL,
    remco_alternate_reference text NOT NULL,
    -- transaction fields
    txn_type text NOT NULL,
    source_currency text NOT NULL,
    destination_currency text NOT NULL,
    exchange_rate text NOT NULL DEFAULT '',
    source_gross_amount text NOT NULL,
    destination_remit_amount text NOT NULL,
    additional_charge_amount text NOT NULL,
    additional_charge_currency text NOT NULL,
    taxes jsonb NOT NULL,
    charges jsonb NOT NULL,
    promo_code text NOT NULL DEFAULT '',
    promo_description text NOT NULL DEFAULT '',
    message text[],
    -- destination fields
    origin_city text NOT NULL,
    origin_state text NOT NULL,
    destination_city text NOT NULL,
    destination_state text NOT NULL,
    bank_bic text NOT NULL DEFAULT '',
    bank_location text NOT NULL DEFAULT '',
    bank_account_no text NOT NULL DEFAULT '',
    bank_account_suffix text NOT NULL DEFAULT '',
    -- kyc fields
    sending_reason text NOT NULL,
    sender_relationship text NOT NULL,
    -- remitter fields
    remitter_fname text NOT NULL,
    remitter_mname text NOT NULL,
    remitter_lname text NOT NULL,
    remitter_gender text NOT NULL,
    remitter_address1 text NOT NULL,
    remitter_address2 text NOT NULL DEFAULT '',
    remitter_city text NOT NULL,
    remitter_state text NOT NULL,
    remitter_postal_code text NOT NULL,
    remitter_country text NOT NULL,
    remitter_phone_country text NOT NULL,
    remitter_phone_number text NOT NULL,
    remitter_mobile_country text NOT NULL,
    remitter_mobile_number text NOT NULL,
    remitter_email text NOT NULL DEFAULT '',
    -- receiver fields
    receiver_fname text NOT NULL,
    receiver_mname text NOT NULL,
    receiver_lname text NOT NULL,
    receiver_address1 text NOT NULL,
    receiver_address2 text NOT NULL DEFAULT '',
    receiver_city text NOT NULL,
    receiver_state text NOT NULL,
    receiver_postal_code text NOT NULL,
    receiver_country text NOT NULL,
    receiver_phone_country text NOT NULL DEFAULT '',
    receiver_phone_number text NOT NULL DEFAULT '',
    receiver_mobile_country text NOT NULL DEFAULT '',
    receiver_mobile_number text NOT NULL DEFAULT '',
    -- error fields
    error_code text NOT NULL DEFAULT '',
    error_message text NOT NULL DEFAULT '',
    error_details text NOT NULL DEFAULT '',
    errored timestamptz,
    -- timestamps
    paid timestamptz,
    processed timestamptz,
    updated timestamptz NOT NULL DEFAULT now(),
    created timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
