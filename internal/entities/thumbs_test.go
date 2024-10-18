package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"passkeeper/internal/entities/structs"
)

func TestAddCard(t *testing.T) {
	type args struct {
		cards   []*structs.Card
		newCard *structs.Card
	}
	var (
		cards = []*structs.Card{
			&structs.Card{
				Type:        structs.BlobCard,
				ID:          "123123",
				Name:        "name",
				Bank:        "bank",
				Person:      "person",
				Number:      123,
				CVC:         4432,
				Expiration:  time.Now(),
				PIN:         111,
				Description: "some description",
			},
		}
		card1 = &structs.Card{
			Type:        structs.BlobCard,
			ID:          "so-------me-new",
			Name:        "name123",
			Bank:        "bank444",
			Person:      "person333",
			Number:      111122223333444,
			CVC:         4432,
			Expiration:  time.Now(),
			PIN:         111,
			Description: "som description",
		}
	)

	tests := []struct {
		name         string
		args         args
		wantNewCards []*structs.Card
		wantErr      bool
	}{
		{
			name: "Test 1: add card",
			args: args{
				cards:   cards,
				newCard: card1,
			},
			wantNewCards: []*structs.Card{card1, cards[0]},
			wantErr:      false,
		},
		{
			name: "Test 2: add card",
			args: args{
				cards:   []*structs.Card{},
				newCard: card1,
			},
			wantNewCards: []*structs.Card{card1},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewCards, err := AddCard(tt.args.cards, tt.args.newCard)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotNewCards) != len(tt.wantNewCards) {
				t.Errorf("AddCard() len gotNewCards = %v, want %v", len(gotNewCards), len(tt.wantNewCards))
				return
			}
			for ind := 0; ind < len(tt.wantNewCards); ind++ {
				require.Equal(t, tt.wantNewCards[ind].Type, gotNewCards[ind].Type)
				require.Equal(t, tt.wantNewCards[ind].ID, gotNewCards[ind].ID)
				require.Equal(t, tt.wantNewCards[ind].Name, gotNewCards[ind].Name)
				require.Equal(t, tt.wantNewCards[ind].Bank, gotNewCards[ind].Bank)
				require.Equal(t, tt.wantNewCards[ind].Person, gotNewCards[ind].Person)
				require.Equal(t, tt.wantNewCards[ind].Number, gotNewCards[ind].Number)
				require.Equal(t, tt.wantNewCards[ind].CVC, gotNewCards[ind].CVC)
				require.Equal(t, tt.wantNewCards[ind].Expiration, gotNewCards[ind].Expiration)
				require.Equal(t, tt.wantNewCards[ind].PIN, gotNewCards[ind].PIN)
				require.Equal(t, tt.wantNewCards[ind].Description, gotNewCards[ind].Description)

			}

		})
	}
}

func TestAddCred(t *testing.T) {
	type args struct {
		creds   []*structs.Credential
		newCard *structs.Credential
	}
	var (
		creds = []*structs.Credential{
			{
				Type: structs.BlobCard,
				ID:   "123123",
				Date: time.Now(),

				Login:       "name",
				Password:    "person",
				Resource:    "bank",
				Description: "some description",
			},
		}
		cred1 = &structs.Credential{
			Type:        structs.BlobCard,
			ID:          "so-------me-new",
			Login:       "name123",
			Password:    "bank444",
			Resource:    "person333",
			Date:        time.Now(),
			Description: "som description",
		}
	)

	tests := []struct {
		name      string
		args      args
		wantCreds []*structs.Credential
		wantErr   bool
	}{
		{
			name: "Test 1: add cred",
			args: args{
				creds:   creds,
				newCard: cred1,
			},
			wantCreds: []*structs.Credential{cred1, creds[0]},
			wantErr:   false,
		},
		{
			name: "Test 2: add cred",
			args: args{
				creds:   []*structs.Credential{},
				newCard: cred1,
			},
			wantCreds: []*structs.Credential{cred1},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddCred(tt.args.creds, tt.args.newCard)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.wantCreds) {
				t.Errorf("AddCard() len got = %v, want %v", len(got), len(tt.wantCreds))
				return
			}
			for ind := 0; ind < len(tt.wantCreds); ind++ {
				require.Equal(t, tt.wantCreds[ind].Type, got[ind].Type)
				require.Equal(t, tt.wantCreds[ind].ID, got[ind].ID)
				require.Equal(t, tt.wantCreds[ind].Login, got[ind].Login)
				require.Equal(t, tt.wantCreds[ind].Password, got[ind].Password)
				require.Equal(t, tt.wantCreds[ind].Resource, got[ind].Resource)
				require.Equal(t, tt.wantCreds[ind].Date, got[ind].Date)
				require.Equal(t, tt.wantCreds[ind].Description, got[ind].Description)

			}

		})
	}
}

