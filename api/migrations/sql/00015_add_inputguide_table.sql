-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS inputguide (
	partner text NOT NULL,
   inputguide jsonb NULL DEFAULT '{}',
   updated timestamptz NOT NULL DEFAULT NOW(),
   created timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER inputguide_updated
    BEFORE UPDATE ON inputguide
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS inputguide;
