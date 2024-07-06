-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remit_history
	ALTER COLUMN txn_staged_time SET DEFAULT NULL,
	ALTER COLUMN txn_completed_time SET DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
