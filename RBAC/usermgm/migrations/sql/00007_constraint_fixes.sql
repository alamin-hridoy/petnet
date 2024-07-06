-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE service_assignment
    DROP CONSTRAINT service_assignment_service_id_org_id_key,
    ADD CONSTRAINT duplicate_service_assignment UNIQUE (service_id, org_id, environment);

ALTER TABLE org_permission
    DROP CONSTRAINT org_permission_org_id_permission_id_key,
    ADD CONSTRAINT duplicate_org_permission UNIQUE (org_id, permission_id, environment);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
