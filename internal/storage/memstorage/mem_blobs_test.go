package memstorage

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"passkeeper/internal/entities/structs"
)

func TestMemStorage_AddBlob(t *testing.T) {
	var (
		ctx = context.Background()

		bl1 = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     "123",
				UserID: "1",
				Blob:   "22222",
			},
			&structs.CryptoBlob{
				ID:     "1234",
				UserID: "1",
				Blob:   "33333",
			},
		}
		blob1 = &structs.CryptoBlob{
			ID:     "123",
			UserID: "1",
			Blob:   "22222",
		}

		bl2 = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     "123333",
				UserID: "2",
				Blob:   "2222322",
			},
			&structs.CryptoBlob{
				ID:     "12234334",
				UserID: "2",
				Blob:   "33333",
			},
			&structs.CryptoBlob{
				ID:     "1234444",
				UserID: "2",
				Blob:   "33333555",
			},
		}

		invalidBlob = &structs.CryptoBlob{
			ID:     "66666",
			UserID: "ffff",
			Blob:   "777777777",
		}
	)

	type args struct {
		blobs map[string][]*structs.CryptoBlob
		blob  *structs.CryptoBlob
	}

	tests := []struct {
		name      string
		args      args
		wantBlobs map[string][]*structs.CryptoBlob
		wantErr   bool
	}{
		{
			name: "Test Add blob 1: valid",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1,
					"2": bl2,
				},
				blob: blob1,
			},
			wantBlobs: map[string][]*structs.CryptoBlob{
				"1": append(bl1, blob1),
				"2": bl2,
			},
			wantErr: false,
		},
		{
			name: "Test Add blob 2: invalid userID",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1,
					"2": bl2,
				},
				blob: invalidBlob,
			},
			wantBlobs: map[string][]*structs.CryptoBlob{
				"1": append(bl1, blob1),
				"2": bl2,
				invalidBlob.UserID: []*structs.CryptoBlob{
					invalidBlob,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				blobs: tt.args.blobs,
			}
			m.Init(ctx)

			err := m.AddBlob(ctx, tt.args.blob)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddBlob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			require.Equal(t, len(tt.args.blobs), len(tt.wantBlobs))
			require.Equal(t, tt.args.blobs[tt.args.blob.UserID], tt.wantBlobs[tt.args.blob.UserID])

		})
	}
}

