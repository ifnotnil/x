package tst

import (
	"errors"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

type ErrorAssertionFunc func(t TestingT, err error) bool

func (e ErrorAssertionFunc) AsRequire() require.ErrorAssertionFunc {
	return func(tt require.TestingT, err error, _ ...any) {
		if suc := e(tt, err); !suc {
			tt.FailNow()
		}
	}
}

func (e ErrorAssertionFunc) AsAssert() assert.ErrorAssertionFunc {
	return func(tt assert.TestingT, err error, _ ...any) bool {
		t, is := tt.(TestingT)
		if is {
			return e(t, err)
		}

		// not possible
		tt.Errorf("wrong TestingT type %T", tt)
		return false
	}
}

func All(expected ...ErrorAssertionFunc) ErrorAssertionFunc {
	return func(t TestingT, err error) bool {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}

		ret := true
		for _, fn := range expected {
			ok := fn(t, err)
			if !ok {
				ret = ok
			}
		}

		return ret
	}
}

func NoError() ErrorAssertionFunc {
	return func(t TestingT, err error) bool {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}

		if err != nil {
			t.Errorf("expected nil error but received : %T(%s)", err, err.Error())
			return false
		}

		return true
	}
}

func Error() ErrorAssertionFunc {
	return func(t TestingT, err error) bool {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}

		if err == nil {
			t.Errorf("expected error but none received")
			return false
		}

		return true
	}
}

// ErrorIs returns an ErrorAssertionFunc that checks if the given error matches
// any of the expected errors using errors.Is.
// If no expected errors are provided, it simply checks that an error is present (similar to Error()).
// Returns false if the error is nil or doesn't match any expected errors.
func ErrorIs(allExpectedErrors ...error) ErrorAssertionFunc {
	return func(t TestingT, err error) bool {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}
		if err == nil {
			t.Errorf("expected error but none received")
			return false
		}

		ret := true
		for _, expected := range allExpectedErrors {
			if !errors.Is(err, expected) {
				t.Errorf("error unexpected.\nExpected error: %T(%s) \nGot           : %T(%s)", expected, expected.Error(), err, err.Error())
				ret = false
			}
		}

		return ret
	}
}

func ErrorOfType[T error](typedAsserts ...func(TestingT, T)) ErrorAssertionFunc {
	return func(t TestingT, err error) bool {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}

		if err == nil {
			t.Errorf("expected error but none received")
			return false
		}

		var wantErr T
		if !errors.As(err, &wantErr) {
			var tErr T
			t.Errorf("Error type check failed.\nExpected error type: %T\nGot                : %T(%s)", tErr, err, err)
			return false
		}

		for _, e := range typedAsserts {
			e(t, wantErr)
		}

		return true
	}
}

func ErrorStringContains(s string) ErrorAssertionFunc {
	return func(t TestingT, err error) bool {
		if h, ok := t.(interface{ Helper() }); ok {
			h.Helper()
		}

		if err == nil {
			t.Errorf("expected error but none received")
			return false
		}

		if !strings.Contains(err.Error(), s) {
			t.Errorf("error string check failed. \nExpected to contain: %s\nGot                : %s\n", s, err.Error())
			return false
		}

		return true
	}
}
