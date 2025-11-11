package id

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ifnotnil/x/tst"
	"github.com/stretchr/testify/assert"
)

type id = [uuidSize]byte

func TestUUID_JSON(t *testing.T) {
	type Foo struct {
		ID ID[id, Base64] `json:"id"`
	}

	type FooPointer struct {
		ID *ID[id, Base64] `json:"id"`
	}

	type FooOZ struct {
		ID ID[id, Base64] `json:"id,omitzero"`
	}

	type FooPointerOZ struct {
		ID *ID[id, Base64] `json:"id,omitzero"`
	}

	const js = `{"id":"AZomiURKfF6MYQcDvjFM_A"}`
	uuid := id{0x01, 0x9a, 0x26, 0x89, 0x44, 0x4a, 0x7c, 0x5e, 0x8c, 0x61, 0x07, 0x03, 0xbe, 0x31, 0x4c, 0xfc}

	unmarshalTests := []struct {
		input         string
		destination   any
		expected      any
		errorAsserter tst.ErrorAssertionFunc
	}{
		{
			input:         js,
			destination:   &Foo{ID: ID[id, Base64]{Value: zeroUUID}},
			expected:      &Foo{ID: ID[id, Base64]{Value: uuid}},
			errorAsserter: tst.NoError(),
		},
		{
			input:         js,
			destination:   &FooPointer{ID: nil},
			expected:      &FooPointer{ID: &ID[id, Base64]{Value: uuid}},
			errorAsserter: tst.NoError(),
		},
	}

	for i, tc := range unmarshalTests {
		t.Run(fmt.Sprintf("unmarshal_%d", i), func(t *testing.T) {
			gotErr := json.Unmarshal([]byte(tc.input), tc.destination)
			tc.errorAsserter(t, gotErr)
			assert.Equal(t, tc.expected, tc.destination)
		})
	}

	marshalTests := []struct {
		input         any
		expectedJSON  string
		errorAsserter tst.ErrorAssertionFunc
	}{
		0: {
			input:         Foo{ID: ID[id, Base64]{Value: uuid}},
			expectedJSON:  js,
			errorAsserter: tst.NoError(),
		},
		1: {
			input:         &Foo{ID: ID[id, Base64]{Value: uuid}},
			expectedJSON:  js,
			errorAsserter: tst.NoError(),
		},
		2: {
			input:         FooPointer{ID: &ID[id, Base64]{Value: uuid}},
			expectedJSON:  js,
			errorAsserter: tst.NoError(),
		},
		3: {
			input:         &FooPointer{ID: &ID[id, Base64]{Value: uuid}},
			expectedJSON:  js,
			errorAsserter: tst.NoError(),
		},
		4: {
			input:         FooOZ{ID: ID[id, Base64]{Value: zeroUUID}},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
		5: {
			input:         &FooOZ{ID: ID[id, Base64]{Value: zeroUUID}},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
		6: {
			input:         FooPointerOZ{ID: &ID[id, Base64]{Value: zeroUUID}},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
		7: {
			input:         &FooPointerOZ{ID: &ID[id, Base64]{Value: zeroUUID}},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
		8: {
			input:         &FooPointerOZ{ID: nil},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
	}

	for i, tc := range marshalTests {
		t.Run(fmt.Sprintf("marshal_%d", i), func(t *testing.T) {
			got, gotErr := json.Marshal(tc.input)
			tc.errorAsserter(t, gotErr)
			assert.Equal(t, tc.expectedJSON, string(got))
		})
	}
}