func TestMemStorage_DeleteBlob(t *testing.T) {
	var (
		ctx   = context.Background()
		blob1 = &structs.CryptoBlob{
			ID:     "123444",
			UserID: "1",
			Blob:   "22222",
		}

		blobID1 = "1234"
		bl1_1   = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     blobID1,
				UserID: "1",
				Blob:   "33333",
			},
			blob1,
		}

		bl1_2 = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     "123333",
				UserID: "2",
				Blob:   "2222322",
			},
			&structs.CryptoBlob{
				ID:     "12234334",
				UserID: "2",
				Blob:   "33333",
			},
			&structs.CryptoBlob{
				ID:     "1234444",
				UserID: "2",
				Blob:   "33333555",
			},
		}

		bl2_1 = []*structs.CryptoBlob{
			blob1,
			&structs.CryptoBlob{
				ID:     blobID1,
				UserID: "1",
				Blob:   "33333",
			},
		}

		bl2_2 = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     "123333",
				UserID: "2",
				Blob:   "2222322",
			},
			&structs.CryptoBlob{
				ID:     "12234334",
				UserID: "2",
				Blob:   "33333",
			},
			&structs.CryptoBlob{
				ID:     "1234444",
				UserID: "2",
				Blob:   "33333555",
			},
		}

		bl3_1 = []*structs.CryptoBlob{
			blob1,
			&structs.CryptoBlob{
				ID:     blobID1,
				UserID: "1",
				Blob:   "33333",
			},
		}

		bl3_2 = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     "123333",
				UserID: "2",
				Blob:   "2222322",
			},
			&structs.CryptoBlob{
				ID:     "12234334",
				UserID: "2",
				Blob:   "33333",
			},
			&structs.CryptoBlob{
				ID:     "1234444",
				UserID: "2",
				Blob:   "33333555",
			},
		}
	)

	type args struct {
		blobs  map[string][]*structs.CryptoBlob
		userID string
		blobID string
	}

	tests := []struct {
		name      string
		args      args
		wantBlobs map[string][]*structs.CryptoBlob
		wantErr   bool
	}{
		{
			name: "Test Del blob 1: valid",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1_1,
					"2": bl1_2,
				},
				userID: "1",
				blobID: blobID1,
			},
			wantBlobs: map[string][]*structs.CryptoBlob{
				"1": []*structs.CryptoBlob{
					blob1,
				},
				"2": bl1_2,
			},
			wantErr: false,
		},
		{
			name: "Test Del blob 2: invalid blobID",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl2_1,
					"2": bl2_2,
				},
				userID: "1",
				blobID: blobID1 + "ffff",
			},
			wantErr: true,
		},
		{
			name: "Test Del blob 3: invalid userID",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl3_1,
					"2": bl3_2,
				},
				userID: "ffff",
				blobID: blobID1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				blobs: tt.args.blobs,
			}
			m.Init(ctx)

			err := m.DeleteBlob(ctx, tt.args.userID, tt.args.blobID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBlob() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			require.Equal(t, len(tt.args.blobs), len(tt.wantBlobs))
			require.Equal(t, tt.args.blobs[tt.args.userID], tt.wantBlobs[tt.args.userID])
		})
	}
}

func TestMemStorage_GetAllBlobs(t *testing.T) {

	type args struct {
		blobs  map[string][]*structs.CryptoBlob
		userID string
	}

	ctx := context.Background()

	var (
		bl1 = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     "123",
				UserID: "1",
				Blob:   "22222",
			},
			&structs.CryptoBlob{
				ID:     "1234",
				UserID: "1",
				Blob:   "33333",
			},
		}

		bl2 = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     "123333",
				UserID: "2",
				Blob:   "2222322",
			},
			&structs.CryptoBlob{
				ID:     "12234334",
				UserID: "2",
				Blob:   "33333",
			},
			&structs.CryptoBlob{
				ID:     "1234444",
				UserID: "2",
				Blob:   "33333555",
			},
		}
	)
	tests := []struct {
		name      string
		args      args
		wantBlobs []*structs.CryptoBlob
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name: "Test 1: get all blobs valid",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1,
					"2": bl2,
				},
				userID: "1",
			},
			wantBlobs: bl1,
			wantErr:   false,
		},
		{
			name: "Test 2: err get all blobs valid",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1,
					"2": bl2,
				},
				userID: "3",
			},
			wantBlobs: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				blobs: tt.args.blobs,
			}
			m.Init(ctx)

			gotBlobs, err := m.GetAllBlobs(ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllBlobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBlobs, tt.wantBlobs) {
				t.Errorf("GetAllBlobs() gotBlobs = %v, want %v", gotBlobs, tt.wantBlobs)
			}
		})
	}
}

