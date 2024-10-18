package structs

import (
	"time"
)

type CredInf interface {
	GetID() string
	SetID(id string)
}

// BlobType using for identify CryptoBlob
type BlobType int

const (
	BlobCred BlobType = iota + 1
	BlobCard
	BlobNote
	BlobFile
)

// Credential is one of blob struct type
type Credential struct {
	Type        BlobType  `json:"type"`
	ID          string    `json:"-"`
	Date        time.Time `json:"date"` // date for last changing this line
	Resource    string    `json:"resource"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
	Description string    `json:"desc"`
}

// Card is one of blob struct type
type Card struct {
	Type        BlobType  `json:"type"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Bank        string    `json:"bank"`
	Person      string    `json:"person"`
	Number      int       `json:"number"`
	CVC         int       `json:"cvc"`
	Expiration  time.Time `json:"exp"`
	PIN         int       `json:"pin"`
	Description string    `json:"desc"`
}

// Note is one of blob struct type
type Note struct {
	Type BlobType  `json:"type"`
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Date time.Time `json:"date"` // date for last changing this line
	Body string    `json:"body"`
}

// File is one of blob struct type
type File struct {
	ID   string
	Type BlobType
	Name string
	Body []byte
}

// CryptoBlob encrypted entity that stores on server side
type CryptoBlob struct {
	ID     string // blob ID
	UserID string // user's ID who owns credential
	Blob   string // encrypted saved resource (credentials)
}

// GetID GeID return ID of entity
func (c *Credential) GetID() string {
	return c.ID
}

// SetID set ID to entity
func (c *Credential) SetID(id string) {
	c.ID = id
}

// GetID return ID of entity
func (c *Card) GetID() string {
	return c.ID
}

// SetID set ID to entity
func (c *Card) SetID(id string) {
	c.ID = id
}

// GetID return ID of entity
func (f *File) GetID() string {
	return f.ID
}

// SetID set ID to entity
func (f *File) SetID(id string) {
	f.ID = id
}

// GetID return ID of entity
func (n *Note) GetID() string {
	return n.ID
}

// SetID set ID to entity
func (c *Note) SetID(id string) {
	c.ID = id
}
