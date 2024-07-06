-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE micro_insurance_history
    ADD COLUMN IF NOT EXISTS org_id text NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE micro_insurance_history
DROP COLUMN org_id;
