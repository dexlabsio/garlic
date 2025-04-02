package errors

type OptInserter interface {
	InsertOpts([]Opt)
}

type using struct {
	opts []Opt
}

func Using(opts ...Opt) *using {
	return &using{opts}
}

func (u *using) Add(opts ...Opt) *using {
	u.opts = append(u.opts, opts...)
	return u
}

func (u *using) NewUserError(message string, opts ...Opt) *UserError {
	return Enrich(NewUserError(message, opts...), u)
}

func (u *using) NewSystemError(message string, opts ...Opt) *SystemError {
	return Enrich(NewSystemError(message, opts...), u)
}

func (u *using) New(message string, opts ...Opt) *ErrorT {
	return Enrich(New(message, opts...), u)
}

func (u *using) Propagate(message string, err error, opts ...Opt) *ErrorT {
	return Enrich(Propagate(message, err, opts...), u)
}

func (u *using) Raw(message string, opts ...Opt) *ErrorT {
	return Enrich(Raw(message, opts...), u)
}

// Enrich takes an OptInserter and a using instance, and inserts the options
// from the using instance into the OptInserter. This function is useful for
// enriching an existing error or object with additional options or metadata
// that are encapsulated within the using instance.
func Enrich[T OptInserter](inserter T, u *using) T {
	inserter.InsertOpts(u.opts)
	return inserter
}
