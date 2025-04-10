package errors

import (
	"maps"
	"slices"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type context struct {
	entries map[string]Entry
}

func Context(entries ...Entry) *context {
	ectx := &context{
		entries: map[string]Entry{},
	}

	for _, entry := range entries {
		ectx.Insert(entry)
	}

	return ectx
}

func (u *context) Add(entries ...Entry) *context {
	for _, entry := range entries {
		u.Insert(entry)
	}

	return u
}

func (u *context) Opt() Opt {
	entries := slices.Collect(maps.Values(u.entries))
	return func(e *ErrorT) {
		e.Insert(userScope(entries))
		e.Insert(systemScope(entries))
	}
}

func (u *context) Insert(entry Entry) {
	if entry == nil {
		return
	}

	key := entry.Key()
	current, ok := u.entries[key]
	if !ok || current == nil {
		u.entries[key] = entry
	} else {
		u.entries[key] = current.Insert(entry)
	}
}

func (u *context) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, v := range u.entries {
		if v.Value() != nil {
			enc.AddReflected(v.Key(), v.Value())
		}
	}

	return nil
}

func (u *context) Zap() zap.Field {
	return zap.Object("context", u)
}
