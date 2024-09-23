package hashes

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"github.com/w1ns3c/go-examples/crypto"

	"golang.org/x/crypto/bcrypt"
)

// GenerateHash func gen sha256 hash of (password with salt)
func GenerateHash(password, salt string) string {
	password = fmt.Sprintf("%s-%s.%s.%s", string(salt), string(password), string(password), string(salt))
	h := sha256.Sum256([]byte(password))

	return fmt.Sprintf("%x", h)
}

// GenerateCryptoHash func
func GenerateCryptoHash(password, salt string) (hash string, err error) {
	password = GenerateHash(password, salt)
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

func ComparePassAndCryptoHash(password, hash string, salt string) bool {
	genHash := GenerateHash(password, salt)

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(genHash)); err != nil {
		return false
	}
	return true
}

// GenerateUserID func to generate ID for any user before save it to storage
func GenerateUserID(secret, salt string) string {
	h := md5.New()

	return genID(secret, salt, h)
}

// GeneratePassID func to generate ID for any credential before save it to storage
func GeneratePassID() string {
	h := sha256.New()
	n := 32
	uid, err := crypto.GenRandStr(n)
	if err != nil {
		// TODO
		return ""
	}

	_, _ = io.WriteString(h, uid)

	return hex.EncodeToString(h.Sum(nil))
}

// genID main scheme for generating ID, but with variable hash type
func genID(secret, salt string, h hash.Hash) string {
	_, err := io.WriteString(h, fmt.Sprintf("%s.%s.%s", salt, secret, salt))
	if err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}