func TestAddFile(t *testing.T) {
	type args struct {
		Files   []*structs.File
		newCard *structs.File
	}
	var (
		Files = []*structs.File{
			&structs.File{
				Type: structs.BlobCard,
				ID:   "123123",
				Body: []byte("some description"),
				Name: "name",
			},
		}
		File1 = &structs.File{
			Type: structs.BlobCard,
			ID:   "so-------me-new",
			Name: "name123",
			Body: []byte("bank444"),
		}
	)

	tests := []struct {
		name      string
		args      args
		wantFiles []*structs.File
		wantErr   bool
	}{
		{
			name: "Test 1: add file",
			args: args{
				Files:   Files,
				newCard: File1,
			},
			wantFiles: []*structs.File{File1, Files[0]},
			wantErr:   false,
		},
		{
			name: "Test 2: add  file",
			args: args{
				Files:   []*structs.File{},
				newCard: File1,
			},
			wantFiles: []*structs.File{File1},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddFile(tt.args.Files, tt.args.newCard)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.wantFiles) {
				t.Errorf("AddCard() len got = %v, want %v", len(got), len(tt.wantFiles))
				return
			}
			for ind := 0; ind < len(tt.wantFiles); ind++ {
				require.Equal(t, tt.wantFiles[ind].Type, got[ind].Type)
				require.Equal(t, tt.wantFiles[ind].ID, got[ind].ID)
				require.Equal(t, tt.wantFiles[ind].Name, got[ind].Name)
				require.Equal(t, tt.wantFiles[ind].Body, got[ind].Body)

			}

		})
	}
}

func TestAddNote(t *testing.T) {
	type args struct {
		Notes   []*structs.Note
		newCard *structs.Note
	}
	var (
		Notes = []*structs.Note{
			{
				Type: structs.BlobCard,
				ID:   "123123",
				Date: time.Now(),
				Name: "name",
				Body: "person",
			},
		}
		Note1 = &structs.Note{
			Type: structs.BlobCard,
			ID:   "so-------me-new",
			Name: "name123",
			Date: time.Now(),
			Body: "som description",
		}
	)

	tests := []struct {
		name      string
		args      args
		wantNotes []*structs.Note
		wantErr   bool
	}{
		{
			name: "Test 1: add note",
			args: args{
				Notes:   Notes,
				newCard: Note1,
			},
			wantNotes: []*structs.Note{Note1, Notes[0]},
			wantErr:   false,
		},
		{
			name: "Test 2: add note",
			args: args{
				Notes:   []*structs.Note{},
				newCard: Note1,
			},
			wantNotes: []*structs.Note{Note1},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddNote(tt.args.Notes, tt.args.newCard)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.wantNotes) {
				t.Errorf("AddCard() len got = %v, want %v", len(got), len(tt.wantNotes))
				return
			}
			for ind := 0; ind < len(tt.wantNotes); ind++ {
				require.Equal(t, tt.wantNotes[ind].Type, got[ind].Type)
				require.Equal(t, tt.wantNotes[ind].ID, got[ind].ID)
				require.Equal(t, tt.wantNotes[ind].Name, got[ind].Name)
				require.Equal(t, tt.wantNotes[ind].Date, got[ind].Date)
				require.Equal(t, tt.wantNotes[ind].Body, got[ind].Body)

			}

		})
	}
}

