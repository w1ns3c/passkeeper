package hashes

import (
	"crypto/sha512"
	"encoding/hex"
)

// Hash func return sha512 hash of password
// Client function
func Hash(password string) string {
	h := sha512.Sum512([]byte(password))
	return hex.EncodeToString(h[:])
}
