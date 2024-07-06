-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE partner_list
ADD COLUMN is_provider BOOLEAN DEFAULT false;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE partner_list 
DROP COLUMN is_provider;
