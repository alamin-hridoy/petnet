-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE roles ADD COLUMN IF NOT EXISTS updateduid text NOT NULL DEFAULT '';
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
