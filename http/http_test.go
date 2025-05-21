package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/ifnotnil/x/http/internal/testingx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDrain(t *testing.T) {
	// ensure no panic
	t.Run("ensure no panic on nils", func(t *testing.T) {
		require.NoError(t, DrainAndCloseRequest(nil))
		require.NoError(t, DrainAndCloseRequest(&http.Request{Body: nil}))
		require.NoError(t, DrainAndCloseRequest(&http.Request{Body: http.NoBody}))

		require.NoError(t, DrainAndCloseResponse(nil))
		require.NoError(t, DrainAndCloseResponse(&http.Response{Body: nil}))
		require.NoError(t, DrainAndCloseResponse(&http.Response{Body: http.NoBody}))
	})

	tests := map[string]struct {
		bodyFN        func(*testing.T) io.ReadCloser
		errorAsserter testingx.ErrorAssertion
	}{
		"happy path": {
			bodyFN: func(t *testing.T) io.ReadCloser {
				t.Helper()
				return newIOReadCloserMock(t, nil)
			},
			errorAsserter: testingx.NoError,
		},
		"hapy path with buffer": {
			bodyFN: func(t *testing.T) io.ReadCloser {
				t.Helper()
				m := newIOReadCloserMock(t, nil)
				m.buf.WriteString("abcd1234")
				return m
			},
			errorAsserter: testingx.NoError,
		},
		"error on close": {
			bodyFN: func(t *testing.T) io.ReadCloser {
				t.Helper()
				m := newIOReadCloserMock(t, errors.New("close error"))
				m.buf.WriteString("abcd1234")
				return m
			},
			errorAsserter: func(tt testingx.TestingT, err error) bool { return assert.Equal(tt, "close error", err.Error()) },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b := tc.bodyFN(t)
			err := DrainAndCloseBody(b)
			tc.errorAsserter(t, err)
		})
	}
}

type ioReadCloserMock struct {
	mock.Mock
	buf bytes.Buffer
}

func (rc *ioReadCloserMock) Close() error {
	return rc.Called().Error(0)
}

func (rc *ioReadCloserMock) Read(p []byte) (n int, err error) {
	return rc.buf.Read(p)
}

func newIOReadCloserMock(t *testing.T, closeErr error) *ioReadCloserMock {
	t.Helper()
	m := &ioReadCloserMock{}
	t.Cleanup(func() { m.AssertExpectations(t) })
	m.On("Close").Return(closeErr).Once()

	return m
}

var (
	errMock1 = errors.New("error1")
	errMock2 = errors.New("error2")
)
