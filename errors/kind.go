package errors

import (
	"fmt"
	"net/http"
)

const (
	KIND_FQN_SEPARATOR      = "::"
	HTTP_STATUS_NOT_DEFINED = 0
)

var (
	registeredCodes = map[string]*Kind{}
	registeredNames = map[string]*Kind{}
)

type Kind struct {
	Name           string
	Code           string
	Description    string
	HTTPStatusCode int
	Parent         *Kind
}

// Register adds a new Kind instance to the global registry of error kinds.
// It checks if a Kind with the same code has already been registered, and if so,
// it panics to prevent duplicate registrations. This function ensures that each
// Kind is uniquely identified by its code within the application, allowing for
// consistent error categorization and handling.
func Register(kinds ...*Kind) {
	for _, kind := range kinds {
		if _, ok := registeredCodes[kind.Code]; ok {
			panic("another kind with the same code was already registered")
		}

		if _, ok := registeredNames[kind.Name]; ok {
			panic("another kind with the same name was already registered")
		}

		registeredCodes[kind.Code] = kind
		registeredNames[kind.Name] = kind
	}
}

// GetByCode retrieves a Kind instance from the global registry using the provided code.
// If the code does not correspond to any registered Kind, the function panics,
// indicating that the requested error kind does not exist. This function is
// essential for accessing predefined error kinds based on their unique codes,
// facilitating error handling and categorization within the application.
func GetByCode(code string) *Kind {
	kind, ok := registeredCodes[code]
	if !ok {
		panic(fmt.Errorf("error kind with code `%s` doesn't exist", code))
	}

	return kind
}

// Get retrieves a Kind instance from the global registry using the provided name.
// If the name does not correspond to any registered Kind, the function panics,
// indicating that the requested error kind does not exist. This function is
// crucial for accessing predefined error kinds based on their unique names,
// facilitating error handling and categorization within the application.
func Get(name string) *Kind {
	kind, ok := registeredNames[name]
	if !ok {
		panic(fmt.Errorf("error kind with name `%s` doesn't exist", name))
	}

	return kind
}

// FQN returns a string representation of the Kind's hierarchy.
// It constructs the hierarchy by concatenating the Kind's name with its
// parent's hierarchy, separated by the KIND_FQN_SEPARATOR. If the Kind has
// no parent, it simply returns its name. This method is useful for
// understanding the hierarchical structure of error kinds.
func (k *Kind) FQN() string {
	if k.Parent == nil {
		return k.Name
	}

	return fmt.Sprintf("%s%s%s", k.Name, KIND_FQN_SEPARATOR, k.Parent.FQN())
}

// StatusCode returns the HTTP status code associated with the Kind instance.
// It traverses up the hierarchy of the Kind, checking each ancestor for an
// assigned HTTP status code. If a Kind in the hierarchy has an HTTP status
// code, it returns that code. If no HTTP status code is found in the hierarchy,
// it defaults to returning http.StatusInternalServerError.
func (k *Kind) StatusCode() int {
	for current := k; current != nil; current = current.Parent {
		if current.HTTPStatusCode != HTTP_STATUS_NOT_DEFINED {
			return current.HTTPStatusCode
		}
	}

	return http.StatusInternalServerError
}

// Is checks if the current Kind instance matches the specified other Kind instance
// by comparing their codes. It traverses up the hierarchy of the current Kind,
// checking each ancestor's code against the code of the other Kind. If a match
// is found, it returns true, indicating that the two Kinds are equivalent or
// related in the hierarchy. Otherwise, it returns false.
func (k *Kind) Is(other *Kind) bool {
	for current := k; current != nil; current = current.Parent {
		if current.Code == other.Code {
			return true
		}
	}
	return false
}
