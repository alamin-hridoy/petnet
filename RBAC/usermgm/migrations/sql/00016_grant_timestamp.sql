-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE consent_grant RENAME COLUMN updated TO timestamp;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE consent_grant RENAME COLUMN timestamp TO updated;

