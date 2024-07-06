-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS service_request (
   id uuid NOT NULL DEFAULT uuid_generate_v4(),
   org_id uuid NOT NULL,
	partner text NOT NULL,
	service_name text NOT NULL,
	company_name text NOT NULL,
   remarks text NOT NULL DEFAULT '',
   status text NOT NULL DEFAULT '',
   enabled bool NOT NULL DEFAULT false,
	updated_by text NOT NULL DEFAULT '',
   applied timestamptz DEFAULT NULL,
   updated timestamptz NOT NULL DEFAULT NOW(),
   created timestamptz NOT NULL DEFAULT NOW()
);
ALTER TABLE service_request ADD UNIQUE (org_id, partner, service_name);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS service_request;
