package cli

import (
	"bufio"
	"context"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"github.com/w1ns3c/go-examples/crypto"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/compress"
	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
	"passkeeper/internal/entities/structs"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/cli/filesUC"
	mocks "passkeeper/mocks/gservice"
)

func TestSortCredsByDate(t *testing.T) {
	var (
		t1     = time.Now()
		creds1 = []*structs.Credential{
			{
				ID:          "1111",
				Resource:    "res_1111",
				Password:    "pass_1111",
				Description: "desc_1111",
				Date:        t1,
			},
			{
				ID:          "2222",
				Resource:    "res_2222",
				Password:    "pass_2222",
				Description: "desc_2222",
				Date:        t1.Add(time.Hour * 4), // t1 + 4 hour
			},
			{
				ID:          "3333",
				Resource:    "res_3333",
				Password:    "pass_3333",
				Description: "desc_3333",
				Date:        t1.Add(time.Hour * -2), // t1 - 2 hour
			},
		}
		order1 = []string{"2222", "1111", "3333"}
	)
	tests := []struct {
		name    string
		creds   []*structs.Credential
		orderID []string
	}{
		{
			name:    "Check 1",
			creds:   creds1,
			orderID: order1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.orderID) != len(tt.creds) {
				t.Errorf("wrong test, len of order: %d and len of notes: %d, should be the same",
					len(tt.orderID), len(tt.creds))
			}
			SortCredsByDate(tt.creds)

			for ind, id := range tt.orderID {
				require.Equal(t, tt.creds[ind].ID, id, "id not the same")
			}

		})
	}
}

func TestSortNotesByDate(t *testing.T) {
	var (
		t1    = time.Now()
		notes = []*structs.Note{
			{
				ID:   "1111",
				Body: "res_1111",
				Date: t1,
			},
			{
				ID:   "2222",
				Body: "res_2222",
				Date: t1.Add(time.Hour * 4), // t1 + 4 hour
			},
			{
				ID:   "3333",
				Body: "res_3333",
				Date: t1.Add(time.Hour * -2), // t1 - 2 hour
			},
		}
		order1 = []string{"2222", "1111", "3333"}
	)
	tests := []struct {
		name    string
		notes   []*structs.Note
		orderID []string
	}{
		{
			name:    "Check 1",
			notes:   notes,
			orderID: order1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.orderID) != len(tt.notes) {
				t.Errorf("wrong test, len of order: %d and len of notes: %d, should be the same",
					len(tt.orderID), len(tt.notes))
			}
			SortNotesByDate(tt.notes)

			for ind, id := range tt.orderID {
				require.Equal(t, tt.notes[ind].ID, id, "id not the same")
			}

		})
	}
}

type invType struct {
}

func (i invType) GetID() string {
	return ""
}

func (i invType) SetID(id string) {
	return
}

