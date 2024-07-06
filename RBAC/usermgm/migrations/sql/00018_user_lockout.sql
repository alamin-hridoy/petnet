-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE user_account
    ADD COLUMN IF NOT EXISTS locked timestamptz,
    ADD COLUMN IF NOT EXISTS last_login timestamptz,
    ADD COLUMN IF NOT EXISTS last_failed timestamptz,
    ADD COLUMN IF NOT EXISTS fail_count int NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reset_required timestamptz;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE mfa_event
    DROP COLUMN IF EXISTS locked,
    DROP COLUMN IF EXISTS last_login,
    DROP COLUMN IF EXISTS last_failed,
    DROP COLUMN IF EXISTS fail_count,
    DROP COLUMN IF EXISTS reset_required;

