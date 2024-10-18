package hashes

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/w1ns3c/go-examples/crypto"

	"passkeeper/internal/entities/structs"
)

// EncryptBlob func just encrypt Credential to CryptoBlob with key
func EncryptBlob(cred structs.CredInf, key string) (blob *structs.CryptoBlob, err error) {
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

	return &structs.CryptoBlob{
		ID:   cred.GetID(),
		Blob: cryptoStr,
	}, nil
}

// DecryptBlob func just decrypt CryptoBlob back to Credential with key
func DecryptBlob(blob *structs.CryptoBlob, key string) (cred structs.CredInf, err error) {
	cryptoCred, err := hex.DecodeString(blob.Blob)
	if err != nil {
		return nil, fmt.Errorf("can't decode from hex: %v", err)
	}

	filledKey := crypto.FillKeyWithHash([]byte(key))
	jsonCred, err := crypto.DecryptAES(cryptoCred, filledKey)
	if err != nil {
		return nil, fmt.Errorf("can't decrypt with aes: %v", err)
	}

	type blobType struct {
		BlobType structs.BlobType `json:"type"`
	}
	var tmp blobType

	err = json.Unmarshal(jsonCred, &tmp)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal from json: %v", err)
	}

	switch tmp.BlobType {
	case structs.BlobCred:
		var tmpCred *structs.Credential
		err = json.Unmarshal(jsonCred, &tmpCred)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal Cred from json: %v", err)
		}

		cred = tmpCred

	case structs.BlobCard:
		var tmpCard *structs.Card
		err = json.Unmarshal(jsonCred, &tmpCard)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal Card from json: %v", err)
		}

		cred = tmpCard

	case structs.BlobNote:
		var tmpNote *structs.Note
		err = json.Unmarshal(jsonCred, &tmpNote)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal Card from json: %v", err)
		}

		cred = tmpNote

	case structs.BlobFile:
		var tmpFile *structs.File
		err = json.Unmarshal(jsonCred, &tmpFile)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal Card from json: %v", err)
		}

		cred = tmpFile

	default:
		return nil, fmt.Errorf("unknown blob type: %v", tmp.BlobType)
	}

	cred.SetID(blob.ID)

	return cred, nil
}
