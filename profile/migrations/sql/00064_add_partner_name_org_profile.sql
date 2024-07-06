-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile
ADD COLUMN partner text DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE org_profile 
DROP COLUMN partner;
