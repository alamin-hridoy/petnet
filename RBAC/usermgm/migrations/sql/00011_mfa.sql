-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS mfa (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL,
    mfa_type text NOT NULL,
    token text NOT NULL DEFAULT '',
    active boolean NOT NULL DEFAULT FALSE,
    revoked timestamptz,
    deadline timestamptz NOT NULL,
    confirmed timestamptz,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS mfa_event (
    event_id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL,
    mfa_id uuid NOT NULL REFERENCES mfa (id),
    mfa_type text NOT NULL,
    active boolean NOT NULL DEFAULT TRUE,
    token text NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '',
    initiated timestamptz NOT NULL DEFAULT now(),
    deadline timestamptz NOT NULL,
    confirmed timestamptz
);

CREATE INDEX user_mfa ON mfa (user_id);

CREATE INDEX user_mfa_type ON mfa (user_id, mfa_type);

CREATE INDEX user_mfa_active ON mfa (user_id, active);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