func TestDeleteCard(t *testing.T) {
	type args struct {
		Cards []*structs.Card
		ind   int
	}
	var (
		Cards = []*structs.Card{
			{
				Type:        structs.BlobCard,
				ID:          "123123",
				Name:        "name",
				Bank:        "bank",
				Person:      "person",
				Number:      123,
				CVC:         4432,
				Expiration:  time.Now(),
				PIN:         111,
				Description: "some description",
			},
			{
				Type:        structs.BlobCard,
				ID:          "so-------me-newffffffffffffff",
				Name:        "name123fffffffffffff",
				Bank:        "bank444fffffffffffff",
				Person:      "person333fffffffffffff",
				Number:      123444,
				CVC:         3333,
				Expiration:  time.Now().Add(time.Hour * -10),
				PIN:         666,
				Description: "som description",
			},
		}
		Cards2 = []*structs.Card{
			{
				Type:        structs.BlobCard,
				ID:          "123123",
				Name:        "name",
				Bank:        "bank",
				Person:      "person",
				Number:      123,
				CVC:         4432,
				Expiration:  time.Now(),
				PIN:         111,
				Description: "some description",
			},
			{
				Type:        structs.BlobCard,
				ID:          "so-------me-newffffffffffffff",
				Name:        "name123fffffffffffff",
				Bank:        "bank444fffffffffffff",
				Person:      "person333fffffffffffff",
				Number:      123444,
				CVC:         3333,
				Expiration:  time.Now().Add(time.Hour * -10),
				PIN:         666,
				Description: "som description",
			},
		}
	)

	tests := []struct {
		name      string
		args      args
		wantCards []*structs.Card
		wantErr   bool
	}{
		{
			name: "Test 1: delete Card (the first)",
			args: args{
				Cards: Cards,
				ind:   0,
			},
			wantCards: []*structs.Card{
				{
					Type:        structs.BlobCard,
					ID:          Cards[1].ID,
					Name:        Cards[1].Name,
					Bank:        Cards[1].Bank,
					Person:      Cards[1].Person,
					Number:      Cards[1].Number,
					CVC:         Cards[1].CVC,
					Expiration:  Cards[1].Expiration,
					PIN:         Cards[1].PIN,
					Description: Cards[1].Description,
				},
			},
			wantErr: false,
		},
		{
			name: "Test 2: delete Card (the last)",
			args: args{
				Cards: Cards2,
				ind:   1,
			},
			wantCards: []*structs.Card{
				{
					Type:        structs.BlobCard,
					ID:          Cards[0].ID,
					Name:        Cards[0].Name,
					Bank:        Cards[0].Bank,
					Person:      Cards[0].Person,
					Number:      Cards[0].Number,
					CVC:         Cards[0].CVC,
					Expiration:  Cards[0].Expiration,
					PIN:         Cards[0].PIN,
					Description: Cards[0].Description,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCards, err := DeleteCard(tt.args.Cards, tt.args.ind)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.wantCards) != len(gotCards) {
				t.Errorf("AddCard() len tt.args.Cards = %v, want %v", len(tt.wantCards), len(gotCards))
				return
			}
			for ind := 0; ind < len(gotCards); ind++ {
				require.Equal(t, gotCards[ind].Type, tt.wantCards[ind].Type)
				require.Equal(t, gotCards[ind].ID, tt.wantCards[ind].ID)
				require.Equal(t, gotCards[ind].Name, tt.wantCards[ind].Name)
				require.Equal(t, gotCards[ind].Expiration.Format(time.DateTime), tt.wantCards[ind].Expiration.Format(time.DateTime))
				require.Equal(t, gotCards[ind].Description, tt.wantCards[ind].Description)

			}

		})
	}
}

