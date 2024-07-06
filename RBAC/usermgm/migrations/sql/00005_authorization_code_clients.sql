-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS oauth2_client (
    org_id text NOT NULL,
    client_id text NOT NULL,
    client_name text NOT NULL,
    created_user_id text NOT NULL,
    updated_user_id text NOT NULL,
    deleted_user_id text NOT NULL DEFAULT '',
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now(),
    deleted timestamptz
);

CREATE TRIGGER set_timestamp_oauth2_client
    BEFORE UPDATE ON oauth2_client
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
