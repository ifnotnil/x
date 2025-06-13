package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRFC5987ExtendedNotationParameterValue(t *testing.T) {
	tests := map[string]struct {
		input           string
		expectedCharset string
		expectedLang    string
		expectedValue   string
		errorAssertion  require.ErrorAssertionFunc
	}{
		"happy path UTF-8 with lang": {
			input:           `UTF-8'en'file%20name.jpg`,
			expectedCharset: "UTF-8",
			expectedLang:    "en",
			expectedValue:   "file name.jpg",
			errorAssertion:  require.NoError,
		},
		"happy path UTF-8 no lang": {
			input:           `UTF-8''file%20name.jpg`,
			expectedCharset: "UTF-8",
			expectedLang:    "",
			expectedValue:   "file name.jpg",
			errorAssertion:  require.NoError,
		},
		"happy path UTF-8 no lang and special characters": {
			input:           `UTF-8''%c2%a3%20and%20%e2%82%ac%20rates.txt`,
			expectedCharset: "UTF-8",
			expectedLang:    "",
			expectedValue:   "£ and € rates.txt",
			errorAssertion:  require.NoError,
		},
		"happy path iso-8859-7 no lang": {
			input:           `iso-8859-7''%EA%E1%EB%E7%EC%DD%F1%E1+%DE%EB%E9%E5%2C+%EA%E1%EB%E7%EC%DD%F1%E1`,
			expectedCharset: "iso-8859-7",
			expectedLang:    "",
			expectedValue:   "καλημέρα ήλιε, καλημέρα",
			errorAssertion:  require.NoError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			charset, lang, value, err := RFC5987ExtendedNotationParameterValue(tc.input)
			tc.errorAssertion(t, err)
			require.Equal(t, tc.expectedCharset, charset)
			require.Equal(t, tc.expectedLang, lang)
			require.Equal(t, tc.expectedValue, value)
		})
	}
}
