package http

import "encoding/base64"

// URLSafeBase64 returns a [base64.Encoding] based on [base64.URLEncoding] replacing the default padding character ('=') padding character to a url safe one ('~').
// In URL parameters, the following characters are considered safe and do not need encoding [rfc3986](https://www.rfc-editor.org/rfc/rfc3986.html#section-3.1):
// Alphabetic characters: A-Z, a-z
// Digits: 0-9
// Hyphen: -
// Underscore: _
// Period: .
// Tilde: ~
func URLSafeBase64() *base64.Encoding {
	return base64.URLEncoding.WithPadding('~')
}

var urlSafeBase64 = URLSafeBase64()
