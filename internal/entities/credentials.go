package entities

import (
	"time"
)

// TODO set ID to unexported field
type Credential struct {
	ID          string    `json:"-"`
	Date        time.Time `json:"date"` // date for last changing this line
	Resource    string    `json:"resource"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
	Description string    `json:"desc"`
}

type CredBlob struct {
	ID     string // blob ID
	UserID string // user's ID who owns credential
	Blob   string // encrypted saved resource (credentials)
}
