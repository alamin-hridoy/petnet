-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE service ADD COLUMN status text default 'DISABLED' NOT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE service DROP COLUMN status;
