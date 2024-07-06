-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE file_upload ADD COLUMN file_names text NOT NULL;
ALTER TABLE file_upload RENAME COLUMN bucket_url TO bucket_name;
ALTER TABLE file_upload ADD UNIQUE (org_id, upload_type);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
