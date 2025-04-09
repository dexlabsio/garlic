package errors

import (
	"go.uber.org/zap/zapcore"
)

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

func (u *context) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, v := range u.opts {
		if v.Value() != nil {
			enc.AddReflected(v.Key(), v.Value())
		}
	}

	return nil
}
