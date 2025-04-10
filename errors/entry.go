package errors

type Entry interface {
	// Insert merges the current Entry with another Entry, allowing for
	// the combination or replacement of values. This method is useful
	// for updating or aggregating entries in a structured way.
	Insert(Entry) Entry

	// Key controls the name of the key that will appear for the
	// user and developer in DTOs and Logs.
	Key() string

	// Value is any object suitable to be present in DTOs or Logs.
	Value() any

	// Visibility says if a specific Opt should be returned in DTO.
	Visibility() Visibility
}

type EntryOpt func(Entry)

type VisibilitySetter interface {
	SetVisibility(Visibility)
}

func Restrict(e Entry) {
	if vs, ok := e.(VisibilitySetter); ok {
		vs.SetVisibility(RESTRICT)
	}
}

func Public(e Entry) {
	if vs, ok := e.(VisibilitySetter); ok {
		vs.SetVisibility(PUBLIC)
	}
}
