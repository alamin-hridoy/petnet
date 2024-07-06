-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE revenue_sharing 
DROP COLUMN IF EXISTS start_date,
DROP COLUMN IF EXISTS end_date;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

