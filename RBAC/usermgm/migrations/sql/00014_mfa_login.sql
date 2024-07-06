-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE user_account
    ADD COLUMN IF NOT EXISTS mfa_login boolean DEFAULT FALSE;

ALTER TABLE organization_information
    ADD COLUMN IF NOT EXISTS mfa_login boolean DEFAULT FALSE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
