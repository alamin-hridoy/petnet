-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remit_history
	ADD COLUMN error_type TEXT NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
