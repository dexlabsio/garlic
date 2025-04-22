package database

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
)

type Filter struct {
	key   string
	value any
}

func (f *Filter) Statement() string {
	return fmt.Sprintf("%s='%s'", f.key, f.value)
}

// ExtractFilters inspects any struct (or pointer to a struct) and returns a map
// where the keys are the values of the "filter" tag and the values are the underlying
// values pointed to by the field. For each field with the tag:
//   - If the field is not a pointer, the function panics.
//   - If the pointer is nil, the field is skipped.
//   - If the field has no "filter" tag, it is skipped.
func ExtractFilters(input interface{}) []*Filter {
	// Use reflection to obtain the input’s value.
	val := reflect.ValueOf(input)

	// If the input is a pointer, make sure it isn't nil and then dereference it.
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return []*Filter{}
		}
		val = val.Elem()
	}

	// The input must now be a struct; otherwise, we panic.
	if val.Kind() != reflect.Struct {
		panic("input must be a struct or a pointer to a struct")
	}

	// Prepare a map to hold the filters.
	uniqueFilters := map[string]*Filter{}
	typ := val.Type()

	// Iterate over all fields in the struct.
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		filterTag := field.Tag.Get("filter")
		// Skip if the tag "filter" is not set.
		if filterTag == "" {
			continue
		}

		fieldValue := val.Field(i)

		// Enforce that the field must be a pointer.
		if fieldValue.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("field %q is tagged with filter but is not a pointer", field.Name))
		}

		// If the pointer is nil, skip this field.
		if fieldValue.IsNil() {
			continue
		}

		filter := &Filter{
			key:   filterTag,
			value: fieldValue.Elem().Interface(),
		}

		// Dereference the pointer to get the actual value and add it to the map.
		uniqueFilters[filterTag] = filter
	}

	filters := slices.Collect(maps.Values(uniqueFilters))
	return filters
}
