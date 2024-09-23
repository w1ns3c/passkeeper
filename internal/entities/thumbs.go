package entities

import (
	"crypto/md5"
	"encoding/hex"

	"passkeeper/internal/entities/structs"
)

// SaveCred save cred blob on client side
func SaveCred(creds []*structs.Credential, ind int, cred *structs.Credential) error {
	creds[ind] = cred
	return nil
}

// SaveCard save card blob on client side
func SaveCard(cards []*structs.Card, ind int, card *structs.Card) error {
	cards[ind] = card
	return nil
}

// SaveNote save note blob on client side
func SaveNote(notes []*structs.Note, ind int, note *structs.Note) error {
	notes[ind] = note
	return nil
}

// SaveFile save file blob on client side
func SaveFile(files []*structs.File, ind int, file *structs.File) error {
	files[ind] = file
	return nil
}

// AddCred add cred blob on client side
func AddCred(creds []*structs.Credential, newCred *structs.Credential) (newCreds []*structs.Credential, err error) {
	if len(creds) == 0 {

		return append(creds, newCred), nil
	}

	creds = append(creds, &structs.Credential{})

	for i := len(creds) - 1; i > 0; i-- {
		creds[i] = creds[i-1]
	}

	creds[0] = newCred

	return creds, nil
}

// AddCard add card blob on client side
func AddCard(cards []*structs.Card, newCard *structs.Card) (newCards []*structs.Card, err error) {
	if len(cards) == 0 {

		return append(cards, newCard), nil
	}

	cards = append(cards, &structs.Card{})

	for i := len(cards) - 1; i > 0; i-- {
		cards[i] = cards[i-1]
	}

	cards[0] = newCard

	return cards, nil
}

// AddNote add note blob on client side
func AddNote(notes []*structs.Note, newNote *structs.Note) (newNotes []*structs.Note, err error) {
	if len(notes) == 0 {

		return append(notes, newNote), nil
	}

	notes = append(notes, &structs.Note{})

	for i := len(notes) - 1; i > 0; i-- {
		notes[i] = notes[i-1]
	}

	notes[0] = newNote

	return notes, nil
}

// AddFile add file blob on client side
func AddFile(files []*structs.File, newFile *structs.File) (newFiles []*structs.File, err error) {
	if len(files) == 0 {

		return append(files, newFile), nil
	}

	files = append(files, &structs.File{})

	for i := len(files) - 1; i > 0; i-- {
		files[i] = files[i-1]
	}

	files[0] = newFile

	return files, nil
}

// DeleteCred delete cred blob on client side
func DeleteCred(creds []*structs.Credential, ind int) (newCreds []*structs.Credential, err error) {
	for i := ind; i < len(creds)-1; i++ {
		creds[i] = creds[i+1]
	}
	creds = creds[:len(creds)-1]
	return creds, nil
}

// DeleteCard delete card blob on client side
func DeleteCard(cards []*structs.Card, ind int) (newCards []*structs.Card, err error) {
	for i := ind; i < len(cards)-1; i++ {
		cards[i] = cards[i+1]
	}
	cards = cards[:len(cards)-1]
	return cards, nil
}

// DeleteNote delete note blob on client side
func DeleteNote(notes []*structs.Note, ind int) (newCreds []*structs.Note, err error) {
	for i := ind; i < len(notes)-1; i++ {
		notes[i] = notes[i+1]
	}
	notes = notes[:len(notes)-1]
	return notes, nil
}

// DeleteFile delete file blob on client side
func DeleteFile(files []*structs.File, ind int) (newFiles []*structs.File, err error) {
	for i := ind; i < len(files)-1; i++ {
		files[i] = files[i+1]
	}
	files = files[:len(files)-1]
	return files, nil
}

// GenHash just generate md5 hash of string
func GenHash(s string) string {
	s1 := md5.Sum([]byte(s))

	return hex.EncodeToString(s1[:])
}
