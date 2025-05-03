package errors

import (
	"encoding/json"
)

type Transferable interface {
	ErrorDTO() *DTO
}

type DTO struct {
	Name    string         `json:"name" mapstructure:"name"`
	Error   string         `json:"error" mapstructure:"error"`
	Code    string         `json:"kind" mapstructure:"kind"`
	Details map[string]any `json:"details,omitempty" mapstructure:"details,omitempty"`
}

func NewDTO(err error) *DTO {
	e, ok := err.(Transferable)
	if !ok {
		e = Raw(KindError, err.Error())
	}

	return e.ErrorDTO()
}

func (dto *DTO) Decode() *ErrorT {
	return &ErrorT{
		kind:    GetByCode(dto.Code),
		message: dto.Error,
		Details: dto.Details,
	}
}

// JSON serializes the DTO struct into a JSON formatted byte slice.
// It returns the serialized data as json.RawMessage, which is a type alias for []byte.
// If an error occurs during the marshaling process, the function will panic.
func (dto *DTO) JSON() json.RawMessage {
	b, err := json.Marshal(dto)
	if err != nil {
		panic(err)
	}

	return json.RawMessage(b)
}
