package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	Beginx() (*sqlx.Tx, error)
	Rebind(string) string
	MustExec(string, ...interface{}) sql.Result
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(string, interface{}) (sql.Result, error)
	Get(interface{}, string, ...interface{}) error
	Select(interface{}, string, ...interface{}) error
}

type Database struct {
	config *Config
	*sqlx.DB
}

func New(config *Config) *Database {
	return &Database{config: config}
}

// BuildConnectionString converts connection options to a format
// that the database library understands.
func (db *Database) BuildConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.config.Host,
		db.config.Port,
		db.config.Username,
		db.config.Password,
		db.config.Database,
		db.config.SSLMode,
	)
}

// BuildConnectionURL converts connection options to a format
// that can be used to connect from CLI to postgres.
func (db *Database) BuildConnectionURL() string {
	return fmt.Sprintf(
		"pgx5://%s:%s@%s:%d/%s?sslmode=%s",
		db.config.Username,
		db.config.Password,
		db.config.Host,
		db.config.Port,
		db.config.Database,
		db.config.SSLMode,
	)
}

// Connect tries to connect to the database
// using options that describe the necessary
// information.
func (db *Database) Connect() error {
	engine, err := sqlx.Open("pgx", db.BuildConnectionString())
	if err != nil {
		return err
	}

	db.DB = engine
	return nil
}
