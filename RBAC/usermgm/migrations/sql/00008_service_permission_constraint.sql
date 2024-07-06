-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE service_permissions
    DROP CONSTRAINT service_permissions_resource_action_key;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
