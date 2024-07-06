-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS file_upload (
    file_id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    org_id uuid NOT NULL REFERENCES org_profile (org_id),
    user_id uuid NOT NULL,
    upload_type text NOT NULL,
    bucket_url text NOT NULL,
    created timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
