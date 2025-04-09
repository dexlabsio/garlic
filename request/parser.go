package request

import (
	"encoding/json"
	"net/http"

	"github.com/dexlabsio/garlic/crypto"
	"github.com/dexlabsio/garlic/errors"
	"github.com/dexlabsio/garlic/rest"
	"github.com/dexlabsio/garlic/validator"
)

type Form[T any] interface {
	ToModel() (T, error)
}

type UnsafeForm[T any] interface {
	ToModel(crpt crypto.Manager) (T, error)
}

// DecodeRequestBody applies the JSON decoder into the request body and
// validate the struct formatting requirements using the validator package.
func DecodeRequestBody[T any](r *http.Request, form T) error {
	l := GetLogger(r)

	if r.ContentLength == 0 {
		l.Warn("Empty request body")
	} else if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		return errors.PropagateAs(
			InvalidRequestError,
			err,
			"invalid request body",
			errors.Hint(
				"something may be wrong with formatting or the content of the request body",
			),
			rest.StatusCode(http.StatusBadRequest),
		)
	}

	if err := ValidateForm(form); err != nil {
		return errors.Propagate(
			err,
			"failed to validate form",
			rest.StatusCode(http.StatusBadRequest),
		)
	}

	return nil
}

func ValidateForm[T any](form T) error {
	if err := validator.Global().Struct(form); err != nil {
		return errors.Propagate(
			validator.ParseValidationErrors(err),
			"failed to validate form",
			errors.Hint(
				"please, verify the correctness of the fields",
			),
		)
	}

	return nil
}

// ParseForm handles decoding and validation of request bodies
// into generic forms.
func ParseForm[T any, F Form[T]](r *http.Request, form F) (T, error) {
	var model T

	if err := DecodeRequestBody(r, form); err != nil {
		return model, errors.Propagate(err, "failed to decode request body")
	}

	model, err := form.ToModel()
	if err != nil {
		return model, errors.Propagate(err, "failed to parse form")
	}

	return model, nil
}

// ParseUnsafeForm handles decoding and validation of request bodies
// into generic forms with decrypted values that should be encrypted.
func ParseUnsafeForm[T any, F UnsafeForm[T]](r *http.Request, form F, crpt crypto.Manager) (T, error) {
	var model T

	if err := DecodeRequestBody(r, form); err != nil {
		return model, errors.Propagate(err, "failed to decode unsafe request body into a form")
	}

	model, err := form.ToModel(crpt)
	if err != nil {
		return model, errors.Propagate(err, "failed to parse unsafe form into model")
	}

	return model, nil
}
