-- +goose Up
-- SQL in this section is executed when the migration is applied.
UPDATE remit_history SET txn_staged_time = now(), txn_completed_time = now();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
