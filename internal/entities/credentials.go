package entities

import (
	"time"
)

type CredInf interface {
	GetID() string
	SetID(id string)
}

// TODO set ID to unexported field
type Credential struct {
	Type        BlobType  `json:"type"`
	ID          string    `json:"-"`
	Date        time.Time `json:"date"` // date for last changing this line
	Resource    string    `json:"resource"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
	Description string    `json:"desc"`
}

func (c *Credential) GetID() string {
	return c.ID
}

func (c *Credential) SetID(id string) {
	c.ID = id
}

type CryptoBlob struct {
	ID     string // blob ID
	UserID string // user's ID who owns credential
	Blob   string // encrypted saved resource (credentials)
}

type Card struct {
	Type        BlobType `json:"type"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Bank        string   `json:"bank"`
	Person      string   `json:"person"`
	Number      int      `json:"number"`
	CVC         int      `json:"cvc"`
	Expiration  string   `json:"exp"`
	PIN         int      `json:"pin"`
	Description string   `json:"desc"`
}

func (c *Card) GetID() string {
	return c.ID
}

func (c *Card) SetID(id string) {
	c.ID = id
}

type Note struct {
	Type BlobType  `json:"type"`
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Date time.Time `json:"date"` // date for last changing this line
	Body string    `json:"body"`
}

func (n *Note) GetID() string {
	return n.ID
}

func (c *Note) SetID(id string) {
	c.ID = id
}

// BlobType using for identify CryptoBlob type
type BlobType int

const (
	UserCred BlobType = iota + 1
	UserCard
	UserNote
)
