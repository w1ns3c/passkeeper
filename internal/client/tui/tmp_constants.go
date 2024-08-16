package tui

import (
	"passkeeper/internal/entities"
	"time"
)

var (
	users = map[string]string{
		"user":  "pass",
		"user2": "pass",
		"empty": "pass",
		"error": "pass",
	}
	id1 = entities.GenHash("1")
	id2 = entities.GenHash("2")
	id3 = entities.GenHash("3")

	DateTime = "2006-01-02 15:05:05"
	time1, _ = time.Parse(DateTime, "2000-10-01 15:02:01")
	time2, _ = time.Parse(DateTime, "2002-10-11 12:12:01")
	time3, _ = time.Parse(DateTime, "2010-10-22 13:22:01")

	user1 = "username1"
	user2 = "username2"
	user3 = "max"

	pass1 = "P@#ss"
	pass2 = "password123"
	pass3 = "s1mple"

	res1 = "mysite.com"
	res2 = "https://site333.com/"
	res3 = "simple.org.com"

	des1 = "Some simple password \n for my site"
	des2 = "12312"
	des3 = "simple description"

	CredsList = []entities.Credential{
		{
			ID:          string(id1[:]),
			Date:        time1,
			Resource:    res1,
			Login:       user1,
			Password:    pass1,
			Description: des1,
		},
		{
			ID:          string(id2[:]),
			Date:        time2,
			Resource:    res2,
			Login:       user2,
			Password:    pass2,
			Description: des2,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    res3,
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    res3,
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    "deleting",
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    res3,
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    res3,
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    res3,
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    res3,
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    res3,
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    res3,
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    res3,
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    res3,
			Login:       user3,
			Password:    pass3,
			Description: des3,
		},
		{
			ID:          string(id3[:]),
			Date:        time3,
			Resource:    "Last Resource",
			Login:       "last " + user3,
			Password:    "last " + pass3,
			Description: "last " + des3,
		},
	}
)
