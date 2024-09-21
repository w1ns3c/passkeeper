package hashes

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/w1ns3c/go-examples/crypto"
)

var (
// ErrorDecode = fmt.Errorf("can't decode secret from hex")
)

// GenerateSecret func for generate random hex string
func GenerateSecret(secretLen int) (secret string, err error) {
	s, err := crypto.GenRandStr(secretLen)
	if err != nil {
		return "", nil
	}

	return s, nil
}

// EncryptSecret func will encrypt user's secret with his utils.hashpass.Hash(password)
// via AES, return encrypted hex value
func EncryptSecret(secret, key string) (cryptSecret string, err error) {
	keyH := sha256.Sum256([]byte(key))
	filledKey := crypto.FillKeyWithHash(keyH[:])
	cryptSecretSl, err := crypto.EncryptAES([]byte(secret), filledKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(cryptSecretSl), nil
}

// DecryptSecret func will decrypt user's secret with his utils.hashpass.Hash(password)
// via AES, return original secret value
func DecryptSecret(cryptSecret, key string) (secret string, err error) {
	keyH := sha256.Sum256([]byte(key))
	filledKey := crypto.FillKeyWithHash(keyH[:])
	cryptSecretSl, err := hex.DecodeString(cryptSecret)
	if err != nil {
		return "", err
	}

	secretsl, err := crypto.DecryptAES(cryptSecretSl, filledKey)
	if err != nil {
		return "", err
	}

	return string(secretsl), nil
}
