package errors

type scope struct {
	key        string
	opts       map[string]Opt
	visibility Visibility
}

// Scope creates a new scope with the specified key and options,
// setting the visibility to PUBLIC. It returns a pointer to the
// newly created scope object.
func Scope(key string, opts ...Opt) *scope {
	return _scope(key, PUBLIC, opts...)
}

// RScope creates a new scope with the specified key and options,
// setting the visibility to RESTRICT. It returns a pointer to the
// newly created scope object.
func RScope(key string, opts ...Opt) *scope {
	return _scope(key, RESTRICT, opts...)
}

// UserScope creates a new scope with the key "context" and the provided options,
// setting the visibility to PUBLIC. It returns a pointer to the newly created scope object.
func UserScope(opts ...Opt) *scope {
	return Scope("context", opts...)
}

// SystemScope creates a new scope with the key "debug" and the provided options,
// setting the visibility to RESTRICT. It uses the caller's name as part of the scope
// key by attempting to find the caller's information. This function is typically used
// to create a restricted scope for system-level debugging and tracing purposes.
func SystemScope(opts ...Opt) *scope {
	caller := "unknown"
	if _, _, name, ok := findCaller(); ok {
		caller = name
	}

	call := RScope(caller, opts...)
	return RScope("debug", call)
}

func _scope(key string, visibility Visibility, opts ...Opt) *scope {
	s := &scope{
		key:        key,
		opts:       make(map[string]Opt, len(opts)),
		visibility: visibility,
	}

	for _, opt := range opts {
		s.insertOpt(opt)
	}

	return s
}

func (s *scope) insertOpt(opt Opt) {
	if opt == nil {
		return
	}

	key := opt.Key()

	current, ok := s.opts[key]
	if !ok || current == nil {
		s.opts[key] = opt
	} else {
		s.opts[key] = current.Insert(opt)
	}
}

func (s *scope) Key() string {
	return s.key
}

func (s *scope) Value() any {
	res := map[string]any{}
	for _, opt := range s.opts {
		if s.visibility == PUBLIC {
			if opt.Visibility() == PUBLIC {
				res[opt.Key()] = opt.Value()
			}
		} else {
			res[opt.Key()] = opt.Value()
		}
	}

	if len(res) == 0 {
		return nil
	}

	return res
}

func (s *scope) Visibility() Visibility {
	return s.visibility
}

func (s *scope) Append(others ...Opt) *scope {
	for _, opt := range others {
		s.insertOpt(opt)
	}

	return s
}

func (s *scope) Insert(other Opt) Opt {
	o, ok := other.(*scope)
	if !ok {
		panic("trying to insert scope unmatch opt")
	}

	for _, opt := range o.opts {
		s.insertOpt(opt)
	}

	return s
}
