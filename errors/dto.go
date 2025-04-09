package errors

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
