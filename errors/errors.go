package errors

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Visibility uint8

const (
	PUBLIC Visibility = iota
	RESTRICT
)

type Opt func(*ErrorT)

type Opter interface {
	Opt() Opt
}

type Error interface {
	Kind() *Kind
}

type ErrorT struct {
	kind    *Kind
	message string
	cause   error
	entries map[string]Entry
}

// Propagate creates a new ErrorT instance with a base error kind, message, and options,
// and wraps an existing error with this new instance. It appends additional options for
// reverse trace and stack trace to the provided options, ensuring that the error context
// is enriched with detailed tracing information. This function is useful for propagating
// errors while maintaining comprehensive error tracking and debugging capabilities.
func Propagate(err error, message string, opts ...Opter) *ErrorT {
	kind := KindUnknownError
	if kinder, ok := err.(Error); ok {
		kind = kinder.Kind()
	}

	return PropagateAs(kind, err, message, opts...)
}

// PropagateAs creates a new ErrorT instance with a specified error kind, message, and options,
// and wraps an existing error with this new instance. It appends additional options for
// reverse trace and stack trace to the provided options, ensuring that the error context
// is enriched with detailed tracing information. This function is useful for propagating
// errors with a specific kind while maintaining comprehensive error tracking and debugging capabilities.
func PropagateAs(kind *Kind, err error, message string, opts ...Opter) *ErrorT {
	e := New(kind, message, opts...)
	return e.wrap(err)
}

// New creates a new instance of ErrorT with the specified kind, message, and options.
// It initializes the ErrorT structure, sets the kind and message, and processes the provided
// options by inserting them into the opts map using the insert method. This function is
// essential for constructing error objects with additional context and metadata, which can
// be used for detailed error reporting and handling.
func New(kind *Kind, message string, opts ...Opter) *ErrorT {
	opts = append(opts, StackTrace(), RevTrace())
	return Raw(kind, message, opts...)
}

// Raw creates a new instance of ErrorT with the specified kind, message, and options.
// Unlike the New function, it does not append additional options for stack trace or reverse trace.
// This function is useful when you want to create an error object without automatically
// adding tracing information, allowing for more control over the error's metadata and context.
func Raw(kind *Kind, message string, opts ...Opter) *ErrorT {
	e := &ErrorT{
		kind:    kind,
		message: message,
		entries: map[string]Entry{},
	}

	for _, opter := range opts {
		opter.Opt()(e)
	}

	return e
}

// From creates a new ErrorT instance from an existing error, using the error's
// message and kind (if available) as the basis for the new instance. It allows
// for additional options to be specified, which are inserted into the new ErrorT
// instance. This function is useful for converting standard errors into ErrorT
// instances, enabling enhanced error tracking and handling with additional context
// and metadata.
func From(err error, opts ...Opter) *ErrorT {
	kind := KindUnknownError
	if kinder, ok := err.(Error); ok {
		kind = kinder.Kind()
	}

	e := Raw(kind, err.Error(), opts...)
	return e.wrap(err)
}

// Find retrieves an option from the ErrorT instance's options map using the specified key.
// It returns the option and a boolean indicating whether the option was found. This method
// is useful for accessing specific error metadata or context that has been stored in the
// options map, allowing for detailed inspection and handling of error-related information.
func (e *ErrorT) Find(key string) (entry Entry, ok bool) {
	entry, ok = e.entries[key]
	return
}

// Kind returns the kind of the ErrorT instance.
// This method provides access to the error kind, which is used to
// categorize and identify the nature of the error. It is useful for
// error handling and reporting, allowing developers to determine the
// specific type of error encountered.
func (e *ErrorT) Kind() *Kind {
	return e.kind
}

// As sets the kind of the ErrorT instance to the specified kind.
// This method allows for modifying the error kind of an existing ErrorT
// instance, enabling dynamic categorization of errors based on their
// nature or origin. It returns the modified ErrorT instance, allowing
// for method chaining.
func (e *ErrorT) As(kind *Kind) *ErrorT {
	err := e.Copy()
	err.kind = kind
	return err
}

// With creates a new ErrorT instance by copying the current ErrorT instance
// and merging it with the options from the provided context. This method
// allows for the augmentation of an existing error with additional context
// options, facilitating more detailed error reporting and handling. It
// returns a new ErrorT instance that includes both the original and
// context-specific options.
func (e *ErrorT) With(opts ...Opter) *ErrorT {
	err := e.Copy()
	for _, opter := range opts {
		opter.Opt()(err)
	}

	return err
}

