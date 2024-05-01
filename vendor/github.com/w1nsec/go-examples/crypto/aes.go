package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

const (
	AES_KEY_LEN = 32
)

func EncryptAES(data []byte, key []byte) (cData []byte, err error) {
	keyBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("can't create new cipher from KEY: %v", err)

	}

	gcm, err := cipher.NewGCM(keyBlock)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	cData = gcm.Seal(nonce, nonce, data, nil)

	return cData, err
}

func DecryptAES(cData []byte, key []byte) (data []byte, err error) {
	keyBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("can't create new cipher from KEY: %v", err)
	}

	gcm, err := cipher.NewGCM(keyBlock)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, cData := cData[:nonceSize], cData[nonceSize:]

	return gcm.Open(nil, nonce, cData, nil)
}

func FillAESKey(key []byte) []byte {
	if len(key) == AES_KEY_LEN {
		return key
	}

	filledKey := make([]byte, AES_KEY_LEN)
	for i := 0; i < len(key); i++ {
		filledKey[i] = key[i]
	}

	for i := len(key); i < len(filledKey); i++ {
		filledKey[i] = 0
	}

	return filledKey
}

func FillKeyWithHash(key []byte) []byte {
	hash := sha256.Sum256(key)
	return hash[:]
}
