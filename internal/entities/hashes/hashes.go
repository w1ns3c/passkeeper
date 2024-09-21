package hashes

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"github.com/w1ns3c/go-examples/crypto"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken = fmt.Errorf("token sign is not valid")
)

// GenerateHash func gen sha256 hash of (password with salt)
func GenerateHash(password, salt string) string {
	password = fmt.Sprintf("%s-%s.%s.%s", string(salt), string(password), string(password), string(salt))
	hash := sha256.Sum256([]byte(password))

	return fmt.Sprintf("%x", hash)
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
	//hash := md5.Sum([]byte(fmt.Sprintf("%s.%s.%s", salt, secret, salt)))
	//return hex.EncodeToString(hash[:])

	h := md5.New()

	return genID(secret, salt, h)
}

// GeneratePassID func to generate ID for any credential before save it to storage
func GeneratePassID(secret, salt string) string {
	h := sha512.New()

	return genID(secret, salt, h)
}

func GeneratePassID2() string {
	h := sha256.New()
	n := 32
	uid, err := crypto.GenRandStr(n)
	if err != nil {
		// TODO

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
