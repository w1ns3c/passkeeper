package cli

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	pb "github.com/w1ns3c/passkeeper/internal/transport/grpc/protofiles/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrPassNotSame = fmt.Errorf("pass and pass repeat are not the same")
)

type Client struct {
	ctx   context.Context
	Token string

	userSvc pb.UserSvcClient
	passSvc pb.UserPassSvcClient
}

func NewClient(ctx context.Context, addr string) (cli *Client, err error) {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		ctx:     ctx,
		userSvc: pb.NewUserSvcClient(conn),
		passSvc: pb.NewUserPassSvcClient(conn),
	}, nil
}

func (c *Client) Login(login, password string) error {
	hash := Hash(password)

	req := &pb.UserLoginRequest{
		Login:    login,
		Password: hash,
	}

	resp, err := c.userSvc.LoginUser(c.ctx, req)
	if err != nil {
		return err
	}

	c.Token = resp.Token

	return nil
}

func (c *Client) Register(login, password, repeat, email string) error {
	hash1 := Hash(password)
	hash2 := Hash(repeat)
	if hash1 != hash2 {
		return ErrPassNotSame
	}

	req := &pb.UserRegisterRequest{
		Login:      login,
		Password:   hash1,
		RePassword: hash2,
	}

	resp, err := c.userSvc.RegisterUser(c.ctx, req)
	if err != nil {
		return err
	}

	c.Token = resp.Token

	return nil
}

func (c *Client) ChangePass(curPass, newPass, repeat string) error {
	oldHash := Hash(curPass)
	newHash := Hash(newPass)
	repeatH := Hash(repeat)

	if newHash != repeatH {
		return ErrPassNotSame
	}

	req := &pb.UserChangePassReq{
		OldPass: oldHash,
		NewPass: newPass,
		Repeat:  repeatH,
	}

	_, err := c.passSvc.ChangePass(c.ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func Hash(password string) string {
	h := sha512.Sum512([]byte(password))
	return hex.EncodeToString(h[:])
}
