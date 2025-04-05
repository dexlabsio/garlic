package errors

type context struct {
	opts []Opt
}

func Context(opts ...Opt) *context {
	return &context{opts}
}

func (u *context) Add(opts ...Opt) *context {
	u.opts = append(u.opts, opts...)
	return u
}
