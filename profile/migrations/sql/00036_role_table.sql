-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS role (
    id uuid NOT NULL,
    role_data jsonb NOT NULL DEFAULT '{}',
    created_by text NOT NULL,
    updated_by text NOT NULL,
    updated timestamptz NOT NULL DEFAULT NOW(),
    created timestamptz NOT NULL DEFAULT NOW(),
    deleted timestamptz DEFAULT NULL 
);
CREATE TRIGGER role_updated
    BEFORE UPDATE ON role 
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS role;
