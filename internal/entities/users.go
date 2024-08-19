package entities

type User struct {
	ID    string
	Login string // nickname
	Hash  string

	Salt   string // for encrypt user password
	Secret string // for encrypt/decrypt saved passwords

	Phone string
	Email string
}
