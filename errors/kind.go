package errors

import "fmt"

const KIND_SEPARATOR = "::"

var (
	KindUnknownError = &Kind{Name: "Unknown Error"}
	KindUserError    = &Kind{Name: "User Error"}
	KindSystemError  = &Kind{Name: "System Error"}

	KindExternalUnknownError = &Kind{Name: "External Unknown Error", Parent: KindUnknownError}
	KindExternalUserError    = &Kind{Name: "External User Error", Parent: KindUserError}
	KindExternalSystemError  = &Kind{Name: "External System Error", Parent: KindSystemError}
)

type Kind struct {
	Name   string
	Parent *Kind
}

func NewKind(name string, parent *Kind) *Kind {
	return &Kind{
		Name:   name,
		Parent: parent,
	}
}

func (k *Kind) Hierarchy() string {
	if k.Parent == nil {
		return k.Name
	}

	return fmt.Sprintf("%s%s%s", k.Name, KIND_SEPARATOR, k.Parent.Hierarchy())
}

func (k *Kind) Is(other *Kind) bool {
	for current := k; current != nil; current = current.Parent {
		if current == other {
			return true
		}
	}
	return false
}
