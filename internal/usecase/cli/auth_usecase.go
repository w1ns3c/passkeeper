package cli

import (
	"context"
	"fmt"
	"passkeeper/internal/entities"
	"regexp"
	"strings"

	"passkeeper/internal/entities/hashes"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
)

var (
	ErrInvalidEmail        = fmt.Errorf("email is not valid by new regexp")
	ErrPassDiff            = fmt.Errorf("passwords are not the same")
	ErrEmptyUsername       = fmt.Errorf("username is empty")
	ErrEmptyEmail          = fmt.Errorf("email is empty")
	ErrEmptyPassword       = fmt.Errorf("pass is empty")
	ErrEmptyPasswordRepeat = fmt.Errorf("pass repeat is empty")
)

// Login func is client login logic for tui app
// to interact with server login logic
func (c *ClientUC) Login(ctx context.Context, login, password string) (err error) {

	hash := hashes.Hash(password)

	req := &pb.UserLoginRequest{
		Login:    login,
		Password: hash,
	}

	resp, err := c.userSvc.LoginUser(ctx, req)
	if err != nil {
		return err
	}
	c.Token = resp.Token

	userID, err := hashes.ExtractUserID(resp.Token)
	if err != nil {
		return err
	}
	//c.UserID = userID
	cryptSecret := resp.SrvSecret

	fullSecret, err := hashes.GenerateCredsSecret(password, userID, cryptSecret)
	if err != nil {
		return err
	}
	//c.CredsSecret = fullSecret

	c.User = &entities.User{
		ID:     userID,
		Login:  login,
		Hash:   password,
		Secret: fullSecret,
	}
	c.Authed = true

	return nil
}

// Register func is client register logic for tui app
// to interact with server register logic
func (c *ClientUC) Register(ctx context.Context, email, login, password, repeat string) error {
	email = strings.TrimSpace(email)
	login = strings.TrimSpace(login)

	err := c.FilterUserRegValues(login, password, repeat, email)
	if err != nil {
		return err
	}

	hash1 := hashes.Hash(password)
	hash2 := hashes.Hash(repeat)
	if hash1 != hash2 {
		return ErrPassNotSame
	}

	req := &pb.UserRegisterRequest{
		Email:      email,
		Login:      login,
		Password:   hash1,
		RePassword: hash2,
	}

	_, err = c.userSvc.RegisterUser(ctx, req)

	return err
}

// FilterUserRegValues func filter user input values from tui app
func (c *ClientUC) FilterUserRegValues(username, password, passRepeat, email string) error {
	// check if user accidentally add space
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	// passwords should not trim space, because it can contain spaces on the ends
	password = strings.Trim(password, "\n\t")

	if username == "" {
		return ErrEmptyUsername
	}
	if email == "" {
		return ErrEmptyEmail
	}

	if password == "" {
		return ErrEmptyPassword
	}

	if passRepeat == "" {
		return ErrEmptyPasswordRepeat
	}

	if err := FilterEmail(email); err != nil {
		return err
	}

	if password != passRepeat {
		return ErrPassDiff
	}

	// TODO Uncomment it
	//if len(password) < app.MinPassLen {
	//	return fmt.Errorf("password len should be a least %d signs", app.MinPassLen)
	//}

	// TODO Change User check
	//if _, ok := usersUC[username]; ok {
	//	return fmt.Errorf("user already exist")
	//}

	return nil
}

// FilterEmail check if email string is valid email
// now this pattern tested by myself and modified from
// original https://emaillistvalidation.com/blog/mastering-email-validation-in-golang-crafting-robust-regex-patterns/
func FilterEmail(email string) error {
	pattern := "^[a-zA-Z0-9._-]*[a-zA-Z0-9]+@[a-zA-Z0-9-.]+[a-zA-A0-9].[a-zA-Z]{2,}$"
	if result, _ := regexp.MatchString(pattern, email); !result {
		return ErrInvalidEmail
	}

	username := strings.Split(email, "@")[0]
	if strings.Contains(username, "..") {
		return ErrInvalidEmail // double dots in username
	}

	return nil
}

// Logout func filter user input values from tui app
func (c *ClientUC) Logout() {
	c.Token = ""
	c.User = nil
	c.Creds = nil
	c.Authed = false

	return
}
