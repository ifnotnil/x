package testingx

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AsRequireErrorAssertionFunc(ea ErrorAssertion) require.ErrorAssertionFunc {
	return func(t require.TestingT, err error, _ ...any) {
		if h, ok := t.(tHelper); ok {
			h.Helper()
		}
		if ok := ea.Assert(t, err); !ok {
			if f, ok := t.(tFailNow); ok {
				f.FailNow()
			}
		}
	}
}

func AsAssertErrorAssertionFunc(ea ErrorAssertion) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, _ ...any) bool {
		if h, ok := t.(tHelper); ok {
			h.Helper()
		}
		return ea.Assert(t, err)
	}
}

type TestingT interface {
	Errorf(format string, args ...any)
}

type tFailNow interface {
	FailNow()
}

type tHelper interface {
	Helper()
}

type ErrorAssertion func(TestingT, error) bool

func (e ErrorAssertion) Assert(t TestingT, err error) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	switch {
	case err != nil && e != nil:
		return e(t, err)
	case err != nil && e == nil:
		t.Errorf("unexpected error returned.\nError: %T(%s)", err, err.Error())
		return false
	case err == nil && e != nil:
		t.Errorf("expected error but none received")
		return false
	}

	return true
}

func NoError(t TestingT, err error) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if err != nil {
		t.Errorf("unexpected error returned.\nError: %T(%s)", err, err.Error())
		return false
	}

	return true
}

func ExpectedErrorChecks(expected ...ErrorAssertion) ErrorAssertion {
	return func(t TestingT, err error) bool {
		if h, ok := t.(tHelper); ok {
			h.Helper()
		}

		ok := true
		for _, fn := range expected {
			if fnOk := fn(t, err); !fnOk {
				ok = false
			}
		}

		return ok
	}
}

func ErrorIs(allExpectedErrors ...error) ErrorAssertion {
	return func(t TestingT, err error) bool {
		if h, ok := t.(tHelper); ok {
			h.Helper()
		}
		if err == nil {
			t.Errorf("expected error but none received")
			return false
		}

		ok := true
		for _, expected := range allExpectedErrors {
			if !errors.Is(err, expected) {
				t.Errorf("error unexpected.\nExpected error: %T(%s) \nGot           : %T(%s)", expected, expected.Error(), err, err.Error())
				ok = false
			}
		}

		return ok
	}
}

func ErrorOfType[T error](assertsOfType ...func(T) bool) ErrorAssertion {
	return func(t TestingT, err error) bool {
		if h, ok := t.(tHelper); ok {
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

		ok := true
		for _, e := range assertsOfType {
			if a := e(wantErr); !a {
				ok = false
			}
		}

		return ok
	}
}

func ErrorStringContains(s string) ErrorAssertion {
	return func(t TestingT, err error) bool {
		if h, ok := t.(tHelper); ok {
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

func ReadAndClose(t *testing.T, body io.ReadCloser) []byte {
	t.Helper()
	defer func() {
		t.Helper()
		if body == nil {
			return
		}
		err := body.Close()
		if err != nil {
			t.Errorf("error while closing body: %s", err.Error())
		}
	}()

	b, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("error while reading body: %s", err.Error())
	}

	return b
}
