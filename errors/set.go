package errors

type SetT struct {
	entries []Entry
}

func Set(entries ...Entry) *SetT {
	set := &SetT{
		entries: []Entry{},
	}

	set.Extend(entries...)

	return set
}

func (s *SetT) Values() []Entry {
	return s.entries
}

func (s *SetT) Insert(e Entry) {
	for i, existing := range s.entries {
		if existing.Key() == e.Key() {
			s.entries[i] = s.entries[i].Insert(e)
			return
		}
	}

	s.entries = append(s.entries, e)
}

func (s *SetT) Extend(es ...Entry) {
	for _, e := range es {
		s.Insert(e)
	}
}
