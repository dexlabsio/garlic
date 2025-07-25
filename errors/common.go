package errors

import "net/http"

var (
	KindError = &Kind{
		Name:           "Error",
		Code:           "U00000",
		Description:    "Any error that has not been mapped in the application.",
		HTTPStatusCode: HTTP_STATUS_NOT_DEFINED,
		Parent:         nil,
	}

	KindUserError = &Kind{
		Name:           "UserError",
		Code:           "E00001",
		Description:    "Any error that was caused by some incorrect user action.",
		HTTPStatusCode: http.StatusBadRequest,
		Parent:         KindError,
	}

	KindSystemError = &Kind{
		Name:           "SystemError",
		Code:           "S00001",
		Description:    "Any error that was caused by some unexpected system failure.",
		HTTPStatusCode: http.StatusInternalServerError,
		Parent:         KindError,
	}

	KindInvalidRequestError = &Kind{
		Name:           "InvalidRequestError",
		Code:           "E00002",
		Description:    "The request is incorrectly formatted or has errors in the request body.",
		HTTPStatusCode: http.StatusBadRequest,
		Parent:         KindUserError,
	}

	KindNotFoundError = &Kind{
		Name:           "NotFoundError",
		Code:           "E00003",
		Description:    "The requested resource was not found in our system or external services.",
		HTTPStatusCode: http.StatusNotFound,
		Parent:         KindUserError,
	}

	KindValidationError = &Kind{
		Name:           "ValidationError",
		Code:           "E00004",
		Description:    "Some field on a form was filled incorrectly by the user or is missing.",
		HTTPStatusCode: http.StatusBadRequest,
		Parent:         KindInvalidRequestError,
	}

	KindAuthError = &Kind{
		Name:           "AuthError",
		Code:           "E00005",
		Description:    "An error occurred during authentication, such as invalid credentials.",
		HTTPStatusCode: http.StatusUnauthorized,
		Parent:         KindUserError,
	}

	KindForbiddenError = &Kind{
		Name:           "ForbiddenError",
		Code:           "E00006",
		Description:    "The user does not have permission to access the requested resource.",
		HTTPStatusCode: http.StatusForbidden,
		Parent:         KindUserError,
	}

	KindContextError = &Kind{
		Name:        "ContextError",
		Code:        "S00002",
		Description: "An error occurred due to a problem with the context.",
		Parent:      KindSystemError,
	}

	KindContextValueNotFoundError = &Kind{
		Name:        "ContextValueNotFoundError",
		Code:        "S00003",
		Description: "A required value was not found in the context.",
		Parent:      KindContextError,
	}

	KindDatabaseRecordNotFoundError = &Kind{
		Name:           "DatabaseRecordNotFoundError",
		Code:           "E00007",
		Description:    "The requested database record was not found.",
		HTTPStatusCode: http.StatusNotFound,
		Parent:         KindNotFoundError,
	}

	KindDatabaseTransactionError = &Kind{
		Name:           "DatabaseTransactionError",
		Code:           "S00005",
		Description:    "An error occurred during a database transaction.",
		HTTPStatusCode: http.StatusInternalServerError,
		Parent:         KindSystemError,
	}
)

func init() {
	Register(
		KindError,
		KindUserError,
		KindSystemError,
		KindInvalidRequestError,
		KindValidationError,
		KindNotFoundError,
		KindAuthError,
		KindForbiddenError,
		KindContextError,
		KindContextValueNotFoundError,
		KindDatabaseRecordNotFoundError,
		KindDatabaseTransactionError,
	)
}
