-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE file_upload 
ADD COLUMN submitted smallint DEFAULT 0,
ADD COLUMN checked timestamptz DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