func TestClientUC_GetBlobs(t *testing.T) {
	filename1 := "testfile1"
	file, err := os.CreateTemp(config.TmpDir, filename1)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file.Name())

	buf := bufio.NewWriter(file)
	buf.WriteString(strings.Repeat("123123", 1000))
	err = buf.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename2 := "testfile2"
	file2, err := os.CreateTemp(config.TmpDir, filename2)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file2.Name())

	buf2 := bufio.NewWriter(file2)
	buf2.WriteString(strings.Repeat("222222errerr", 1000))
	err = buf2.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file2.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename3 := "testfile3"
	file3, err := os.CreateTemp(config.TmpDir, filename3)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file3.Name())

	buf3 := bufio.NewWriter(file3)
	buf3.WriteString(strings.Repeat("3333abc", 1000))
	err = buf3.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file3.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	var (
		zipData1, _ = compress.CompressFile(file.Name())
		zipData2, _ = compress.CompressFile(file2.Name())
		zipData3, _ = compress.CompressFile(file3.Name())

		login1  = "user1"
		pass1   = "123"
		userid1 = login1 + "_ID"

		hash1 = hashes.Hash(pass1)

		userSalt1, _    = crypto.GenRandStr(config.UserPassSaltLen)
		cryptHash1, _   = hashes.GenerateCryptoHash(hash1, userSalt1)
		secret1, _      = hashes.GenerateSecret(config.UserPassSaltLen)
		cryptSecret1, _ = hashes.EncryptSecret(secret1, hash1)

		user1 = &structs.User{
			ID:     userid1,
			Login:  login1,
			Hash:   cryptHash1,
			Salt:   userSalt1,
			Secret: cryptSecret1,
		}

		// Passwords
		password1 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID1111",
			Date:        time.Now().Add(time.Second * -200),
			Resource:    "localhost1111",
			Login:       "my_favorite_username1111",
			Password:    "my_favorite_password1111",
			Description: "some description1111",
		}
		password2 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID2222",
			Date:        time.Now().Add(time.Second * -500),
			Resource:    "localhost2222",
			Login:       "admin2222",
			Password:    "secret password2222",
			Description: "some new description2222",
		}
		password3 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "superID3333",
			Date:        time.Now(),
			Resource:    "localhost3333",
			Login:       "my_favorite_username333",
			Password:    "my_favorite_password3333",
			Description: "some description3333",
		}

		testCreds = []*structs.Credential{
			password1,
			password2,
			password3,
		}

		time1, _  = time.Parse("01/06", "33/44")
		time2, _  = time.Parse("01/06", "1/11")
		testCards = []*structs.Card{
			{
				ID:          "ID_CARD_1111",
				Type:        structs.BlobCard,
				Name:        "test1",
				Bank:        entities.Banks[0],
				Person:      "string",
				Number:      122222222222,
				CVC:         232,
				Expiration:  time1,
				PIN:         3333,
				Description: "test description only",
			},
			{
				ID:          "ID_CARD_22222",
				Type:        structs.BlobCard,
				Name:        "test333331",
				Bank:        entities.Banks[2],
				Person:      "Major Tom",
				Number:      234872398472,
				CVC:         23244444,
				Expiration:  time2,
				PIN:         11111,
				Description: "test description2",
			},
			{
				ID:          "ID_CARD_33",
				Type:        structs.BlobCard,
				Name:        "super secret card",
				Bank:        entities.Banks[4],
				Person:      "Major Jerry",
				Number:      2348723984721111,
				CVC:         232444443333,
				Expiration:  time2.Add(time.Second * 500),
				PIN:         2323,
				Description: "test myself",
			},
		}

		testNotes = []*structs.Note{
			{
				ID:   "ID_NOTE_1",
				Type: structs.BlobNote,
				Name: "test1",
				Date: time.Now().Add(time.Second * -300000),
				Body: "Hello\nWorld!",
			},
			{
				ID:   "ID_NOTE_2",
				Type: structs.BlobNote,
				Name: "HELLO 222222",
				Date: time.Now().Add(time.Second * -3000010),
				Body: "Hello\nWorld! 9234928309482390480298340923809840",
			},
			{
				ID:   "ID_NOTE_3",
				Type: structs.BlobNote,
				Name: "New Test Blob",
				Date: time.Now().Add(time.Second * -500000),
				Body: "Hello\nWorld! Amigo",
			},
		}

		testFiles = []*structs.File{
			{
				ID:   "FILE_ID_1",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename1),
				Body: zipData1,
			},
			{
				ID:   "FILE_ID_2",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename2),
				Body: zipData2,
			},
			{
				ID:   "FILE_ID_3",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename3),
				Body: zipData3,
			},
		}
	)

	key1, _ := hashes.GenerateCredsSecret(pass1, user1.ID, cryptSecret1)

	blob1, _ := hashes.EncryptBlob(password1, key1)
	blob1.UserID = userid1
	blob2, _ := hashes.EncryptBlob(password2, key1)
	blob2.UserID = userid1
	blob3, _ := hashes.EncryptBlob(password3, key1)
	blob3.UserID = userid1
	blob4, _ := hashes.EncryptBlob(testCards[0], key1)
	blob4.UserID = userid1
	blob5, _ := hashes.EncryptBlob(testCards[1], key1)
	blob5.UserID = userid1
	blob6, _ := hashes.EncryptBlob(testCards[2], key1)
	blob6.UserID = userid1
	blob7, _ := hashes.EncryptBlob(testNotes[0], key1)
	blob7.UserID = userid1
	blob8, _ := hashes.EncryptBlob(testNotes[1], key1)
	blob8.UserID = userid1
	blob9, _ := hashes.EncryptBlob(testNotes[2], key1)
	blob9.UserID = userid1

	blob10, _ := hashes.EncryptBlob(testFiles[0], key1)
	blob10.UserID = userid1
	blob11, _ := hashes.EncryptBlob(testFiles[1], key1)
	blob11.UserID = userid1
	blob12, _ := hashes.EncryptBlob(testFiles[2], key1)
	blob12.UserID = userid1

	invalidBlob1, _ := hashes.EncryptBlob(testNotes[0], key1+"123")
	inv := &structs.Note{
		Type: 0,
		ID:   "1213",
		Name: "name",
		Date: time.Now(),
		Body: "993243",
	}
	invalidBlob2, _ := hashes.EncryptBlob(inv, key1)

	blobs := []*structs.CryptoBlob{
		blob1, blob2, blob3, blob4, blob5, blob6,
		blob7, blob8, blob9, blob10, blob11, blob12,
	}
	pbBlobs := make([]*pb.CryptoBlob, len(blobs))

	for ind, bl := range blobs {
		pbBlobs[ind] = new(pb.CryptoBlob)
		pbBlobs[ind].ID = bl.ID
		pbBlobs[ind].Blob = bl.Blob
	}

	tests := []struct {
		name    string
		prepare func(*mocks.MockBlobSvcClient)
		creds   []*structs.Credential
		cards   []*structs.Card
		notes   []*structs.Note
		files   []*structs.File
		wantErr bool
	}{
		{
			name: "Test GetBlobs 1: success",
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobList(gomock.Any(), new(empty.Empty)).
					Return(&pb.BlobListResponse{Blobs: pbBlobs}, nil)
			},

			creds: testCreds,
			cards: testCards,
			notes: testNotes,
			files: testFiles,

			wantErr: false,
		},
		{
			name: "Test GetBlobs 2: success but empty results",
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobList(gomock.Any(), new(empty.Empty)).
					Return(&pb.BlobListResponse{Blobs: []*pb.CryptoBlob{
						{
							ID:   invalidBlob1.ID,
							Blob: invalidBlob1.Blob,
						},
						{
							ID:   invalidBlob2.ID,
							Blob: invalidBlob2.Blob,
						},
					}}, nil)
			},

			wantErr: false,
		},
		{
			name: "Test GetBlobs 3: err",
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobList(gomock.Any(), new(empty.Empty)).
					Return(&pb.BlobListResponse{Blobs: []*pb.CryptoBlob{}},
						myerrors.ErrBlobList)
			},

			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mock       = mocks.NewMockBlobSvcClient(ctrl)
				usecase, _ = NewClientUC()
			)

			if tt.prepare != nil {
				tt.prepare(mock)
			}

			usecase.blobsSvc = mock
			usecase.User = user1
			usecase.User.Secret = key1

			err := usecase.GetBlobs(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, len(tt.creds), len(usecase.Creds))
			require.Equal(t, len(tt.cards), len(usecase.Cards))
			require.Equal(t, len(tt.files), len(usecase.Files))
			require.Equal(t, len(tt.notes), len(usecase.Notes))

		})
	}
}

