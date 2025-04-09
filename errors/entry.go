package errors

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
}

type EntryList []Entry

func Entries(entries ...Entry) EntryList {
	return entries
}

func (el EntryList) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, v := range el {
		if v.Value() != nil {
			enc.AddReflected(v.Key(), v.Value())
		}
	}

	return nil
}

func (el EntryList) Zap() zap.Field {
	return zap.Object("context", el)
}

func (el EntryList) AsUserScope() *ScopeT {
	return UserScope(el...)
}

func (el EntryList) AsSystemScope() *ScopeT {
	return SystemScope(el...)
}
