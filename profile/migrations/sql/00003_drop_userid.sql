-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile DROP COLUMN user_id;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