func TestClientUC_AddBlob(t *testing.T) {
	filename1 := "testfile1"
	file, err := os.CreateTemp(config.TmpDir, filename1)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file.Name())

	buf := bufio.NewWriter(file)
	buf.WriteString(strings.Repeat("123123", 1000))
	err = buf.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename2 := "testfile2"
	file2, err := os.CreateTemp(config.TmpDir, filename2)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file2.Name())

	buf2 := bufio.NewWriter(file2)
	buf2.WriteString(strings.Repeat("222222errerr", 1000))
	err = buf2.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file2.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename3 := "testfile3"
	file3, err := os.CreateTemp(config.TmpDir, filename3)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file3.Name())

	buf3 := bufio.NewWriter(file3)
	buf3.WriteString(strings.Repeat("3333abc", 1000))
	err = buf3.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file3.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	var (
		login1  = "user1"
		pass1   = "123"
		userid1 = login1 + "_ID"

		hash1 = hashes.Hash(pass1)

		userSalt1, _    = crypto.GenRandStr(config.UserPassSaltLen)
		cryptHash1, _   = hashes.GenerateCryptoHash(hash1, userSalt1)
		secret1, _      = hashes.GenerateSecret(config.UserPassSaltLen)
		cryptSecret1, _ = hashes.EncryptSecret(secret1, hash1)

		user1 = &structs.User{
			ID:     userid1,
			Login:  login1,
			Hash:   cryptHash1,
			Salt:   userSalt1,
			Secret: cryptSecret1,
		}

		// Passwords
		password1 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID1111",
			Date:        time.Now().Add(time.Second * -200),
			Resource:    "localhost1111",
			Login:       "my_favorite_username1111",
			Password:    "my_favorite_password1111",
			Description: "some description1111",
		}
		password2 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID2222",
			Date:        time.Now().Add(time.Second * -500),
			Resource:    "localhost2222",
			Login:       "admin2222",
			Password:    "secret password2222",
			Description: "some new description2222",
		}
		password3 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "superID3333",
			Date:        time.Now(),
			Resource:    "localhost3333",
			Login:       "my_favorite_username333",
			Password:    "my_favorite_password3333",
			Description: "some description3333",
		}

		testCreds = []*structs.Credential{
			password1,
			password2,
			password3,
		}

		addPass = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "new_superID3333",
			Date:        time.Now(),
			Resource:    "new localhost3333",
			Login:       "new my_favorite_username333",
			Password:    "new my_favorite_password3333",
			Description: "new some description3333",
		}

		time1, _  = time.Parse("01/06", "33/44")
		time2, _  = time.Parse("01/06", "1/11")
		testCards = []*structs.Card{
			{
				ID:          "ID_CARD_1111",
				Type:        structs.BlobCard,
				Name:        "test1",
				Bank:        entities.Banks[0],
				Person:      "string",
				Number:      122222222222,
				CVC:         232,
				Expiration:  time1,
				PIN:         3333,
				Description: "test description only",
			},
			{
				ID:          "ID_CARD_22222",
				Type:        structs.BlobCard,
				Name:        "test333331",
				Bank:        entities.Banks[2],
				Person:      "Major Tom",
				Number:      234872398472,
				CVC:         23244444,
				Expiration:  time2,
				PIN:         11111,
				Description: "test description2",
			},
			{
				ID:          "ID_CARD_33",
				Type:        structs.BlobCard,
				Name:        "super secret card",
				Bank:        entities.Banks[4],
				Person:      "Major Jerry",
				Number:      2348723984721111,
				CVC:         232444443333,
				Expiration:  time2.Add(time.Second * 500),
				PIN:         2323,
				Description: "test myself",
			},
		}
		addCard = &structs.Card{
			ID:          "new_ID_CARD_33",
			Type:        structs.BlobCard,
			Name:        "super secret card",
			Bank:        entities.Banks[5],
			Person:      "new Major Jerry",
			Number:      11111111111111111,
			CVC:         222,
			Expiration:  time2.Add(time.Second * 8000),
			PIN:         333,
			Description: "new test myself",
		}

		testNotes = []*structs.Note{
			{
				ID:   "ID_NOTE_1",
				Type: structs.BlobNote,
				Name: "test1",
				Date: time.Now().Add(time.Second * -300000),
				Body: "Hello\nWorld!",
			},
			{
				ID:   "ID_NOTE_2",
				Type: structs.BlobNote,
				Name: "HELLO 222222",
				Date: time.Now().Add(time.Second * -3000010),
				Body: "Hello\nWorld! 9234928309482390480298340923809840",
			},
			{
				ID:   "ID_NOTE_3",
				Type: structs.BlobNote,
				Name: "New Test Blob",
				Date: time.Now().Add(time.Second * -500000),
				Body: "Hello\nWorld! Amigo",
			},
		}
		addNote = &structs.Note{
			ID:   "new_ID_NOTE_3",
			Type: structs.BlobNote,
			Name: "--New Test Blob--",
			Date: time.Now().Add(time.Second * -5000),
			Body: "Hello\nWorld! Amigo! I should finish it...",
		}

		zipData1, _ = compress.CompressFile(file.Name())
		zipData2, _ = compress.CompressFile(file2.Name())
		zipData3, _ = compress.CompressFile(file3.Name())

		testFiles = []*structs.File{
			{
				ID:   "FILE_ID_1",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename1),
				Body: zipData1,
			},
			{
				ID:   "FILE_ID_2",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename2),
				Body: zipData2,
			},
			{
				ID:   "FILE_ID_3",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename3),
				Body: zipData3,
			},
		}

		addFile = &structs.File{
			ID:   "new_FILE_ID_2",
			Type: structs.BlobFile,
			Name: filesUC.GenFileShortName(filename2),
			Body: zipData3,
		}
	)

	key1, _ := hashes.GenerateCredsSecret(pass1, user1.ID, cryptSecret1)

	blob1, _ := hashes.EncryptBlob(testCreds[0], key1)
	blob1.UserID = userid1
	blob2, _ := hashes.EncryptBlob(testCreds[1], key1)
	blob2.UserID = userid1
	blob3, _ := hashes.EncryptBlob(testCreds[2], key1)
	blob3.UserID = userid1
	blob4, _ := hashes.EncryptBlob(testCards[0], key1)
	blob4.UserID = userid1
	blob5, _ := hashes.EncryptBlob(testCards[1], key1)
	blob5.UserID = userid1
	blob6, _ := hashes.EncryptBlob(testCards[2], key1)
	blob6.UserID = userid1
	blob7, _ := hashes.EncryptBlob(testNotes[0], key1)
	blob7.UserID = userid1
	blob8, _ := hashes.EncryptBlob(testNotes[1], key1)
	blob8.UserID = userid1
	blob9, _ := hashes.EncryptBlob(testNotes[2], key1)
	blob9.UserID = userid1

	blob10, _ := hashes.EncryptBlob(testFiles[0], key1)
	blob10.UserID = userid1
	blob11, _ := hashes.EncryptBlob(testFiles[1], key1)
	blob11.UserID = userid1
	blob12, _ := hashes.EncryptBlob(testFiles[2], key1)
	blob12.UserID = userid1

	type args struct {
		cred structs.CredInf
	}
	type fiels struct {
		cards []*structs.Card
		creds []*structs.Credential
		notes []*structs.Note
		files []*structs.File
	}
	tests := []struct {
		name    string
		args    args
		f       fiels
		prepare func(*mocks.MockBlobSvcClient)
		wantErr bool
	}{
		{
			name: "Test AddBlob 1: valid cred",
			args: args{
				cred: addPass,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobAdd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test AddBlob 2: valid card",
			args: args{
				cred: addCard,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobAdd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test AddBlob 3: valid note",
			args: args{
				cred: addNote,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobAdd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test AddBlob 4: valid file",
			args: args{
				cred: addFile,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobAdd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test AddBlob 5: err",
			args: args{
				cred: addFile,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobAdd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), myerrors.ErrBlobAdd)
			},
			wantErr: true,
		},
		{
			name: "Test AddBlob 6: invalid blob type",
			args: args{
				cred: &invType{},
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobAdd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mock       = mocks.NewMockBlobSvcClient(ctrl)
				usecase, _ = NewClientUC()
			)

			if tt.prepare != nil {
				tt.prepare(mock)
			}

			usecase.blobsSvc = mock
			usecase.Creds = tt.f.creds
			usecase.Cards = tt.f.cards
			usecase.Notes = tt.f.notes
			usecase.Files = tt.f.files

			usecase.User = user1
			usecase.User.Secret = key1

			err := usecase.AddBlob(context.Background(), tt.args.cred)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddBlob() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			switch tt.args.cred.(type) {
			case *structs.Credential:
				found := false
				for _, bl := range usecase.Creds {
					if reflect.DeepEqual(bl, tt.args.cred) {
						found = true
						break
					}
				}
				require.True(t, found)

				return

			case *structs.Card:
				found := false
				for _, bl := range usecase.Cards {
					if reflect.DeepEqual(bl, tt.args.cred) {
						found = true
						break
					}
				}
				require.True(t, found)

				return

			case *structs.Note:
				found := false
				for _, bl := range usecase.Notes {
					if reflect.DeepEqual(bl, tt.args.cred) {
						found = true
						break
					}
				}
				require.True(t, found)

				return

			case *structs.File:
				found := false
				for _, bl := range usecase.Files {
					if reflect.DeepEqual(bl, tt.args.cred) {
						found = true
						break
					}
				}
				require.True(t, found)

				return

			default:
				t.Errorf("unknown type of blob")

			}

		})
	}
}

