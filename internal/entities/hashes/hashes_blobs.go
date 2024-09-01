package hashes

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/w1ns3c/go-examples/crypto"
	"passkeeper/internal/entities"
)

// EncryptBlob func just encrypt Credential to CredBlob with key
func EncryptBlob(cred *entities.Credential, key string) (blob *entities.CredBlob, err error) {
	jsonCred, err := json.Marshal(cred)
	if err != nil {
		return nil, fmt.Errorf("can't marshal cred: %v", err)
	}

	filledKey := crypto.FillKeyWithHash([]byte(key))
	cryptoCred, err := crypto.EncryptAES(jsonCred, filledKey)
	if err != nil {
		return nil, fmt.Errorf("can't encrypt json with aes: %v", err)
	}

	cryptoStr := hex.EncodeToString(cryptoCred)

	return &entities.CredBlob{
		ID:   cred.ID,
		Blob: cryptoStr,
	}, nil
}

// DecryptBlob func just decrypt CredBlob back to Credential with key
func DecryptBlob(blob *entities.CredBlob, key string) (cred *entities.Credential, err error) {
	cryptoCred, err := hex.DecodeString(blob.Blob)
	if err != nil {
		return nil, fmt.Errorf("can't decode from hex: %v", err)
	}

	filledKey := crypto.FillKeyWithHash([]byte(key))
	jsonCred, err := crypto.DecryptAES(cryptoCred, filledKey)
	if err != nil {
		return nil, fmt.Errorf("can't decrypt with aes: %v", err)
	}

	err = json.Unmarshal(jsonCred, &cred)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal from json: %v", err)
	}

	cred.ID = blob.ID

	return cred, nil
}
