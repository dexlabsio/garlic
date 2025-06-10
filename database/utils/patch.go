package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

func Named(query string, resource any) (string, []any) {
	query, args, err := sqlx.Named(query, resource)
	if err != nil {
		panic(fmt.Errorf("fatal failure trying to get named query: %w", err))
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		panic(fmt.Errorf("fatal failure trying to get expand named query: %w", err))
	}

	// Specific to postgres
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	return query, args
}

func JoinedPatchResourceBindings(resource any) string {
	return strings.Join(NamedResourceBindings(resource), ", ")
}

func NamedResourceBindings(resource any) []string {
	params := []string{}
	for k := range ResourceIter(resource) {
		params = append(params, fmt.Sprintf("%s = :%s", k, k))
	}

	return params
}

func ResourceIter(resource any) func(func(string, any) bool) {
	v := reflect.ValueOf(resource)
	if v.Kind() == reflect.Ptr {
		if v.Elem().Kind() == reflect.Struct {
			v = v.Elem() // Dereference the pointer to get the struct
		} else {
			panic("pointer does not point to a struct")
		}
	}

	t := v.Type()
	return func(yield func(string, any) bool) {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			dbTag := field.Tag.Get("db")

			if dbTag == "" {
				continue // Skip fields without db tags
			}

			if value.Kind() == reflect.Ptr {
				if value.IsNil() {
					continue // Skip nil pointers, it means user didn't provide a value for this field
				}

				val := value.Elem().Interface()
				if !yield(dbTag, val) {
					return
				}
			}
		}
	}
}
