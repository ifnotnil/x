package tst

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNoError(t *testing.T) {
	tests := []errorAssertionFuncTestCase{
		{
			name:           "nil error should pass",
			input:          nil,
			asserterToTest: NoError(),
			expectedResult: true,
			initMock:       func(mt *MockTestingT) {},
		},
		{
			name:           "non-nil error should fail",
			input:          errors.New("test error"),
			asserterToTest: NoError(),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Expected nil error but received"), mock.Anything) },
		},
		{
			name:           "wrapped error should fail",
			input:          fmt.Errorf("wrapped: %w", errors.New("base error")),
			asserterToTest: NoError(),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Expected nil error but received"), mock.Anything) },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.Test)
	}
}

func TestError(t *testing.T) {
	tests := []errorAssertionFuncTestCase{
		{
			name:           "nil error should fail",
			input:          nil,
			asserterToTest: Error(),
			expectedResult: false,
			initMock: func(mt *MockTestingT) {
				mt.EXPECT().Errorf(contains("Expected error but none received"), mock.Anything)
			},
		},
		{
			name:           "non-nil error should pass",
			input:          errors.New("test error"),
			asserterToTest: Error(),
			expectedResult: true,
		},
		{
			name:           "wrapped error should pass",
			input:          fmt.Errorf("wrapped: %w", errors.New("base error")),
			asserterToTest: Error(),
			expectedResult: true,
		},
		{
			name:           "*os.PathError should pass",
			input:          &os.PathError{Op: "open", Path: "/test", Err: errors.New("base")},
			asserterToTest: Error(),
			expectedResult: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.Test)
	}
}

func TestErrorIs(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)
	anotherErr := errors.New("another error")

	tests := []errorAssertionFuncTestCase{
		{
			name:           "nil error should fail", // works like IsError
			input:          nil,
			asserterToTest: ErrorIs(baseErr),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Expected error but none received")).Once() },
		},
		{
			name:           "nil error without arguments should fail", // works like IsError
			input:          nil,
			asserterToTest: ErrorIs(),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Expected error but none received")).Once() },
		},
		{
			name:           "empty expected errors list", // works like IsError
			input:          baseErr,
			asserterToTest: ErrorIs(),
			expectedResult: true,
			initMock:       func(mt *MockTestingT) {},
		},
		{
			name:           "single matching error",
			input:          baseErr,
			asserterToTest: ErrorIs(baseErr),
			expectedResult: true,
			initMock:       func(mt *MockTestingT) {},
		},
		{
			name:           "single wrapped matching error",
			input:          wrappedErr,
			asserterToTest: ErrorIs(baseErr),
			expectedResult: true,
			initMock:       func(mt *MockTestingT) {},
		},
		{
			name:           "single non-matching error",
			input:          baseErr,
			asserterToTest: ErrorIs(anotherErr),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Error is unexpected")).Once() },
		},
		{
			name:           "multiple matching errors",
			input:          wrappedErr,
			asserterToTest: ErrorIs(baseErr, wrappedErr),
			expectedResult: true,
			initMock:       func(mt *MockTestingT) {},
		},
		{
			name:           "multiple errors - some match",
			input:          wrappedErr,
			asserterToTest: ErrorIs(anotherErr, wrappedErr),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Error is unexpected")).Once() },
		},
		{
			name:           "multiple errors - none match",
			input:          baseErr,
			asserterToTest: ErrorIs(anotherErr, errors.New("yet another")),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Error is unexpected")).Once() },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.Test)
	}
}

func TestErrorOfType(t *testing.T) {
	baseErr := errors.New("base error")
	pathErr := &os.PathError{Op: "open", Path: "/test", Err: baseErr}
	wrappedPathErr := fmt.Errorf("wrapped: %w", pathErr)

	tests := []errorAssertionFuncTestCase{
		{
			name:           "nil error should fail",
			input:          nil,
			asserterToTest: ErrorOfType[*os.PathError](),
			expectedResult: false,
			initMock: func(mt *MockTestingT) {
				mt.EXPECT().Errorf(contains("Expected error but none received")).Once()
			},
		},
		{
			name:           "nil *os.PathError should fail",
			input:          (*os.PathError)(nil),
			asserterToTest: ErrorOfType[*os.PathError](),
			expectedResult: false,
			initMock: func(mt *MockTestingT) {
				mt.EXPECT().Errorf(contains("Expected not nill error value"), mock.Anything).Once()
			},
		},
		{
			name:           "matching error type with no assertions",
			input:          pathErr,
			asserterToTest: ErrorOfType[*os.PathError](),
			expectedResult: true,
			initMock:       func(_ *MockTestingT) {},
		},
		{
			name:           "wrapped matching error type with no assertions",
			input:          wrappedPathErr,
			asserterToTest: ErrorOfType[*os.PathError](),
			expectedResult: true,
			initMock:       func(_ *MockTestingT) {},
		},
		{
			name:           "non-matching error type",
			input:          baseErr,
			asserterToTest: ErrorOfType[*os.PathError](),
			expectedResult: false,
			initMock: func(mt *MockTestingT) {
				mt.EXPECT().Errorf(contains("Error type check failed."), mock.Anything).Once()
			},
		},
		{
			name:  "matching error type with passing assertion",
			input: pathErr,
			asserterToTest: ErrorOfType[*os.PathError](
				func(t TestingT, pe *os.PathError) {
					assert.Equal(t, "open", pe.Op)
				},
			),
			expectedResult: true,
			initMock:       func(mt *MockTestingT) {},
		},
		{
			name:  "matching error type with multiple passing assertions",
			input: pathErr,
			asserterToTest: ErrorOfType[*os.PathError](
				func(t TestingT, pe *os.PathError) {
					assert.Equal(t, "open", pe.Op)
				},
				func(t TestingT, pe *os.PathError) {
					assert.Equal(t, "/test", pe.Path)
				},
			),
			expectedResult: true,
			initMock:       func(mt *MockTestingT) {},
		},
		{
			name:  "wrapped matching error type with passing assertion",
			input: wrappedPathErr,
			asserterToTest: ErrorOfType[*os.PathError](
				func(t TestingT, pe *os.PathError) {
					assert.Equal(t, "open", pe.Op)
				},
			),
			expectedResult: true,
			initMock:       func(mt *MockTestingT) {},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.Test)
	}

	t.Run("ensure sub-asserts calls", func(t *testing.T) {
		baseErr := errors.New("base error")
		pathErr := &os.PathError{Op: "open", Path: "/test", Err: baseErr}

		mt := NewMockTestingT(t)
		mt.EXPECT().Helper().Maybe()
		ma := &mockErrorTypedAssertionFunc{}
		ma.OnAssert(mt, pathErr).Return(true).Times(3)

		got := ErrorOfType[*os.PathError](ma.Assert, ma.Assert, ma.Assert)(mt, pathErr)

		assert.True(t, got)
	})
}

func TestErrorStringContains(t *testing.T) {
	testErr := errors.New("this is a test error message")

	tests := []errorAssertionFuncTestCase{
		{
			name:           "nil error should fail",
			input:          nil,
			asserterToTest: ErrorStringContains("test"),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Expected error but none received")) },
		},
		{
			name:           "exact match",
			input:          testErr,
			asserterToTest: ErrorStringContains("this is a test error message"),
			expectedResult: true,
		},
		{
			name:           "partial match",
			input:          testErr,
			asserterToTest: ErrorStringContains("test error"),
			expectedResult: true,
		},
		{
			name:           "empty string should match any error",
			input:          testErr,
			asserterToTest: ErrorStringContains(""),
			expectedResult: true,
		},
		{
			name:           "no match - different content",
			input:          testErr,
			asserterToTest: ErrorStringContains("not found"),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Error string check failed."), mock.Anything) },
		},
		{
			name:           "no match - case sensitive",
			input:          testErr,
			asserterToTest: ErrorStringContains("TEST ERROR"),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Error string check failed."), mock.Anything) },
		},
		{
			name:           "no match - extra characters",
			input:          testErr,
			asserterToTest: ErrorStringContains("this is a test error message!"),
			expectedResult: false,
			initMock:       func(mt *MockTestingT) { mt.EXPECT().Errorf(contains("Error string check failed."), mock.Anything) },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.Test)
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

type errorAssertionFuncTestCase struct {
	name           string
	input          error
	asserterToTest ErrorAssertionFunc
	expectedResult bool
	initMock       func(mt *MockTestingT)
}

func (tc errorAssertionFuncTestCase) Test(t *testing.T) {
	mt := NewMockTestingT(t)
	mt.EXPECT().Helper().Maybe()
	if tc.initMock != nil {
		tc.initMock(mt)
	}
	result := tc.asserterToTest(mt, tc.input)
	assert.Equal(t, tc.expectedResult, result, "error asserter returned bool mismatch")
}

func contains(s string) any {
	return mock.MatchedBy(func(format string) bool {
		return strings.Contains(format, s)
	})
}

//nolint:thelper
func TestTestifyIntegration(t *testing.T) {
	tests := []struct {
		name               string
		mockAsserterReturn bool
		mockTInit          func(*MockTestingT)
		run                func(t *testing.T, mt *MockTestingT, e ErrorAssertionFunc)
	}{
		{
			name:               "AsRequire pass",
			mockAsserterReturn: true,
			mockTInit:          func(mt *MockTestingT) {},
			run:                func(t *testing.T, mt *MockTestingT, e ErrorAssertionFunc) { e.AsRequire()(mt, nil) },
		},
		{
			name:               "AsRequire fail",
			mockAsserterReturn: false,
			mockTInit:          func(mt *MockTestingT) { mt.EXPECT().FailNow().Once() },
			run:                func(t *testing.T, mt *MockTestingT, e ErrorAssertionFunc) { e.AsRequire()(mt, nil) },
		},
		{
			name:               "AsAssert pass",
			mockAsserterReturn: true,
			mockTInit:          func(mt *MockTestingT) {},
			run: func(t *testing.T, mt *MockTestingT, e ErrorAssertionFunc) {
				got := e.AsAssert()(mt, nil)
				assert.True(t, got)
			},
		},
		{
			name:               "AsAssert fail",
			mockAsserterReturn: false,
			mockTInit:          func(mt *MockTestingT) {},
			run: func(t *testing.T, mt *MockTestingT, e ErrorAssertionFunc) {
				got := e.AsAssert()(mt, nil)
				assert.False(t, got)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// expect asserter to be called
			m := &mock.Mock{}
			t.Cleanup(func() { m.AssertExpectations(t) })
			m.On("Assert", mock.Anything, mock.Anything).Return(tc.mockAsserterReturn).Once()
			var f ErrorAssertionFunc = func(t TestingT, err error) bool {
				return m.MethodCalled("Assert", t, err).Bool(0)
			}

			// mock T init
			mt := NewMockTestingT(t)
			mt.EXPECT().Helper().Maybe()
			tc.mockTInit(mt)

			tc.run(t, mt, f)
		})
	}

	t.Run("mismatch of T", func(t *testing.T) {
		t.Run("require testingT", func(t *testing.T) {
			mt := &testifyRequireTestingTMock{}
			mt.On("Errorf", "Wrong TestingT type %T", mock.Anything).Once()
			mt.On("FailNow").Once()

			NoError().AsRequire()(mt, nil)
		})
		t.Run("assert testingT", func(t *testing.T) {
			mt := &testifyAssertTestingTMock{}
			mt.On("Errorf", "Wrong TestingT type %T", mock.Anything).Once()

			got := NoError().AsAssert()(mt, nil)
			assert.False(t, got)
		})
	})
}

func TestAll(t *testing.T) {
	tests := []struct {
		name                string
		assertionFuncs      []ErrorAssertionFunc
		expectResult        bool
		expectedErrorfCalls int
	}{
		{
			name:           "no assertion functions",
			assertionFuncs: []ErrorAssertionFunc{},
			expectResult:   true,
		},
		{
			name: "all assertions pass",
			assertionFuncs: []ErrorAssertionFunc{
				func(t TestingT, err error) bool { return true },
				func(t TestingT, err error) bool { return true },
			},
			expectResult: true,
		},
		{
			name: "one assertion fails",
			assertionFuncs: []ErrorAssertionFunc{
				func(t TestingT, err error) bool { return true },
				func(t TestingT, err error) bool { return false },
			},
			expectResult: false,
		},
		{
			name: "all assertions fail",
			assertionFuncs: []ErrorAssertionFunc{
				func(t TestingT, err error) bool { return false },
				func(t TestingT, err error) bool { return false },
			},
			expectResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := NewMockTestingT(t)
			mt.EXPECT().Helper().Maybe()
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

func TestFail(t *testing.T) {
	t.SkipNow()

	err := errors.New("error one")

	tests := []struct {
		name     string
		input    error
		asserter ErrorAssertionFunc
	}{
		{
			name:     "NoError",
			input:    err,
			asserter: NoError(),
		},
		{
			name:     "Error",
			input:    nil,
			asserter: Error(),
		},
		{
			name:     "ErrorIs one",
			input:    err,
			asserter: ErrorIs(io.ErrClosedPipe),
		},
		{
			name:     "ErrorIs many",
			input:    err,
			asserter: ErrorIs(io.ErrClosedPipe, io.ErrNoProgress),
		},
		{
			name:     "ErrorOfType",
			input:    err,
			asserter: ErrorOfType[*os.PathError](),
		},
		{
			name:     "ErrorOfType nil",
			input:    (*os.PathError)(nil),
			asserter: ErrorOfType[*os.PathError](),
		},
		{
			name:     "ErrorStringContains",
			input:    err,
			asserter: ErrorStringContains("just a string"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.asserter(t, tc.input)
		})
	}
}

var _ (assert.TestingT) = (*testifyAssertTestingTMock)(nil)

type testifyAssertTestingTMock struct {
	mock.Mock
}

func (m *testifyAssertTestingTMock) Errorf(format string, args ...any) {
	if len(args) > 0 {
		m.Called(format, args)
	} else {
		m.Called(format)
	}
}

var _ (require.TestingT) = (*testifyRequireTestingTMock)(nil)

type testifyRequireTestingTMock struct {
	mock.Mock
}

func (m *testifyRequireTestingTMock) Errorf(format string, args ...any) {
	if len(args) > 0 {
		m.Called(format, args)
	} else {
		m.Called(format)
	}
}

func (m *testifyRequireTestingTMock) FailNow() {
	m.Called()
}
