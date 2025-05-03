package validator

import (
	"reflect"
	"strings"

	"github.com/dexlabsio/garlic/errors"
	"github.com/dexlabsio/garlic/logging"
	val "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var singleton *Validator

type Field = val.FieldLevel

type Validator struct {
	*val.Validate
}

type FieldValidator interface {
	Key() string
	Validate(Field) bool
}

type Validation struct {
	key string
	fn  func(Field) bool
}

func NewValidation(key string, fn func(Field) bool) *Validation {
	return &Validation{
		key: key,
		fn:  fn,
	}
}

func (v *Validation) Key() string {
	return v.key
}

func (v *Validation) Validate(field Field) bool {
	return v.fn(field)
}

func New() *Validator {
	v := &Validator{val.New()}

	// Using the names which have been specified for JSON representations of structs,
	// rather than normal Go field names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	v.Extend(defaultExtendedValidations...)
	return v
}

func (v *Validator) Extend(validations ...FieldValidator) {
	for _, validation := range validations {
		if err := v.RegisterValidation(validation.Key(), validation.Validate); err != nil {
			logging.Global().Fatal(
				"Failed to register field validator",
				zap.String("validator_key", validation.Key()),
				zap.Error(err),
			)
		}
	}
}

func ParseValidationErrors(err error) error {
	if err == nil {
		return nil
	}

	valErrs, ok := err.(val.ValidationErrors)
	if !ok {
	}

	return errors.PropagateAs(
		KindValidationError,
		err,
		"validation error",
		errors.Hint(
			"One or more fields of the form were completed incorrectly. Please, fix the errors and try again.",
		),
		ValidationErrors(valErrs),
	)
}

// Global returns the singleton instance of the Validator.
// If the singleton is not yet initialized, it creates a new Validator instance
// and returns it. This ensures that the same Validator instance is used
// throughout the application, allowing for consistent validation logic.
func Global() *Validator {
	if singleton == nil {
		singleton = New()
	}

	return singleton
}

// Init initializes the global Validator instance with the provided
// field validators. If the singleton Validator is already set, it logs
// a fatal error to prevent reinitialization. This function ensures that
// the application uses a consistent set of validation rules by extending
// the default validations with any additional ones provided as arguments.
func Init(validations ...FieldValidator) {
	if singleton != nil {
		logging.Global().Fatal("Failed to initialize new global validator: this is already set")
	}

	v := New()
	v.Extend(validations...)
	singleton = v
}
