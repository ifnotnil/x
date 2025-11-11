package id

import (
	"fmt"
	"math"
	"math/big"
)

type Base62Encoding struct{}

func (e *Base62Encoding) AppendEncode(dst, src []byte) []byte {
	if len(src) == 0 {
		return nil
	}

	num := big.Int{}
	num.SetBytes(src)
	return num.Append(dst, 62)
}

func base62Decode(b []byte) ([]byte, error) {
	num := big.Int{}
	n, ok := num.SetString(string(b), 62)
	if !ok {
		return nil, fmt.Errorf("tbd")
	}

	return n.Bytes(), nil
}

func base62EncodedLen(n int) int {
	// log2(62) ≈ 5.954196310386875
	return int(math.Ceil(float64(n) * 8 / 5.954196310386875))
}

func base62DecodedLen(m int) int {
	// log2(62) ≈ 5.954196310386875
	return int(math.Floor(float64(m) * 5.954196310386875 / 8))
}
