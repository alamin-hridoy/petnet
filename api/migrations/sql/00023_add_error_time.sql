-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remit_history
	ADD COLUMN error_time timestamptz NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
