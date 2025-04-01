package errors

const REDACTION_PLACEHOLDER = "****"

type field struct {
	key        string
	value      any
	visibility Visibility
}

func newField(key string, value any, v Visibility) *field {
	return &field{
		key:        key,
		value:      value,
		visibility: v,
	}
}

// Field creates a new field with the given key and value,
// setting its visibility to PUBLIC. This function is used to
// create key-value pairs that can be attached to errors for
// additional context or metadata.
func Field(key string, value any) *field {
	return newField(key, value, PUBLIC)
}

// RField creates a new field with the given key and value,
// setting its visibility to RESTRICT. This function is used to
// create key-value pairs that can be attached to errors for
// additional context or metadata, but with restricted visibility.
func RField(key string, value any) *field {
	return newField(key, value, RESTRICT)
}

// RedactedString creates a partially visible string value for debugging purposes,
// adaptively showing approximately 1/3 of the value while protecting sensitive data.
func RedactedString(key, value string) *field {
	length := len(value)
	if length < 5 {
		return newField(key, REDACTION_PLACEHOLDER, RESTRICT)
	}

	// shows 1/3 of the content, half in the beggining, half in the end
	visibleChars := length / (3 * 2)
	if visibleChars < 1 {
		visibleChars = 1
	}

	prefix := value[:visibleChars]
	suffix := value[length-visibleChars:]
	redactedValue := prefix + REDACTION_PLACEHOLDER + suffix

	return newField(key, redactedValue, RESTRICT)
}

func (h *field) Key() string {
	return h.key
}

func (h *field) Value() any {
	return h.value
}

func (h *field) Visibility() Visibility {
	return h.visibility
}

func (h *field) Insert(other Opt) Opt {
	return h
}
