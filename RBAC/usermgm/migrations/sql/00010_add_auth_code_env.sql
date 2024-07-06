-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE oauth2_client ADD COLUMN environment text NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE oauth2_client DROP COLUMN environment;
