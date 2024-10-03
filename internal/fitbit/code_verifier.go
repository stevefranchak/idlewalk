package fitbit

import (
	"crypto/sha256"
	"encoding/base64"
	"math/big"
	"strings"

	"github.com/stevefranchak/idlewalk/internal"
)

const (
	CodeVerifierLength  int    = 96
	CodeVerifierCharset string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	CodeChallengeMethod string = "S256"
)

type CodeVerifier struct {
	code []byte
}

func (verifier *CodeVerifier) SecretCode() string {
	return string(verifier.code)
}

func (verifier *CodeVerifier) ChallengeCodeHash() []byte {
	hasher := sha256.New()
	hasher.Write(verifier.code)
	return hasher.Sum(nil)
}

func (verifier *CodeVerifier) ChallengeCode() string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(verifier.ChallengeCodeHash()), "=")
}

func NewCodeVerifier(gen internal.RandIntGenerator) (*CodeVerifier, error) {
	result := make([]byte, CodeVerifierLength)
	charsetLength := big.NewInt(int64(len(CodeVerifierCharset)))

	for i := 0; i < CodeVerifierLength; i++ {
		randomIndex, err := gen.Int(charsetLength)
		if err != nil {
			return nil, err
		}
		result[i] = CodeVerifierCharset[randomIndex.Int64()]
	}

	verifier := CodeVerifier{
		code: result,
	}
	return &verifier, nil
}
