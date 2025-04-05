package errors

import (
	stderrors "errors"
)

// Is checks whether the error 'err' is equivalent to the 'target' error.
// It uses the standard library's errors.Is function to perform this check,
// which compares errors based on their types and values, including any
// wrapped errors in the chain.
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// As attempts to set the target to the first error in the error chain
// that matches the target's type. It uses the AsEmbedded function to
// perform this check, which not only considers the error itself but
// also any embedded fields within the error that might match the target's type.
// Returns true if a match is found and the target is set, otherwise false.
func As(err error, target any) bool {
	return stderrors.As(err, target)
}

// AsKind checks if the provided error 'err' or any error in its chain
// is of the specified 'kind'. It unwraps the error chain and looks for
// an error of type *ErrorT that matches the given kind. If a match is
// found, it returns true along with the matched *ErrorT. Otherwise, it
// returns false and nil, indicating no match was found in the error chain.
func AsKind(err error, kind *Kind) (*ErrorT, bool) {
	for current := err; current != nil; current = stderrors.Unwrap(current) {
		if e, ok := current.(*ErrorT); ok {
			if e.kind.Is(kind) {
				return e, true
			}
		}
	}

	return nil, false
}

// IsKind checks if the provided error 'err' or any error in its chain
// is of the specified 'kind'. It utilizes the AsKind function to determine
// if there is an error of type *ErrorT in the chain that matches the given kind.
// Returns true if such an error is found, otherwise false.
func IsKind(err error, kind *Kind) bool {
	_, ok := AsKind(err, kind)
	return ok
}
