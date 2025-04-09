package request

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dexlabsio/garlic/errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ParseResourceUUID reads the resource id from the request path and tries to parse it into a valid UUID.
func ParseResourceUUID(r *http.Request, param string) (uuid.UUID, error) {
	rawResourceId := chi.URLParam(r, param)

	resourceId, err := uuid.Parse(rawResourceId)
	if err != nil {
		return uuid.Nil, errors.PropagateAs(
			InvalidRequestError,
			err,
			"failed to parse resource id (int)",
			errors.Hint(
				fmt.Sprintf("Something is wrong with the request field '%s'", param),
			),
		)
	}

	return resourceId, nil
}

// ParseResourceInt reads the resource id from the request path and attempts to parse it into an integer.
// If the parsing fails, it logs a warning and returns an error. Otherwise, it returns the parsed integer.
func ParseResourceInt(r *http.Request, param string) (int, error) {
	rawResourceId := chi.URLParam(r, param)

	resourceId, err := strconv.Atoi(rawResourceId)
	if err != nil {
		return 0, errors.PropagateAs(
			InvalidRequestError,
			err,
			"failed to parse resource id (int)",
			errors.Hint(
				fmt.Sprintf("Something is wrong with the request field '%s'", param),
			),
		)
	}

	return resourceId, nil
}

// ParseResourceString extracts a string parameter from the request path, ensuring it is not empty,
// and attempts to unescape it. If the parameter is empty or cannot be unescaped, it logs a warning
// and returns an error. Otherwise, it returns the unescaped string.
func ParseResourceString(r *http.Request, param string) (string, error) {
	str := chi.URLParam(r, param)
	if str == "" {
		return "", errors.New(
			InvalidRequestError,
			"path string is empty",
			errors.Hint(
				fmt.Sprintf("String path '%s' can't be empty", param),
			),
		)
	}

	unescapedPath, err := url.PathUnescape(str)
	if err != nil {
		return "", errors.PropagateAs(
			InvalidRequestError,
			err,
			"failed to unescape path string",
			errors.Hint(
				fmt.Sprintf("We couldn't unescape the path string '%s'", param),
			),
		)
	}

	return unescapedPath, nil
}

// ParseParamPagination extracts pagination parameters 'limit' and 'start' from the request query string.
// It attempts to convert these parameters to integers. If the conversion fails or the parameters are not
// set, it defaults both 'limit' and 'start' to 0. This function logs debug messages if the parameters
// are not set or cannot be converted.
func ParseParamPagination(r *http.Request) (limit, start int) {
	l := GetLogger(r)
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
func ParseParamUUID(r *http.Request, param string) (uuid.UUID, error) {
	rawParam := r.URL.Query().Get(param)
	if rawParam == "" {
		err := errors.New(
			InvalidRequestError,
			"required request param is missing",
			errors.Hint(
				fmt.Sprintf("Something is wrong with the request param '%s'", param),
			),
		)

		return uuid.Nil, err
	}

	paramUUID, err := uuid.Parse(rawParam)
	if err != nil {
		return uuid.Nil, errors.PropagateAs(
			InvalidRequestError,
			err,
			"malformed required request param",
			errors.Hint(
				fmt.Sprintf("Something is wrong with the request param '%s'", param),
			),
		)
	}

	return paramUUID, nil
}

// ParseOptionalParamUUID attempts to retrieve an optional UUID parameter from the request query string.
// If the parameter is not present, it returns uuid.Nil without an error. If the parameter is present
// but cannot be parsed into a valid UUID, it logs a warning and returns an error. Otherwise, it returns
// the parsed UUID.
func ParseOptionalParamUUID(r *http.Request, param string) (uuid.UUID, error) {
	rawParam := r.URL.Query().Get(param)
	if rawParam == "" {
		return uuid.Nil, nil
	}

	paramUUID, err := uuid.Parse(rawParam)
	if err != nil {
		return uuid.Nil, errors.PropagateAs(
			InvalidRequestError,
			err,
			"malformed optional request param",
			errors.Hint(
				fmt.Sprintf("Something is wrong with the optional request param '%s'", param),
			),
		)
	}

	return paramUUID, nil
}

// ParseParamString retrieves the specified query parameter from the request URL.
// If the parameter is missing, it logs a warning and returns an error indicating
// that the required request parameter is missing. Otherwise, it returns the
// parameter value as a string.
func ParseParamString(r *http.Request, param string) (string, error) {
	rawParam := r.URL.Query().Get(param)
	if rawParam == "" {
		return "", errors.New(
			InvalidRequestError,
			"missing required request param",
			errors.Hint(
				fmt.Sprintf("Request param '%s' is missing", param),
			),
		)
	}

	return rawParam, nil
}

// ParseOptionalParamBool attempts to retrieve an optional boolean parameter from the request query string.
// If the parameter is not present, it returns false without an error. If the parameter is present but cannot
// be parsed into a valid boolean, it logs a debug message and returns an error. Otherwise, it returns the
// parsed boolean value.
func ParseOptionalParamBool(r *http.Request, param string) (bool, error) {
	l := GetLogger(r)

	rawParam := r.URL.Query().Get(param)

	if rawParam == "" {
		l.Debug("Optional request param not set", zap.String("param", param))
		return false, nil
	}

	val, err := strconv.ParseBool(rawParam)
	if err != nil {
		return false, errors.PropagateAs(
			InvalidRequestError,
			err,
			"Optional request param bool was provided but cannot be parsed",
			errors.Hint(fmt.Sprintf("make sure '%s' request param is either 'true' or 'false', no other value is accepted", param)),
		)
	}

	return val, err
}
