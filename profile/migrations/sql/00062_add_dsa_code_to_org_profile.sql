-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile
ADD COLUMN dsa_code text NOT NULL DEFAULT '',
ADD COLUMN terminal_id_otc text NOT NULL DEFAULT '',
ADD COLUMN terminal_id_digital text NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE org_profile 
DROP COLUMN dsa_code,
DROP COLUMN terminal_id_otc,
DROP COLUMN terminal_id_digital;