package keygenerator

import (
	"math/rand"
	"time"
)

type (
	Generator struct {
	}

	KeyGenerator interface {
		Generate() string
	}
)

const base62Len uint64 = 62

var base62set = []rune{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

func (*Generator) Generate() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Uint64()

	return encode62(id)
}

func encode62(in uint64) string {
	if in == 0 {
		return ""
	}
	out := make([]rune, 0)
	tmp := in
	for tmp > 0 {
		if tmp == base62Len {
			break
		}
		m := tmp % base62Len
		if m == 0 {
			m = base62Len
		}
		tmp /= base62Len
		c := base62set[m-1]
		out = append(out, c)
	}

	return string(out)
}
