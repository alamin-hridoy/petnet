-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE service DROP CONSTRAINT IF EXISTS service_org_id_key;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