func TestDeleteCred(t *testing.T) {
	type args struct {
		Credentials []*structs.Credential
		ind         int
	}
	var (
		Creds = []*structs.Credential{
			{
				Type:        structs.BlobCred,
				ID:          "123123",
				Login:       "name",
				Password:    "bank",
				Resource:    "person",
				Date:        time.Now(),
				Description: "some description",
			},
			{
				Type:        structs.BlobCred,
				ID:          "so-------me-newffffffffffffff",
				Login:       "name123fffffffffffff",
				Password:    "bank444fffffffffffff",
				Resource:    "person333fffffffffffff",
				Date:        time.Now().Add(time.Hour * -10),
				Description: "som description",
			},
		}
		Creds2 = []*structs.Credential{
			{
				Type:        structs.BlobCred,
				ID:          "123123",
				Login:       "name",
				Password:    "bank",
				Resource:    "person",
				Date:        time.Now(),
				Description: "some description",
			},
			{
				Type:        structs.BlobCred,
				ID:          "so-------me-newffffffffffffff",
				Login:       "name123fffffffffffff",
				Password:    "bank444fffffffffffff",
				Resource:    "person333fffffffffffff",
				Date:        time.Now().Add(time.Hour * -10),
				Description: "som description",
			},
		}
	)

	tests := []struct {
		name            string
		args            args
		wantCredentials []*structs.Credential
		wantErr         bool
	}{
		{
			name: "Test 1: delete Credential (the first)",
			args: args{
				Credentials: Creds,
				ind:         0,
			},
			wantCredentials: []*structs.Credential{
				Creds[1],
			},
			wantErr: false,
		},
		{
			name: "Test 2: delete Credential (the last)",
			args: args{
				Credentials: Creds2,
				ind:         1,
			},
			wantCredentials: []*structs.Credential{
				Creds2[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCredentials, err := DeleteCred(tt.args.Credentials, tt.args.ind)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCredential() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.wantCredentials) != len(gotCredentials) {
				t.Errorf("AddCredential() len tt.args.Credentials = %v, want %v", len(tt.wantCredentials), len(gotCredentials))
				return
			}
			for ind := 0; ind < len(gotCredentials); ind++ {
				require.Equal(t, gotCredentials[ind].Type, tt.wantCredentials[ind].Type)
				require.Equal(t, gotCredentials[ind].ID, tt.wantCredentials[ind].ID)
				require.Equal(t, gotCredentials[ind].Login, tt.wantCredentials[ind].Login)
				require.Equal(t, gotCredentials[ind].Password, tt.wantCredentials[ind].Password)
				require.Equal(t, gotCredentials[ind].Resource, tt.wantCredentials[ind].Resource)
				require.Equal(t, gotCredentials[ind].Date.Format(time.DateTime), tt.wantCredentials[ind].Date.Format(time.DateTime))
				require.Equal(t, gotCredentials[ind].Description, tt.wantCredentials[ind].Description)

			}

		})
	}
}

func TestDeleteFile(t *testing.T) {
	type args struct {
		Files []*structs.File
		ind   int
	}
	var (
		Creds = []*structs.File{
			{
				Type: structs.BlobCred,
				ID:   "123123",
				Name: "name",
				Body: []byte("bank"),
			},
			{
				Type: structs.BlobCred,
				ID:   "so-------me-newffffffffffffff",
				Name: "name123fffffffffffff",
				Body: []byte("bank444fffffffffffff"),
			},
		}
		Creds2 = []*structs.File{
			{
				Type: structs.BlobCred,
				ID:   "123123",
				Name: "name",
				Body: []byte("bank"),
			},
			{
				Type: structs.BlobCred,
				ID:   "so-------me-newffffffffffffff",
				Name: "name123fffffffffffff",
				Body: []byte("bank444fffffffffffff"),
			},
		}
	)

	tests := []struct {
		name      string
		args      args
		wantFiles []*structs.File
		wantErr   bool
	}{
		{
			name: "Test 1: delete File (the first)",
			args: args{
				Files: Creds,
				ind:   0,
			},
			wantFiles: []*structs.File{
				Creds[1],
			},
			wantErr: false,
		},
		{
			name: "Test 2: delete File (the last)",
			args: args{
				Files: Creds2,
				ind:   1,
			},
			wantFiles: []*structs.File{
				Creds2[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiles, err := DeleteFile(tt.args.Files, tt.args.ind)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.wantFiles) != len(gotFiles) {
				t.Errorf("AddFile() len tt.args.Files = %v, want %v", len(tt.wantFiles), len(gotFiles))
				return
			}
			for ind := 0; ind < len(gotFiles); ind++ {
				require.Equal(t, gotFiles[ind].Type, tt.wantFiles[ind].Type)
				require.Equal(t, gotFiles[ind].ID, tt.wantFiles[ind].ID)
				require.Equal(t, gotFiles[ind].Body, tt.wantFiles[ind].Body)
				require.Equal(t, gotFiles[ind].Name, tt.wantFiles[ind].Name)
			}

		})
	}
}

