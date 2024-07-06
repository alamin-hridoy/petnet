-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS change_password (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id uuid REFERENCES user_account (id) ON DELETE CASCADE NOT NULL,
    event_id uuid NOT NULL,
    new_password TEXT NOT NULL,
    created timestamptz NOT NULL DEFAULT NOW()
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.