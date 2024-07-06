-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE permissions DROP CONSTRAINT permissions_org_id_service_permission_id_key; 
ALTER TABLE permissions ADD CONSTRAINT org_id_svc_perm_del_key UNIQUE (org_id, service_permission_id, deleted);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
