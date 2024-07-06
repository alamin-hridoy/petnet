-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE IF EXISTS remit_history
	ADD COLUMN txn_status TEXT NOT NULL DEFAULT '',
	ADD COLUMN txn_step TEXT NOT NULL DEFAULT '',
	ADD COLUMN error_code TEXT NOT NULL DEFAULT '',
	ADD COLUMN error_message TEXT NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
