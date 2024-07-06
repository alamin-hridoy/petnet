package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const driver = "postgres"

// NewDBStringFromConfig build database connection string from config file.
func NewDBStringFromConfig(config *viper.Viper) (string, error) {
	var dbParams []string
	dbParams = append(dbParams, fmt.Sprintf("user=%s", config.GetString("database.user")))
	dbParams = append(dbParams, fmt.Sprintf("host=%s", config.GetString("database.host")))
	dbParams = append(dbParams, fmt.Sprintf("port=%s", config.GetString("database.port")))
	dbParams = append(dbParams, fmt.Sprintf("dbname=%s", config.GetString("database.dbname")))
	if password := config.GetString("database.password"); password != "" {
		dbParams = append(dbParams, fmt.Sprintf("password=%s", password))
	}
	dbParams = append(dbParams, fmt.Sprintf("sslmode=%s",
		config.GetString("database.sslMode")))
	dbParams = append(dbParams, fmt.Sprintf("connect_timeout=%d",
		config.GetInt("database.connectionTimeout")))
	dbParams = append(dbParams, fmt.Sprintf("statement_timeout=%d",
		config.GetInt("database.statementTimeout")))
	dbParams = append(dbParams, fmt.Sprintf("idle_in_transaction_session_timeout=%d",
		config.GetInt("database.idleInTransactionSessionTimeout")))

	return strings.Join(dbParams, " "), nil
}

// Open opens a connection to database with given connection string.
func Open(config *viper.Viper) (*sql.DB, error) {
	dbString, err := NewDBStringFromConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driver, dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Open opens a connection to database with given connection string, using sqlx opener.
func Openx(config *viper.Viper) (*sqlx.DB, error) {
	dbString, err := NewDBStringFromConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open(driver, dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Connectx opens a connection to database with given connection string using sqlx opener
// and verify the connection with a ping.
func Connectx(config *viper.Viper) (*sqlx.DB, error) {
	dbString, err := NewDBStringFromConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect(driver, dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
