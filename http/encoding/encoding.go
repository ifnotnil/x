package encoding

import (
	"errors"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
)

// RFC5987ExtendedNotationParameterValue decodes RFC 5987 encoded filenames expecting the extended notation
// (charset  "'" [ language ] "'" value-chars)
// example: UTF-8'en'file%20name.jpg
// https://datatracker.ietf.org/doc/html/rfc5987#section-3.2
func RFC5987ExtendedNotationParameterValue(parameterValue string) (charset string, lang string, value string, err error) {
	parts := strings.Split(parameterValue, "'")
	if len(parts) != 3 {
		return "", "", "", ErrRFC5987ParameterValueMalformed
	}

	charset, lang, value = parts[0], parts[1], parts[2]

	// unescape value
	decodedValue, er := url.QueryUnescape(value)
	if er != nil {
		return "", "", "", errors.Join(ErrRFC5987ParameterValueMalformed, er)
	}
	value = decodedValue

	if strings.ToUpper(charset) == "UTF-8" {
		return charset, lang, value, err
	}

	enc, er := encodingFromCharset(charset)
	if er != nil {
		return "", "", "", er
	}

	value, er = enc.NewDecoder().String(value)
	if er != nil {
		return "", "", "", errors.Join(ErrRFC5987ParameterValueMalformed, er)
	}

	return charset, lang, value, err
}

var (
	ErrRFC5987ParameterValueMalformed = errors.New("RFC5987 Parameter Value Malformed")
	ErrCharsetNotSupported            = errors.New("charset is not supported")
)

// FromCharset maps the official names (plus preferred mime names and aliases) for character sets to the equivalent golang [encoding.Encoding].
// https://www.iana.org/assignments/character-sets/character-sets.xhtml
// https://github.com/unicode-org/icu-data
// https://encoding.spec.whatwg.org/
func FromCharset(mimeName string) (encoding.Encoding, error) {
	m := strings.TrimSpace(strings.ToUpper(mimeName))
	return encodingFromCharset(m)
}

//go:generate go run -tags=generators mktable.go

var (
	encoderPerMIB     map[uint16]encoding.Encoding
	encoderPerMIBOnce sync.Once
)

func encodingFromCharset(mimeName string) (encoding.Encoding, error) {
	encoderPerMIBOnce.Do(initEncoderPerMID)

	mid, midFound := toMIB[strings.ToUpper(mimeName)]
	if !midFound {
		return nil, ErrCharsetNotSupported
	}
	enc, encFound := encoderPerMIB[mid]
	if !encFound || enc == nil {
		return nil, ErrCharsetNotSupported
	}

	return enc, nil
}

func initEncoderPerMID() {
	encoderPerMIB = map[uint16]encoding.Encoding{}
	for _, c := range charmap.All {
		if cm, is := c.(*charmap.Charmap); is {
			id, _ := cm.ID()
			encoderPerMIB[uint16(id)] = c
		}
	}

	encoderPerMIB[3] = charmap.Windows1252

	encoderPerMIB[106] = unicode.UTF8                                            // UTF-8
	encoderPerMIB[1013] = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)    // UTF-16BE
	encoderPerMIB[1014] = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM) // UTF-16LE
	encoderPerMIB[1015] = unicode.UTF16(unicode.BigEndian, unicode.UseBOM)       // UTF-16

	encoderPerMIB[17] = japanese.ShiftJIS
	encoderPerMIB[18] = japanese.EUCJP
	encoderPerMIB[39] = japanese.ISO2022JP

	encoderPerMIB[38] = korean.EUCKR

	encoderPerMIB[113] = simplifiedchinese.GBK
	encoderPerMIB[114] = simplifiedchinese.GB18030
	encoderPerMIB[2085] = simplifiedchinese.HZGB2312

	encoderPerMIB[2026] = traditionalchinese.Big5
}

// TODO: identify the utf-16 based on bom:
// var (
// 	utf16BEBOM = []byte{0xFE, 0xFF}
// 	utf16LEBOM = []byte{0xFF, 0xFE}
// )
