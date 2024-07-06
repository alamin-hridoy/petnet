-- +goose Up
-- SQL in this section is executed when the migration is applied.
DELETE from wu_id_types;
DELETE from wu_occupation;
DELETE from wu_employment_position_level;
DELETE from wu_purpose_of_transaction;
DELETE from wu_relationship;
DELETE from wu_source_of_funds;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