func TestDeleteNote(t *testing.T) {
	type args struct {
		Notes []*structs.Note
		ind   int
	}
	var (
		Notes = []*structs.Note{
			{
				Type: structs.BlobNote,
				ID:   "123123",
				Name: "name",
				Body: "bank",
			},
			{
				Type: structs.BlobNote,
				ID:   "so-------me-newffffffffffffff",
				Name: "name123fffffffffffff",
				Body: "bank444fffffffffffff",
			},
		}
		Notes2 = []*structs.Note{
			{
				Type: structs.BlobNote,
				ID:   "123123",
				Name: "name",
				Body: "bank",
			},
			{
				Type: structs.BlobNote,
				ID:   "so-------me-newffffffffffffff",
				Name: "name123fffffffffffff",
				Body: "bank444fffffffffffff",
			},
		}
	)

	tests := []struct {
		name      string
		args      args
		wantNotes []*structs.Note
		wantErr   bool
	}{
		{
			name: "Test 1: delete Note (the first)",
			args: args{
				Notes: Notes,
				ind:   0,
			},
			wantNotes: []*structs.Note{
				Notes[1],
			},
			wantErr: false,
		},
		{
			name: "Test 2: delete Note (the last)",
			args: args{
				Notes: Notes2,
				ind:   1,
			},
			wantNotes: []*structs.Note{
				Notes2[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotes, err := DeleteNote(tt.args.Notes, tt.args.ind)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.wantNotes) != len(gotNotes) {
				t.Errorf("AddNote() len tt.args.Notes = %v, want %v", len(tt.wantNotes), len(gotNotes))
				return
			}
			for ind := 0; ind < len(gotNotes); ind++ {
				require.Equal(t, gotNotes[ind].Type, tt.wantNotes[ind].Type)
				require.Equal(t, gotNotes[ind].ID, tt.wantNotes[ind].ID)
				require.Equal(t, gotNotes[ind].Body, tt.wantNotes[ind].Body)
				require.Equal(t, gotNotes[ind].Name, tt.wantNotes[ind].Name)
			}

		})
	}
}

func TestGenHash(t *testing.T) {
	tests := []struct {
		name string
		sec  string
		hash string
	}{
		{
			name: "Test GenHash",
			sec:  "123456",
			hash: "e10adc3949ba59abbe56e057f20f883e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenHash(tt.sec); got != tt.hash {
				t.Errorf("GenHash() = %v, want %v", got, tt.hash)
			}
		})
	}
}

func TestSaveCard(t *testing.T) {
	type args struct {
		Cards   []*structs.Card
		newCard *structs.Card
		ind     int
	}
	var (
		Cards = []*structs.Card{
			{
				Type:        structs.BlobCard,
				ID:          "123123",
				Name:        "name",
				Bank:        "bank",
				Person:      "person",
				Number:      123,
				CVC:         4432,
				Expiration:  time.Now(),
				PIN:         111,
				Description: "some description",
			},
			{
				Type:        structs.BlobCard,
				ID:          "so-------me-newffffffffffffff",
				Name:        "name123fffffffffffff",
				Bank:        "bank444fffffffffffff",
				Person:      "person333fffffffffffff",
				Number:      123444,
				CVC:         3333,
				Expiration:  time.Now().Add(time.Hour * -10),
				PIN:         666,
				Description: "som description",
			},
		}
		Cards2 = []*structs.Card{
			{
				Type:        structs.BlobCard,
				ID:          "123123",
				Name:        "name",
				Bank:        "bank",
				Person:      "person",
				Number:      123,
				CVC:         4432,
				Expiration:  time.Now(),
				PIN:         111,
				Description: "some description",
			},
			{
				Type:        structs.BlobCard,
				ID:          "so-------me-newffffffffffffff",
				Name:        "name123fffffffffffff",
				Bank:        "bank444fffffffffffff",
				Person:      "person333fffffffffffff",
				Number:      123444,
				CVC:         3333,
				Expiration:  time.Now().Add(time.Hour * -10),
				PIN:         666,
				Description: "som description",
			},
		}
		Card1 = &structs.Card{
			Type:        structs.BlobCard,
			ID:          "so-------me-new",
			Name:        "name123",
			Bank:        "bank444",
			Person:      "person333",
			Number:      111122223333444,
			CVC:         4432,
			Expiration:  time.Now(),
			PIN:         111,
			Description: "som description",
		}
	)

	tests := []struct {
		name      string
		args      args
		wantCards []*structs.Card
		wantErr   bool
	}{
		{
			name: "Test 1: save Card (the first)",
			args: args{
				Cards:   Cards,
				newCard: Card1,
				ind:     0,
			},
			wantCards: []*structs.Card{Card1, Cards[1]},
			wantErr:   false,
		},
		{
			name: "Test 2: save Card (the last)",
			args: args{
				Cards:   Cards2,
				newCard: Card1,
				ind:     1,
			},
			wantCards: []*structs.Card{Cards2[0], Card1},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveCard(tt.args.Cards, tt.args.ind, tt.args.newCard)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.args.Cards) != len(tt.wantCards) {
				t.Errorf("AddCard() len tt.args.Cards = %v, want %v", len(tt.args.Cards), len(tt.wantCards))
				return
			}
			for ind := 0; ind < len(tt.wantCards); ind++ {
				require.Equal(t, tt.wantCards[ind].Type, tt.args.Cards[ind].Type)
				require.Equal(t, tt.wantCards[ind].ID, tt.args.Cards[ind].ID)
				require.Equal(t, tt.wantCards[ind].Name, tt.args.Cards[ind].Name)
				require.Equal(t, tt.wantCards[ind].Expiration, tt.args.Cards[ind].Expiration)
				require.Equal(t, tt.wantCards[ind].Description, tt.args.Cards[ind].Description)

			}

		})
	}
}

