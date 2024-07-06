-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE wu_employment_position_level
DROP COLUMN id,
DROP COLUMN last_updated,
ADD COLUMN updated timestamptz NOT NULL DEFAULT NOW(),
ADD COLUMN id SERIAL PRIMARY KEY;

ALTER TABLE wu_id_types
DROP COLUMN bo_version,
ADD COLUMN updated timestamptz NOT NULL DEFAULT NOW(),
ADD COLUMN id SERIAL PRIMARY KEY;

ALTER TABLE wu_occupation
DROP COLUMN id,
DROP COLUMN position_level,
DROP COLUMN last_updated,
ADD COLUMN updated timestamptz NOT NULL DEFAULT NOW(),
ADD COLUMN id SERIAL PRIMARY KEY;

ALTER TABLE wu_purpose_of_transaction
DROP COLUMN id,
DROP COLUMN last_updated,
ADD COLUMN updated timestamptz NOT NULL DEFAULT NOW(),
ADD COLUMN id SERIAL PRIMARY KEY;

ALTER TABLE wu_relationship
DROP COLUMN id,
DROP COLUMN last_updated,
ADD COLUMN updated timestamptz NOT NULL DEFAULT NOW(),
ADD COLUMN id SERIAL PRIMARY KEY;

ALTER TABLE wu_source_of_funds
DROP COLUMN id,
DROP COLUMN last_updated,
ADD COLUMN updated timestamptz NOT NULL DEFAULT NOW(),
ADD COLUMN id SERIAL PRIMARY KEY;

CREATE TRIGGER wu_employment_position_level_last_updated
    BEFORE UPDATE ON wu_employment_position_level
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER wu_id_types_last_updated
    BEFORE UPDATE ON wu_id_types
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER wu_occupation_last_updated
    BEFORE UPDATE ON wu_occupation
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER wu_purpose_of_transaction_last_updated
    BEFORE UPDATE ON wu_purpose_of_transaction
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER wu_relationship_last_updated
    BEFORE UPDATE ON wu_relationship
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER wu_source_of_funds_last_updated
    BEFORE UPDATE ON wu_source_of_funds
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
