-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile ADD COLUMN acc_info_currency smallint DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE org_profile DROP COLUMN acc_info_currency;
