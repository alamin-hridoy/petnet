-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remit_cache RENAME COLUMN remco_reference TO remco_control_number;
ALTER TABLE remit_cache RENAME COLUMN remco_alternate_reference TO remco_alternate_control_number;

ALTER TABLE remit_history RENAME COLUMN remco_ref TO remco_control_number;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
