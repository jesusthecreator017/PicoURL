package shortcode

import (
	"crypto/sha256"
	"math/big"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateShortCode creates a short code from the original URL using SHA-256 hashing and base62 encoding
func GenerateShortCode(url string, length int) string {
	hash := sha256.Sum256([]byte(url))
	num := new(big.Int).SetBytes(hash[:8])

	encoded := make([]byte, 0, length)
	base := big.NewInt(62)
	mod := new(big.Int)

	for num.Sign() > 0 && len(encoded) < length {
		num.DivMod(num, base, mod)
		encoded = append(encoded, base62Chars[mod.Int64()])
	}

	// Pad if the encoded string is shorter than the desired length
	for len(encoded) < length {
		encoded = append(encoded, []byte(base62Chars)[0])
	}

	return string(encoded)
}
