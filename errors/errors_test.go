//go:build unit
// +build unit

package errors

import (
	"fmt"
	"testing"
)

func TestPropagationKind(t *testing.T) {
	var (
		rootErr            = fmt.Errorf("Root Error")
		KindBaseError      = NewKind("Base Error", nil)
		KindOtherBaseError = NewKind("Other Base Error", nil)
	)

	cases := []struct {
		title string
		fn    func() error
		kind  *Kind
	}{
		{
			title: "standard errors should be propagated as KindUnknownError by default",
			fn: func() error {
				return Propagate(rootErr, "test error")
			},
			kind: KindUnknownError,
		},
		{
			title: "base errors should be propagated as the given kind",
			fn: func() error {
				return PropagateAs(KindBaseError, rootErr, "test error")
			},
			kind: KindBaseError,
		},
		{
			title: "secondary errors should be propagated as the given kind",
			fn: func() error {
				err := PropagateAs(KindBaseError, rootErr, "test error")
				return PropagateAs(KindOtherBaseError, err, "test other error")
			},
			kind: KindOtherBaseError,
		},
		{
			title: "errors should be propagated as the given kind override",
			fn: func() error {
				err := PropagateAs(KindBaseError, rootErr, "test error")
				return Propagate(err, "test other error").As(KindOtherBaseError)
			},
			kind: KindOtherBaseError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			err := tc.fn()
			if !IsKind(err, tc.kind) {
				t.Errorf("expected IsKind to return true for error `%s`, but got false", tc.kind.Name)
			}
		})
	}
}

func TestPropagationOfUserScopeEntries(t *testing.T) {
	var (
		rootErr            = fmt.Errorf("Root Error")
		KindBaseError      = NewKind("Base Error", nil)
		KindOtherBaseError = NewKind("Other Base Error", nil)
		FieldA             = Field("testA", "A")
		FieldAPlus         = Field("testA", "A+")
		FieldB             = Field("testB", "B")
		FieldC             = Field("testC", "C")
		ScopeA             = UserScope(FieldA)
		ScopeAPlus         = UserScope(FieldAPlus)
		ScopeBC            = UserScope(FieldB, FieldC)
	)

	cases := []struct {
		title          string
		fn             func() error
		expectedFields []Entry
	}{
		{
			title: "first level errors should propagate their fields",
			fn: func() error {
				return PropagateAs(KindBaseError, rootErr, "test error", ScopeA)
			},
			expectedFields: []Entry{FieldA},
		},
		{
			title: "override errors should propagate previous error fields",
			fn: func() error {
				err := PropagateAs(KindBaseError, rootErr, "test error", ScopeA)
				return PropagateAs(KindOtherBaseError, err, "other test error")
			},
			expectedFields: []Entry{FieldA},
		},
		{
			title: "secondary level errors should propagate previous error fields",
			fn: func() error {
				err := PropagateAs(KindBaseError, rootErr, "test error", ScopeA)
				return Propagate(err, "other test error")
			},
			expectedFields: []Entry{FieldA},
		},
		{
			title: "secondary level errors should propagate previous error fields plus its own fields",
			fn: func() error {
				err := PropagateAs(KindBaseError, rootErr, "test error", ScopeA)
				return Propagate(err, "other test error", ScopeBC)
			},
			expectedFields: []Entry{FieldA, FieldB, FieldC},
		},
		{
			title: "secondary level errors should override previous error fields with the same key",
			fn: func() error {
				err := PropagateAs(KindBaseError, rootErr, "test error", ScopeA)
				return Propagate(err, "other test error", ScopeAPlus)
			},
			expectedFields: []Entry{FieldAPlus},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			err := tc.fn()
			errt, ok := err.(*ErrorT)
			if !ok {
				t.Errorf("tried to parse *ErrorT but this is not valid")
			}

			dto := errt.DTO()
			detailsDTO, ok := dto["details"]
			if !ok {
				t.Errorf("dto doesn't have an expected details field")
			}

			details, ok := detailsDTO.(map[string]any)
			if !ok {
				t.Errorf("dto.details should be a `map[string]any`")
			}

			userDetailsDTO, ok := details["user"]
			if !ok {
				t.Errorf("dto.details doesn't have an expected user field")
			}

			userDetails, ok := userDetailsDTO.(map[string]any)
			if !ok {
				t.Errorf("dto.details.user should be a `map[string]any`")
			}

			if len(userDetails) != len(tc.expectedFields) {
				t.Errorf("the number of expected fields `%d` is different from what we got `%d`", len(tc.expectedFields), len(userDetails))
			}

			for _, field := range tc.expectedFields {
				val, ok := userDetails[field.Key()]
				if !ok {
					t.Errorf("expected field `%s` was not found in the returned error", field.Key())
				}

				if val != field.Value() {
					t.Errorf("expected field `%s` value `%v` is different from the returned value `%v`", field.Key(), field.Value(), val)
				}
			}
		})
	}
}
