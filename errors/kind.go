package errors

import "fmt"

var (
	KindUnknownError = &Kind{Name: "Unknown Error"}
	KindUserError    = &Kind{Name: "User Error"}
	KindSystemError  = &Kind{Name: "System Error"}
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

	return fmt.Sprintf("%s::%s", k.Name, k.Parent.Hierarchy())
}

func (k *Kind) Is(other *Kind) bool {
	for current := k; current != nil; current = current.Parent {
		if current == other {
			return true
		}
	}
	return false
}
