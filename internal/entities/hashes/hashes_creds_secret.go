package hashes

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateCredsSecret this key will encrypt/decrypt all user's credential
// --------------------------------------------------------------------------------------------------
// | 									Full Client Secret				 							|
// --------------------------------------------------------------------------------------------------
// |	 			Client Side					| 				Server Side		 					|
// --------------------------------------------------------------------------------------------------
// | 		sha1(sha1(clearPass) + userID)		| 	decrypt(secret)	with hashes.Hash(clearPass)		| // userID -->> from server side
// --------------------------------------------------------------------------------------------------
// | 							SHA256(client_side + server_side)									|
// --------------------------------------------------------------------------------------------------
func GenerateCredsSecret(clearPass, userID, secret string) (key string, err error) {

	// client side
	clientSecret := sha1.Sum([]byte(clearPass))
	clientSecret = sha1.Sum(append(clientSecret[:], []byte(userID)...))

	// srv side
	srvHash := Hash(clearPass)
	clearSRVSecret, err := DecryptSecret(secret, srvHash)
	if err != nil {
		return "", err
	}

	keySl := sha256.Sum256(append(clientSecret[:], []byte(clearSRVSecret)...))

	return hex.EncodeToString(keySl[:]), nil
}