func TestClientUC_DelBlob(t *testing.T) {
	filename1 := "testfile1"
	file, err := os.CreateTemp(config.TmpDir, filename1)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file.Name())

	buf := bufio.NewWriter(file)
	buf.WriteString(strings.Repeat("123123", 1000))
	err = buf.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename2 := "testfile2"
	file2, err := os.CreateTemp(config.TmpDir, filename2)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file2.Name())

	buf2 := bufio.NewWriter(file2)
	buf2.WriteString(strings.Repeat("222222errerr", 1000))
	err = buf2.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file2.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename3 := "testfile3"
	file3, err := os.CreateTemp(config.TmpDir, filename3)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file3.Name())

	buf3 := bufio.NewWriter(file3)
	buf3.WriteString(strings.Repeat("3333abc", 1000))
	err = buf3.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file3.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	var (
		zipData1, _ = compress.CompressFile(file.Name())
		zipData2, _ = compress.CompressFile(file2.Name())
		zipData3, _ = compress.CompressFile(file3.Name())

		login1  = "user1"
		pass1   = "123"
		userid1 = login1 + "_ID"

		hash1 = hashes.Hash(pass1)

		userSalt1, _    = crypto.GenRandStr(config.UserPassSaltLen)
		cryptHash1, _   = hashes.GenerateCryptoHash(hash1, userSalt1)
		secret1, _      = hashes.GenerateSecret(config.UserPassSaltLen)
		cryptSecret1, _ = hashes.EncryptSecret(secret1, hash1)

		user1 = &structs.User{
			ID:     userid1,
			Login:  login1,
			Hash:   cryptHash1,
			Salt:   userSalt1,
			Secret: cryptSecret1,
		}

		// Passwords
		password1 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID1111",
			Date:        time.Now().Add(time.Second * -200),
			Resource:    "localhost1111",
			Login:       "my_favorite_username1111",
			Password:    "my_favorite_password1111",
			Description: "some description1111",
		}
		password2 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID2222",
			Date:        time.Now().Add(time.Second * -500),
			Resource:    "localhost2222",
			Login:       "admin2222",
			Password:    "secret password2222",
			Description: "some new description2222",
		}
		password3 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "superID3333",
			Date:        time.Now(),
			Resource:    "localhost3333",
			Login:       "my_favorite_username333",
			Password:    "my_favorite_password3333",
			Description: "some description3333",
		}

		testCreds = []*structs.Credential{
			password1,
			password2,
			password3,
		}

		time1, _  = time.Parse("01/06", "33/44")
		time2, _  = time.Parse("01/06", "1/11")
		testCards = []*structs.Card{
			{
				ID:          "ID_CARD_1111",
				Type:        structs.BlobCard,
				Name:        "test1",
				Bank:        entities.Banks[0],
				Person:      "string",
				Number:      122222222222,
				CVC:         232,
				Expiration:  time1,
				PIN:         3333,
				Description: "test description only",
			},
			{
				ID:          "ID_CARD_22222",
				Type:        structs.BlobCard,
				Name:        "test333331",
				Bank:        entities.Banks[2],
				Person:      "Major Tom",
				Number:      234872398472,
				CVC:         23244444,
				Expiration:  time2,
				PIN:         11111,
				Description: "test description2",
			},
			{
				ID:          "ID_CARD_33",
				Type:        structs.BlobCard,
				Name:        "super secret card",
				Bank:        entities.Banks[4],
				Person:      "Major Jerry",
				Number:      2348723984721111,
				CVC:         232444443333,
				Expiration:  time2.Add(time.Second * 500),
				PIN:         2323,
				Description: "test myself",
			},
		}

		testNotes = []*structs.Note{
			{
				ID:   "ID_NOTE_1",
				Type: structs.BlobNote,
				Name: "test1",
				Date: time.Now().Add(time.Second * -300000),
				Body: "Hello\nWorld!",
			},
			{
				ID:   "ID_NOTE_2",
				Type: structs.BlobNote,
				Name: "HELLO 222222",
				Date: time.Now().Add(time.Second * -3000010),
				Body: "Hello\nWorld! 9234928309482390480298340923809840",
			},
			{
				ID:   "ID_NOTE_3",
				Type: structs.BlobNote,
				Name: "New Test Blob",
				Date: time.Now().Add(time.Second * -500000),
				Body: "Hello\nWorld! Amigo",
			},
		}

		testFiles = []*structs.File{
			{
				ID:   "FILE_ID_1",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename1),
				Body: zipData1,
			},
			{
				ID:   "FILE_ID_2",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename2),
				Body: zipData2,
			},
			{
				ID:   "FILE_ID_3",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename3),
				Body: zipData3,
			},
		}
	)

	key1, _ := hashes.GenerateCredsSecret(pass1, user1.ID, cryptSecret1)

	blob1, _ := hashes.EncryptBlob(testCreds[0], key1)
	blob1.UserID = userid1
	blob2, _ := hashes.EncryptBlob(testCreds[1], key1)
	blob2.UserID = userid1
	blob3, _ := hashes.EncryptBlob(testCreds[2], key1)
	blob3.UserID = userid1
	blob4, _ := hashes.EncryptBlob(testCards[0], key1)
	blob4.UserID = userid1
	blob5, _ := hashes.EncryptBlob(testCards[1], key1)
	blob5.UserID = userid1
	blob6, _ := hashes.EncryptBlob(testCards[2], key1)
	blob6.UserID = userid1
	blob7, _ := hashes.EncryptBlob(testNotes[0], key1)
	blob7.UserID = userid1
	blob8, _ := hashes.EncryptBlob(testNotes[1], key1)
	blob8.UserID = userid1
	blob9, _ := hashes.EncryptBlob(testNotes[2], key1)
	blob9.UserID = userid1

	blob10, _ := hashes.EncryptBlob(testFiles[0], key1)
	blob10.UserID = userid1
	blob11, _ := hashes.EncryptBlob(testFiles[1], key1)
	blob11.UserID = userid1
	blob12, _ := hashes.EncryptBlob(testFiles[2], key1)
	blob12.UserID = userid1

	type args struct {
		blobType structs.BlobType
		ind      int
	}
	type fiels struct {
		cards []*structs.Card
		creds []*structs.Credential
		notes []*structs.Note
		files []*structs.File
	}
	tests := []struct {
		name    string
		args    args
		f       fiels
		prepare func(*mocks.MockBlobSvcClient)
		wantErr bool
	}{
		{
			name: "Test EditBlob 1: valid ind cred",
			args: args{
				blobType: structs.BlobCred,
				ind:      2,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobDel(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test DelBlob 2: valid ind card",
			args: args{
				blobType: structs.BlobCard,
				ind:      2,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobDel(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test DelBlob 3: valid ind note",
			args: args{
				blobType: structs.BlobNote,
				ind:      2,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobDel(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test DelBlob 4: valid ind file",
			args: args{
				blobType: structs.BlobFile,
				ind:      2,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobDel(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test DelBlob 5: invalid ind cred",
			args: args{
				blobType: structs.BlobCred,
				ind:      25,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
			},
			wantErr: true,
		},
		{
			name: "Test DelBlob 6: invalid ind card",
			args: args{
				blobType: structs.BlobCard,
				ind:      25,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
			},
			wantErr: true,
		},
		{
			name: "Test DelBlob 7: invalid ind note",
			args: args{
				blobType: structs.BlobNote,
				ind:      25,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
			},
			wantErr: true,
		},
		{
			name: "Test DelBlob 8: invalid ind file",
			args: args{
				blobType: structs.BlobFile,
				ind:      25,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
			},
			wantErr: true,
		},
		{
			name: "Test DelBlob 9: invalid blob type",
			args: args{
				blobType: 0,
				ind:      1,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
			},
			wantErr: true,
		},
		{
			name: "Test DelBlob 10: err on server side",
			args: args{
				blobType: structs.BlobFile,
				ind:      1,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobDel(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), myerrors.ErrBlobDel)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mock       = mocks.NewMockBlobSvcClient(ctrl)
				usecase, _ = NewClientUC()
			)

			if tt.prepare != nil {
				tt.prepare(mock)
			}

			usecase.blobsSvc = mock
			usecase.Creds = tt.f.creds
			usecase.Cards = tt.f.cards
			usecase.Notes = tt.f.notes
			usecase.Files = tt.f.files

			usecase.User = user1
			usecase.User.Secret = key1

			err := usecase.DelBlob(context.Background(), tt.args.ind, tt.args.blobType)
			if (err != nil) != tt.wantErr {
				t.Errorf("DelBlob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientUC_EditBlob(t *testing.T) {
	filename1 := "testfile1"
	file, err := os.CreateTemp(config.TmpDir, filename1)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file.Name())

	buf := bufio.NewWriter(file)
	buf.WriteString(strings.Repeat("123123", 1000))
	err = buf.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename2 := "testfile2"
	file2, err := os.CreateTemp(config.TmpDir, filename2)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file2.Name())

	buf2 := bufio.NewWriter(file2)
	buf2.WriteString(strings.Repeat("222222errerr", 1000))
	err = buf2.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file2.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename3 := "testfile3"
	file3, err := os.CreateTemp(config.TmpDir, filename3)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file3.Name())

	buf3 := bufio.NewWriter(file3)
	buf3.WriteString(strings.Repeat("3333abc", 1000))
	err = buf3.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file3.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	var (
		zipData1, _ = compress.CompressFile(file.Name())
		zipData2, _ = compress.CompressFile(file2.Name())
		zipData3, _ = compress.CompressFile(file3.Name())

		login1  = "user1"
		pass1   = "123"
		userid1 = login1 + "_ID"

		hash1 = hashes.Hash(pass1)

		userSalt1, _    = crypto.GenRandStr(config.UserPassSaltLen)
		cryptHash1, _   = hashes.GenerateCryptoHash(hash1, userSalt1)
		secret1, _      = hashes.GenerateSecret(config.UserPassSaltLen)
		cryptSecret1, _ = hashes.EncryptSecret(secret1, hash1)

		user1 = &structs.User{
			ID:     userid1,
			Login:  login1,
			Hash:   cryptHash1,
			Salt:   userSalt1,
			Secret: cryptSecret1,
		}

		// Passwords
		password1 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID1111",
			Date:        time.Now().Add(time.Second * -200),
			Resource:    "localhost1111",
			Login:       "my_favorite_username1111",
			Password:    "my_favorite_password1111",
			Description: "some description1111",
		}
		password2 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID2222",
			Date:        time.Now().Add(time.Second * -500),
			Resource:    "localhost2222",
			Login:       "admin2222",
			Password:    "secret password2222",
			Description: "some new description2222",
		}
		password3 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "superID3333",
			Date:        time.Now(),
			Resource:    "localhost3333",
			Login:       "my_favorite_username333",
			Password:    "my_favorite_password3333",
			Description: "some description3333",
		}

		testCreds = []*structs.Credential{
			password1,
			password2,
			password3,
		}

		editPass = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "superID3333",
			Date:        time.Now(),
			Resource:    "new localhost3333",
			Login:       "new my_favorite_username333",
			Password:    "new my_favorite_password3333",
			Description: "new some description3333",
		}

		time1, _  = time.Parse("01/06", "33/44")
		time2, _  = time.Parse("01/06", "1/11")
		testCards = []*structs.Card{
			{
				ID:          "ID_CARD_1111",
				Type:        structs.BlobCard,
				Name:        "test1",
				Bank:        entities.Banks[0],
				Person:      "string",
				Number:      122222222222,
				CVC:         232,
				Expiration:  time1,
				PIN:         3333,
				Description: "test description only",
			},
			{
				ID:          "ID_CARD_22222",
				Type:        structs.BlobCard,
				Name:        "test333331",
				Bank:        entities.Banks[2],
				Person:      "Major Tom",
				Number:      234872398472,
				CVC:         23244444,
				Expiration:  time2,
				PIN:         11111,
				Description: "test description2",
			},
			{
				ID:          "ID_CARD_33",
				Type:        structs.BlobCard,
				Name:        "super secret card",
				Bank:        entities.Banks[4],
				Person:      "Major Jerry",
				Number:      2348723984721111,
				CVC:         232444443333,
				Expiration:  time2.Add(time.Second * 500),
				PIN:         2323,
				Description: "test myself",
			},
		}
		editCard = &structs.Card{
			ID:          "ID_CARD_33",
			Type:        structs.BlobCard,
			Name:        "super secret card",
			Bank:        entities.Banks[5],
			Person:      "new Major Jerry",
			Number:      11111111111111111,
			CVC:         222,
			Expiration:  time2.Add(time.Second * 8000),
			PIN:         333,
			Description: "new test myself",
		}

		testNotes = []*structs.Note{
			{
				ID:   "ID_NOTE_1",
				Type: structs.BlobNote,
				Name: "test1",
				Date: time.Now().Add(time.Second * -300000),
				Body: "Hello\nWorld!",
			},
			{
				ID:   "ID_NOTE_2",
				Type: structs.BlobNote,
				Name: "HELLO 222222",
				Date: time.Now().Add(time.Second * -3000010),
				Body: "Hello\nWorld! 9234928309482390480298340923809840",
			},
			{
				ID:   "ID_NOTE_3",
				Type: structs.BlobNote,
				Name: "New Test Blob",
				Date: time.Now().Add(time.Second * -500000),
				Body: "Hello\nWorld! Amigo",
			},
		}
		editNote = &structs.Note{
			ID:   "ID_NOTE_3",
			Type: structs.BlobNote,
			Name: "--New Test Blob--",
			Date: time.Now().Add(time.Second * -5000),
			Body: "Hello\nWorld! Amigo! I should finish it...",
		}

		testFiles = []*structs.File{
			{
				ID:   "FILE_ID_1",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename1),
				Body: zipData1,
			},
			{
				ID:   "FILE_ID_2",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename2),
				Body: zipData2,
			},
			{
				ID:   "FILE_ID_3",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename3),
				Body: zipData3,
			},
		}

		editFile = &structs.File{
			ID:   "FILE_ID_2",
			Type: structs.BlobFile,
			Name: filesUC.GenFileShortName(filename2),
			Body: zipData3,
		}
	)

	key1, _ := hashes.GenerateCredsSecret(pass1, user1.ID, cryptSecret1)

	blob1, _ := hashes.EncryptBlob(testCreds[0], key1)
	blob1.UserID = userid1
	blob2, _ := hashes.EncryptBlob(testCreds[1], key1)
	blob2.UserID = userid1
	blob3, _ := hashes.EncryptBlob(testCreds[2], key1)
	blob3.UserID = userid1
	blob4, _ := hashes.EncryptBlob(testCards[0], key1)
	blob4.UserID = userid1
	blob5, _ := hashes.EncryptBlob(testCards[1], key1)
	blob5.UserID = userid1
	blob6, _ := hashes.EncryptBlob(testCards[2], key1)
	blob6.UserID = userid1
	blob7, _ := hashes.EncryptBlob(testNotes[0], key1)
	blob7.UserID = userid1
	blob8, _ := hashes.EncryptBlob(testNotes[1], key1)
	blob8.UserID = userid1
	blob9, _ := hashes.EncryptBlob(testNotes[2], key1)
	blob9.UserID = userid1

	blob10, _ := hashes.EncryptBlob(testFiles[0], key1)
	blob10.UserID = userid1
	blob11, _ := hashes.EncryptBlob(testFiles[1], key1)
	blob11.UserID = userid1
	blob12, _ := hashes.EncryptBlob(testFiles[2], key1)
	blob12.UserID = userid1

	//editPassBlob, _ := hashes.EncryptBlob(editPass, key1)
	//editCardBlob, _ := hashes.EncryptBlob(editCard, key1)
	//editNoteBlob, _ := hashes.EncryptBlob(editNote, key1)
	//editFileBlob, _ := hashes.EncryptBlob(editFile, key1)

	//invalidBlob1, _ := hashes.EncryptBlob(testNotes[0], key1+"123")

	//invalidBlob2, _ := hashes.EncryptBlob(inv, key1)

	type args struct {
		cred structs.CredInf
		ind  int
	}
	type fiels struct {
		cards []*structs.Card
		creds []*structs.Credential
		notes []*structs.Note
		files []*structs.File
	}
	tests := []struct {
		name    string
		args    args
		f       fiels
		prepare func(*mocks.MockBlobSvcClient)
		wantErr bool
	}{
		{
			name: "Test EditBlob 1: valid ind cred",
			args: args{
				cred: editPass,
				ind:  2,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobUpd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test EditBlob 2: valid ind card",
			args: args{
				cred: editCard,
				ind:  2,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobUpd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test EditBlob 3: valid ind note",
			args: args{
				cred: editNote,
				ind:  2,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobUpd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test EditBlob 4: valid ind file",
			args: args{
				cred: editFile,
				ind:  2,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobUpd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), nil)
			},
			wantErr: false,
		},
		{
			name: "Test EditBlob 5: invalid ind cred",
			args: args{
				cred: editPass,
				ind:  25,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
			},
			wantErr: true,
		},
		{
			name: "Test EditBlob 6: invalid ind card",
			args: args{
				cred: editCard,
				ind:  25,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
			},
			wantErr: true,
		},
		{
			name: "Test EditBlob 7: invalid ind note",
			args: args{
				cred: editNote,
				ind:  25,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
			},
			wantErr: true,
		},
		{
			name: "Test EditBlob 8: invalid ind file",
			args: args{
				cred: editFile,
				ind:  25,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
			},
			wantErr: true,
		},
		{
			name: "Test EditBlob 9: invalid blob type",
			args: args{
				cred: &invType{},
				ind:  1,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
			},
			wantErr: true,
		},
		{
			name: "Test EditBlob 10: err on server side",
			args: args{
				cred: editFile,
				ind:  1,
			},
			f: fiels{
				testCards,
				testCreds,
				testNotes,
				testFiles,
			},
			prepare: func(cli *mocks.MockBlobSvcClient) {
				cli.EXPECT().BlobUpd(gomock.Any(), gomock.Any()).
					Return(new(empty.Empty), myerrors.ErrBlobUpd)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mock       = mocks.NewMockBlobSvcClient(ctrl)
				usecase, _ = NewClientUC()
			)

			if tt.prepare != nil {
				tt.prepare(mock)
			}

			usecase.blobsSvc = mock
			usecase.Creds = tt.f.creds
			usecase.Cards = tt.f.cards
			usecase.Notes = tt.f.notes
			usecase.Files = tt.f.files

			usecase.User = user1
			usecase.User.Secret = key1

			err := usecase.EditBlob(context.Background(), tt.args.cred, tt.args.ind)
			if (err != nil) != tt.wantErr {
				t.Errorf("EditBlob() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			switch tt.args.cred.(type) {
			case *structs.Credential:
				if tt.args.ind < 0 || tt.args.ind >= len(usecase.Creds) {
					t.Errorf("invalid index of notes")
					return
				}
				require.Equal(t, tt.args.cred, usecase.Creds[tt.args.ind])

			case *structs.Card:
				if tt.args.ind < 0 || tt.args.ind >= len(usecase.Cards) {
					t.Errorf("invalid index of cards")
					return
				}
				require.Equal(t, tt.args.cred, usecase.Cards[tt.args.ind])
			case *structs.Note:
				if tt.args.ind < 0 || tt.args.ind >= len(usecase.Notes) {
					t.Errorf("invalid index of notes")
					return
				}
				require.Equal(t, tt.args.cred, usecase.Notes[tt.args.ind])
			case *structs.File:
				if tt.args.ind < 0 || tt.args.ind >= len(usecase.Files) {
					t.Errorf("invalid index of files")
					return
				}
				require.Equal(t, tt.args.cred, usecase.Files[tt.args.ind])
			default:
				t.Errorf("unknown type of blob")

			}
		})
	}
}

func TestClientUC_GetCredByIND(t *testing.T) {
	var (
		password1 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID1111",
			Date:        time.Now().Add(time.Second * -200),
			Resource:    "localhost1111",
			Login:       "my_favorite_username1111",
			Password:    "my_favorite_password1111",
			Description: "some description1111",
		}
		password2 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID2222",
			Date:        time.Now().Add(time.Second * -500),
			Resource:    "localhost2222",
			Login:       "admin2222",
			Password:    "secret password2222",
			Description: "some new description2222",
		}
		password3 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "superID3333",
			Date:        time.Now(),
			Resource:    "localhost3333",
			Login:       "my_favorite_username333",
			Password:    "my_favorite_password3333",
			Description: "some description3333",
		}

		tmpCreds = []*structs.Credential{
			password1,
			password2,
			password3,
		}
	)

	type args struct {
		ind      int
		tmpCreds []*structs.Credential
	}
	tests := []struct {
		name     string
		args     args
		wantCred *structs.Credential
		wantErr  bool
	}{
		{
			name: "Test GetByInd Cred 1: valid ind",
			args: args{
				ind:      1,
				tmpCreds: tmpCreds,
			},
			wantCred: password2,
			wantErr:  false,
		},
		{
			name: "Test GetByInd Cred 2: ind < 0",
			args: args{
				ind:      -10,
				tmpCreds: tmpCreds,
			},
			wantErr: true,
		},
		{
			name: "Test GetByInd Cred 3: ind > len(creds)",
			args: args{
				ind:      100,
				tmpCreds: tmpCreds,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase, _ := NewClientUC()
			usecase.Creds = tt.args.tmpCreds

			gotCred, err := usecase.GetCredByIND(tt.args.ind)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredByIND() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCred, tt.wantCred) {
				t.Errorf("GetCredByIND() gotCred = %v, want %v", gotCred, tt.wantCred)
			}
		})
	}
}

