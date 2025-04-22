package toolkit

import "reflect"

func IsValueNil(i interface{}) bool {
	if i == nil {
		return true
	}
	value := reflect.ValueOf(i)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice:
		return value.IsNil()
	}
	return false
}

func PointerOf[T any](value T) *T {
	return &value
}

func ValueOrDefault[T any](value *T) T {
	if value == nil {
		var zero T
		return zero
	}
	return *value
}
