#!/bin/sh

export PGPASSWORD=$(DATABASE_PASSWORD)

echo "About to get client_ids to delete hydra clients data"
CLIENTIDS=$(psql -qtX --host '$(DATABASE_HOST)' -U '$(DATABASE_USER)' -d 'cebcorp-usermgm' -p '$(DATABASE_PORT)' -c "SELECT client_id FROM service_account WHERE client_name LIKE '%E2E%';")
SVCPERMIDS=$(psql -qtX --host '$(DATABASE_HOST)' -U '$(DATABASE_USER)' -d 'cebcorp-usermgm' -p '$(DATABASE_PORT)' -c "SELECT service_permission_id FROM permissions WHERE permission_name like '%E2E%' GROUP BY service_permission_id;")

CLIENTIDS=`echo $CLIENTIDS | sed -e "s/\(^\|$\)/'/" -e "s/$/'/g" -e "s/ /','/g"`
SVCPERMIDS=`echo $SVCPERMIDS | sed -e "s/\(^\|$\)/'/" -e "s/$/'/g" -e "s/ /','/g"`

echo "About to delete e2e resource in USERMGM DB"
psql --host '$(DATABASE_HOST)' -U '$(DATABASE_USER)' -d 'cebcorp-usermgm' -p '$(DATABASE_PORT)' -v svc_perm_ids=$SVCPERMIDS <<-EOSQL
	DELETE FROM organization_information WHERE org_name like '%e2e%';
	DELETE FROM roles WHERE role_name like '%e2e%';
	DELETE FROM permissions WHERE permission_name like '%E2E%';
	DELETE FROM service_account WHERE client_name like '%E2E%';
	DELETE FROM service_permissions WHERE id IN (:svc_perm_ids);
	DELETE FROM user_account WHERE username like '%e2e%';
	DELETE FROM user_account WHERE username like '%email.webhook.site%';
EOSQL
	
echo "About to delete e2e resource in KETO DB"
psql --host '$(DATABASE_HOST)' -U '$(DATABASE_USER)' -d 'cebcorp-keto' -p '$(DATABASE_PORT)' -v client_ids=$CLIENTIDS <<-EOSQL
	DELETE FROM rego_data where document ->> 'description' LIKE '%E2E%';
	DELETE FROM rego_data where document ->> 'members' IN (:client_ids);
EOSQL

echo "About to delete e2e resource in HYDRA DB"
psql --host '$(DATABASE_HOST)' -U '$(DATABASE_USER)' -d 'cebcorp-hydra' -p '$(DATABASE_PORT)' -v client_ids=$CLIENTIDS <<-EOSQL
	DELETE FROM hydra_client where id IN (:client_ids);
EOSQL
