package internal

import (
	"crypto/rand"
	"math/big"
)

type RandIntGenerator interface {
	Int(max *big.Int) (*big.Int, error)
}

type cryptoRandIntGenerator struct{}

func (c *cryptoRandIntGenerator) Int(max *big.Int) (*big.Int, error) {
	return rand.Int(rand.Reader, max)
}

var CryptoRandIntGenerator RandIntGenerator = new(cryptoRandIntGenerator)
