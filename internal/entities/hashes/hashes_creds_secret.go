package hashes

import (
	"crypto/sha256"
	"encoding/hex"
)

// GenerateCredsSecret will generate key to encrypt/decrypt all user's credentials
// ----------------------------------------------------------------------------------------------------------
// | 										Full Client CredsSecret				 							|
// ----------------------------------------------------------------------------------------------------------
// |	 				Client Side						| 					Server Side	 					|
// ----------------------------------------------------------------------------------------------------------
// | 	sha256(sha256(clearPass) + userID_from_srv)		| 	 decrypt(secret) with hashes.Hash(clearPass)	|
// ----------------------------------------------------------------------------------------------------------
// |	 								SHA256(client_side + server_side)									|
// ----------------------------------------------------------------------------------------------------------
func GenerateCredsSecret(clearPass, userID, cryptSecret string) (key string, err error) {

	// client side
	clientSecret := sha256.Sum256([]byte(clearPass))
	clientSecret = sha256.Sum256(append(clientSecret[:], []byte(userID)...))

	// srv side
	srvHash := Hash(clearPass)
	clearSRVSecret, err := DecryptSecret(cryptSecret, srvHash)
	if err != nil {
		return "", err
	}

	keySl := sha256.Sum256(append(clientSecret[:], []byte(clearSRVSecret)...))

	return hex.EncodeToString(keySl[:]), nil
}
