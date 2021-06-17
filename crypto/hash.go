package crypto

import (
	"math/rand"

	"golang.org/x/crypto/sha3"
)

var (
	chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func Hash(key []byte) []byte {
	hasher := sha3.New256()
	hasher.Write(key)
	return hasher.Sum(nil)
}

func RandomStr(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func RandomBytes(n int) []byte {
	return []byte(RandomStr(n))
}
