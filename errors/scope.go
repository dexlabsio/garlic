package errors

type ScopeT struct {
	*DictT
	visibility Visibility
}

func _scope(key string, visibility Visibility, entries ...Entry) *ScopeT {
	s := &ScopeT{
		DictT:      Dict(key, entries...),
		visibility: visibility,
	}

	return s
}

func Scope(key string, entries ...Entry) *ScopeT {
	return _scope(key, PUBLIC, entries...)
}

func RScope(key string, entries ...Entry) *ScopeT {
	return _scope(key, RESTRICT, entries...)
}

func UserScope(entries ...Entry) *ScopeT {
	return Scope("user", entries...)
}

func SystemScope(entries ...Entry) *ScopeT {
	caller := "unknown"
	if _, _, name, ok := findCaller(); ok {
		caller = name
	}

	call := Dict(caller, entries...)
	return RScope("system", call)
}

func (s *ScopeT) Insert(other Opt) Opt {
	o, ok := other.(*ScopeT)
	if !ok {
		panic("trying to insert unmatching opt")
	}

	s.DictT.Insert(o.DictT)
	return s
}

func (s *ScopeT) Visibility() Visibility {
	return s.visibility
}
