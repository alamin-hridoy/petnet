-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile ADD CONSTRAINT org_id_unique UNIQUE (org_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
