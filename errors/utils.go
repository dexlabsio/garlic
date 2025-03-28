package errors

import (
	stderrors "errors"
	"reflect"
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
	return AsEmbedded(err, target)
}

// findEmbedded recursively searches for a value whose type is assignable to targetType.
// It uses reflect.Type.AssignableTo instead of unconditionally dereferencing pointers.
func findEmbedded(v reflect.Value, targetType reflect.Type) (reflect.Value, bool) {
	// If v's type is already assignable to targetType, return it.
	if v.Type().AssignableTo(targetType) {
		return v, true
	}
	// If v is a pointer, try its element.
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Value{}, false
		}
		if res, ok := findEmbedded(v.Elem(), targetType); ok {
			return res, true
		}
	}
	// If v is a struct, iterate through its anonymous fields.
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldInfo := v.Type().Field(i)
			if fieldInfo.Anonymous {
				if res, ok := findEmbedded(field, targetType); ok {
					return res, true
				}
			}
		}
	}
	return reflect.Value{}, false
}

// AsEmbedded acts like errors.As but also checks whether any of the unwrapped errors
// have an embedded field whose type matches the target. If a match is found,
// it assigns the value to target and returns true.
func AsEmbedded(err error, target interface{}) bool {
	if err == nil {
		return false
	}
	// target must be a pointer to a type; get the type that we’re looking for.
	targetType := reflect.TypeOf(target).Elem()

	// Walk the error chain.
	for err != nil {
		// First, try the standard errors.As.
		if stderrors.As(err, target) {
			return true
		}
		// Use reflection to inspect the concrete value.
		val := reflect.ValueOf(err)
		if embeddedVal, found := findEmbedded(val, targetType); found {
			// Assign the found embedded value to the target pointer.
			reflect.ValueOf(target).Elem().Set(embeddedVal)
			return true
		}
		// Move to the next error in the chain.
		err = stderrors.Unwrap(err)
	}
	return false
}