func TestClientUC_GetCardByIND(t *testing.T) {
	var (
		tmpCards = []*structs.Card{
			{
				ID:          "ID_CARD_1111",
				Type:        structs.BlobCard,
				Name:        "test1",
				Bank:        entities.Banks[0],
				Person:      "string",
				Number:      122222222222,
				CVC:         232,
				Expiration:  time.Now(),
				PIN:         3333,
				Description: "test description only",
			},
			{
				ID:          "ID_CARD_22222",
				Type:        structs.BlobCard,
				Name:        "test333331",
				Bank:        entities.Banks[2],
				Person:      "Major Tom",
				Number:      234872398472,
				CVC:         23244444,
				Expiration:  time.Now().Add(time.Second * 100),
				PIN:         11111,
				Description: "test description2",
			},
			{
				ID:          "ID_CARD_33",
				Type:        structs.BlobCard,
				Name:        "super secret card",
				Bank:        entities.Banks[4],
				Person:      "Major Jerry",
				Number:      2348723984721111,
				CVC:         232444443333,
				Expiration:  time.Now().Add(time.Second * -500),
				PIN:         2323,
				Description: "test myself",
			},
		}
	)

	type args struct {
		ind      int
		tmpCards []*structs.Card
	}
	tests := []struct {
		name     string
		args     args
		wantCard *structs.Card
		wantErr  bool
	}{
		{
			name: "Test GetByInd Card 1: valid ind",
			args: args{
				ind:      1,
				tmpCards: tmpCards,
			},
			wantCard: tmpCards[1],
			wantErr:  false,
		},
		{
			name: "Test GetByInd Card 2: ind < 0",
			args: args{
				ind:      -10,
				tmpCards: tmpCards,
			},
			wantErr: true,
		},
		{
			name: "Test GetByInd Card 3: ind > len(cards)",
			args: args{
				ind:      100,
				tmpCards: tmpCards,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			usecase, _ := NewClientUC()
			usecase.Cards = tt.args.tmpCards

			gotCard, err := usecase.GetCardByIND(tt.args.ind)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCardByIND() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCard, tt.wantCard) {
				t.Errorf("GetCardByIND() gotCard = %v, want %v", gotCard, tt.wantCard)
			}
		})
	}
}

