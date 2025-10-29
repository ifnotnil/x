package id

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// In URL parameters, the following characters are considered safe and do not need encoding [rfc3986](https://www.rfc-editor.org/rfc/rfc3986.html#section-3.1):
// Alphabetic characters: A-Z, a-z
// Digits: 0-9
// Hyphen: -
// Underscore: _
// Period: .
// Tilde: ~

// Base64 is a [base64.Encoding] based on [base64.URLEncoding] without padding character.
// alphabet of base64: ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_
var (
	Base64            = base64.URLEncoding.WithPadding(base64.NoPadding)
	Base64WithPadding = base64.URLEncoding.WithPadding('~')
)

var (
	base64UUIDEncodedLen     = Base64.EncodedLen(uuidSize)
	base64UUIDEncodedLenJSON = base64UUIDEncodedLen + 2
	zeroUUID                 = [uuidSize]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

const uuidSize = 16

type Base64UUID[U ~[uuidSize]byte] struct {
	Value U
}

func (u Base64UUID[U]) IsZero() bool {
	return u.Value == zeroUUID
}

func (u Base64UUID[U]) MarshalJSON() ([]byte, error) {
	b := make([]byte, 1, base64UUIDEncodedLenJSON)
	b[0] = '"'
	sub := b[1:][:base64UUIDEncodedLen]
	Base64.Encode(sub, u.Value[:])
	b = b[0 : base64UUIDEncodedLenJSON-1]
	b = append(b, '"')

	return b, nil
}

func (u *Base64UUID[U]) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	decBytes, err := Base64.DecodeString(s)
	if err != nil {
		return err
	}

	if len(decBytes) != uuidSize {
		return ErrMalformedUUID
	}

	copy(u.Value[:], decBytes)

	return nil
}

var ErrMalformedUUID = errors.New("malformed uuid")

var (
	_ json.Marshaler   = (*Base64UUID[[uuidSize]byte])(nil)
	_ json.Unmarshaler = (*Base64UUID[[uuidSize]byte])(nil)
)
