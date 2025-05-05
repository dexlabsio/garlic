//go:build unit
// +build unit

package errors

import (
	"fmt"
	"testing"
)

func TestPropagationKind(t *testing.T) {
	var (
		rootErr       = fmt.Errorf("Root Error")
		KindBaseError = &Kind{
			Name:           "BaseError",
			Code:           "E1",
			Description:    "Base error.",
			HTTPStatusCode: HTTP_STATUS_NOT_DEFINED,
			Parent:         nil,
		}

		KindOtherBaseError = &Kind{
			Name:           "OtherBaseError",
			Code:           "E2",
			Description:    "Other base error.",
			HTTPStatusCode: HTTP_STATUS_NOT_DEFINED,
			Parent:         nil,
		}
	)

	cases := []struct {
		title string
		fn    func() error
		kind  *Kind
	}{
		{
			title: "standard errors should be propagated as standard KindError by default",
			fn: func() error {
				return Propagate(rootErr, "test error")
			},
			kind: KindError,
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
