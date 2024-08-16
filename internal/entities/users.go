package entities

import (
	"time"
)

type User struct {
	ID            string
	Login         string // nickname
	Hash          string
	Challenge     string
	ChallengeTime time.Time
	Secret        string // for encrypt/decrypt passwords
	Phone         string
	Email         string
}
