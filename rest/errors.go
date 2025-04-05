package rest

import "github.com/dexlabsio/garlic/errors"

var (
	NotFoundError = errors.NewKind("Not Found Error", errors.KindUserError)
)
