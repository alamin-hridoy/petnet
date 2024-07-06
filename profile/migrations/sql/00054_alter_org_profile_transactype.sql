-- +goose Up
-- SQL in this section is executed when the migration is applied.
UPDATE org_profile set transaction_types = 'OTC' WHERE transaction_types = '';
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
