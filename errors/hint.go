package errors

type hint struct {
	message string
}

func Hint(message string) *hint {
	return &hint{
		message: message,
	}
}

func (h *hint) Key() string {
	return "hint"
}

func (h *hint) Value() any {
	return h.message
}

func (h *hint) Visibility() Visibility {
	return PUBLIC
}

func (h *hint) Insert(other Opt) Opt {
	return h
}