func TestSaveCred(t *testing.T) {
	type args struct {
		Credentials   []*structs.Credential
		newCredential *structs.Credential
		ind           int
	}
	var (
		Credentials = []*structs.Credential{
			{
				Type:        structs.BlobCard,
				ID:          "123123",
				Date:        time.Now(),
				Login:       "name",
				Password:    "person",
				Resource:    "bank",
				Description: "some description",
			},
			{
				Type:        structs.BlobCard,
				ID:          "ggggggggggggg",
				Date:        time.Now().Add(time.Hour * -10),
				Login:       "111111111111",
				Password:    "22222222222",
				Resource:    "ban11111k",
				Description: "some descriptio2222222222n",
			},
		}
		Credentials3 = []*structs.Credential{
			{
				Type:        structs.BlobCard,
				ID:          "123123",
				Date:        time.Now(),
				Login:       "name",
				Password:    "person",
				Resource:    "bank",
				Description: "some description",
			},
			{
				Type:        structs.BlobCard,
				ID:          "ggggggggggggg",
				Date:        time.Now().Add(time.Hour * -10),
				Login:       "111111111111",
				Password:    "22222222222",
				Resource:    "ban11111k",
				Description: "some descriptio2222222222n",
			},
		}
		cred1 = &structs.Credential{
			Type:        structs.BlobCard,
			ID:          "so-------me-new",
			Login:       "name123",
			Password:    "bank444",
			Resource:    "person333",
			Date:        time.Now(),
			Description: "som description",
		}
	)

	tests := []struct {
		name            string
		args            args
		wantCredentials []*structs.Credential
		wantErr         bool
	}{
		{
			name: "Test 1: save Credential (the first)",
			args: args{
				Credentials:   Credentials,
				newCredential: cred1,
				ind:           0,
			},
			wantCredentials: []*structs.Credential{cred1, Credentials[1]},
			wantErr:         false,
		},
		{
			name: "Test 2: save Credential (the last)",
			args: args{
				Credentials:   Credentials3,
				newCredential: cred1,
				ind:           1,
			},
			wantCredentials: []*structs.Credential{Credentials3[0], cred1},
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveCred(tt.args.Credentials, tt.args.ind, tt.args.newCredential)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCredential() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.args.Credentials) != len(tt.wantCredentials) {
				t.Errorf("AddCredential() len tt.args.Credentials = %v, want %v", len(tt.args.Credentials), len(tt.wantCredentials))
				return
			}
			for ind := 0; ind < len(tt.wantCredentials); ind++ {
				require.Equal(t, tt.wantCredentials[ind].Type, tt.args.Credentials[ind].Type)
				require.Equal(t, tt.wantCredentials[ind].ID, tt.args.Credentials[ind].ID)
				require.Equal(t, tt.wantCredentials[ind].Login, tt.args.Credentials[ind].Login)
				require.Equal(t, tt.wantCredentials[ind].Password, tt.args.Credentials[ind].Password)
				require.Equal(t, tt.wantCredentials[ind].Resource, tt.args.Credentials[ind].Resource)
				require.Equal(t, tt.wantCredentials[ind].Date, tt.args.Credentials[ind].Date)
				require.Equal(t, tt.wantCredentials[ind].Description, tt.args.Credentials[ind].Description)

			}

		})
	}
}

