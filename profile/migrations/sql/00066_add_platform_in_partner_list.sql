-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE partner_list
ADD COLUMN platform text DEFAULT 'Perahub';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE partner_list 
DROP COLUMN platform;
