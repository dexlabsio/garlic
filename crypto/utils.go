package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashSHA256 calculates the SHA-256 hash of the given string and returns it as a hex string.
func HashSHA256(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
