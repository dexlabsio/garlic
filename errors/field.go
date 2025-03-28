package errors

type field struct {
	key   string
	value any
}

func Field(key string, value any) *field {
	return &field{
		key:   key,
		value: value,
	}
}

func (h *field) Key() string {
	return h.key
}

func (h *field) Value() any {
	return h.value
}

func (h *field) Visibility() Visibility {
	return PUBLIC
}

func (h *field) Insert(other Opt) Opt {
	return h
}
