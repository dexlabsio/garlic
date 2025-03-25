package request

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/dexlabsio/garlic/pkg/errors"
	"go.uber.org/zap"
)

// ParseResourceUUID reads the resource id from the request path and tries to parse it into a valid UUID.
func ParseResourceUUID(l *zap.Logger, r *http.Request, param string) (uuid.UUID, error) {
	rawResourceId := chi.URLParam(r, param)

	resourceId, err := uuid.Parse(rawResourceId)
	if err != nil {
		err = NewInvalidRequestError("failed to parse resource id (int)", errors.Hint(
			fmt.Sprintf("Something is wrong with the request field '%s'", param),
		)).Wrap(err)
		l.Warn("Failed to parse resource id", errors.Zap(err), zap.String("param", param))
		return uuid.Nil, err
	}

	return resourceId, nil
}

// ParseResourceString reads the resource id from the request path and returns it as it's since it's
// already in the string format
func ParseResourceString(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

// ParseResourceInt reads the resource id from the request path and tries to parse it into a valid integer.
func ParseResourceInt(l *zap.Logger, r *http.Request, param string) (int, error) {
	rawResourceId := chi.URLParam(r, param)

	resourceId, err := strconv.Atoi(rawResourceId)
	if err != nil {
		err = NewInvalidRequestError("failed to parse resource id (int)", errors.Hint(
			fmt.Sprintf("Something is wrong with the request field '%s'", param),
		)).Wrap(err)
		l.Warn("Failed to parse resource id", errors.Zap(err), zap.String("param", param))
		return 0, err
	}

	return resourceId, nil
}

// ParseParamPagination is used to parse pagination parameters passed limit and start
// defaults to 0 when they're not found
func ParseParamPagination(l *zap.Logger, r *http.Request) (limit, start int) {
	var err error

	rawPaginationLimit := r.URL.Query().Get("limit")
	rawPaginationStart := r.URL.Query().Get("start")

	limit, err = strconv.Atoi(rawPaginationLimit)
	if err != nil {
		l.Debug("Pagination limit not set, defaults to 0")
		limit = 0 // explicit zero value
	}

	start, err = strconv.Atoi(rawPaginationStart)
	if err != nil {
		l.Debug("Pagination start not set, defaults to 0")
		start = 0 // explicit zero value
	}

	return
}

// ParseParamUUID takes the request query string and tries to find the given param. Then it tries to parse
// the respective value into an UUID. If it breaks the function returns uuid.Nil and a false checker. It also
// returns a common error message to the user.
func ParseParamUUID(l *zap.Logger, r *http.Request, param string) (uuid.UUID, error) {
	rawParam := r.URL.Query().Get(param)
	if rawParam == "" {
		err := NewInvalidRequestError("required request param is missing", errors.Hint(
			fmt.Sprintf("Something is wrong with the request param '%s'", param),
		))

		l.Warn("Missing required request param", errors.Zap(err), zap.String("param", param))
		return uuid.Nil, err
	}

	paramUUID, err := uuid.Parse(rawParam)
	if err != nil {
		err = NewInvalidRequestError("malformed required request param", errors.Hint(
			fmt.Sprintf("Something is wrong with the request param '%s'", param),
		)).Wrap(err)

		l.Warn("Malformed mandatory request param", errors.Zap(err), zap.String("param", param))
		return uuid.Nil, err
	}

	return paramUUID, nil
}

func ParseOptionalParamUUID(l *zap.Logger, r *http.Request, param string) (uuid.UUID, error) {
	rawParam := r.URL.Query().Get(param)
	if rawParam == "" {
		return uuid.Nil, nil
	}

	paramUUID, err := uuid.Parse(rawParam)
	if err != nil {
		err = NewInvalidRequestError("malformed optional request param", errors.Hint(
			fmt.Sprintf("Something is wrong with the optional request param '%s'", param),
		)).Wrap(err)

		l.Warn("Malformed optional request param", errors.Zap(err), zap.String("param", param))
		return uuid.Nil, err
	}

	return paramUUID, nil
}

// ParseStringPath tries to find the given path param in the given request and tries to unescape it
func ParseStringPath(l *zap.Logger, w http.ResponseWriter, r *http.Request, param string) (string, error) {
	str := chi.URLParam(r, param)
	if str == "" {
		err := NewInvalidRequestError("path string is empty", errors.Hint(
			fmt.Sprintf("String path '%s' can't be empty", param),
		))

		l.Warn("Path string is empty", errors.Zap(err), zap.String("param", param))
		return "", err
	}

	unescapedPath, err := url.PathUnescape(str)
	if err != nil {
		err := NewInvalidRequestError("failed to unescape path string", errors.Hint(
			fmt.Sprintf("We couldn't unescape the path string '%s'", param),
		)).Wrap(err)

		l.Warn("Failed to unescape path string", errors.Zap(err), zap.String("param", param))
		return "", err
	}

	return unescapedPath, nil
}

// ParseParamString takes the request query string and tries to find the given param.
func ParseParamString(l *zap.Logger, w http.ResponseWriter, r *http.Request, param string) (string, error) {
	rawParam := r.URL.Query().Get(param)
	if rawParam == "" {
		err := NewInvalidRequestError("missing required request param", errors.Hint(
			fmt.Sprintf("Request param '%s' is missing", param),
		))

		l.Warn("Missing required request param", errors.Zap(err), zap.String("param", param))

		return "", err
	}

	return rawParam, nil
}
