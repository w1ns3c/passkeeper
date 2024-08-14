package cli

import (
	"context"
	"fmt"
	"github.com/w1ns3c/passkeeper/internal/config"
	pb "github.com/w1ns3c/passkeeper/internal/transport/grpc/protofiles/proto"
	"github.com/w1ns3c/passkeeper/internal/utils/hashpass"
	"google.golang.org/grpc/metadata"
	"net/mail"
	"regexp"
	"strings"
)

// Login func is client login logic for tui app
// to interact with server login logic
func (c *ClientUC) Login(ctx context.Context, login, password string) error {

	hash := hashpass.Hash(password)

	req := &pb.UserLoginRequest{
		Login:    login,
		Password: hash,
	}

	resp, err := c.userSvc.LoginUser(ctx, req)
	if err != nil {
		return err
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		tokens := md.Get(config.TokenHeader)
		if len(tokens) < 1 {
			return fmt.Errorf("no token in response")
		}
		c.Token = tokens[0]
	}

	c.SecretHash = resp.GetSecret()

	return nil
}

// Register func is client register logic for tui app
// to interact with server register logic
func (c *ClientUC) Register(ctx context.Context, login, password, repeat, email string) error {
	email = strings.TrimSpace(email)
	login = strings.TrimSpace(login)

	err := c.FilterUserRegValues(login, password, repeat, email)
	if err != nil {
		return err
	}

	hash1 := hashpass.Hash(password)
	hash2 := hashpass.Hash(repeat)
	if hash1 != hash2 {
		return ErrPassNotSame
	}

	req := &pb.UserRegisterRequest{
		Login:      login,
		Password:   hash1,
		RePassword: hash2,
		Email:      email,
	}

	resp, err := c.userSvc.RegisterUser(ctx, req)
	if err != nil {
		return err
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		tokens := md.Get(config.TokenHeader)
		if len(tokens) < 1 {
			return fmt.Errorf("no token in response")
		}
		c.Token = tokens[0]
	}

	c.SecretHash = resp.GetSecret()

	return nil
}

// FilterUserRegValues func filter user input values from tui app
func (c *ClientUC) FilterUserRegValues(username, password, passRepeat, email string) error {
	if username == "" {
		return fmt.Errorf("username is empty")
	}
	if email == "" {
		return fmt.Errorf("email is empty")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("email is not valid")
	}
	// from here https://emaillistvalidation.com/blog/mastering-email-validation-in-golang-crafting-robust-regex-patterns/
	if result, _ := regexp.MatchString("^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", email); !result {
		return fmt.Errorf("email is not valid by new regexp")
	}

	if password != passRepeat {
		return fmt.Errorf("passwords are not the same")
	}

	// TODO Uncomment it
	//if len(password) < app.MinPassLen {
	//	return fmt.Errorf("password len should be a least %d signs", app.MinPassLen)
	//}

	// TODO Change User check
	//if _, ok := users[username]; ok {
	//	return fmt.Errorf("user already exist")
	//}

	return nil
}
