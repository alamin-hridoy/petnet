// Package postgres provides common function to creating database connection and migration from config file.
// Consumers of this package are expected to copy following lines to envcfg config file.
// [database]
// host="$DATABASE_HOST||localhost"
// port="$DATABASE_PORT||5432"
// user="$DATABASE_USER||postgres"
// password="$DATABASE_PASSWORD||"
// dbname="$DATABASE_NAME||postgres"
// migrationDir="$DATABASE_MIGRATION_DIR||migrations/sql"
// sslMode="$DATABASE_SSL_MODE||disable"
// connectionTimeout="$DATABASE_CONNECTION_TIMEOUT||30000"
// statementTimeout="$DATABASE_STATEMENT_TIMEOUT||30000"
// idleInTransactionSessionTimeout="$DATABASE_IDLE_IN_TRANSACTION_SESSION_TIMEOUT||30000"
package postgres
