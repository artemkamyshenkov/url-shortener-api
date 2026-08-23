package shortener

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

type RandomCodeGenerator struct {
}

func (generator *RandomCodeGenerator) Generate(codeLength int) (string, error) {
	if codeLength <= 0 {
		return "", errors.New("code length must be positive")
	}

	result := make([]byte, codeLength)
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for i := 0; i < codeLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))

		if err != nil {
			return "", fmt.Errorf("generate short code: %w", err)
		}

		symbol := alphabet[int(n.Int64())]
		result[i] = symbol
	}

	return string(result), nil
}
