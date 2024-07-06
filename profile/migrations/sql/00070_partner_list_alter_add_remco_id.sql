-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE partner_list
ADD COLUMN remco_id text NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE partner_list 
DROP COLUMN remco_id;
