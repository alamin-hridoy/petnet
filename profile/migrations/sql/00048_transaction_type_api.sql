-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS api_key_transaction_type (
   id uuid NOT NULL DEFAULT uuid_generate_v4(),
   org_id uuid NOT NULL,
   user_id uuid NOT NULL,
	client_id text NOT NULL,
	environment text NOT NULL,
   transaction_type text NOT NULL,
   created timestamptz NOT NULL DEFAULT NOW(),
   deleted timestamptz DEFAULT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS api_key_transaction_type;
