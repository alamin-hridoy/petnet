-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS upload_service_request (
   id uuid NOT NULL DEFAULT uuid_generate_v4(),
   org_id uuid NOT NULL,
	partner text NOT NULL,
	service_name text NOT NULL,
   status text NOT NULL DEFAULT '',
   file_type text NOT NULL DEFAULT '',
   file_id text NOT NULL DEFAULT '',
	create_by text NOT NULL DEFAULT '',
	verify_by text NOT NULL DEFAULT '',
   verified timestamptz DEFAULT NULL,
   created timestamptz NOT NULL DEFAULT NOW()
);
ALTER TABLE upload_service_request ADD UNIQUE (org_id, service_name, file_type);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS upload_service_request;
