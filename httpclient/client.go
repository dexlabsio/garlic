package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/dexlabsio/garlic/errors"
	"github.com/dexlabsio/garlic/logging"
)

// Post sends a HTTP POST request to the given url
func Post(ctx context.Context, url string, data any) (*http.Response, error) {
	return request(ctx, http.MethodPost, url, data)
}

// Put sends a HTTP PUT request to the given url
func Put(ctx context.Context, url string, data any) (*http.Response, error) {
	return request(ctx, http.MethodPut, url, data)
}

// Patch sends a HTTP PATCH request to the given url
func Patch(ctx context.Context, url string, data any) (*http.Response, error) {
	return request(ctx, http.MethodPatch, url, data)
}

// Get sends a HTTP GET request to the given url
func Get(ctx context.Context, url string) (*http.Response, error) {
	return request(ctx, http.MethodGet, url, nil)
}

// Delete sends a HTTP DELETE request to the given url
func Delete(ctx context.Context, url string) (*http.Response, error) {
	return request(ctx, http.MethodDelete, url, nil)
}

func request(ctx context.Context, method, url string, data any) (*http.Response, error) {
	ectx := errors.Context(
		errors.Field("http_method", method, errors.Restrict),
		errors.Field("http_url", url, errors.Restrict),
	)

	l := logging.GetLoggerFromContext(ctx).With(ectx.Zap())

	// encode data into json bytes
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, errors.PropagateAs(errors.KindSystemError, err, "failed to unmarshal data into JSON", ectx)
	}

	// create net_http request with given method and request body
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.PropagateAs(errors.KindSystemError, err, "failed to create HTTP request", ectx)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", logging.GetSessionIdFromContext(ctx))
	req.Header.Set("X-Request-ID", logging.GetRequestIdFromContext(ctx))

	// create net_http client and send request
	client := &http.Client{}
	var res *http.Response

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 100 * time.Millisecond
	expBackoff.MaxInterval = 2 * time.Second
	expBackoff.MaxElapsedTime = 1 * time.Minute // Overall retry duration limit

	operation := func() error {
		res, err = client.Do(req)
		if err != nil {
			return err
		}

		return nil
	}

	// notify is called after each failed attempt.
	notify := func(err error, delay time.Duration) {
		l.Error(fmt.Sprintf("Failed to send request: %v. Retrying in %v...\n", err, delay))
	}

	l.Info("Sending request to remote API")
	if err := backoff.RetryNotify(operation, expBackoff, notify); err != nil {
		return nil, errors.Propagate(err, "failed to make request", ectx)
	}

	l.Debug("Successfully performed external request")
	return res, nil
}
