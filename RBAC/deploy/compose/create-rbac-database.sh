export PGUSER=postgres
for db in usermgmdb hydradb ketodb profiledb; do
	psql <<-EOSQL
		    	 CREATE USER ${db}_user with password 'secret';
		    	 CREATE DATABASE ${db};
		    	 GRANT ALL PRIVILEGES ON DATABASE $db TO ${db}_user;
	EOSQL

	export PGDATABASE=${db}

	for e in ${CREATE_EXTENSION}; do
		psql <<-EOSQL
			     CREATE EXTENSION IF NOT EXISTS "${e}";
		EOSQL
	done
done