// Copy creates a deep copy of the current ErrorT instance, including its kind,
// message, cause, and options. This method is useful for duplicating an error
// object while preserving its original state and metadata, allowing for
// independent modifications or further processing without affecting the
// original error instance.
func (e *ErrorT) Copy() *ErrorT {
	err := &ErrorT{
		kind:    e.kind,
		message: e.message,
		cause:   e.cause,
		entries: make(map[string]Entry, len(e.entries)),
	}

	for k, entry := range e.entries {
		err.entries[k] = entry
	}

	return err
}

// wrap takes an existing error and wraps it with the current ErrorT instance,
// incorporating any options from the existing error into the current instance.
// If the existing error is of type ErrorT, its options are merged into the current
// instance using the insert method. This allows for the aggregation of error
// context and metadata, facilitating enhanced error tracking and debugging.
func (e *ErrorT) wrap(other error) *ErrorT {
	if other == nil {
		return e
	}

	if o, ok := other.(*ErrorT); ok {
		for _, entry := range o.entries {
			e.Insert(entry)
		}
	}

	e.cause = other
	return e
}

// Unwrap returns the wrapped error from the ErrorT instance.
// This method is used to retrieve the original error that was wrapped
// by the ErrorT instance, enabling error unwrapping and inspection
// in error handling workflows.
func (e *ErrorT) Unwrap() error {
	return e.cause
}

// Insert controls when to call opt.Insert for a new opt.
// If the key is missing or its value is nil, we just insert
// the given opt directly. If it exists, we call the respective
// Insert function to handle the newly introduced object.
func (e *ErrorT) Insert(entry Entry) {
	if entry == nil {
		return
	}

	key := entry.Key()
	current, ok := e.entries[key]
	if !ok || current == nil {
		e.entries[key] = entry
	} else {
		e.entries[key] = current.Insert(entry)
	}
}

// Error returns the error message for the ErrorT instance.
// If the ErrorT instance wraps another error, this method
// appends the wrapped error's message to the current error
// message, providing a complete error description. This is
// useful for error reporting and logging, as it gives a
// comprehensive view of the error chain.
func (e *ErrorT) Error() string {
	message := e.message
	if e.cause != nil {
		message = fmt.Sprintf("%s: %s", message, e.cause.Error())
	}

	return message
}

// DTO converts the ErrorT instance into a map suitable for data transfer objects (DTOs).
// This method constructs a DTO object containing the error message and kind hierarchy, and
// includes additional details from the options map if they are marked as PUBLIC visibility.
// It ensures that only relevant and non-sensitive information is exposed, making it
// suitable for returning error details in API responses or logs.
func (e *ErrorT) DTO() *DTO {
	content := &DTO{
		Error: e.message,
		Kind:  e.kind.Hierarchy(),
	}

	details := map[string]any{}
	for k, v := range e.entries {
		if v.Visibility() == PUBLIC {
			if v.Value() != nil {
				details[k] = v.Value()
			}
		}
	}

	if len(details) > 0 {
		content.Details = details
	}

	return content
}

// MarshalLogObject encodes the ErrorT instance into a zapcore.ObjectEncoder for structured logging.
// This method adds the error message, kind, and any additional details from the options map to the
// encoder. It ensures that all relevant error information is captured in the log, facilitating
// comprehensive error tracking and debugging when using the zap logging library.
func (e *ErrorT) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("message", e.message)
	enc.AddString("error", e.Error())
	enc.AddString("kind", e.kind.Hierarchy())

	details := map[string]any{}
	for k, v := range e.entries {
		if v.Value() != nil {
			details[k] = v.Value()
		}
	}

	enc.AddReflected("details", details)
	return nil
}

// Zap creates a zap.Field for logging an error using the zap logging library.
// If the provided error is of type ErrorT, it logs the error as a zap object,
// which includes detailed error information such as the message, kind, and options.
// Otherwise, it logs the error using zap.Error, which captures the error message
// and stack trace. This function is useful for integrating structured error logging
// into applications using the zap logging framework.
func Zap(err error) zap.Field {
	if e, ok := err.(*ErrorT); ok {
		return zap.Object("error", e)
	}

	return zap.Error(err)
}
