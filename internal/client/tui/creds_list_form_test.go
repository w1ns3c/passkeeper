package tui

import (
	"crypto/md5"
	"encoding/hex"
	"testing"
	"time"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/structs"

	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	sum1 := md5.Sum([]byte("1"))
	sum2 := md5.Sum([]byte("2"))
	sum3 := md5.Sum([]byte("3"))

	creds := []*structs.Credential{
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
		{
			ID:          entities.GenHash("3"),
			Resource:    "wiki.org",
			Login:       "juice",
			Password:    "secinfo",
			Date:        time.Now(),
			Description: "main info from site",
		},
	}

	creds1 := make([]*structs.Credential, len(creds))
	creds2 := make([]*structs.Credential, len(creds))
	creds3 := make([]*structs.Credential, len(creds))

	_ = copy(creds1, creds)
	_ = copy(creds2, creds)
	_ = copy(creds3, creds)

	tmpCred := structs.Credential{
		ID:          hex.EncodeToString(sum3[:]),
		Date:        time.Time{},
		Resource:    "new_example.com",
		Login:       "james",
		Password:    "spiderman",
		Description: "Your favorite neighbour",
	}

	type args struct {
		creds    []*structs.Credential
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
		{
			name: "Check 1: save new the first cred",
			args: args{
				creds:    creds1,
				ind:      0,
				res:      tmpCred.Resource,
				login:    tmpCred.Login,
				password: tmpCred.Password,
				desc:     tmpCred.Description,
			},
			wantErr: false,
		},
		{
			name: "Check 2: save new middle cred",
			args: args{
				creds:    creds2,
				ind:      1,
				res:      tmpCred.Resource,
				login:    tmpCred.Login,
				password: tmpCred.Password,
				desc:     tmpCred.Description,
			},
			wantErr: false,
		},
		{
			name: "Check 3: save new the last cred",
			args: args{
				creds:    creds3,
				ind:      2,
				res:      tmpCred.Resource,
				login:    tmpCred.Login,
				password: tmpCred.Password,
				desc:     tmpCred.Description,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := entities.Save(tt.args.creds, tt.args.ind, tt.args.res, tt.args.login, tt.args.password, tt.args.desc); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			tmp := tt.args.creds[tt.args.ind]
			//require.Equal(t, tt.args.ind, tmp.ID)
			require.Equal(t, tt.args.res, tmp.Resource)
			require.Equal(t, tt.args.login, tmp.Login)
			require.Equal(t, tt.args.password, tmp.Password)
			require.Equal(t, tt.args.desc, tmp.Description)

		})
	}
}

func TestAdd(t *testing.T) {
	cred1 := &structs.Credential{
		ID:          entities.GenHash("1"),
		Resource:    "contoso.local",
		Login:       "username",
		Password:    "SomeSecret",
		Date:        time.Now(),
		Description: "Current Description",
	}

	cred2 := &structs.Credential{
		ID:          entities.GenHash("2"),
		Resource:    "example.com",
		Login:       "mike",
		Password:    "password",
		Date:        time.Now(),
		Description: "Zip zip zip",
	}

	cred3 := &structs.Credential{
		ID:          entities.GenHash("3"),
		Resource:    "wiki.org",
		Login:       "juice",
		Password:    "secinfo",
		Date:        time.Now(),
		Description: "main info from site",
	}

	newCred := &structs.Credential{
		ID:          entities.GenHash("new"),
		Resource:    "new",
		Login:       "new_login",
		Password:    "new_pass",
		Date:        time.Now(),
		Description: "new_description",
	}

	creds0 := []*structs.Credential{}
	creds1 := []*structs.Credential{cred1}
	creds2 := []*structs.Credential{cred1, cred2}
	creds3 := []*structs.Credential{cred1, cred2, cred3}

	newCreds0 := []*structs.Credential{newCred}
	newCreds1 := []*structs.Credential{newCred, cred1}
	newCreds2 := []*structs.Credential{newCred, cred1, cred2}
	newCreds3 := []*structs.Credential{newCred, cred1, cred2, cred3}

	type args struct {
		creds []*structs.Credential
		new   *structs.Credential
	}
	tests := []struct {
		name         string
		args         args
		wantNewCreds []*structs.Credential
		wantErr      bool
	}{
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
