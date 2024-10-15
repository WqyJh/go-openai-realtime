package openairt

import (
	"crypto/rand"
	"math/big"
)

func GenerateId(prefix string, length int) string {
	const chars = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	if length <= len(prefix) {
		return prefix
	}

	result := make([]byte, length-len(prefix))
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[num.Int64()]
	}

	return prefix + string(result)
}
