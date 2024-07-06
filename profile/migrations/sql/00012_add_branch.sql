-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS branch (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    org_profile_id uuid NOT NULL,
    title text NOT NULL DEFAULT '',
    address1 text NOT NULL DEFAULT '',
    city text NOT NULL DEFAULT '',
    state text NOT NULL DEFAULT '',
    postal_code text NOT NULL DEFAULT '',
    phone_number text NOT NULL DEFAULT '',
    fax_number text NOT NULL DEFAULT '',
    contact_person text NOT NULL DEFAULT '',
    updated timestamptz NOT NULL DEFAULT now(),
    created timestamptz NOT NULL DEFAULT now(),
    deleted timestamptz DEFAULT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS branch;
