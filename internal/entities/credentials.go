package entities

import (
	"time"
)

type Cred interface {
	GetID() string
}

// TODO set ID to unexported field
type Credential struct {
	ID          string    `json:"-"`
	Date        time.Time `json:"date"` // date for last changing this line
	Resource    string    `json:"resource"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
	Description string    `json:"desc"`
}

func (c Credential) GetID() string {
	return c.ID
}

type CredBlob struct {
	ID     string // blob ID
	UserID string // user's ID who owns credential
	Blob   string // encrypted saved resource (credentials)
}

type Card struct {
	ID          string
	Name        string
	Bank        string
	Person      string
	Number      int
	CVC         int
	Expiration  string
	PIN         int
	Description string
}

func (c Card) GetID() string {
	return c.ID
}

type Note struct {
	ID   string
	Name string
	Body string
}

func (n Note) GetID() string {
	return n.ID
}
