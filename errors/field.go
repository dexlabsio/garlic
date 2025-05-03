package errors

const REDACTION_PLACEHOLDER = "****"

type FieldT struct {
	key   string
	value any
}

func Field(key string, value any) *FieldT {
	f := &FieldT{
		key:   key,
		value: value,
	}

	return f
}

func (f *FieldT) Key() string {
	return f.key
}

func (f *FieldT) Value() any {
	return f.value
}

func (f *FieldT) Insert(other Entry) Entry {
	return f
}

// RedactedString creates a partially visible string value for debugging purposes,
// adaptively showing approximately 1/3 of the value while protecting sensitive data.
func RedactedString(key, value string) Entry {
	length := len(value)
	if length < 5 {
		return Field(key, REDACTION_PLACEHOLDER)
	}

	// shows 1/3 of the content, half in the beggining, half in the end
	visibleChars := length / (3 * 2)
	if visibleChars < 1 {
		visibleChars = 1
	}

	prefix := value[:visibleChars]
	suffix := value[length-visibleChars:]
	redactedValue := prefix + REDACTION_PLACEHOLDER + suffix

	return Field(key, redactedValue)
}
