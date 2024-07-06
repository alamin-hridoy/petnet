-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS user_session (
    customer_code text NOT NULL UNIQUE,
    session_data jsonb NOT NULL,
    updated timestamptz NOT NULL DEFAULT now(),
    created timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER remittances_updated
    BEFORE UPDATE ON remittances
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

CREATE TRIGGER remittance_history_updated
    BEFORE UPDATE ON remittance_history
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

CREATE TRIGGER user_session_updated
    BEFORE UPDATE ON user_session
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
