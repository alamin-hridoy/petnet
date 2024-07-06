package random

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
)

func CIntn(n int) int {
	for i := 0; i < 10; i++ {
		num, err := crand.Int(crand.Reader, big.NewInt(int64(n)))
		if err != nil {
			continue
		}
		return int(num.Int64())
	}
	return 0
}

func String(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	read, err := crand.Read(bytes)
	// Note that err == nil iff we read len(b) bytes.
	if err != nil {
		return "", err
	}
	if read != n {
		return "", fmt.Errorf("failed to read %d bytes from random, bytes read %d", n, read)
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func NumberString(n int) string {
	letters := []rune("123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

func InvitationCode(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[CIntn(len(letters))]
	}
	return string(b)
}
