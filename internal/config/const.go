package config

import "time"

const (
	TokenLifeTime = time.Duration(time.Hour * 10)
	UserSecretLen = 32

	ChallengeLen      = 16
	ChallengeLifeTime = 5 // minutes

	TokenHeader = "token"

	MinPassLen = 8
)
