-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE user_profile 
ADD COLUMN profile_picture text NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
