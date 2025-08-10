package tst

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const expectedError = "Expected error but none received"

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
}

type ErrorAssertionFunc func(t TestingT, err error) bool

func (e ErrorAssertionFunc) AsRequire() require.ErrorAssertionFunc {
	return func(tt require.TestingT, err error, _ ...any) {
		t, is := tt.(TestingT)
		if !is { // not possible
			tt.Errorf("Wrong TestingT type %T", tt)
			tt.FailNow()
			return
		}

		if suc := e(t, err); !suc {
			tt.FailNow()
		}
	}
}

func (e ErrorAssertionFunc) AsAssert() assert.ErrorAssertionFunc {
	return func(tt assert.TestingT, err error, _ ...any) bool {
		t, is := tt.(TestingT)
		if !is { // not possible
			tt.Errorf("Wrong TestingT type %T", tt)
			return false
		}

		return e(t, err)
	}
}

func NoError() ErrorAssertionFunc {
	return func(t TestingT, err error) bool {
		t.Helper()

		if err != nil {
			t.Errorf("Expected nil error but received : %T(%s)", err, err.Error())
			return false
		}

		return true
	}
}

func Error() ErrorAssertionFunc {
	return func(t TestingT, err error) bool {
		t.Helper()

		if err == nil {
			t.Errorf(expectedError)
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
		t.Helper()

		if err == nil {
			t.Errorf(expectedError)
			return false
		}

		if len(allExpectedErrors) == 0 {
			return true
		}

		suc := true
		notMatched := make([]error, 0, len(allExpectedErrors))
		for _, expected := range allExpectedErrors {
			if !errors.Is(err, expected) {
				notMatched = append(notMatched, expected)
				suc = false
			}
		}

		if !suc {
			sb := strings.Builder{}
			sb.WriteString("Error is unexpected.\n")
			sb.WriteString(fmt.Sprintf("Got error      : %T(%s)\n", err, err.Error()))

			if len(notMatched) == 1 {
				sb.WriteString(fmt.Sprintf("Expected error : %T(%s)\n", notMatched[0], notMatched[0].Error()))
				t.Errorf(sb.String())
				return suc
			}

			sb.WriteString("Expected errors:\n")
			for _, e := range notMatched {
				sb.WriteString(fmt.Sprintf("        -> %T(%s)\n", e, e.Error()))
			}
			t.Errorf(sb.String())
		}

		return suc
	}
}

func ErrorOfType[T error](typedAsserts ...func(TestingT, T)) ErrorAssertionFunc {
	return func(t TestingT, err error) bool {
		t.Helper()

		if err == nil {
			t.Errorf(expectedError)
			return false
		}

		var wantErr T
		if !errors.As(err, &wantErr) {
			var tErr T
			t.Errorf("Error type check failed.\nExpected error type: %T\nGot                : %T(%s)", tErr, err, err)
			return false
		}

		if v := reflect.ValueOf(wantErr); v.Kind() == reflect.Pointer && v.IsNil() {
			t.Errorf("Error check failed.\nExpected not nill error value: %T\nGot                          : %T(nil)", wantErr, wantErr)
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
		t.Helper()

		if err == nil {
			t.Errorf(expectedError)
			return false
		}

		// consider case insensitive?
		if !strings.Contains(err.Error(), s) {
			t.Errorf("Error string check failed. \nExpected to contain: %s\nGot                : %s\n", s, err.Error())
			return false
		}

		return true
	}
}

func All(expected ...ErrorAssertionFunc) ErrorAssertionFunc {
	return func(t TestingT, err error) bool {
		t.Helper()

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
