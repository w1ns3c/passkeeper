package blobsUC

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"passkeeper/mocks"
)

func TestNewBlobUCWithOpts(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx     = context.Background()
		lg      = zerolog.New(io.Discard)
		storage = mocks.NewMockStorage(ctrl)
		uc1     = NewBlobUCWithOpts(
			WithContext(ctx),
			WithLogger(&lg),
			WithStorage(storage))
	)

	tests := []struct {
		name string
		opts []BlobsOption
		want *BlobUsecase
	}{
		// TODO: Add test cases.
		{
			name: "Test 1: success",
			want: uc1,
			opts: []BlobsOption{
				WithLogger(&lg),
				WithStorage(storage),
				WithContext(ctx),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBlobUCWithOpts(tt.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlobUCWithOpts() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_newBlobUC(t *testing.T) {
	var want = &BlobUsecase{}
	tests := []struct {
		name string
		want *BlobUsecase
	}{
		// TODO: Add test cases.
		{
			name: "Test: empty blobsUC",
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newBlobUC(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newBlobUC() = %v, want %v", got, tt.want)
			}
		})
	}
}
