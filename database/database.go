package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dexlabsio/garlic/errors"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Executor interface {
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
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

func (db *Database) BeginContext(ctx context.Context) (ctxTx context.Context, commit, rollback func() error, err error) {
	ctxTx, commit, rollback, err = BeginContext(ctx, db.DB)
	if err != nil {
		err = errors.Propagate(err, "failed to begin database transaction")
	}

	return
}

func (db *Database) Executor(ctx context.Context) Executor {
	tx := Transaction(ctx)
	if tx != nil {
		return tx
	}

	return db
}

func (db *Database) Create(ctx context.Context, query string, resource any) error {
	ectx := errors.Context(
		errors.Field("query", query),
		errors.Field("resource_name", resource),
	)

	executor := db.Executor(ctx)

	rows, err := executor.NamedQuery(query, resource)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505": // UNIQUE violation
				// @TODO include the unique parameter violated in the error user scope.
				return errors.PropagateAs(
					errors.KindUserError,
					err,
					"resource already exists",
					errors.Hint(
						"Another resource with similar parameters already exists in our system. "+
							"Please change the parameters and try again.",
					),
					ectx,
				)
			case "23502": // NOT NULL constraint violation
				return errors.PropagateAs(
					errors.KindUserError,
					err,
					"missing required field",
					errors.Hint(
						"You're trying to create a new resource but some parameters are missing. "+
							"Please, check the documentation to clarify which parameters are necessary and try again.",
					),
					ectx,
				)
			default:
				return errors.PropagateAs(errors.KindSystemError, err, "missing required field", ectx)
			}
		}

		return errors.PropagateAs(errors.KindSystemError, err, "failed to insert resource", ectx)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(resource); err != nil {
			return errors.PropagateAs(errors.KindSystemError, err, "failed to scan returned resource", ectx)
		}
	} else {
		return errors.PropagateAs(errors.KindSystemError, err, "no rows returned while scanning resource during creation", ectx)
	}

	return nil
}

func (db *Database) List(ctx context.Context, query string, resourceList any, args ...any) error {
	ectx := errors.Context(
		errors.Field("query", query),
	)

	executor := db.Executor(ctx)
	if err := executor.Select(resourceList, query, args...); err != nil {
		return errors.PropagateAs(errors.KindSystemError, err, "failed to select resources", ectx)
	}

	return nil
}

func (db *Database) Delete(ctx context.Context, query string, args ...any) error {
	ectx := errors.Context(
		errors.Field("query", query),
		errors.Field("args", args),
	)

	executor := db.Executor(ctx)
	res, err := executor.Exec(query, args...)
	if err != nil {
		return errors.PropagateAs(errors.KindSystemError, err, "failed to execute delete query")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return errors.PropagateAs(errors.KindSystemError, err, "failed to get affected rows while deleting resource", ectx)
	}

	if rows < 1 {
		return errors.New(
			errors.KindNotFoundError,
			"resource not found",
			errors.Hint("Check if the reference of this resource is right and if exists."),
			ectx,
		)
	}

	return nil
}

func (db *Database) Update(ctx context.Context, query string, args ...any) error {
	ectx := errors.Context(
		errors.Field("query", query),
		errors.Field("args", args),
	)

	executor := db.Executor(ctx)
	res, err := executor.Exec(query, args...)
	if err != nil {
		return errors.PropagateAs(errors.KindSystemError, err, "failed to execute query while updating resource", ectx)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return errors.PropagateAs(errors.KindSystemError, err, "failed to get affected rows while updating resource", ectx)
	}

	if rows < 1 {
		return errors.New(
			errors.KindNotFoundError,
			"resource not found",
			errors.Hint("Check if the reference of this resource is right and if exists."),
			ectx,
		)
	}

	return nil
}

func (db *Database) Read(ctx context.Context, query string, resource any, args ...any) error {
	ectx := errors.Context(
		errors.Field("query", query),
		errors.Field("args", args),
	)

	executor := db.Executor(ctx)
	err := executor.Get(resource, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New(
			errors.KindNotFoundError,
			"resource not found",
			errors.Hint("Check if the reference of this resource is right and if exists."),
			ectx,
		)
	}

	if err != nil {
		return errors.PropagateAs(errors.KindSystemError, err, "failed to read dataset from database", ectx)
	}

	return nil
}
