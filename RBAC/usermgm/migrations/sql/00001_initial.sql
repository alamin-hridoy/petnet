-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS organization_information (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    org_name text NOT NULL DEFAULT '',
    contact_email text NOT NULL DEFAULT '',
    contact_phone text NOT NULL DEFAULT '',
    active boolean NOT NULL DEFAULT FALSE,
    created timestamptz NOT NULL DEFAULT NOW(),
    updated timestamptz NOT NULL DEFAULT NOW(),
    deleted timestamptz
);

CREATE TABLE IF NOT EXISTS user_account (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    org_id uuid NOT NULL,
    username text NOT NULL DEFAULT '' UNIQUE,
    password TEXT NOT NULL,
    first_name text NOT NULL DEFAULT '',
    last_name text NOT NULL DEFAULT '',
    email text NOT NULL UNIQUE,
    email_verified boolean NOT NULL DEFAULT FALSE,
    invite_status text NOT NULL DEFAULT '',
    invite_sender text NOT NULL DEFAULT '',
    invite_code text NOT NULL DEFAULT '' UNIQUE,
    invite_expiry timestamptz,
    created timestamptz NOT NULL DEFAULT NOW(),
    updated timestamptz NOT NULL DEFAULT NOW(),
    deleted timestamptz
);

CREATE TABLE IF NOT EXISTS password_reset (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id uuid REFERENCES user_account (id) ON DELETE CASCADE NOT NULL,
    expiry timestamptz NOT NULL DEFAULT NOW(),
    created timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE email_confirmation (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1mc (),
    user_id uuid REFERENCES user_account (id) ON DELETE CASCADE NOT NULL,
    created timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS service_account (
    auth_type text NOT NULL DEFAULT '',
    org_id uuid NOT NULL,
    environment text NOT NULL,
    client_name text NOT NULL DEFAULT '',
    client_id text NOT NULL,
    create_user_id text NOT NULL,
    disable_user_id text NOT NULL,
    created timestamptz DEFAULT now(),
    disabled timestamptz
);

CREATE TABLE IF NOT EXISTS google_token (
    id varchar(255) PRIMARY KEY,
    user_id uuid REFERENCES user_account (id) ON DELETE CASCADE NOT NULL,
    hosted_domain text,
    email text,
    email_verified bool NOT NULL DEFAULT TRUE,
    name text,
    picture text,
    given_name text,
    family_name text,
    locale text,
    created timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS github_token (
    id bigint PRIMARY KEY,
    user_id uuid REFERENCES user_account (id) ON DELETE CASCADE NOT NULL,
    email text,
    name text,
    avatar_url text
);

CREATE TABLE IF NOT EXISTS permissions (
    id text PRIMARY KEY,
    org_id uuid NOT NULL,
    service_permission_id uuid NOT NULL,
    permission_name text NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '',
    create_user_id text NOT NULL,
    created timestamptz NOT NULL DEFAULT NOW(),
    updated timestamptz NOT NULL DEFAULT NOW(),
    delete_user_id text,
    deleted timestamptz,

    UNIQUE(org_id, service_permission_id)
);

CREATE TABLE IF NOT EXISTS roles (
    id text PRIMARY KEY,
    org_id uuid NOT NULL,
    role_name text NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '',
    create_user_id uuid NOT NULL,
    created timestamptz NOT NULL DEFAULT NOW(),
    updated timestamptz NOT NULL DEFAULT NOW(),
    delete_user_id uuid,
    deleted timestamptz,

    UNIQUE(org_id, role_name)
);

CREATE TABLE IF NOT EXISTS service_permissions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    resource text NOT NULL DEFAULT '',
    action text NOT NULL DEFAULT '',
    created timestamptz NOT NULL DEFAULT NOW(),
    updated timestamptz NOT NULL DEFAULT NOW(),
    
    UNIQUE(resource, action)
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

CREATE TRIGGER set_update_permission
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

CREATE TRIGGER set_update_role
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

CREATE TRIGGER set_update_service_permission
    BEFORE UPDATE ON service_permissions
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS user_account CASCADE;

DROP TABLE IF EXISTS google_token;

DROP TABLE IF EXISTS organization_information CASCADE;

DROP TABLE IF EXISTS github_token;

DROP TABLE IF EXISTS email_confirmation;

DROP TABLE IF EXISTS permissions CASCADE;

DROP TABLE IF EXISTS roles CASCADE;

DROP TABLE IF EXISTS service_permissions CASCADE;

DROP FUNCTION IF EXISTS trigger_set_timestamp ();

DROP EXTENSION IF EXISTS "uuid-ossp";

