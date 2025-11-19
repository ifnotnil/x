package id

import (
	"encoding"
	"encoding/json"
	"errors"
)

const uuidSize = 16

type Encoding interface {
	encode(dst, src []byte)
	appendEncode(dst, src []byte) []byte
	// encodeToString(src []byte) string
	encodedLen(n int) int
	// appendDecode(dst, src []byte) ([]byte, error)
	// decodeString(s string) ([]byte, error)
	decode(dst, src []byte) (n int, err error)
	decodedLen(n int) int
}

type ID[UUID ~[uuidSize]byte, Enc Encoding] struct {
	Value UUID
}

func (u ID[UUID, Enc]) IsZero() bool {
	return u.Value == zeroUUID
}

func (u ID[UUID, Enc]) MarshalJSON() ([]byte, error) {
	var enc Enc
	ln := enc.encodedLen(uuidSize)
	b := make([]byte, 0, ln+2)
	b = append(b, '"')

	b = enc.appendEncode(b, u.Value[:])

	b = append(b, '"')

	return b, nil
}

func (u ID[UUID, Enc]) AppendText(b []byte) ([]byte, error) {
	var enc Enc
	return enc.appendEncode(b, u.Value[:]), nil
}

func (u ID[UUID, Enc]) MarshalText() ([]byte, error) {
	return u.AppendText(nil)
}

func (u *ID[UUID, Enc]) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	return u.UnmarshalText([]byte(s))
}

func (u *ID[UUID, Enc]) UnmarshalText(text []byte) error {
	var enc Enc
	dec := [uuidSize]byte{}

	n, err := enc.decode(dec[:], text)
	if err != nil {
		return err
	}

	if n != uuidSize {
		return ErrMalformedUUID
	}

	u.Value = dec

	return nil
}

var (
	_ json.Marshaler           = (*ID[[uuidSize]byte, Base64])(nil)
	_ json.Unmarshaler         = (*ID[[uuidSize]byte, Base64])(nil)
	_ encoding.TextAppender    = (*ID[[uuidSize]byte, Base64])(nil)
	_ encoding.TextMarshaler   = (*ID[[uuidSize]byte, Base64])(nil)
	_ encoding.TextUnmarshaler = (*ID[[uuidSize]byte, Base64])(nil)
)

var ErrMalformedUUID = errors.New("malformed uuid")

var zeroUUID = [uuidSize]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