func TestMemStorage_GetBlob(t *testing.T) {
	var (
		ctx     = context.Background()
		blobID1 = "11111111fff"
		blob1   = &structs.CryptoBlob{
			ID:     blobID1,
			UserID: "1",
			Blob:   "22222",
		}

		blobID2 = "1234"
		bl1     = []*structs.CryptoBlob{
			blob1,
			&structs.CryptoBlob{
				ID:     blobID2,
				UserID: "1",
				Blob:   "33333",
			},
		}

		bl2 = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     "123333",
				UserID: "2",
				Blob:   "2222322",
			},
			&structs.CryptoBlob{
				ID:     "12234334",
				UserID: "2",
				Blob:   "33333",
			},
			&structs.CryptoBlob{
				ID:     "1234444",
				UserID: "2",
				Blob:   "33333555",
			},
		}
	)

	type args struct {
		blobs  map[string][]*structs.CryptoBlob
		userID string
		blobID string
	}

	tests := []struct {
		name     string
		args     args
		wantBlob *structs.CryptoBlob
		wantErr  bool
	}{
		{
			name: "Test Get blob 1: valid",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1,
					"2": bl2,
				},
				userID: "1",
				blobID: blobID1,
			},
			wantBlob: blob1,
			wantErr:  false,
		},
		{
			name: "Test Get blob 2: invalid blobID",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1,
					"2": bl2,
				},
				userID: "1",
				blobID: blobID1 + "ffff",
			},
			wantErr: true,
		},
		{
			name: "Test Get blob 3: invalid userID",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1,
					"2": bl2,
				},
				userID: "3",
				blobID: blobID1 + "ffff",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				blobs: tt.args.blobs,
			}
			m.Init(ctx)

			gotBlob, err := m.GetBlob(ctx, tt.args.userID, tt.args.blobID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlob() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			require.Equal(t, tt.wantBlob, gotBlob)
		})
	}
}

func TestMemStorage_UpdateBlob(t *testing.T) {
	var (
		ctx = context.Background()

		blob1 = &structs.CryptoBlob{
			ID:     "123",
			UserID: "1",
			Blob:   "1231231212312312",
		}
		blob2 = &structs.CryptoBlob{
			ID:     "33333",
			UserID: "1",
			Blob:   "11111111",
		}

		bl1 = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     "123",
				UserID: "1",
				Blob:   "22222",
			},
			blob2,
		}
		want1 = []*structs.CryptoBlob{
			blob1, blob2,
		}

		blobInvID = &structs.CryptoBlob{
			ID:     "invalid",
			UserID: "1",
			Blob:   "44444444444444",
		}

		blobInvUserID = &structs.CryptoBlob{
			ID:     blob1.ID,
			UserID: "4444444",
			Blob:   "44444444444444",
		}

		bl2 = []*structs.CryptoBlob{
			&structs.CryptoBlob{
				ID:     "123333",
				UserID: "2",
				Blob:   "2222322",
			},
			&structs.CryptoBlob{
				ID:     "12234334",
				UserID: "2",
				Blob:   "33333",
			},
			&structs.CryptoBlob{
				ID:     "1234444",
				UserID: "2",
				Blob:   "33333555",
			},
		}
	)

	type args struct {
		blobs map[string][]*structs.CryptoBlob
		blob  *structs.CryptoBlob
	}

	tests := []struct {
		name      string
		args      args
		wantBlobs map[string][]*structs.CryptoBlob
		wantErr   bool
	}{
		{
			name: "Test Update blob 1: valid",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1,
					"2": bl2,
				},
				blob: blob1,
			},
			wantBlobs: map[string][]*structs.CryptoBlob{
				"1": want1,
				"2": bl2,
			},
			wantErr: false,
		},
		{
			name: "Test Update blob 2: invalid blobID",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1,
					"2": bl2,
				},
				blob: blobInvID,
			},
			wantBlobs: map[string][]*structs.CryptoBlob{
				"1": bl1,
				"2": bl2,
			},
			wantErr: true,
		},
		{
			name: "Test Update blob 3: invalid userID",
			args: args{
				blobs: map[string][]*structs.CryptoBlob{
					"1": bl1,
					"2": bl2,
				},
				blob: blobInvUserID,
			},
			wantBlobs: map[string][]*structs.CryptoBlob{
				"1": bl1,
				"2": bl2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				blobs: tt.args.blobs,
			}
			m.Init(ctx)

			err := m.UpdateBlob(ctx, tt.args.blob)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBlob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, len(tt.args.blobs), len(tt.wantBlobs))
			require.Equal(t, tt.args.blobs[tt.args.blob.UserID], tt.wantBlobs[tt.args.blob.UserID])

		})
	}
}
