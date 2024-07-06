// Package postgres provides common function to creating database connection and migration from config file.
// Consumers of this package are expected to copy following lines to ini config file.
// [database]
// host="localhost"
// port="5432"
// user="postgres"
// password="password"
// dbname="postgres"
// migrationDir="migrations/sql"
// sslMode="disable"
// connectionTimeout="30000"
// statementTimeout="30000"
// idleInTransactionSessionTimeout="30000"
package postgres
