package tui

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"passkeeper/internal/entities"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	sum1 := md5.Sum([]byte("1"))
	sum2 := md5.Sum([]byte("2"))
	sum3 := md5.Sum([]byte("3"))

	creds := []entities.Credential{
		{
			ID:          hex.EncodeToString(sum1[:]),
			Resource:    "contoso.local",
			Login:       "username",
			Password:    "SomeSecret",
			Date:        time.Now(),
			Description: "Current Description",
		},
		{
			ID:          hex.EncodeToString(sum2[:]),
			Resource:    "example.com",
			Login:       "mike",
			Password:    "password",
			Date:        time.Now(),
			Description: "Zip zip zip",
		},
	}

	tmpCred := entities.Credential{
		ID:          hex.EncodeToString(sum3[:]),
		Date:        time.Time{},
		Resource:    "new_example.com",
		Login:       "james",
		Password:    "spiderman",
		Description: "Your favorite neighbour",
	}

	t.Log(creds)
	t.Log("123")
	err := entities.Save(creds, 0, tmpCred.Resource, tmpCred.Login, tmpCred.Password, tmpCred.Description)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(creds)
	//return

	type args struct {
		creds    []entities.Credential
		ind      int
		res      string
		login    string
		password string
		desc     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := entities.Save(tt.args.creds, tt.args.ind, tt.args.res, tt.args.login, tt.args.password, tt.args.desc); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	cred1 := entities.Credential{
		ID:          entities.GenHash("1"),
		Resource:    "contoso.local",
		Login:       "username",
		Password:    "SomeSecret",
		Date:        time.Now(),
		Description: "Current Description",
	}

	cred2 := entities.Credential{
		ID:          entities.GenHash("2"),
		Resource:    "example.com",
		Login:       "mike",
		Password:    "password",
		Date:        time.Now(),
		Description: "Zip zip zip",
	}

	cred3 := entities.Credential{
		ID:          entities.GenHash("3"),
		Resource:    "wiki.org",
		Login:       "juice",
		Password:    "secinfo",
		Date:        time.Now(),
		Description: "main info from site",
	}

	newCred := entities.Credential{
		ID:          entities.GenHash("new"),
		Resource:    "new",
		Login:       "new_login",
		Password:    "new_pass",
		Date:        time.Now(),
		Description: "new_description",
	}

	creds0 := []entities.Credential{}
	creds1 := []entities.Credential{cred1}
	creds2 := []entities.Credential{cred1, cred2}
	creds3 := []entities.Credential{cred1, cred2, cred3}

	newCreds0 := []entities.Credential{newCred}
	newCreds1 := []entities.Credential{newCred, cred1}
	newCreds2 := []entities.Credential{newCred, cred1, cred2}
	newCreds3 := []entities.Credential{newCred, cred1, cred2, cred3}

	type args struct {
		creds []entities.Credential
		new   entities.Credential
	}
	tests := []struct {
		name         string
		args         args
		wantNewCreds []entities.Credential
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name:         "No elements",
			args:         args{creds: creds0, new: newCred},
			wantNewCreds: newCreds0,
		},
		{
			name:         "One element",
			args:         args{creds: creds1, new: newCred},
			wantNewCreds: newCreds1,
		},
		{
			name:         "Two elements",
			args:         args{creds: creds2, new: newCred},
			wantNewCreds: newCreds2,
		},

		{
			name:         "Many elements",
			args:         args{creds: creds3, new: newCred},
			wantNewCreds: newCreds3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewCreds, err := entities.Add(tt.args.creds, tt.args.new.Resource, tt.args.new.Login, tt.args.new.Password, tt.args.new.Description)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotNewCreds) != len(tt.wantNewCreds) {
				t.Errorf("Add() len(gotNewCreds) = %v, want %v", len(gotNewCreds), len(tt.wantNewCreds))
			}
			for ind := 0; ind < len(gotNewCreds); ind++ {
				if gotNewCreds[ind].ID != tt.wantNewCreds[ind].ID {
					t.Errorf("Add() gotNewCreds[%d] = %v, want %v", ind, gotNewCreds[ind].ID, tt.wantNewCreds[ind].ID)
				}
				if gotNewCreds[ind].Resource != tt.wantNewCreds[ind].Resource {
					t.Errorf("Add() gotNewCreds[%d] = %v, want %v", ind, gotNewCreds[ind].Resource, tt.wantNewCreds[ind].Resource)
				}
				if gotNewCreds[ind].Login != tt.wantNewCreds[ind].Login {
					t.Errorf("Add() gotNewCreds[%d] = %v, want %v", ind, gotNewCreds[ind].Login, tt.wantNewCreds[ind].Login)
				}
				if gotNewCreds[ind].Password != tt.wantNewCreds[ind].Password {
					t.Errorf("Add() gotNewCreds[%d] = %v, want %v", ind, gotNewCreds[ind].Password, tt.wantNewCreds[ind].Password)
				}
				if gotNewCreds[ind].Description != tt.wantNewCreds[ind].Description {
					t.Errorf("Add() gotNewCreds[%d] = %v, want %v", ind, gotNewCreds[ind].Description, tt.wantNewCreds[ind].Description)
				}
			}

			//if !reflect.DeepEqual(gotNewCreds, tt.wantNewCreds) {
			//	t.Errorf("Add() gotNewCreds = %v, want %v", gotNewCreds, tt.wantNewCreds)
			//}
		})
	}
}
