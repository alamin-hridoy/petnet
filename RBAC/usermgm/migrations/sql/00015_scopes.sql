-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS scopes (
    id text NOT NULL PRIMARY KEY,
    name text NOT NULL,
    group_name text NOT NULL,
    description text NOT NULL,
    updated timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS scope_group (
    name text NOT NULL UNIQUE,
    description text NOT NULL DEFAULT '',
    updated timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS consent_grant (
    grant_id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id text NOT NULL,
    client_id text NOT NULL,
    owner_id text NOT NULL,
    scopes text[] NOT NULL,
    updated timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS consent_grant;

DROP TABLE IF EXISTS scope_group;

DROP TABLE IF EXISTS scopes;
