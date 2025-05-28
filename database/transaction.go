package database

import (
	"context"

	"github.com/dexlabsio/garlic/errors"
	"github.com/jmoiron/sqlx"
)

type key int

const (
	TransactionKey key = iota
)

// BeginContext starts a new database transaction within the provided context.
// It returns a new context containing the transaction, along with commit and rollback functions.
// If a transaction already exists in the context, it returns the existing context and no-op functions.
// If starting the transaction fails, it returns an error wrapped with KindDatabaseTransactionError.
func BeginContext(ctx context.Context, db *sqlx.DB) (ctxTx context.Context, commit, rollback func() error, err error) {
	if Transaction(ctx) != nil {
		return ctx, Nop(), Nop(), nil
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return ctx, Nop(), Nop(), errors.PropagateAs(
			KindDatabaseTransactionError,
			err,
			"failed to begin transaction",
		)
	}

	ctx = context.WithValue(ctx, TransactionKey, tx)
	commit = Commit(tx)
	rollback = Rollback(tx)
	return
}

// Rollback attempts to roll back the given transaction. If the rollback
// operation fails, it returns an error wrapped with KindDatabaseTransactionError.
// Otherwise, it returns nil, indicating the rollback was successful.
func Rollback(tx *sqlx.Tx) func() error {
	return func() error {
		if err := tx.Rollback(); err != nil {
			return errors.PropagateAs(
				KindDatabaseTransactionError,
				err,
				"failed to rollback transaction",
			)
		}

		return nil
	}

}

// Commit attempts to commit the given transaction. If the commit
// operation fails, it returns an error wrapped with KindDatabaseTransactionError.
// Otherwise, it returns nil, indicating the commit was successful.
func Commit(tx *sqlx.Tx) func() error {
	return func() error {
		if err := tx.Commit(); err != nil {
			return errors.PropagateAs(
				KindDatabaseTransactionError,
				err,
				"failed to commit transaction",
			)
		}

		return nil
	}
}

// Transaction retrieves the current database transaction from the provided context.
// If no transaction is found in the context, it returns nil. This function is useful
// for checking whether a transaction is already active within a given context.
func Transaction(ctx context.Context) *sqlx.Tx {
	tx, ok := ctx.Value(TransactionKey).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx
}

// Nop returns a no-operation function that always returns nil.
// This is useful as a placeholder for commit or rollback functions
// when no actual operation is needed, such as when a transaction
// is not active or has already been handled.
func Nop() func() error {
	return func() error {
		return nil
	}
}
