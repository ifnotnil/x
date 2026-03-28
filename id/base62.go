package id

import (
	"fmt"
	"math"
	"math/big"
)

type Base62 struct{}

func (e Base62) appendEncode(dst, src []byte) []byte {
	if len(src) == 0 {
		return nil
	}

	num := big.Int{}
	num.SetBytes(src)
	return num.Append(dst, 62)
}

func (e Base62) encodedLen(n int) int {
	return base62EncodedLen(n)
}

func (e Base62) decode(dst, src []byte) (n int, err error) {
	if len(src) == 0 {
		return 0, nil
	}

	num := big.Int{}
	nn, ok := num.SetString(string(src), 62)
	if !ok {
		return 0, fmt.Errorf("base62 error while parsing")
	}

	bb := nn.Bytes()
	// if len(bb) != uuidSize {
	// 	return 0, fmt.Errorf("base62 error while parsing")
	// }

	return copy(dst, bb), nil
}

// lg262 : log2(62) ≈ 5.954196310386875
const lg262 = 5.954196310386875

func base62EncodedLen(n int) int {
	return int(math.Ceil(float64(n) * 8 / lg262))
}

func base62DecodedLen(m int) int {
	return int(math.Floor(float64(m) * lg262 / 8))
}

var _ encoding = Base62{}
