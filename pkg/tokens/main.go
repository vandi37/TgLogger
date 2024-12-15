package tokens

import (
	"crypto/rand"
	"math/big"
	"regexp"
)

var length = 25

func Ok(token string) bool {
	if ok, err := regexp.MatchString(`^[A-Z0-9]*$`, token); !ok || err != nil || len(token) != length {
		return false
	}
	return true
}

var allowed = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func NewToken() string {
	alwLen := big.NewInt(int64(len(allowed)))
	res := make([]byte, length)

	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, alwLen)
		res[i] = allowed[n.Int64()]
	}

	return string(res)
}
