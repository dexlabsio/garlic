package validator

import (
	"github.com/dexlabsio/garlic/errors"
)

var (
	ValidationError = errors.NewKind("Validation Error", errors.KindUserError)
)
