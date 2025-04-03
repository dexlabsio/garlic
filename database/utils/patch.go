package utils

import (
	"fmt"
	"reflect"

	"github.com/dexlabsio/garlic/errors"
	"github.com/google/uuid"
)

// ParsePatchFields parses the fields of a resource model used for partial updates.
// The majority of the complexity is due to the fact that we don't know which fields
// will be present in the patch operation. Therefore, the manipulations need to be
// generic enough to support the variety of those cases.
func ParsePatchFields(resource any) ([]string, []any, error) {
	v := reflect.ValueOf(resource)
	if v.Kind() == reflect.Ptr {
		if v.Elem().Kind() == reflect.Struct {
			v = v.Elem() // Dereference the pointer to get the struct
		} else {
			panic("Pointer does not point to a struct")
		}
	}

	t := v.Type()

	numFields := t.NumField()
	params := make([]string, 0, numFields)
	values := make([]any, 0, numFields+1) // need +1 because we're going to append the ID later

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "Id" {
			continue // Skip the ID field
		}

		value := v.Field(i)

		dbTag := field.Tag.Get("db")

		if dbTag == "" {
			continue // Skip fields without db tags
		}

		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue // Skip nil pointers, it means user didn't provide a value for this field
			}

			val := value.Elem() // Dereference the pointer
			values = append(values, val.Interface())
			params = append(params, fmt.Sprintf("%s = $%d", dbTag, len(params)+1))
		} else {
			panic(fmt.Sprintf("Patch structs can only have pointer fields. Field %s is invalid", field.Name))
		}
	}

	// Find and append the ID field
	idField := v.FieldByName("Id")
	if !idField.IsValid() {
		return nil, nil, errors.NewSystemError("missing required ID field")
	}

	id := idField.Interface().(*uuid.UUID)
	if id == nil {
		panic("the resource ID cannot be nil")
	}
	values = append(values, id)

	return params, values, nil
}
