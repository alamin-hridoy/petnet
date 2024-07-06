-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE service ADD COLUMN updated_by text default '' NOT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE service DROP COLUMN updated_by;
