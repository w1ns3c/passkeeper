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

func SaveCard(cards []*Card, ind int, card *Card) error {
	cards[ind] = card
	return nil
}

func SaveNote(notes []*Note, ind int, note *Note) error {
	notes[ind] = note
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

func AddCard(cards []*Card, newCard *Card) (newCards []*Card, err error) {
	if len(cards) == 0 {

		return append(cards, newCard), nil
	}

	cards = append(cards, &Card{})

	for i := len(cards) - 1; i > 0; i-- {
		cards[i] = cards[i-1]
	}

	cards[0] = newCard

	return cards, nil
}

func AddNote(notes []*Note, newNote *Note) (newNotes []*Note, err error) {
	if len(notes) == 0 {

		return append(notes, newNote), nil
	}

	notes = append(notes, &Note{})

	for i := len(notes) - 1; i > 0; i-- {
		notes[i] = notes[i-1]
	}

	notes[0] = newNote

	return notes, nil
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

func DeleteCred(creds []*Credential, ind int) (newCreds []*Credential, err error) {
	for i := ind; i < len(creds)-1; i++ {
		creds[i] = creds[i+1]
	}
	creds = creds[:len(creds)-1]
	return creds, nil
}

func DeleteCard(cards []*Card, ind int) (newCards []*Card, err error) {
	for i := ind; i < len(cards)-1; i++ {
		cards[i] = cards[i+1]
	}
	cards = cards[:len(cards)-1]
	return cards, nil
}

func DeleteNote(notes []*Note, ind int) (newCreds []*Note, err error) {
	for i := ind; i < len(notes)-1; i++ {
		notes[i] = notes[i+1]
	}
	notes = notes[:len(notes)-1]
	return notes, nil
}

func GenHash(s string) string {
	s1 := md5.Sum([]byte(s))

	return hex.EncodeToString(s1[:])
}
