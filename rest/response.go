package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dexlabsio/garlic/errors"
)

type PayloadMessage struct {
	Message string `json:"message"`
}

type Response struct {
	StatusCode int
	Payload    any
}

var (
	// We filter internal server errors to provide a standard
	// response and prevent leaking sensitive information
	internalServerErrorResponse = WriteResponse(
		http.StatusInternalServerError,
		errors.Raw(
			errors.KindSystemError,
			"internal server error",
			errors.Hint("internal server error, please contact the support"),
		).ErrorDTO(),
	)

	// This is a generic response for unknown errors
	unknownErrorResponse = WriteResponse(
		http.StatusInternalServerError,
		errors.Raw(
			errors.KindSystemError,
			"unknown error",
			errors.Hint("unknown error, please contact the support"),
		).ErrorDTO(),
	)
)

func (r *Response) Must(w http.ResponseWriter) {
	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	if err := json.NewEncoder(w).Encode(r.Payload); err != nil {
		panic(fmt.Sprintf("Failed to encode response %s", err))
	}
}

// WriteResponse is a generic function to create a response with a payload
func WriteResponse(statusCode int, payload any) *Response {
	return &Response{
		StatusCode: statusCode,
		Payload:    payload,
	}
}

// WriteMessage is a helper function to create a response with a message
func WriteMessage(statusCode int, message string) *Response {
	return WriteResponse(statusCode, PayloadMessage{Message: message})
}

// WriteError is a helper function to create a response with a service error
// or a generic error response if the error is not a service error
func WriteError(err error) *Response {
	// Return unknown error if the callen didn't provide an error
	if err == nil {
		return unknownErrorResponse
	}

	// Return internal server error if the error is not a service error
	usrErr, ok := errors.AsKind(err, errors.KindUserError)
	if !ok {
		return internalServerErrorResponse
	}

	statusCode := usrErr.Kind().StatusCode()
	return WriteResponse(statusCode, usrErr.ErrorDTO())
}
