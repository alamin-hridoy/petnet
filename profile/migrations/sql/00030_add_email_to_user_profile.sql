-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE user_profile 
ADD COLUMN email text NOT NULL unique;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
