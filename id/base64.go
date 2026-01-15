package id

import (
	"encoding/base64"
)

// In URL parameters, the following characters are considered safe and do not need encoding [rfc3986](https://www.rfc-editor.org/rfc/rfc3986.html#section-3.1):
// Alphabetic characters: A-Z, a-z
// Digits: 0-9
// Hyphen: -
// Underscore: _
// Period: .
// Tilde: ~

// stdBase64 is a [base64.Encoding] based on [base64.URLEncoding] without padding character.
// alphabet of base64: ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_
var (
	stdBase64            = base64.URLEncoding.WithPadding(base64.NoPadding)
	stdBase64WithPadding = base64.URLEncoding.WithPadding('~')
)

type Base64 struct{}

func (b Base64) encode(dst, src []byte) {
	stdBase64.Encode(dst, src)
}

func (b Base64) appendEncode(dst, src []byte) []byte {
	return stdBase64.AppendEncode(dst, src)
}

func (b Base64) encodeToString(src []byte) string {
	return stdBase64.EncodeToString(src)
}

func (b Base64) encodedLen(n int) int {
	return stdBase64.EncodedLen(n)
}

func (b Base64) appendDecode(dst, src []byte) ([]byte, error) {
	return stdBase64.AppendDecode(dst, src)
}

func (b Base64) decodeString(s string) ([]byte, error) {
	return stdBase64.DecodeString(s)
}

func (b Base64) decode(dst, src []byte) (n int, err error) {
	return stdBase64.Decode(dst, src)
}

func (b Base64) decodedLen(n int) int {
	return stdBase64.DecodedLen(n)
}

type Base64WithPadding struct{}

func (b Base64WithPadding) encode(dst, src []byte) {
	stdBase64WithPadding.Encode(dst, src)
}

func (b Base64WithPadding) appendEncode(dst, src []byte) []byte {
	return stdBase64WithPadding.AppendEncode(dst, src)
}

func (b Base64WithPadding) encodeToString(src []byte) string {
	return stdBase64WithPadding.EncodeToString(src)
}

func (b Base64WithPadding) encodedLen(n int) int {
	return stdBase64WithPadding.EncodedLen(n)
}

func (b Base64WithPadding) appendDecode(dst, src []byte) ([]byte, error) {
	return stdBase64WithPadding.AppendDecode(dst, src)
}

func (b Base64WithPadding) decodeString(s string) ([]byte, error) {
	return stdBase64WithPadding.DecodeString(s)
}

func (b Base64WithPadding) decode(dst, src []byte) (n int, err error) {
	return stdBase64WithPadding.Decode(dst, src)
}

func (b Base64WithPadding) decodedLen(n int) int {
	return stdBase64WithPadding.DecodedLen(n)
}

var (
	_ encoding = Base64{}
	_ encoding = Base64WithPadding{}
)
