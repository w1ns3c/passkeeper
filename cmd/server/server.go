package server

import (
	"time"

	"passkeeper/internal/entities"
	"passkeeper/internal/storage/memstorage"
	"passkeeper/internal/utils/hashpass"
)

func main() {

	listenAddr := "localhost:8001"

	var (
		login1 = "user1"
		login2 = "user2"
		hash1  = hashpass.Hash("123")
		hash2  = hashpass.Hash("password")

		user1 = &entities.User{
			ID:    login1,
			Login: login1,
			Hash:  hash1,
		}

		user2 = &entities.User{
			ID:    login2,
			Login: login2,
			Hash:  hash2,
		}

		// Passwords
		password1 = &entities.Credential{
			ID:          "superID",
			Date:        time.Now(),
			Resource:    "localhost",
			Login:       "my_favorite_username",
			Password:    "my_favorite_password",
			Description: "some description",
		}
		password2 = &entities.Credential{
			ID:          "superID",
			Date:        time.Now().Add(time.Second * 500),
			Resource:    "localhost",
			Login:       "admin",
			Password:    "secret password",
			Description: "some new description",
		}
	)

	users := map[string]*entities.User{
		login1: user1,
		login2: user2,
	}

	passwords := map[string]*entities.Credential{
		login1: password1,
		login2: password2,
	}

	storage := memstorage.NewMemStorage(memstorage.WithUsers(users),
		memstorage.WithPasswords(passwords))

}
