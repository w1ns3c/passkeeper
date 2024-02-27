package config

import "time"

const (
	TokenLifeTime = time.Duration(time.Hour * 10)
	UserSecretLen = 32

	TokenHeader = "token"
)
