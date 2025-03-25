package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dexlabsio/garlic/pkg/errors"
)

func httpStatusFromError(err *errors.ErrorT) int {
	if err == nil {
		panic("can't create an status for nil error")
	}

	var notFoundErr *errors.NotFoundError
	if errors.As(err, &notFoundErr) {
		return http.StatusNotFound
	}

	var userErr *errors.UserError
	if errors.As(err, &userErr) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

type PayloadError struct {
	Error map[string]any `json:"error"`
}

type PayloadMessage struct {
	Message string `json:"message"`
}

type Response[T any] struct {
	StatusCode int
	Payload    T
}

var (
	// We filter internal server errors to provide a standard
	// response and prevent leaking sensitive information
	internalServerErrorResponse = WriteResponse(
		http.StatusInternalServerError,
		PayloadError{
			Error: errors.New(
				"internal server error",
				errors.Hint("internal server error, please contact the support"),
			).DTO(),
		},
	)

	// This is a generic response for unknown errors
	unknownErrorResponse = WriteResponse(
		http.StatusInternalServerError,
		PayloadError{
			Error: errors.New(
				"unknown error",
				errors.Hint("unknown error, please contact the support"),
			).DTO(),
		},
	)
)

func (r *Response[_]) Must(w http.ResponseWriter) {
	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	if err := json.NewEncoder(w).Encode(r.Payload); err != nil {
		panic(fmt.Sprintf("Failed to encode response %s", err))
	}
}

// WriteResponse is a generic function to create a response with a payload
func WriteResponse[T any](statusCode int, payload T) *Response[T] {
	return &Response[T]{
		StatusCode: statusCode,
		Payload:    payload,
	}
}

// WriteMessage is a helper function to create a response with a message
func WriteMessage(statusCode int, message string) *Response[PayloadMessage] {
	return WriteResponse(statusCode, PayloadMessage{Message: message})
}

// WriteError is a helper function to create a response with a service error
// or a generic error response if the error is not a service error
func WriteError(err error) *Response[PayloadError] {
	// Return unknown error if the callen didn't provide an error
	if err == nil {
		return unknownErrorResponse
	}

	// Return internal server error if the error is not a service error
	var errt *errors.ErrorT
	if !errors.As(err, &errt) {
		return internalServerErrorResponse
	}

	// Standardize the error output if the status of the error is 500, 401 or 403,
	// this way we avoid leaking potential dangerous internal information
	responseStatus := httpStatusFromError(errt)
	switch responseStatus {
	case http.StatusInternalServerError:
		return internalServerErrorResponse
	}

	return WriteResponse(responseStatus, PayloadError{Error: errt.DTO()})
}
