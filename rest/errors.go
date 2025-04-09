package rest

import (
	"github.com/dexlabsio/garlic/errors"
)

const (
	StatusCodeOptKey = "statuscode"
)

var (
	NotFoundError = errors.NewKind("Not Found Error", errors.KindUserError)
)

// StatusCode is a custom type that implements several methods
// to implement the Opt iface of the errors package. It provides
// a way to represent and handle status codes with specific behaviors,
// such as returning a key, value, and visibility level, and
// inserting itself into error options.
type StatusCode int

func (sc StatusCode) Key() string {
	return StatusCodeOptKey
}

func (sc StatusCode) Value() any {
	return sc
}

func (sc StatusCode) Visibility() errors.Visibility {
	return errors.RESTRICT
}

func (sc StatusCode) Insert(opt errors.Opt) errors.Opt {
	return sc
}
