package fitbit_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stevefranchak/idlewalk/internal/fitbit"
	"github.com/stretchr/testify/assert"
)

type mockRandIntGenerator struct {
	values []*big.Int
	pos    int
}

func (m *mockRandIntGenerator) Int(max *big.Int) (*big.Int, error) {
	if m.pos >= len(m.values) {
		return nil, fmt.Errorf("exhausted values in mock rand int generator")
	}
	result := m.values[m.pos]
	m.pos++

	return result, nil
}

func toBigIntPointerSlice(nums []int) []*big.Int {
	b := make([]*big.Int, len(nums))
	for i, num := range nums {
		b[i] = big.NewInt(int64(num))
	}
	return b
}

func TestCodeVerifier(t *testing.T) {
	gen := &mockRandIntGenerator{
		// indices into CodeVerifierCharset
		values: toBigIntPointerSlice([]int{
			48, 23, 17, 51, 56, 43, 26, 53, 38, 32, 53, 47, 53, 37, 23, 58, 0, 16, 31, 56, 32, 41, 49, 57, 53, 61, 46, 20,
			25, 56, 15, 54, 26, 3, 18, 57, 50, 46, 13, 23, 20, 23, 61, 7, 14, 51, 0, 14, 37, 39, 46, 8, 32, 1, 23, 31, 30,
			20, 29, 46, 30, 45, 59, 61, 27, 37, 29, 20, 49, 35, 8, 30, 46, 57, 5, 41, 45, 0, 29, 35, 46, 4, 28, 60, 21, 22,
			41, 6, 36, 13, 59, 61, 38, 34, 56, 41,
		}),
	}
	assert := assert.New(t)

	verifier, err := fitbit.NewCodeVerifier(gen)
	assert.Nilf(err, "NewCodeVerifier errored with: %s", err)

	got := verifier.ChallengeCode()
	assert.Equal("5l-F55mJw3lGVNbXYyTsIytpxW_SI9qaGATCNUemCmw", got)
}
