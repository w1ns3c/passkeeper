package hashes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/structs"
)

func TestEncryptDecryptBlob(t *testing.T) {
	type args struct {
		cred structs.CredInf
		key  string
	}

	var (
		pass   = "123"
		hash   = Hash(pass)
		userID = "new-user-id"
		s, _   = GenerateSecret(config.UserSecretLen)
		ss, _  = EncryptSecret(s, hash)
		key, _ = GenerateCredsSecret(pass, userID, ss)

		cred1 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          GeneratePassID(),
			Date:        time.Now(),
			Resource:    "res1",
			Login:       "login1",
			Password:    "pass1",
			Description: "simple description",
		}
		cred2 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          GeneratePassID(),
			Date:        time.Now(),
			Resource:    "res1",
			Login:       "login1",
			Password:    "pass1",
			Description: "Long description -----------------------------------------------------------------------------------------------------\nsidfosdouioewaifsdjsdljfalkfdjalkdsjfklasjdfuiahsdfiua\nsdjfsiodfoiwueroisj sdajkfalkj-*(@(HIUH jsdfkldsfkj",
		}

		testCards = []*structs.Card{
			{
				Type:        structs.BlobCard,
				Name:        "test1",
				Bank:        entities.Banks[0],
				Person:      "string",
				Number:      122222222222,
				CVC:         232,
				Expiration:  time.Now().Add(time.Second * -19999),
				PIN:         3333,
				Description: "test description only",
			},
			{
				Type: structs.BlobCard,
				Name: "test333331",
			},
		}

		testNotes = []*structs.Note{
			{
				Type: structs.BlobNote,
				Name: "New Test Blob",
				Body: "Hello\nWorld! Amigo",
			},
			{
				Type: structs.BlobNote,
			},
			{},
		}
	)

	tests := []struct {
		name        string
		args        args
		wantErrEnc  bool
		wantErrDecr bool
	}{
		// TODO: Add test cases.
		{
			name: "Creds: Simple1",
			args: args{
				cred: cred1,
				key:  key,
			},
			wantErrEnc: false,
		},
		{
			name: "Creds: Long Description",
			args: args{
				cred: cred2,
				key:  key,
			},
			wantErrEnc: false,
		},
		{
			name: "Cards: Valid",
			args: args{
				cred: testCards[0],
				key:  key,
			},
			wantErrEnc: false,
		},
		{
			name: "Cards: Empty fields",
			args: args{
				cred: cred2,
				key:  key,
			},
			wantErrEnc: false,
		},
		{
			name: "Notes: Valid",
			args: args{
				cred: testNotes[0],
				key:  key,
			},
			wantErrEnc: false,
		},
		{
			name: "Notes: Empty fields",
			args: args{
				cred: testNotes[1],
				key:  key,
			},
			wantErrEnc: false,
		},
		{
			name: "Notes: All Empty fields",
			args: args{
				cred: testNotes[2],
				key:  key,
			},
			wantErrEnc:  false,
			wantErrDecr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBlob, err := EncryptBlob(tt.args.cred, tt.args.key)
			if (err != nil) != tt.wantErrEnc {
				t.Errorf("EncryptBlob() error = %v, wantErr %v", err, tt.wantErrEnc)
				return
			}

			got, err := DecryptBlob(gotBlob, tt.args.key)
			if (err != nil) != tt.wantErrDecr {
				t.Errorf("DecryptBlob() error = %v, wantErr %v", err, tt.wantErrDecr)
				return
			}

			switch got.(type) {
			case *structs.Credential:
				gotCred := got.(*structs.Credential)
				putCred := tt.args.cred.(*structs.Credential)
				require.Equal(t, putCred.ID, gotCred.ID, "Encrypt/DecryptBlob() creds ID are not the same")
				//require.Equal(t, putCred.Date, gotCred.Date, "Encrypt/DecryptBlob() creds Date are not the same")
				require.Equal(t, putCred.Resource, gotCred.Resource, "Encrypt/DecryptBlob() creds Resource are not the same")
				require.Equal(t, putCred.Login, gotCred.Login, "Encrypt/DecryptBlob() creds Login are not the same")
				require.Equal(t, putCred.Password, gotCred.Password, "Encrypt/DecryptBlob() creds Password are not the same")
				require.Equal(t, putCred.Description, gotCred.Description, "Encrypt/DecryptBlob() creds Description are not the same")

			case *structs.Card:
				gotCard := got.(*structs.Card)
				putCard := tt.args.cred.(*structs.Card)

				require.Equal(t, putCard.ID, gotCard.ID, "Encrypt/DecryptBlob() card ID are not the same")
				require.Equal(t, putCard.Name, gotCard.Name, "Encrypt/DecryptBlob() card Name are not the same")
				require.Equal(t, putCard.Bank, gotCard.Bank, "Encrypt/DecryptBlob() card Bank are not the same")
				require.Equal(t, putCard.Person, gotCard.Person, "Encrypt/DecryptBlob() card Person are not the same")
				require.Equal(t, putCard.Number, gotCard.Number, "Encrypt/DecryptBlob() card Bank are not the same")
				require.Equal(t, putCard.CVC, gotCard.CVC, "Encrypt/DecryptBlob() card CVC are not the same")
				require.Equal(t, putCard.Expiration.Format(time.DateTime), gotCard.Expiration.Format(time.DateTime), "Encrypt/DecryptBlob() card Expiration are not the same")
				require.Equal(t, putCard.PIN, gotCard.PIN, "Encrypt/DecryptBlob() card PIN are not the same")
				require.Equal(t, putCard.Description, gotCard.Description, "Encrypt/DecryptBlob() card Description are not the same")

			case *structs.Note:
				gotNote := got.(*structs.Note)
				putNote := tt.args.cred.(*structs.Note)

				require.Equal(t, putNote.ID, gotNote.ID, "Encrypt/DecryptBlob() note ID are not the same")
				require.Equal(t, putNote.Name, gotNote.Name, "Encrypt/DecryptBlob() note Name are not the same")
				//require.Equal(t, putNote.Date, gotNote.Date, "Encrypt/DecryptBlob() note Date are not the same")
				require.Equal(t, putNote.Body, gotNote.Body, "Encrypt/DecryptBlob() note Body are not the same")
			case *structs.File:
				gotFile := got.(*structs.File)
				putNote := tt.args.cred.(*structs.File)

				require.Equal(t, putNote.ID, gotFile.ID, "Encrypt/DecryptBlob() file ID are not the same")
				require.Equal(t, putNote.Name, gotFile.Name, "Encrypt/DecryptBlob() file Name are not the same")
				require.Equal(t, putNote.Type, gotFile.Type, "Encrypt/DecryptBlob() file Type are not the same")
				require.Equal(t, putNote.Body, gotFile.Body, "Encrypt/DecryptBlob() file Body are not the same")

			}

		})
	}
}
