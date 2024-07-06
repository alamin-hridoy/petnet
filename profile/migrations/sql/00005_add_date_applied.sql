-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile ADD COLUMN date_applied timestamptz DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE org_profile DROP COLUMN date_applied;
