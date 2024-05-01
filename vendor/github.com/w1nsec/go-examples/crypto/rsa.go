package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

const (
	PRIVATE_KEY = "RSA PRIVATE KEY"
	PUBLIC_KEY  = "RSA PUBLIC KEY"
)

func GenKeys(bits int) (priv, pub []byte, err error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	privKey := x509.MarshalPKCS1PrivateKey(key)
	pubKey := x509.MarshalPKCS1PublicKey(&key.PublicKey)

	privPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  PRIVATE_KEY,
			Bytes: privKey,
		})

	pubPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  PUBLIC_KEY,
			Bytes: pubKey,
		})

	return privPem, pubPem, nil

}

// ReadPrivateKey parse private key from PEM format
func ReadPrivateKey(priv []byte) (privKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(priv)
	if block == nil || block.Type != PRIVATE_KEY {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// ReadPubKey parse public key from PEM format
func ReadPubKey(pub []byte) (pubKey *rsa.PublicKey, err error) {
	block, _ := pem.Decode(pub)
	if block == nil || block.Type != PUBLIC_KEY {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func EncryptRSA(data []byte, pubKey *rsa.PublicKey) (cipherData []byte, err error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
}

func DecryptRSA(cipherData []byte, privKey *rsa.PrivateKey) (data []byte, err error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privKey, cipherData)
}
