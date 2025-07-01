package database

import (
	"context"
	"database/sql"

	"github.com/dexlabsio/garlic/errors"
	"github.com/dexlabsio/garlic/logging"
)

type Store interface {
	BeginContext(ctx context.Context) (ctxTx context.Context, commit, rollback func() error, err error)
	Create(ctx context.Context, query string, resource any) error
	Read(ctx context.Context, query string, resource any, args ...any) error
	Update(ctx context.Context, query string, args ...any) error
	Delete(ctx context.Context, query string, args ...any) error
	List(ctx context.Context, query string, resourceList any, args ...any) error
	RawExec(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedRawExec(ctx context.Context, query string, resource any) (sql.Result, error)
}

type Storer struct {
	Store Store
}

func NewStorer(store Store) *Storer {
	return &Storer{Store: store}
}

func (s *Storer) Transaction(ctx context.Context, fn func(context.Context) error) error {
	var err error

	ctxTx, commit, rollback, err := s.Store.BeginContext(ctx)
	if err != nil {
		return errors.Propagate(err, "storer failed to begin database transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			if err := rollback(); err != nil {
				logging.Global().Error("Failed to rollback transaction during panic handling", errors.Zap(err))
			}
			panic(p)
		}
	}()

	defer func() {
		if err != nil {
			if rerr := rollback(); rerr != nil {
				logging.Global().Error("Failed to rollback transaction during error handling", errors.Zap(rerr))
			}
		}
	}()

	err = fn(ctxTx)
	if err != nil {
		return errors.Propagate(err, "failed to execute transactional function")
	}

	err = commit()
	if err != nil {
		return errors.Propagate(err, "storer failed to commit database transaction")
	}

	return nil
}
