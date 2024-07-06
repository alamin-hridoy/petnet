-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS remit_history (
    remit_id uuid NOT NULL,
    dsa_id text NOT NULL,
    user_id text NOT NULL,
    remco_id text NOT NULL,
    sender_member_id text NOT NULL,
    receiver_member_id text NOT NULL,
    remco_ref text NOT NULL,
    remittance jsonb NOT NULL,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER remit_history_updated
    BEFORE UPDATE ON remit_history
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS remit_history;

