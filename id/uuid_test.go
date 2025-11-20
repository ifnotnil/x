package id

import (
	"encoding/json"
	"testing"

	"github.com/ifnotnil/x/tst"
	"github.com/stretchr/testify/assert"
)

type id = [uuidSize]byte

type unmarshalTest struct {
	name          string
	input         string
	destination   any
	expected      any
	errorAsserter tst.ErrorAssertionFunc
}

func (tc unmarshalTest) Test(t *testing.T) {
	gotErr := json.Unmarshal([]byte(tc.input), tc.destination)
	tc.errorAsserter(t, gotErr)
	assert.Equal(t, tc.expected, tc.destination)
}

func TestJSONUnmarshal(t *testing.T) {
	type Foo[Enc encoding] struct {
		ID ID[id, Enc] `json:"id"`
	}

	type FooPointer[Enc encoding] struct {
		ID *ID[id, Enc] `json:"id"`
	}

	uuid := id{0x01, 0x9a, 0x26, 0x89, 0x44, 0x4a, 0x7c, 0x5e, 0x8c, 0x61, 0x07, 0x03, 0xbe, 0x31, 0x4c, 0xfc}

	unmarshalTests := []unmarshalTest{
		{
			name:          "Foo[Base64]",
			input:         `{"id":"AZomiURKfF6MYQcDvjFM_A"}`,
			destination:   &Foo[Base64]{ID: ID[id, Base64]{Value: zeroUUID}},
			expected:      &Foo[Base64]{ID: ID[id, Base64]{Value: uuid}},
			errorAsserter: tst.NoError(),
		},
		{
			name:          "FooPointer[Base64]",
			input:         `{"id":"AZomiURKfF6MYQcDvjFM_A"}`,
			destination:   &FooPointer[Base64]{ID: nil},
			expected:      &FooPointer[Base64]{ID: &ID[id, Base64]{Value: uuid}},
			errorAsserter: tst.NoError(),
		},
		{
			name:          "Foo[Base64WithPadding]",
			input:         `{"id":"AZomiURKfF6MYQcDvjFM_A~~"}`,
			destination:   &Foo[Base64WithPadding]{ID: ID[id, Base64WithPadding]{Value: zeroUUID}},
			expected:      &Foo[Base64WithPadding]{ID: ID[id, Base64WithPadding]{Value: uuid}},
			errorAsserter: tst.NoError(),
		},
		{
			name:          "FooPointer[Base64WithPadding]",
			input:         `{"id":"AZomiURKfF6MYQcDvjFM_A~~"}`,
			destination:   &FooPointer[Base64WithPadding]{ID: nil},
			expected:      &FooPointer[Base64WithPadding]{ID: &ID[id, Base64WithPadding]{Value: uuid}},
			errorAsserter: tst.NoError(),
		},
	}

	for _, tc := range unmarshalTests {
		t.Run(tc.name, tc.Test)
	}
}

type marshalTestCase struct {
	name          string
	input         any
	expectedJSON  string
	errorAsserter tst.ErrorAssertionFunc
}

func (tc marshalTestCase) Test(t *testing.T) {
	got, gotErr := json.Marshal(tc.input)
	tc.errorAsserter(t, gotErr)
	assert.Equal(t, tc.expectedJSON, string(got))
}

func TestJSONMarshal(t *testing.T) {
	type Foo[enc encoding] struct {
		ID ID[id, enc] `json:"id"`
	}

	type FooPointer[enc encoding] struct {
		ID *ID[id, enc] `json:"id"`
	}

	type FooOZ[enc encoding] struct {
		ID ID[id, enc] `json:"id,omitzero"`
	}

	type FooPointerOZ[enc encoding] struct {
		ID *ID[id, enc] `json:"id,omitzero"`
	}

	const jsBase64 = `{"id":"AZomiURKfF6MYQcDvjFM_A"}`
	uuid := id{0x01, 0x9a, 0x26, 0x89, 0x44, 0x4a, 0x7c, 0x5e, 0x8c, 0x61, 0x07, 0x03, 0xbe, 0x31, 0x4c, 0xfc}

	marshalTests := []marshalTestCase{
		{
			name:          "FooOZ[Base64] zero",
			input:         FooOZ[Base64]{ID: ID[id, Base64]{Value: zeroUUID}},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
		{
			name:          "*FooOZ[Base64] zero",
			input:         &FooOZ[Base64]{ID: ID[id, Base64]{Value: zeroUUID}},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
		{
			name:          "FooPointerOZ[Base64] zero",
			input:         FooPointerOZ[Base64]{ID: &ID[id, Base64]{Value: zeroUUID}},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
		{
			name:          "*FooPointerOZ[Base64] zero",
			input:         &FooPointerOZ[Base64]{ID: &ID[id, Base64]{Value: zeroUUID}},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
		{
			name:          "FooPointerOZ[Base64] nil",
			input:         FooPointerOZ[Base64]{ID: nil},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
		{
			name:          "*FooPointerOZ[Base64] nil",
			input:         &FooPointerOZ[Base64]{ID: nil},
			expectedJSON:  `{}`,
			errorAsserter: tst.NoError(),
		},
		{
			name:          "Foo[Base64] uuid",
			input:         Foo[Base64]{ID: ID[id, Base64]{Value: uuid}},
			expectedJSON:  jsBase64,
			errorAsserter: tst.NoError(),
		},
		{
			name:          "*Foo[Base64] uuid",
			input:         &Foo[Base64]{ID: ID[id, Base64]{Value: uuid}},
			expectedJSON:  jsBase64,
			errorAsserter: tst.NoError(),
		},
		{
			name:          "FooPointerOZ[Base64] uuid",
			input:         FooPointerOZ[Base64]{ID: &ID[id, Base64]{Value: uuid}},
			expectedJSON:  jsBase64,
			errorAsserter: tst.NoError(),
		},
		{
			name:          "*FooPointerOZ[Base64] uuid",
			input:         &FooPointerOZ[Base64]{ID: &ID[id, Base64]{Value: uuid}},
			expectedJSON:  jsBase64,
			errorAsserter: tst.NoError(),
		},
		{
			name:          "Foo[Base64WithPadding] uuid",
			input:         Foo[Base64WithPadding]{ID: ID[id, Base64WithPadding]{Value: uuid}},
			expectedJSON:  `{"id":"AZomiURKfF6MYQcDvjFM_A~~"}`,
			errorAsserter: tst.NoError(),
		},
	}

	for _, tc := range marshalTests {
		t.Run(tc.name, tc.Test)
	}
}
