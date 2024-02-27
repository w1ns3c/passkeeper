package entities

import "fmt"

var (
	ErrPassNotTheSame = fmt.Errorf("passwords not match")
	ErrUserNotFound   = fmt.Errorf("user with this login not found")
)
