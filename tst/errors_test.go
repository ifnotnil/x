package tst

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	tests := []struct {
		name                string
		assertionFuncs      []ErrorAssertionFunc
		expectResult        bool
		expectedErrorfCalls int
	}{
		{
			name:                "no assertion functions",
			assertionFuncs:      []ErrorAssertionFunc{},
			expectResult:        true,
			expectedErrorfCalls: 0,
		},
		{
			name: "all assertions pass",
			assertionFuncs: []ErrorAssertionFunc{
				func(t TestingT, err error) bool { return true },
				func(t TestingT, err error) bool { return true },
			},
			expectResult:        true,
			expectedErrorfCalls: 0,
		},
		{
			name: "one assertion fails",
			assertionFuncs: []ErrorAssertionFunc{
				func(t TestingT, err error) bool { return true },
				func(t TestingT, err error) bool { return false },
			},
			expectResult:        false,
			expectedErrorfCalls: 1,
		},
		{
			name: "all assertions fail",
			assertionFuncs: []ErrorAssertionFunc{
				func(t TestingT, err error) bool { return false },
				func(t TestingT, err error) bool { return false },
			},
			expectResult:        false,
			expectedErrorfCalls: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := NewMockTestingT(t)

			result := All(tt.assertionFuncs...)(mt, errors.New("test error"))

			assert.Equal(t, tt.expectResult, result)
		})
	}
}

func Test_Readme(t *testing.T) {
	err := errors.New("not found")

	tests := []struct {
		input    error
		asserter ErrorAssertionFunc
	}{
		{
			input:    nil,
			asserter: NoError(),
		},
		{
			input:    err,
			asserter: Error(),
		},
		{
			input:    err,
			asserter: ErrorIs(err),
		},
		{
			input:    fmt.Errorf("wrapped %w", err),
			asserter: ErrorIs(err),
		},
		{
			input:    &os.PathError{},
			asserter: ErrorOfType[*os.PathError](),
		},
		{
			input: &os.PathError{Op: "op", Path: "/abc", Err: os.ErrInvalid},
			asserter: ErrorOfType[*os.PathError](
				func(tt TestingT, pe *os.PathError) { assert.Equal(tt, "op", pe.Op) },
			),
		},
		{
			input: fmt.Errorf("wrapped %w", &os.PathError{Op: "op", Path: "/abc", Err: os.ErrInvalid}),
			asserter: All(
				Error(),
				ErrorIs(os.ErrInvalid),
				ErrorOfType[*os.PathError](
					func(tt TestingT, pe *os.PathError) { assert.Equal(tt, "op", pe.Op) },
				),
				ErrorStringContains("op"),
			),
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc.asserter(t, tc.input)
		})
	}
}
