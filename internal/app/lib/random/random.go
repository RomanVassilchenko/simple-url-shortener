package random

import (
	"math/rand"
	"time"
)

// NewRandomString generates random string with given size
func NewRandomString() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	size := rnd.Intn(4) + 3

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")

	b := make([]rune, size)

	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
