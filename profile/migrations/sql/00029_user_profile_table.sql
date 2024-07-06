-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS user_profile (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL UNIQUE,
    org_id uuid NOT NULL,
    session_expiry timestamptz DEFAULT NULL,
    updated timestamptz NOT NULL DEFAULT now(),
    created timestamptz NOT NULL DEFAULT now(),
    deleted timestamptz DEFAULT NULL 
);
CREATE TRIGGER user_profile_updated
    BEFORE UPDATE ON user_profile 
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS user_profile;
