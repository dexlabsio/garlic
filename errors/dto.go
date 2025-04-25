package errors

import "strings"

type Transferable interface {
	DTO() *DTO
}

type DTO struct {
	Error   string         `json:"error" mapstructure:"error"`
	Kind    string         `json:"kind" mapstructure:"kind"`
	Details map[string]any `json:"details,omitempty" mapstructure:"details,omitempty"`
}

func NewDTO(err error) *DTO {
	e, ok := err.(Transferable)
	if !ok {
		e = From(err)
	}

	return e.DTO()
}

func (dto *DTO) Parse() *ErrorT {
	kind := dto.DecodeKind()
	err := New(kind, dto.Error)
	err.extension = dto.Details

	return err
}

func (dto *DTO) DecodeKind() *Kind {
	kindHierarchy := strings.Split(dto.Kind, KIND_SEPARATOR)

	for _, k := range kindHierarchy {
		switch k {
		case KindUserError.Name:
			return KindExternalUserError
		case KindSystemError.Name:
			return KindExternalSystemError
		}
	}

	return KindExternalUnknownError
}
