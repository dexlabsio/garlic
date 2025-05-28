package database

import "github.com/dexlabsio/garlic/errors"

var (
	KindDatabaseRecordNotFoundError = errors.Get("DatabaseRecordNotFoundError")
	KindDatabaseTransactionError    = errors.Get("DatabaseTransactionError")
)
