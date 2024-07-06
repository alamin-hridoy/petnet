-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile 
	ADD COLUMN wu_coy text NOT NULL DEFAULT '',
	ADD COLUMN wu_operator_id text NOT NULL DEFAULT '',
	ADD COLUMN wu_terminal_id text NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE org_profile 
	DROP COLUMN wu_coy,
	DROP COLUMN wu_operator_id,
	DROP COLUMN wu_terminal_id;