func TestClientUC_GetNoteByIND(t *testing.T) {
	var (
		tmpNotes = []*structs.Note{
			{
				ID:   "ID_NOTE_1",
				Type: structs.BlobNote,
				Name: "test1",
				Date: time.Now().Add(time.Second * -300000),
				Body: "Hello\nWorld!",
			},
			{
				ID:   "ID_NOTE_2",
				Type: structs.BlobNote,
				Name: "HELLO 222222",
				Date: time.Now().Add(time.Second * -3000010),
				Body: "Hello\nWorld! 9234928309482390480298340923809840",
			},
			{
				ID:   "ID_NOTE_3",
				Type: structs.BlobNote,
				Name: "New Test Blob",
				Date: time.Now().Add(time.Second * -500000),
				Body: "Hello\nWorld! Amigo",
			},
		}
	)
	type args struct {
		ind      int
		tmpNotes []*structs.Note
	}
	tests := []struct {
		name     string
		args     args
		wantNote *structs.Note
		wantErr  bool
	}{
		{
			name: "Test GetByInd Note 1: valid ind",
			args: args{
				ind:      1,
				tmpNotes: tmpNotes,
			},
			wantNote: tmpNotes[1],
			wantErr:  false,
		},
		{
			name: "Test GetByInd Note 2: ind < 0",
			args: args{
				ind:      -10,
				tmpNotes: tmpNotes,
			},
			wantErr: true,
		},
		{
			name: "Test GetByInd Note 3: ind > len(notes)",
			args: args{
				ind:      100,
				tmpNotes: tmpNotes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase, _ := NewClientUC()
			usecase.Notes = tt.args.tmpNotes

			gotNote, err := usecase.GetNoteByIND(tt.args.ind)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNoteByIND() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNote, tt.wantNote) {
				t.Errorf("GetNoteByIND() gotNote = %v, want %v", gotNote, tt.wantNote)
			}
		})
	}
}

