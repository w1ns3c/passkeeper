package hashes

import (
	"github.com/stretchr/testify/require"
	"passkeeper/internal/config"
	"passkeeper/internal/entities"
	"testing"
	"time"
)

func TestEncryptDecryptBlob(t *testing.T) {
	type args struct {
		cred *entities.Credential
		key  string
	}

	var (
		pass   = "123"
		hash   = Hash(pass)
		userID = "new-user-id"
		s, _   = GenerateSecret(config.UserSecretLen)
		ss, _  = EncryptSecret(s, hash)
		key, _ = GenerateCredsSecret(pass, userID, ss)

		cred1 = &entities.Credential{
			ID:          GeneratePassID2(),
			Date:        time.Now(),
			Resource:    "res1",
			Login:       "login1",
			Password:    "pass1",
			Description: "simple description",
		}
		cred2 = &entities.Credential{
			ID:          GeneratePassID2(),
			Date:        time.Now(),
			Resource:    "res1",
			Login:       "login1",
			Password:    "pass1",
			Description: "Long description -----------------------------------------------------------------------------------------------------\nsidfosdouioewaifsdjsdljfalkfdjalkdsjfklasjdfuiahsdfiua\nsdjfsiodfoiwueroisj sdajkfalkj-*(@(HIUH jsdfkldsfkj",
		}
	)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Simple1",
			args: args{
				cred: cred1,
				key:  key,
			},
			wantErr: false,
		}, {
			name: "Long Description",
			args: args{
				cred: cred2,
				key:  key,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBlob, err := EncryptBlob(tt.args.cred, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptBlob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotCred, err := DecryptBlob(gotBlob, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptBlob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.args.cred.ID, gotCred.ID, "Encrypt/DecryptBlob() creds ID are not the same")

			require.Equal(t, tt.args.cred.Login, gotCred.Login, "Encrypt/DecryptBlob() creds Login are not the same")

			require.Equal(t, tt.args.cred.Password, gotCred.Password, "Encrypt/DecryptBlob() creds Password are not the same")

			require.Equal(t, tt.args.cred.Description, gotCred.Description, "Encrypt/DecryptBlob() creds Description are not the same")
		})
	}
}
