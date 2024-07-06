-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE IF EXISTS service_permissions
    ADD COLUMN service_id uuid,
    ADD COLUMN name text NOT NULL DEFAULT '',
    ADD COLUMN description text NOT NULL DEFAULT '',
    ADD CONSTRAINT unique_service_res_act UNIQUE (service_id, resource, action);

CREATE TABLE IF NOT EXISTS service (
    service_id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    service_name text UNIQUE NOT NULL,
    description text NOT NULL DEFAULT 'unknown',
    assign_default boolean NOT NULL DEFAULT FALSE,
    updated timestamptz NOT NULL DEFAULT now(),
    created timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER service_updated
    BEFORE UPDATE ON service
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp ();

CREATE TABLE IF NOT EXISTS service_assignment (
    grant_id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    service_id uuid REFERENCES service (service_id),
    org_id uuid REFERENCES organization_information (id),
    environment text NOT NULL,
    assign_default boolean NOT NULL DEFAULT FALSE,
    assigned timestamptz NOT NULL DEFAULT now(),
    assign_user_id uuid NOT NULL,
    revoke_user_id uuid,
    revoked timestamptz,
    created timestamptz NOT NULL DEFAULT now(),
    UNIQUE (service_id, org_id)
);

CREATE TABLE IF NOT EXISTS service_assignment_audit (
    LIKE service_assignment
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_service_assignment_audit ()
    RETURNS TRIGGER
    AS $$
BEGIN
    INSERT INTO service_assignment_audit
    SELECT
        OLD.*;
    RETURN NULL;
END;
$$
LANGUAGE plpgsql;

-- +goose StatementEnd
CREATE TRIGGER service_assignment_audit_tg
    AFTER DELETE ON service_assignment
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_service_assignment_audit ();

CREATE TABLE IF NOT EXISTS org_permission (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    grant_id uuid REFERENCES service_assignment (grant_id),
    org_id uuid REFERENCES organization_information (id),
    service_id uuid NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    resource text NOT NULL,
    action text NOT NULL,
    environment text NOT NULL,
    permission_id uuid REFERENCES service_permissions (id),
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now(),
    UNIQUE (org_id, permission_id)
);

CREATE TABLE IF NOT EXISTS public_service (
    grant_id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    service_id uuid REFERENCES service (service_id),
    environment text NOT NULL,
    published timestamptz NOT NULL DEFAULT now(),
    published_user_id uuid NOT NULL,
    retracted_user_id uuid,
    retracted timestamptz,
    created timestamptz NOT NULL DEFAULT now(),
    UNIQUE (service_id, environment)
);

CREATE TABLE IF NOT EXISTS public_service_audit (
    LIKE public_service
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_public_service_audit ()
    RETURNS TRIGGER
    AS $$
BEGIN
    INSERT INTO public_service_audit
    SELECT
        OLD.*;
    RETURN NULL;
END;
$$
LANGUAGE plpgsql;

-- +goose StatementEnd
CREATE TRIGGER public_service_audit_tg
    AFTER DELETE ON public_service
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_service_assignment_audit ();

WITH svc_id AS (
INSERT INTO service (service_name, assign_default)
        VALUES ('Account Management', TRUE)
    RETURNING
        service_id)
    UPDATE
        service_permissions
    SET
        service_id = (
            SELECT
                service_id
            FROM
                svc_id
            LIMIT 1)
WHERE
    resource LIKE 'ACCOUNT:%%';

WITH svc_id AS (
INSERT INTO service (service_name, assign_default)
        VALUES ('Role Management', TRUE)
    RETURNING
        service_id)
    UPDATE
        service_permissions
    SET
        service_id = (
            SELECT
                service_id
            FROM
                svc_id
            LIMIT 1)
WHERE
    resource LIKE 'RBAC:%%';

ALTER TABLE IF EXISTS service_permissions
    ALTER COLUMN service_id SET NOT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
