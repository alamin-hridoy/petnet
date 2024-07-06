-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile
    ALTER COLUMN user_id TYPE text;

UPDATE org_profile SET user_id = '' WHERE user_id IS NULL;

ALTER TABLE org_profile
    ALTER COLUMN user_id SET NOT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
