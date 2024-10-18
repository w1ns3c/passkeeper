package structs

// User is a main entity of any user account
type User struct {
	ID    string // generated on Hash and Salt
	Login string // nickname
	Hash  string // bcrypt hash of sha512 hash of password | bcrypt(sha512(password))

	Salt   string // for encrypt user's password and sign user token
	Secret string // for encrypt/decrypt saved passwords

	Phone string
	Email string
}
