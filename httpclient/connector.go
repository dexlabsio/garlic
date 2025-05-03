package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/dexlabsio/garlic/errors"
)

type Request struct {
	Method      string
	URI         string
	Data        any
	QueryParams map[string]string
}

type Connector struct {
	config *Config
}

func NewConnector(config *Config) *Connector {
	return &Connector{config}
}

func (c *Connector) Request(ctx context.Context, req *Request, result any) error {
	ectx := errors.Context(
		errors.Field("http_method", req.Method),
		errors.Field("http_url", c.config.URL),
		errors.Field("http_uri", req.URI),
		errors.Field("http_query_params", req.QueryParams),
	)

	target, err := buildURL(c.config.URL, req.URI, req.QueryParams)
	if err != nil {
		return errors.PropagateAs(errors.KindSystemError, err, "failed to build request URL", ectx)
	}

	res, err := request(ctx, req.Method, target, req.Data)
	if err != nil {
		return errors.Propagate(err, "failed to make request", ectx)
	}

	// We only support StatusOK for successful operations
	if res.StatusCode != http.StatusOK {
		err := handleFailure(res)
		return errors.Propagate(err, "bad response from external service", ectx)
	}

	if err := handleSuccess(res, result); err != nil {
		return errors.Propagate(err, "failed to process successful response from external service", ectx)
	}

	return nil
}

// buildURL parses the base, joins the URI path, sets params, and returns the final URL string.
func buildURL(baseURL, uri string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	u.Path, err = url.JoinPath(baseURL, uri)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func handleSuccess(res *http.Response, result any) error {
	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		return errors.PropagateAs(errors.KindSystemError, err, "failed to decode response body")
	}

	return nil
}

func handleFailure(res *http.Response) *errors.ErrorT {
	var body errors.DTO

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return errors.PropagateAs(errors.KindSystemError, err, "failed to decode response body")
	}

	return body.Decode()
}
