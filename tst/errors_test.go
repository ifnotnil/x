package tst

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestErrorIs(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)
	anotherErr := errors.New("another error")

	tests := []struct {
		name              string
		err               error
		errorIsArgs       []error
		expectedResult    bool
		expectedErrorfMsg string
	}{
		{
			name:              "nil error should fail",
			err:               nil,
			errorIsArgs:       []error{baseErr},
			expectedResult:    false,
			expectedErrorfMsg: "expected error but none received",
		},
		{
			name:              "nil error without arguments should fail",
			err:               nil,
			errorIsArgs:       []error{},
			expectedResult:    false,
			expectedErrorfMsg: "expected error but none received",
		},
		{
			name:           "empty expected errors list", // works like IsError
			err:            baseErr,
			errorIsArgs:    []error{},
			expectedResult: true,
		},
		{
			name:           "single matching error",
			err:            baseErr,
			errorIsArgs:    []error{baseErr},
			expectedResult: true,
		},
		{
			name:           "single wrapped matching error",
			err:            wrappedErr,
			errorIsArgs:    []error{baseErr},
			expectedResult: true,
		},
		{
			name:              "single non-matching error",
			err:               baseErr,
			errorIsArgs:       []error{anotherErr},
			expectedResult:    false,
			expectedErrorfMsg: "error unexpected",
		},
		{
			name:           "multiple matching errors",
			err:            wrappedErr,
			errorIsArgs:    []error{baseErr, wrappedErr},
			expectedResult: true,
		},
		{
			name:              "multiple errors - some match, some don't",
			err:               wrappedErr,
			errorIsArgs:       []error{baseErr, anotherErr},
			expectedResult:    false,
			expectedErrorfMsg: "error unexpected",
		},
		{
			name:              "multiple errors - none match",
			err:               baseErr,
			errorIsArgs:       []error{anotherErr, errors.New("yet another")},
			expectedResult:    false,
			expectedErrorfMsg: "error unexpected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := NewMockTestingT(t)

			if !tt.expectedResult {
				mt.EXPECT().Errorf(
					mock.MatchedBy(
						func(format string) bool {
							return strings.Contains(format, tt.expectedErrorfMsg)
						},
					),
					mock.Anything,
				).Return()
			}

			result := ErrorIs(tt.errorIsArgs...)(mt, tt.err)

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestErrorOfType(t *testing.T) {
	baseErr := errors.New("base error")
	pathErr := &os.PathError{Op: "open", Path: "/test", Err: baseErr}
	wrappedPathErr := fmt.Errorf("wrapped: %w", pathErr)

	tests := []struct {
		name           string
		err            error
		typedAsserts   []func(TestingT, *os.PathError)
		mockInit       func(*MockTestingT)
		expectedResult bool
	}{
		{
			name:         "nil error should fail",
			err:          nil,
			typedAsserts: []func(TestingT, *os.PathError){},
			mockInit: func(mt *MockTestingT) {
				mt.EXPECT().Errorf(mock.MatchedBy(func(fmtMsg string) bool {
					return fmtMsg == "expected error but none received"
				})).Once()
			},
			expectedResult: false,
		},
		{
			name:         "nil *os.PathError should fail",
			err:          (*os.PathError)(nil),
			typedAsserts: []func(TestingT, *os.PathError){},
			mockInit: func(mt *MockTestingT) {
				mt.EXPECT().Errorf(mock.MatchedBy(func(fmtMsg string) bool {
					return strings.Contains(fmtMsg, "Expected not nill error type")
				}), mock.Anything).Once()
			},
			expectedResult: false,
		},
		{
			name:           "matching error type with no assertions",
			err:            pathErr,
			typedAsserts:   []func(TestingT, *os.PathError){},
			mockInit:       func(_ *MockTestingT) {},
			expectedResult: true,
		},
		{
			name:           "wrapped matching error type with no assertions",
			err:            wrappedPathErr,
			typedAsserts:   []func(TestingT, *os.PathError){},
			expectedResult: true,
		},
		{
			name:         "non-matching error type",
			err:          baseErr,
			typedAsserts: []func(TestingT, *os.PathError){},
			mockInit: func(mt *MockTestingT) {
				mt.EXPECT().Errorf(mock.MatchedBy(func(fmtMsg string) bool {
					return strings.Contains(fmtMsg, "Error type check failed.")
				}), mock.Anything, mock.Anything).Once()
			},
			expectedResult: false,
		},
		{
			name: "matching error type with passing assertion",
			err:  pathErr,
			typedAsserts: []func(TestingT, *os.PathError){
				func(t TestingT, pe *os.PathError) {
					assert.Equal(t, "open", pe.Op)
				},
			},
			mockInit: func(mt *MockTestingT) {
			},
			expectedResult: true,
		},
		{
			name: "matching error type with multiple passing assertions",
			err:  pathErr,
			typedAsserts: []func(TestingT, *os.PathError){
				func(t TestingT, pe *os.PathError) {
					assert.Equal(t, "open", pe.Op)
				},
				func(t TestingT, pe *os.PathError) {
					assert.Equal(t, "/test", pe.Path)
				},
			},
			expectedResult: true,
		},
		{
			name: "wrapped matching error type with passing assertion",
			err:  wrappedPathErr,
			typedAsserts: []func(TestingT, *os.PathError){
				func(t TestingT, pe *os.PathError) {
					assert.Equal(t, "open", pe.Op)
				},
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := NewMockTestingT(t)

			if tt.mockInit != nil {
				tt.mockInit(mt)
			}

			result := ErrorOfType[*os.PathError](tt.typedAsserts...)(mt, tt.err)

			assert.Equal(t, tt.expectedResult, result)
		})
	}

	t.Run("ensure sub-asserts calls", func(t *testing.T) {
		baseErr := errors.New("base error")
		pathErr := &os.PathError{Op: "open", Path: "/test", Err: baseErr}

		mt := NewMockTestingT(t)
		ma := &mockErrorTypedAssertionFunc{}
		ma.OnAssert(mt, pathErr).Return(true).Times(3)

		got := ErrorOfType[*os.PathError](ma.Assert, ma.Assert, ma.Assert)(mt, pathErr)

		assert.True(t, got)
	})
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

type mockErrorTypedAssertionFunc struct {
	mock.Mock
}

func (m *mockErrorTypedAssertionFunc) OnAssert(t any, err any) *mock.Call {
	return m.On("Assert", t, err)
}

func (m *mockErrorTypedAssertionFunc) Assert(t TestingT, err *os.PathError) {
	m.Called(t, err)
}
