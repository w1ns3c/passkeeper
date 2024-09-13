package config

import "time"

const (
	TokenLifeTime   = time.Duration(time.Hour * 10)
	UserPassSaltLen = 32
	UserSecretLen   = 32

	TokenHeader = "token"

	MinPassLen = 8

	DefaultAddr = "localhost:8000"

	Level = "Debug"

	SyncMax     = time.Second * 1000
	SyncMin     = time.Second * 1
	SyncDefault = time.Second * 5
)
