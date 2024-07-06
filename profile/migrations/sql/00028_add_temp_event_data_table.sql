-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS temp_event_data (
    event_id uuid PRIMARY KEY NOT NULL,
	 resource text NOT NULL,
	 action text NOT NULL,
    data jsonb NULL DEFAULT '{}',
    updated timestamptz NOT NULL DEFAULT NOW(),
    created timestamptz NOT NULL DEFAULT NOW(),
    deleted timestamptz DEFAULT NULL 
);
CREATE TRIGGER event_data_updated
    BEFORE UPDATE ON temp_event_data 
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS temp_event_data;
