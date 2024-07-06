-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS org_profile (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    org_id uuid NOT NULL,
    bus_info_company_name text NOT NULL DEFAULT '',
    bus_info_store_name text NOT NULL DEFAULT '',
    bus_info_phone_number text NOT NULL DEFAULT '',
    bus_info_fax_number text NOT NULL DEFAULT '',
    bus_info_website text NOT NULL DEFAULT '',
    bus_info_company_email text NOT NULL DEFAULT '',
    bus_info_contact_person text NOT NULL DEFAULT '',
    bus_info_position text NOT NULL DEFAULT '',
    bus_info_address1 text NOT NULL DEFAULT '',
    bus_info_city text NOT NULL DEFAULT '',
    bus_info_state text NOT NULL DEFAULT '',
    bus_info_postal_code text NOT NULL DEFAULT '',
    bus_info_id_photo_urls text NOT NULL DEFAULT '',
    bus_info_picture_urls text NOT NULL DEFAULT '',
    bus_info_nbi_clearance_urls text NOT NULL DEFAULT '',
    bus_info_court_clearance_url text NOT NULL DEFAULT '',
    bus_info_incorporation_paper_urls text NOT NULL DEFAULT '',
    bus_info_mayors_permit_url text NOT NULL DEFAULT '',
    fin_info_financial_statement_urls text NOT NULL DEFAULT '',
    fin_info_bank_statement_urls text NOT NULL DEFAULT '',
    acc_info_bank text NOT NULL DEFAULT '',
    acc_info_bank_account_number text NOT NULL DEFAULT '',
    acc_info_bank_account_holder text NOT NULL DEFAULT '',
    drp_info_service text NOT NULL DEFAULT '',
    drp_info_questionnaire_urls text NOT NULL DEFAULT '',
    updated timestamptz NOT NULL DEFAULT now(),
    created timestamptz NOT NULL DEFAULT now(),
    deleted timestamptz DEFAULT NULL 
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS org_profile;
