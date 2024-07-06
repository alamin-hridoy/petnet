-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS cico_partner_list (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
	stype text NOT NULL,
	name text NOT NULL,
    status text default 'DISABLED' NOT NULL,
    updated timestamptz NOT NULL DEFAULT NOW(),
    created timestamptz NOT NULL DEFAULT NOW(),
    deleted timestamptz DEFAULT NULL 
);

ALTER TABLE cico_partner_list ADD UNIQUE (stype);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS cico_partner_list;
