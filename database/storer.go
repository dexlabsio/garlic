package database

import (
	"context"

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
		err = errors.Propagate(err, "storer failed to begin database transaction")
	}

	err = fn(ctxTx)
	if p := recover(); p != nil {
		if rerr := rollback(); rerr != nil {
			logging.Global().Error("Failed to rollback transaction during panic handling", errors.Zap(rerr))
		}
		panic(p)
	}

	if err != nil {
		if rerr := rollback(); rerr != nil {
			logging.Global().Error("Failed to rollback transaction during error handling", errors.Zap(rerr))
		}

		return errors.Propagate(err, "storer transaction failed and rollback was applied")
	}

	err = commit()
	if err != nil {
		return errors.Propagate(err, "storer failed to commit database transaction")
	}

	return nil
}
