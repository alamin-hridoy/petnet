-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_account (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    org_id uuid NOT NULL,
    created timestamptz NOT NULL DEFAULT NOW(),
    updated timestamptz NOT NULL DEFAULT NOW(),
    deleted timestamptz
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_timestamp ()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.updated = NOW();
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- +goose StatementEnd
CREATE TRIGGER set_timestamp_user
    BEFORE UPDATE ON user_account
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS user_account CASCADE;

DROP FUNCTION IF EXISTS trigger_set_timestamp ();

DROP EXTENSION IF EXISTS "uuid-ossp";
