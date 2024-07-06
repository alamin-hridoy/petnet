-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS service (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
	 type text NOT NULL,
    org_id uuid NOT NULL REFERENCES org_profile (org_id) UNIQUE,
    service jsonb NULL DEFAULT '{}',
    updated timestamptz NOT NULL DEFAULT NOW(),
    created timestamptz NOT NULL DEFAULT NOW(),
    deleted timestamptz DEFAULT NULL 
);
ALTER TABLE service ADD UNIQUE (org_id, type);
CREATE TRIGGER service_updated
    BEFORE UPDATE ON service 
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS service;
