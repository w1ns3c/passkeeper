package config

import (
	"path/filepath"
	"time"
)

const (
	TokenLifeTime   = time.Duration(time.Hour * 10) // how often token will be valid
	UserPassSaltLen = 32                            // len of user's salt string
	UserSecretLen   = 32                            // len of user's secret string

	TokenHeader = "token" // header in requests/response to store token

	MinPassLen = 8 // min len of user password during registartion

	DefaultAddr = "localhost:8000" // on this port server will start

	Level = "Debug" // default log level

	SyncMax     = time.Second * 1000 // max blobs sync time
	SyncMin     = time.Second * 1    // min blobs sync time
	SyncDefault = time.Second * 25   // default blobs sync time

	MaxNameLen     = 15   // max len for blob's name in card/note list
	MaxFilenameLen = 40   // max len for blob's name in file list
	MaxTextAreaLen = 200  // max len for blob's description
	MaxNoteLen     = 1000 // max len for note blob body

	CliLogFileFile = "passkeeper.log"
)

var (
	CliLogFilePath = filepath.Join(CliLogDir, CliLogFileFile)
)
