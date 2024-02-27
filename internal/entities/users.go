package entities

type User struct {
	ID    string
	Login string // nickname
	//Password string
	Hash  string
	Token string
	Phone string
}
