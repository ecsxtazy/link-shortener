package service

import (
	"crypto/rand"
	"math/big"
)

const (
	alphabet = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789_"
	length   = 10
)

type Generator interface {
	Generate() (string, error)
}

func Generate() (string, error) {
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		code[i] = alphabet[n.Int64()]
	}
	return string(code), nil
}
