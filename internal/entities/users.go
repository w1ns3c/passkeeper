package entities

type User struct {
	ID    string
	Login string // nickname
	//Credential string
	Hash   string
	Secret string // for encrypt/decrypt passwords
	Phone  string
}
