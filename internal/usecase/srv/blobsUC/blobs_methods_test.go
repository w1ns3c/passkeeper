package blobsUC

import (
	"context"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
	"passkeeper/internal/entities/structs"
	"passkeeper/mocks/mockstorage"
)

func TestBlobUsecase_AddBlob(t *testing.T) {
	type args struct {
		userID string
		blob   *structs.CryptoBlob
	}

	var (
		ctx = context.Background()
		log = zerolog.New(io.Discard)

		userID1 = "1111"
		userID2 = "2222"
		cred1   = &structs.Note{
			Type: structs.BlobNote,
			ID:   "1111111",
			Name: "simple",
			Date: time.Now(),
			Body: "some secret note",
		}

		secret1, _ = hashes.GenerateSecret(config.UserSecretLen)
		blob1, _   = hashes.EncryptBlob(cred1, secret1)
	)

	blob1.UserID = userID1

	tests := []struct {
		name string

		args    args
		prepare func(storage *mocks.MockBlobStorage)
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test 1: success",
			args: args{
				userID: userID1,
				blob:   blob1,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
					storage.EXPECT().AddBlob(ctx, gomock.Any()).Return(nil),
				)
			},
			wantErr: false,
		},
		{
			name: "Test 2: success",
			args: args{
				userID: userID2,
				blob:   blob1,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mocks.NewMockBlobStorage(ctrl)

			if &tt.prepare != nil {
				tt.prepare(store)
			}

			var u = &BlobUsecase{
				ctx:     ctx,
				storage: store,
				log:     &log,
			}

			if err := u.AddBlob(ctx, tt.args.userID, tt.args.blob); (err != nil) != tt.wantErr {
				t.Errorf("AddBlob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlobUsecase_DelBlob(t *testing.T) {
	type args struct {
		userID string
		blobID string
	}

	var (
		ctx = context.Background()
		log = zerolog.New(io.Discard)

		userID1 = "1111"
		userID2 = "2222"
		cred1   = &structs.Note{
			Type: structs.BlobNote,
			ID:   "1111111",
			Name: "simple",
			Date: time.Now(),
			Body: "some secret note",
		}

		secret1, _ = hashes.GenerateSecret(config.UserSecretLen)
		blob1, _   = hashes.EncryptBlob(cred1, secret1)
	)

	blob1.UserID = userID1

	tests := []struct {
		name string

		args    args
		prepare func(storage *mocks.MockBlobStorage)
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test 1 Del: success",
			args: args{
				userID: userID1,
				blobID: blob1.ID,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
					storage.EXPECT().GetBlob(ctx, userID1, blob1.ID).Return(blob1, nil),
					storage.EXPECT().DeleteBlob(ctx, userID1, blob1.ID).Return(nil),
				)
			},
			wantErr: false,
		},
		{
			name: "Test 2 Del: can't return blob (not exist)",
			args: args{
				userID: userID1,
				blobID: blob1.ID,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
					storage.EXPECT().GetBlob(ctx, userID1, blob1.ID).Return(nil, myerrors.ErrBlobNotFound),
				)
			},
			wantErr: true,
		},
		{
			name: "Test 3 Del: store return wrong blob",
			args: args{
				userID: userID2,
				blobID: blob1.ID,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
					storage.EXPECT().GetBlob(ctx, userID2, blob1.ID).Return(blob1, nil),
					//storage.EXPECT().DeleteBlob(ctx, userID1, blob1.ID).Return(nil),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mocks.NewMockBlobStorage(ctrl)

			if &tt.prepare != nil {
				tt.prepare(store)
			}

			var u = &BlobUsecase{
				ctx:     ctx,
				storage: store,
				log:     &log,
			}

			if err := u.DelBlob(ctx, tt.args.userID, tt.args.blobID); (err != nil) != tt.wantErr {
				t.Errorf("DelBlob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlobUsecase_GetBlob(t *testing.T) {
	type args struct {
		userID string
		blobID string
		blob   *structs.CryptoBlob
	}

	var (
		ctx = context.Background()
		log = zerolog.New(io.Discard)

		userID1 = "1111"
		userID2 = "2222"
		cred1   = &structs.Note{
			Type: structs.BlobNote,
			ID:   "1111111",
			Name: "simple",
			Date: time.Now(),
			Body: "some secret note",
		}

		secret1, _ = hashes.GenerateSecret(config.UserSecretLen)
		blob1, _   = hashes.EncryptBlob(cred1, secret1)
	)

	blob1.UserID = userID1

	tests := []struct {
		name string

		args    args
		prepare func(storage *mocks.MockBlobStorage)
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test 1 Get: success",
			args: args{
				userID: userID1,
				blobID: blob1.ID,
				blob:   blob1,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
					storage.EXPECT().GetBlob(ctx, userID1, blob1.ID).Return(blob1, nil),
				)
			},
			wantErr: false,
		},
		{
			name: "Test 2 Get: can't return blob (not exist)",
			args: args{
				userID: userID1,
				blobID: blob1.ID,
				blob:   nil,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
					storage.EXPECT().GetBlob(ctx, userID1, blob1.ID).Return(nil, myerrors.ErrBlobNotFound),
				)
			},
			wantErr: true,
		},
		{
			name: "Test 3 Get: store return wrong blob",
			args: args{
				userID: userID2,
				blobID: blob1.ID,
				blob:   blob1,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
					storage.EXPECT().GetBlob(ctx, userID2, blob1.ID).Return(blob1, nil),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mocks.NewMockBlobStorage(ctrl)

			if &tt.prepare != nil {
				tt.prepare(store)
			}

			var u = &BlobUsecase{
				ctx:     ctx,
				storage: store,
				log:     &log,
			}

			gotBlob, err := u.GetBlob(ctx, tt.args.userID, tt.args.blobID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DelBlob() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotBlob != nil {
				require.Equal(t, tt.args.blob.Blob, gotBlob.Blob)
				require.Equal(t, tt.args.blob.UserID, gotBlob.UserID)
				require.Equal(t, tt.args.blob.ID, gotBlob.ID)
			}

		})
	}
}

func TestBlobUsecase_UpdBlob(t *testing.T) {
	type args struct {
		userID string
		blob   *structs.CryptoBlob
	}

	var (
		ctx = context.Background()
		log = zerolog.New(io.Discard)

		userID1 = "1111"
		userID2 = "2222"
		cred1   = &structs.Note{
			Type: structs.BlobNote,
			ID:   "1111111",
			Name: "simple",
			Date: time.Now(),
			Body: "some secret note",
		}

		secret1, _ = hashes.GenerateSecret(config.UserSecretLen)
		blob1, _   = hashes.EncryptBlob(cred1, secret1)
	)

	blob1.UserID = userID1

	tests := []struct {
		name    string
		args    args
		prepare func(storage *mocks.MockBlobStorage)
		wantErr bool
	}{
		{
			name: "Test 1 Upd: success",
			args: args{
				userID: userID1,
				blob:   blob1,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
					storage.EXPECT().UpdateBlob(ctx, blob1).Return(nil),
				)
			},
			wantErr: false,
		},
		{
			name: "Test 2 Upd: userID != blob.UserID",
			args: args{
				userID: userID2,
				blob:   blob1,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
				//storage.EXPECT().UpdateBlob(ctx, blob1.ID).Return(nil, myerrors.ErrBlobNotFound),
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mocks.NewMockBlobStorage(ctrl)

			if &tt.prepare != nil {
				tt.prepare(store)
			}

			var u = &BlobUsecase{
				ctx:     ctx,
				storage: store,
				log:     &log,
			}

			if err := u.UpdBlob(ctx, tt.args.userID, tt.args.blob); (err != nil) != tt.wantErr {
				t.Errorf("UpdBlob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlobUsecase_ListBlobs(t *testing.T) {
	var (
		ctx = context.Background()
		log = zerolog.New(io.Discard)

		userID1 = "1111"
		userID2 = "2222"
		cred1   = &structs.Note{
			Type: structs.BlobNote,
			ID:   "1111111",
			Name: "simple",
			Date: time.Now(),
			Body: "some secret note",
		}

		cred2 = &structs.Note{
			Type: structs.BlobNote,
			ID:   "222222",
			Name: "new simple",
			Date: time.Now(),
			Body: "some secret note V2",
		}

		secret1, _ = hashes.GenerateSecret(config.UserSecretLen)
		blob1, _   = hashes.EncryptBlob(cred1, secret1)
		blob2, _   = hashes.EncryptBlob(cred2, secret1)
	)

	blob1.UserID = userID1
	blob2.UserID = userID2

	tests := []struct {
		name      string
		userID    string
		wantBlobs []*structs.CryptoBlob
		prepare   func(storage *mocks.MockBlobStorage)
		wantErr   bool
	}{
		{
			name:   "Test 1 List: success",
			userID: userID1,
			wantBlobs: []*structs.CryptoBlob{
				blob1, blob2,
			},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
					storage.EXPECT().GetAllBlobs(ctx, userID1).Return([]*structs.CryptoBlob{
						blob1, blob2,
					}, nil),
				)
			},
			wantErr: false,
		},
		{
			name:      "Test 2 List: error, not exist",
			userID:    userID2,
			wantBlobs: []*structs.CryptoBlob{},
			prepare: func(storage *mocks.MockBlobStorage) {
				gomock.InOrder(
					storage.EXPECT().GetAllBlobs(ctx, userID2).Return(nil, myerrors.ErrBlobNotFound),
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mocks.NewMockBlobStorage(ctrl)

			if &tt.prepare != nil {
				tt.prepare(store)
			}

			var u = &BlobUsecase{
				ctx:     ctx,
				storage: store,
				log:     &log,
			}

			gotBlobs, err := u.ListBlobs(ctx, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListBlobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.wantBlobs) == 0 {
				return
			}
			if !reflect.DeepEqual(gotBlobs, tt.wantBlobs) {
				t.Errorf("ListBlobs() gotBlobs = %v, want %v", gotBlobs, tt.wantBlobs)
			}
		})
	}
}
