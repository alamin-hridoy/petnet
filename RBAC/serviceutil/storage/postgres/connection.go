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

type DBConfig struct {
	User                            string `mapstructure:"user"`
	Host                            string `mapstructure:"host"`
	Port                            string `mapstructure:"port"`
	DBName                          string `mapstructure:"dbname"`
	Password                        string `mapstructure:"password"`
	SSLMode                         string `mapstructure:"sslMode"`
	ConnectionTimeout               int    `mapstructure:"connectionTimeout"`
	StatementTimeout                int    `mapstructure:"statementTimeout"`
	IdleInTransactionSessionTimeout int    `mapstructure:"idleInTransactionSessionTimeout"`
}

func NewDBStringFromDBConfig(config DBConfig) (string, error) {
	var dbParams []string
	dbParams = append(dbParams, fmt.Sprintf("user=%s", config.User))
	dbParams = append(dbParams, fmt.Sprintf("host=%s", config.Host))
	dbParams = append(dbParams, fmt.Sprintf("port=%s", config.Port))
	dbParams = append(dbParams, fmt.Sprintf("dbname=%s", config.DBName))
	if password := config.Password; password != "" {
		dbParams = append(dbParams, fmt.Sprintf("password=%s", password))
	}
	dbParams = append(dbParams, fmt.Sprintf("sslmode=%s",
		config.SSLMode))
	dbParams = append(dbParams, fmt.Sprintf("connect_timeout=%d",
		config.ConnectionTimeout))
	dbParams = append(dbParams, fmt.Sprintf("statement_timeout=%d",
		config.StatementTimeout))
	dbParams = append(dbParams, fmt.Sprintf("idle_in_transaction_session_timeout=%d",
		config.IdleInTransactionSessionTimeout))

	return strings.Join(dbParams, " "), nil
}

// NewDBStringFromConfig build database connection string from config file.
func NewDBStringFromConfig(config *viper.Viper) (string, error) {
	var allConfig struct {
		Database DBConfig `mapstructure:"database"`
	}
	if err := config.Unmarshal(&allConfig); err != nil {
		return "", fmt.Errorf("cannot unmarshal db config: %w", err)
	}

	return NewDBStringFromDBConfig(allConfig.Database)
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

func OpenDBConfig(config DBConfig) (*sql.DB, error) {
	dbString, err := NewDBStringFromDBConfig(config)
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

func OpenxDBConfig(config DBConfig) (*sqlx.DB, error) {
	dbString, err := NewDBStringFromDBConfig(config)
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

func ConnectxDBConfig(config DBConfig) (*sqlx.DB, error) {
	dbString, err := NewDBStringFromDBConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect(driver, dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
