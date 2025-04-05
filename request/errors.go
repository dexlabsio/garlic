package request

import (
	"github.com/dexlabsio/garlic/errors"
)

var (
	InvalidRequestError = errors.NewKind("Invalid Request Error", errors.KindUser)
)
