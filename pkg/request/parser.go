package request

import (
	"encoding/json"
	"net/http"

	"github.com/dexlabsio/garlic/pkg/crypto"
	"github.com/dexlabsio/garlic/pkg/errors"
	"github.com/dexlabsio/garlic/pkg/validator"
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
		return NewInvalidRequestError("invalid request body", errors.Hint(
			"Something may be wrong with formatting or the content of the request body",
		)).Wrap(err)
	}

	if err := ValidateForm(form); err != nil {
		return errors.Propagate("failed to validate form", err, errors.Hint(
			"Please, verify the correctness of the fields",
		))
	}

	return nil
}

func ValidateForm[T any](form T) error {
	if err := validator.SimpleValidator.Struct(form); err != nil {
		return validator.ParseValidationErrors(err)
	}

	return nil
}

// ParseForm handles decoding and validation of request bodies
// into generic forms.
func ParseForm[T any, F Form[T]](r *http.Request, form F) (T, error) {
	var model T
	l := GetLogger(r)

	if err := DecodeRequestBody(r, form); err != nil {
		l.Error("Failed to decode request body into a form", errors.Zap(err))
		return model, err
	}

	model, err := form.ToModel()
	if err != nil {
		l.Error("Failed to convert form into a model", errors.Zap(err))
		return model, err
	}

	return model, nil
}

// ParseUnsafeForm handles decoding and validation of request bodies
// into generic forms with decrypted values that should be encrypted.
func ParseUnsafeForm[T any, F UnsafeForm[T]](r *http.Request, form F, crpt crypto.Manager) (T, error) {
	var model T
	l := GetLogger(r)

	if err := DecodeRequestBody(r, form); err != nil {
		err = NewInvalidRequestError("failed to decode unsafe request body into a form").Wrap(err)
		l.Error("Failed to parse request body", errors.Zap(err))
		return model, err
	}

	model, err := form.ToModel(crpt)
	if err != nil {
		err = NewInvalidRequestError("failed to convert unsafe form into model").Wrap(err)
		l.Error("Failed to convert form", errors.Zap(err))
		return model, err
	}

	return model, nil
}
