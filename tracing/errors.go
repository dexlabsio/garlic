package tracing

import "github.com/dexlabsio/garlic/errors"

var (
	KindContextError              = errors.Get("ContextError")
	KindContextValueNotFoundError = errors.Get("ContextValueNotFoundError")
)