func TestSaveFile(t *testing.T) {
	type args struct {
		Files   []*structs.File
		newFile *structs.File
		ind     int
	}
	var (
		Files = []*structs.File{
			{
				Type: structs.BlobCard,
				ID:   "123123",
				Name: "name",
				Body: []byte("person"),
			},
			{
				Type: structs.BlobCard,
				ID:   "44444444",
				Name: "1111111111111",
				Body: []byte("222222222222222222222"),
			},
		}
		Files2 = []*structs.File{
			{
				Type: structs.BlobCard,
				ID:   "123123",
				Name: "name",
				Body: []byte("person"),
			},
			{
				Type: structs.BlobCard,
				ID:   "44444444",
				Name: "1111111111111",
				Body: []byte("222222222222222222222"),
			},
		}
		File1 = &structs.File{
			Type: structs.BlobCard,
			ID:   "so-------me-new",
			Name: "name123",
			Body: []byte("som description"),
		}
	)

	tests := []struct {
		name      string
		args      args
		wantFiles []*structs.File
		wantErr   bool
	}{
		{
			name: "Test 1: save File (the first)",
			args: args{
				Files:   Files,
				newFile: File1,
				ind:     0,
			},
			wantFiles: []*structs.File{File1, Files[1]},
			wantErr:   false,
		},
		{
			name: "Test 2: save File (the last)",
			args: args{
				Files:   Files2,
				newFile: File1,
				ind:     1,
			},
			wantFiles: []*structs.File{Files2[0], File1},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveFile(tt.args.Files, tt.args.ind, tt.args.newFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.args.Files) != len(tt.wantFiles) {
				t.Errorf("AddFile() len tt.args.Files = %v, want %v", len(tt.args.Files), len(tt.wantFiles))
				return
			}
			for ind := 0; ind < len(tt.wantFiles); ind++ {
				require.Equal(t, tt.wantFiles[ind].Type, tt.args.Files[ind].Type)
				require.Equal(t, tt.wantFiles[ind].ID, tt.args.Files[ind].ID)
				require.Equal(t, tt.wantFiles[ind].Name, tt.args.Files[ind].Name)
				require.Equal(t, tt.wantFiles[ind].Body, tt.args.Files[ind].Body)
			}

		})
	}
}

func TestSaveNote(t *testing.T) {
	type args struct {
		Notes   []*structs.Note
		newNote *structs.Note
		ind     int
	}
	var (
		Notes = []*structs.Note{
			{
				Type: structs.BlobCard,
				ID:   "123123",
				Date: time.Now(),
				Name: "name",
				Body: "person",
			},
			{
				Type: structs.BlobCard,
				ID:   "fffffffffffff",
				Date: time.Now(),
				Name: "2222222222",
				Body: "fffffffffffffffffff",
			},
		}
		Notes2 = []*structs.Note{
			{
				Type: structs.BlobCard,
				ID:   "123123",
				Date: time.Now(),
				Name: "name",
				Body: "person",
			},
			{
				Type: structs.BlobCard,
				ID:   "fffffffffffff",
				Date: time.Now(),
				Name: "2222222222",
				Body: "fffffffffffffffffff",
			},
		}
		Note1 = &structs.Note{
			Type: structs.BlobCard,
			ID:   "so-------me-new",
			Name: "name123",
			Date: time.Now(),
			Body: "som description",
		}
	)

	tests := []struct {
		name      string
		args      args
		wantNotes []*structs.Note
		wantErr   bool
	}{
		{
			name: "Test 1: save Note (the first)",
			args: args{
				Notes:   Notes,
				newNote: Note1,
				ind:     0,
			},
			wantNotes: []*structs.Note{Note1, Notes[1]},
			wantErr:   false,
		},
		{
			name: "Test 2: save Note (the last)",
			args: args{
				Notes:   Notes2,
				newNote: Note1,
				ind:     1,
			},
			wantNotes: []*structs.Note{Notes2[0], Note1},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveNote(tt.args.Notes, tt.args.ind, tt.args.newNote)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.args.Notes) != len(tt.wantNotes) {
				t.Errorf("AddNote() len tt.args.Notes = %v, want %v", len(tt.args.Notes), len(tt.wantNotes))
				return
			}
			for ind := 0; ind < len(tt.wantNotes); ind++ {
				require.Equal(t, tt.wantNotes[ind].Type, tt.args.Notes[ind].Type)
				require.Equal(t, tt.wantNotes[ind].ID, tt.args.Notes[ind].ID)
				require.Equal(t, tt.wantNotes[ind].Name, tt.args.Notes[ind].Name)
				require.Equal(t, tt.wantNotes[ind].Body, tt.args.Notes[ind].Body)

			}

		})
	}
}
