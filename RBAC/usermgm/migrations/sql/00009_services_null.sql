-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE public_service
    ALTER COLUMN retracted_user_id TYPE text,
    ALTER COLUMN retracted_user_id SET DEFAULT '';

ALTER TABLE service_assignment
    ALTER COLUMN revoke_user_id TYPE text,
    ALTER COLUMN revoke_user_id SET DEFAULT '';

ALTER TABLE service_assignment_audit
    ALTER COLUMN revoke_user_id TYPE text,
    ALTER COLUMN revoke_user_id SET DEFAULT '';

UPDATE public_service SET retracted_user_id = '' WHERE retracted_user_id IS NULL;

UPDATE service_assignment SET revoke_user_id = '' WHERE revoke_user_id IS NULL;

UPDATE service_assignment_audit SET revoke_user_id = '' WHERE revoke_user_id IS NULL;

ALTER TABLE public_service
    ALTER COLUMN retracted_user_id SET NOT NULL;

ALTER TABLE service_assignment
    ALTER COLUMN revoke_user_id SET NOT NULL;

ALTER TABLE service_assignment_audit
    ALTER COLUMN revoke_user_id SET NOT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
