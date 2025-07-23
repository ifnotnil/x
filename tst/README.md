# tst
[![ci](https://github.com/ifnotnil/x/actions/workflows/sub_tst.yml/badge.svg)](https://github.com/ifnotnil/x/actions/workflows/sub_tst.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ifnotnil/x/tst)](https://goreportcard.com/report/github.com/ifnotnil/x/tst)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ifnotnit/x/tst)](https://pkg.go.dev/github.com/ifnotnil/x/tst)
[![Version](https://img.shields.io/github/v/tag/ifnotnil/x?filter=tst%2F*)](https://pkg.go.dev/github.com/ifnotnil/x/tst?tab=versions)
[![codecov](https://codecov.io/gh/ifnotnil/x/graph/badge.svg?token=n0t9q5Y3Sf&component=tst)](https://codecov.io/gh/ifnotnil/x)


Error asserting functions for table driven testcases.

Examples:
```golang
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
```
