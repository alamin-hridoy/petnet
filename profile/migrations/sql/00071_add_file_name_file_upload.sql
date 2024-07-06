-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE file_upload 
ADD COLUMN file_name TEXT NOT NULL DEFAULT '';
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE file_upload 
DROP COLUMN file_name;
