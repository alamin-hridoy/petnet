-- +goose Up
-- SQL in this section is executed when the migration is applied.
DROP TABLE IF EXISTS wu_employment_position_level;
DROP TABLE IF EXISTS wu_id_types;
DROP TABLE IF EXISTS wu_occupation;
DROP TABLE IF EXISTS wu_purpose_of_transaction;
DROP TABLE IF EXISTS wu_relationship;
DROP TABLE IF EXISTS wu_source_of_funds;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
