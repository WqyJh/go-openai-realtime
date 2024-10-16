package openairt

import (
	"crypto/rand"
	"math/big"
)

// GenerateId generates a random ID with a prefix and a specified length.
// The length of the returned ID is equal to the length parameter, therefore the prefix must be shorter than the length.
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
