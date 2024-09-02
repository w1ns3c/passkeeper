package entities

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func Save(creds []*Credential, ind int, res, login, password, desc string) error {
	creds[ind].Resource = res
	creds[ind].Login = login
	creds[ind].Password = password
	creds[ind].Description = desc

	return nil
}

func SaveCred(creds []*Credential, ind int, cred *Credential) error {
	creds[ind] = cred
	return nil
}

func AddCred(creds []*Credential, newCred *Credential) (newCreds []*Credential, err error) {
	if len(creds) == 0 {

		return append(creds, newCred), nil
	}

	creds = append(creds, &Credential{})

	for i := len(creds) - 1; i > 0; i-- {
		creds[i] = creds[i-1]
	}

	creds[0] = newCred

	return creds, nil
}

func Add(creds []*Credential, res, login, password, desc string) (newCreds []*Credential, err error) {

	tmpCred := &Credential{
		ID:          GenHash(res),
		Date:        time.Now(),
		Resource:    res,
		Login:       login,
		Password:    password,
		Description: desc,
	}

	if len(creds) == 0 {
		return append(creds, tmpCred), nil
	}

	creds = append(creds, &Credential{})

	for i := len(creds) - 1; i > 0; i-- {
		creds[i] = creds[i-1]
	}

	creds[0] = tmpCred

	return creds, nil
}

func Delete(creds []*Credential, ind int) (newCreds []*Credential, err error) {
	for i := ind; i < len(creds)-1; i++ {
		creds[i] = creds[i+1]
	}
	creds = creds[:len(creds)-1]
	return creds, nil
}

func GenHash(s string) string {
	s1 := md5.Sum([]byte(s))

	return hex.EncodeToString(s1[:])
}
