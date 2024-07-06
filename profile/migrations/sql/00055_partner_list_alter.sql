-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE partner_list
ADD COLUMN service_name text NOT NULL DEFAULT '',
ADD COLUMN updated_by text NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE partner_list 
DROP COLUMN service_name,
DROP COLUMN updated_by;
