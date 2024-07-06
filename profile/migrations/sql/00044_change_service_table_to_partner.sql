-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE service
RENAME TO partner;
ALTER TABLE service_list
RENAME TO partner_list;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
