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

type Opt interface {
	// Key controls the name of the key that will appear for the
	// user and developer in DTOs and Logs.
	Key() string

	// Value is any object suitable to be present in DTOs or Logs.
	Value() any

	// Visibility says if a specific Opt should be returned in DTO.
	Visibility() Visibility

	// Insert controls how a new value is inserted in the map of
	// existing Opts. For example, if a specific opt value is inserted
	// but the key already exists, this function controls if the new
	// should be aggregated or replaced.
	Insert(other Opt) Opt
}

type ErrorT struct {
	message string
	wrap    error
	opts    map[string]Opt
}

// New creates a new ErrorT instance with the specified message and options.
// It appends RevTrace and StackTrace options to the provided options list
// and then calls the Raw function to construct the ErrorT object. This function
// is typically used to create a new error with additional context or metadata.
func New(message string, opts ...Opt) *ErrorT {
	opts = append(opts,
		RevTrace(),
		StackTrace(),
	)

	return Raw(message, opts...)
}

// Propagate creates a new ErrorT instance with the specified message and options,
// appending RevTrace and StackTrace options to the provided options list. It then
// wraps the given error with the newly created ErrorT instance. This function is
// typically used to propagate an existing error with additional context or metadata.
func Propagate(message string, err error, opts ...Opt) *ErrorT {
	opts = append(opts,
		RevTrace(),
		StackTrace(),
	)

	e := Raw(message, opts...)
	return e.Wrap(err)
}

// Raw constructs a new ErrorT instance with the given message and options.
// It initializes an empty map for options and iterates over the provided options,
// inserting each one into the map using the insert method. This function is used
// internally to create an ErrorT object with specific options, allowing for
// customization and additional context to be added to the error.
func Raw(message string, opts ...Opt) *ErrorT {
	e := ErrorT{
		message: message,
		opts:    map[string]Opt{},
	}

	for _, opt := range opts {
		e.insert(opt)
	}

	return &e
}

func (e *ErrorT) Wrap(other error) *ErrorT {
	if other == nil {
		return e
	}

	if o, ok := other.(*ErrorT); ok {
		for _, opt := range o.opts {
			e.insert(opt)
		}
	}

	e.wrap = other
	return e
}

func (e *ErrorT) Unwrap() error {
	return e.wrap
}

// insert controls when to call opt.Insert for a new opt.
// If the key is missing or its value is nil, we just insert
// the given opt directly. If it exists, we call the respective
// Insert function to handle the newly introduced object.
func (e *ErrorT) insert(opt Opt) {
	if opt == nil {
		return
	}

	key := opt.Key()
	current, ok := e.opts[key]
	if !ok || current == nil {
		e.opts[key] = opt
	} else {
		e.opts[key] = current.Insert(opt)
	}
}

func (e *ErrorT) InsertOpts(opts []Opt) {
	for _, opt := range opts {
		e.insert(opt)
	}
}

func (e *ErrorT) Copy() *ErrorT {
	newError := &ErrorT{
		message: e.message,
		wrap:    e.wrap,
	}

	for _, opt := range e.opts {
		newError.insert(opt)
	}

	return newError
}

func (e *ErrorT) Error() string {
	message := e.message
	if e.wrap != nil {
		message = fmt.Sprintf("%s: %s", message, e.wrap.Error())
	}

	return message
}

func (e *ErrorT) DTO() map[string]any {
	details := map[string]any{}
	for k, v := range e.opts {
		if v.Visibility() == PUBLIC {
			if v.Value() != nil {
				details[k] = v.Value()
			}
		}
	}

	return map[string]any{
		"error":   e.message,
		"details": details,
	}
}

func (e *ErrorT) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("message", e.message)
	enc.AddString("error", e.Error())

	details := map[string]any{}
	for k, v := range e.opts {
		if v.Value() != nil {
			details[k] = v.Value()
		}
	}

	enc.AddReflected("details", details)
	return nil
}

func Zap(err error) zap.Field {
	if e, ok := err.(*ErrorT); ok {
		return zap.Object("error", e)
	}

	return zap.Error(err)
}
