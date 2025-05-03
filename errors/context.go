package errors

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ContextT struct {
	caller  string
	entries *SetT
}

// Context creates a new ContextT instance with the provided entries.
// It determines the caller's function name using findCaller() and sets it as the caller.
// The entries are stored in a SetT, which allows for efficient management of log entries.
// This function is useful for capturing and organizing contextual information
// that can be logged or used for troubleshooting purposes.
func Context(entries ...Entry) *ContextT {
	caller := "unknown"
	if _, _, name, ok := findCaller(); ok {
		caller = name
	}

	return &ContextT{
		caller:  caller,
		entries: Set(entries...),
	}
}

func (c *ContextT) Add(entries ...Entry) *ContextT {
	c.entries.Extend(entries...)
	return c
}

func (c *ContextT) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("function", c.caller)
	for _, v := range c.entries.Values() {
		if v.Value() != nil {
			enc.AddReflected(v.Key(), v.Value())
		}
	}

	return nil
}

func (c *ContextT) Zap() zap.Field {
	return zap.Object("context", c)
}

func (c *ContextT) Opt(e *ErrorT) {
	outputs := make(map[string]any, len(c.entries.Values()))
	for _, entry := range c.entries.Values() {
		outputs[entry.Key()] = entry.Value()
	}

	if e.Troubleshooting.Context == nil {
		e.Troubleshooting.Context = map[string]any{}
	}

	e.Troubleshooting.Context[c.caller] = outputs
}