func TestClientUC_GetFileByIND(t *testing.T) {
	filename1 := "testfile1"
	file, err := os.CreateTemp(config.TmpDir, filename1)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file.Name())

	buf := bufio.NewWriter(file)
	buf.WriteString(strings.Repeat("123123", 1000))
	err = buf.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename2 := "testfile2"
	file2, err := os.CreateTemp(config.TmpDir, filename2)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file2.Name())

	buf2 := bufio.NewWriter(file2)
	buf2.WriteString(strings.Repeat("222222errerr", 1000))
	err = buf2.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file2.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename3 := "testfile3"
	file3, err := os.CreateTemp(config.TmpDir, filename3)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file3.Name())

	buf3 := bufio.NewWriter(file3)
	buf3.WriteString(strings.Repeat("3333abc", 1000))
	err = buf3.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file3.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	var (
		zipData1, _ = compress.CompressFile(file.Name())
		zipData2, _ = compress.CompressFile(file2.Name())
		zipData3, _ = compress.CompressFile(file3.Name())

		tmpFiles = []*structs.File{
			{
				ID:   "FILE_ID_1",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename1),
				Body: zipData1,
			},
			{
				ID:   "FILE_ID_2",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename2),
				Body: zipData2,
			},
			{
				ID:   "FILE_ID_3",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename3),
				Body: zipData3,
			},
		}
	)

	type args struct {
		ind      int
		tmpFiles []*structs.File
	}
	tests := []struct {
		name     string
		args     args
		wantFile *structs.File
		wantErr  bool
	}{
		{
			name: "Test GetByInd File 1: valid ind",
			args: args{
				ind:      1,
				tmpFiles: tmpFiles,
			},
			wantFile: tmpFiles[1],
			wantErr:  false,
		},
		{
			name: "Test GetByInd File 2: ind < 0",
			args: args{
				ind:      -10,
				tmpFiles: tmpFiles,
			},
			wantErr: true,
		},
		{
			name: "Test GetByInd File 3: ind > len(files)",
			args: args{
				ind:      100,
				tmpFiles: tmpFiles,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase, _ := NewClientUC()
			usecase.Files = tt.args.tmpFiles

			gotFile, err := usecase.GetFileByIND(tt.args.ind)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileByIND() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFile, tt.wantFile) {
				t.Errorf("GetFileByIND() gotFile = %v, want %v", gotFile, tt.wantFile)
			}
		})
	}
}

