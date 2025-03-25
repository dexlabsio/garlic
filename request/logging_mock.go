//go:build unit
// +build unit

package request

import (
	"context"
	"net/http"

	"github.com/dexlabsio/garlic/logging"
	"go.uber.org/zap"
)

func WithLogger(r *http.Request, logger *zap.Logger) *http.Request {
	ctx := context.WithValue(r.Context(), logging.LoggerKey, logger)
	return r.WithContext(ctx)
}
