//go:build unit
// +build unit

package errors

import (
	"fmt"
	"testing"
)

func TestAsKind(t *testing.T) {
	var (
		KindTestError = &Kind{
			Name:           "TestError",
			Code:           "E1",
			Description:    "Test error",
			HTTPStatusCode: HTTP_STATUS_NOT_DEFINED,
			Parent:         nil,
		}

		KindChildTestError = &Kind{
			Name:           "ChildTestError",
			Code:           "E2",
			Description:    "Child test error",
			HTTPStatusCode: HTTP_STATUS_NOT_DEFINED,
			Parent:         KindTestError,
		}

		KindOtherTestError = &Kind{
			Name:           "OtherTestError",
			Code:           "E3",
			Description:    "Other test error",
			HTTPStatusCode: HTTP_STATUS_NOT_DEFINED,
			Parent:         nil,
		}
	)

	cases := []struct {
		title string
		fn    func() error
		kind  *Kind
		found bool
	}{
		{
			title: "standard errors should not translate",
			fn:    func() error { return fmt.Errorf("test") },
			kind:  KindTestError,
			found: false,
		},
		{
			title: "first level errors should translate",
			fn:    func() error { return New(KindTestError, "test") },
			kind:  KindTestError,
			found: true,
		},
		{
			title: "first level errors should not translate to other kind",
			fn:    func() error { return New(KindTestError, "test") },
			kind:  KindOtherTestError,
			found: false,
		},
		{
			title: "second level errors should translate to parent kind",
			fn:    func() error { return New(KindChildTestError, "test child") },
			kind:  KindTestError,
			found: true,
		},
		{
			title: "second level errors should translate to self kind",
			fn:    func() error { return New(KindChildTestError, "test child") },
			kind:  KindChildTestError,
			found: true,
		},
		{
			title: "second level errors should not translate to other kind",
			fn:    func() error { return New(KindChildTestError, "test child") },
			kind:  KindOtherTestError,
			found: false,
		},
		{
			title: "first level propagated errors should translate to parent kind",
			fn: func() error {
				root := New(KindChildTestError, "test child")
				return Propagate(root, "propagated test child")
			},
			kind:  KindTestError,
			found: true,
		},
		{
			title: "first level propagated errors should translate to self kind",
			fn: func() error {
				root := New(KindChildTestError, "test child")
				return Propagate(root, "propagated test child")
			},
			kind:  KindChildTestError,
			found: true,
		},
		{
			title: "first level propagated errors should not translate to other kind",
			fn: func() error {
				root := New(KindChildTestError, "test child")
				return Propagate(root, "propagated test child")
			},
			kind:  KindOtherTestError,
			found: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			err := tc.fn()
			if _, ok := AsKind(err, tc.kind); ok != tc.found {
				t.Errorf("expected AsKind to return %v, but got %v", tc.found, ok)
			}
		})
	}
}