func TestClientUC_GetCreds(t *testing.T) {
	var (
		password1 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID1111",
			Date:        time.Now().Add(time.Second * -200),
			Resource:    "localhost1111",
			Login:       "my_favorite_username1111",
			Password:    "my_favorite_password1111",
			Description: "some description1111",
		}
		password2 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "ID2222",
			Date:        time.Now().Add(time.Second * -500),
			Resource:    "localhost2222",
			Login:       "admin2222",
			Password:    "secret password2222",
			Description: "some new description2222",
		}
		password3 = &structs.Credential{
			Type:        structs.BlobCred,
			ID:          "superID3333",
			Date:        time.Now(),
			Resource:    "localhost3333",
			Login:       "my_favorite_username333",
			Password:    "my_favorite_password3333",
			Description: "some description3333",
		}

		tmpCreds = []*structs.Credential{
			password1,
			password2,
			password3,
		}
	)

	usecase, _ := NewClientUC()
	usecase.Creds = tmpCreds

	got := usecase.GetCreds()

	if !reflect.DeepEqual(got, tmpCreds) {
		t.Errorf("GetCreds() = %v, want %v", got, tmpCreds)
	}

}

func TestClientUC_GetCards(t *testing.T) {
	var (
		tmpCards = []*structs.Card{
			{
				ID:          "ID_CARD_1111",
				Type:        structs.BlobCard,
				Name:        "test1",
				Bank:        entities.Banks[0],
				Person:      "string",
				Number:      122222222222,
				CVC:         232,
				Expiration:  time.Now(),
				PIN:         3333,
				Description: "test description only",
			},
			{
				ID:          "ID_CARD_22222",
				Type:        structs.BlobCard,
				Name:        "test333331",
				Bank:        entities.Banks[2],
				Person:      "Major Tom",
				Number:      234872398472,
				CVC:         23244444,
				Expiration:  time.Now().Add(time.Second * 100),
				PIN:         11111,
				Description: "test description2",
			},
			{
				ID:          "ID_CARD_33",
				Type:        structs.BlobCard,
				Name:        "super secret card",
				Bank:        entities.Banks[4],
				Person:      "Major Jerry",
				Number:      2348723984721111,
				CVC:         232444443333,
				Expiration:  time.Now().Add(time.Second * -500),
				PIN:         2323,
				Description: "test myself",
			},
		}
	)
	usecase, _ := NewClientUC()
	usecase.Cards = tmpCards

	if got := usecase.GetCards(); !reflect.DeepEqual(got, tmpCards) {
		t.Errorf("GetCards() = %v, want %v", got, tmpCards)
	}

}

func TestClientUC_GetNotes(t *testing.T) {

	var (
		testNotes = []*structs.Note{
			{
				ID:   "ID_NOTE_1",
				Type: structs.BlobNote,
				Name: "test1",
				Date: time.Now().Add(time.Second * -300000),
				Body: "Hello\nWorld!",
			},
			{
				ID:   "ID_NOTE_2",
				Type: structs.BlobNote,
				Name: "HELLO 222222",
				Date: time.Now().Add(time.Second * -3000010),
				Body: "Hello\nWorld! 9234928309482390480298340923809840",
			},
			{
				ID:   "ID_NOTE_3",
				Type: structs.BlobNote,
				Name: "New Test Blob",
				Date: time.Now().Add(time.Second * -500000),
				Body: "Hello\nWorld! Amigo",
			},
		}
	)

	usecase, _ := NewClientUC()
	usecase.Notes = testNotes

	if got := usecase.GetNotes(); !reflect.DeepEqual(got, testNotes) {
		t.Errorf("GetNotes() = %v, want %v", got, testNotes)
	}

}

func TestClientUC_GetFiles(t *testing.T) {
	filename1 := "testfile1"
	file, err := os.CreateTemp(config.TmpDir, filename1)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file.Name())

	buf := bufio.NewWriter(file)
	buf.WriteString(strings.Repeat("123123", 1000))
	err = buf.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename2 := "testfile2"
	file2, err := os.CreateTemp(config.TmpDir, filename2)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file2.Name())

	buf2 := bufio.NewWriter(file2)
	buf2.WriteString(strings.Repeat("222222errerr", 1000))
	err = buf2.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file2.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename3 := "testfile3"
	file3, err := os.CreateTemp(config.TmpDir, filename3)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file3.Name())

	buf3 := bufio.NewWriter(file3)
	buf3.WriteString(strings.Repeat("3333abc", 1000))
	err = buf3.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file3.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	var (
		zipData1, _ = compress.CompressFile(file.Name())
		zipData2, _ = compress.CompressFile(file2.Name())
		zipData3, _ = compress.CompressFile(file3.Name())

		testFiles = []*structs.File{
			{
				ID:   "FILE_ID_1",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename1),
				Body: zipData1,
			},
			{
				ID:   "FILE_ID_2",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename2),
				Body: zipData2,
			},
			{
				ID:   "FILE_ID_3",
				Type: structs.BlobFile,
				Name: filesUC.GenFileShortName(filename3),
				Body: zipData3,
			},
		}
	)

	usecase, _ := NewClientUC()
	usecase.Files = testFiles

	if got := usecase.GetFiles(); !reflect.DeepEqual(got, testFiles) {
		t.Errorf("GetFiles() = %v, want %v", got, testFiles)
	}

}
