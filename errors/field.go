package errors

const REDACTION_PLACEHOLDER = "****"

type FieldT struct {
	key        string
	value      any
	visibility Visibility
}

// Field creates a new Entry with the specified key and value.
// This function is used to encapsulate key-value pairs, which can be
// utilized for logging or error reporting purposes. The returned Entry
// implements the necessary interface to be compatible with the error
// handling and logging system, allowing for structured data to be
// associated with errors or log entries.
func Field(key string, value any, opts ...EntryOpt) *FieldT {
	f := &FieldT{
		key:        key,
		value:      value,
		visibility: PUBLIC,
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (f *FieldT) Key() string {
	return f.key
}

func (f *FieldT) Value() any {
	return f.value
}

func (f *FieldT) Visibility() Visibility {
	return f.visibility
}

func (f *FieldT) Insert(other Entry) Entry {
	return f
}

// RedactedString creates a partially visible string value for debugging purposes,
// adaptively showing approximately 1/3 of the value while protecting sensitive data.
func RedactedString(key, value string, opts ...EntryOpt) Entry {
	length := len(value)
	if length < 5 {
		return Field(key, REDACTION_PLACEHOLDER, opts...)
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
