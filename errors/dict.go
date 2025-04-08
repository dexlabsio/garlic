package errors

// DictT is a structure that represents a dictionary-like collection of entries,
// each identified by a unique key. It provides methods to add entries, retrieve
// the key of the dictionary, get the value representation of all entries, and
// insert entries from another DictT instance. This structure is useful for
// managing and aggregating key-value pairs in a structured manner.
type DictT struct {
	key     string
	entries map[string]Entry
}

// Dict creates a new instance of DictT with the specified key and entries.
// It initializes the dictionary with the provided entries, allowing for the
// management of key-value pairs. This function is useful for creating a structured
// collection of entries that can be easily accessed and manipulated.
func Dict(key string, entries ...Entry) *DictT {
	e := &DictT{
		key:     key,
		entries: make(map[string]Entry, len(entries)),
	}

	for _, entry := range entries {
		e.Add(entry)
	}

	return e
}

// Add inserts a new entry into the DictT. If an entry with the same key
// already exists, it merges the new entry with the existing one using
// the Insert method of the Entry interface. This function ensures that
// entries are aggregated correctly, allowing for updates or additions
// without overwriting existing data unless explicitly intended.
func (e *DictT) Add(entries ...Entry) {
	for _, entry := range entries {
		if entry == nil {
			continue
		}

		key := entry.Key()

		current, ok := e.entries[key]
		if !ok || current == nil {
			e.entries[key] = entry
		} else {
			e.entries[key] = current.Insert(entry)
		}
	}
}

// Key returns the key associated with the DictT instance. This key
// serves as an identifier for the dictionary and can be used to
// distinguish it from other dictionaries or entries. The key is
// typically used in contexts where the dictionary needs to be
// referenced or accessed by its unique identifier.
func (e *DictT) Key() string {
	return e.key
}

// Value returns a map representation of the entries within the DictT.
// Each entry is represented as a key-value pair, where the key is
// obtained from the entry's Key method and the value is obtained from
// the entry's Value method. If there are no entries, it returns nil.
// This method provides a structured way to access the contents of the
// dictionary, allowing for easy integration with other systems or
// serialization processes.
func (e *DictT) Value() any {
	res := map[string]any{}
	for _, entry := range e.entries {
		res[entry.Key()] = entry.Value()
	}

	if len(res) == 0 {
		return nil
	}

	return res
}

// Insert merges entries from another DictT instance into the current one.
// It takes an Entry as an argument, which is expected to be of type *DictT.
// If the provided entry is not a *DictT, the function will panic, indicating
// a type mismatch. This method iterates over the entries of the other DictT
// and adds them to the current instance using the Add method, ensuring that
// entries are correctly aggregated. This function is useful for combining
// multiple dictionaries into a single cohesive collection.
func (e *DictT) Insert(other Entry) Entry {
	o, ok := other.(*DictT)
	if !ok {
		panic("trying to insert unmatching entry")
	}

	for _, entry := range o.entries {
		e.Add(entry)
	}

	return e
}
