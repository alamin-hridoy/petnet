-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remit_history ADD COLUMN txn_staged_time timestamptz;
ALTER TABLE remit_history ALTER COLUMN txn_staged_time SET DEFAULT now();
ALTER TABLE remit_history ADD COLUMN txn_completed_time timestamptz;
ALTER TABLE remit_history ALTER COLUMN txn_completed_time SET DEFAULT now();
ALTER TABLE remit_history DROP COLUMN created;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
