package handlers

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/myerrors"
	"passkeeper/internal/entities/structs"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	mocks "passkeeper/mocks/usecase/blobs_usecase"
)

func TestBlobsHandler_BlobAdd(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.BlobAddRequest
	}
	var (
		logger = zerolog.New(io.Discard)
	)

	tests := []struct {
		name    string
		args    args
		prepare func(usecase *mocks.MockBlobUsecaseInf)
		want    *empty.Empty
		wantErr error
	}{
		{
			name: "BlobAdd Test1: success",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "user_ID1"})),
				req: &pb.BlobAddRequest{
					Cred: &pb.CryptoBlob{
						ID:   "11123123123",
						Blob: "33333333333",
					},
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {
				gomock.InOrder(
					usecase.EXPECT().AddBlob(gomock.Any(), "user_ID1", &structs.CryptoBlob{
						ID:     "11123123123",
						UserID: "user_ID1",
						Blob:   "33333333333",
					}).Return(nil),
				)
			},
			want:    new(empty.Empty),
			wantErr: nil,
		},
		{
			name: "BlobAdd Test2: invalid user in token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "user_1"})),
				req: &pb.BlobAddRequest{
					Cred: &pb.CryptoBlob{
						ID:   "11123123123",
						Blob: "33333333333",
					},
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {
				gomock.InOrder(
					usecase.EXPECT().AddBlob(gomock.Any(), gomock.Any(), &structs.CryptoBlob{
						ID:     "11123123123",
						UserID: "user_1",
						Blob:   "33333333333",
					}).Return(myerrors.ErrBlobAdd),
				)
			},
			want:    nil,
			wantErr: myerrors.ErrBlobAdd,
		},
		{
			name: "BlobAdd Test3: Empty token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
				req: &pb.BlobAddRequest{
					Cred: &pb.CryptoBlob{
						ID:   "11123123123",
						Blob: "33333333333",
					},
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {

			},
			want:    nil,
			wantErr: myerrors.ErrEmptyToken,
		},
		{
			name: "BlobAdd Test4: No token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), nil),
				req: &pb.BlobAddRequest{
					Cred: &pb.CryptoBlob{
						ID:   "11123123123",
						Blob: "33333333333",
					},
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {

			},
			want:    nil,
			wantErr: myerrors.ErrEmptyToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockBlobUsecaseInf(ctrl)
			if tt.prepare != nil {
				tt.prepare(service)
			}

			h := &BlobsHandler{
				service: service,
				log:     &logger,
			}
			got, err := h.BlobAdd(tt.args.ctx, tt.args.req)
			if err != nil {
				if !errors.Is(tt.wantErr, err) {
					t.Errorf("BlobAdd() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BlobAdd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlobsHandler_BlobDel(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.BlobDelRequest
	}
	var (
		logger = zerolog.New(io.Discard)
	)

	tests := []struct {
		name    string
		args    args
		prepare func(usecase *mocks.MockBlobUsecaseInf)
		want    *empty.Empty
		wantErr error
	}{
		{
			name: "BlobDel Test1: success",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "user_ID1"})),
				req: &pb.BlobDelRequest{
					CredID: "cred_ID1",
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {
				gomock.InOrder(
					usecase.EXPECT().
						DelBlob(gomock.Any(), "user_ID1", "cred_ID1").
						Return(nil),
				)
			},
			want:    new(empty.Empty),
			wantErr: nil,
		},
		{
			name: "BlobDel Test2: error",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "user_ID1"})),
				req: &pb.BlobDelRequest{
					CredID: "cred_ID1",
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {
				gomock.InOrder(
					usecase.EXPECT().
						DelBlob(gomock.Any(), "user_ID1", "cred_ID1").
						Return(myerrors.ErrBlobDel),
				)
			},
			want:    nil,
			wantErr: myerrors.ErrBlobDel,
		},
		{
			name: "BlobDel Test3: Empty token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
				req: &pb.BlobDelRequest{
					CredID: "cred_ID1",
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {

			},
			want:    nil,
			wantErr: myerrors.ErrEmptyToken,
		},
		{
			name: "BlobDel Test4: No token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), nil),
				req: &pb.BlobDelRequest{
					CredID: "cred_ID1",
				},
			},
			want:    nil,
			wantErr: myerrors.ErrEmptyToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockBlobUsecaseInf(ctrl)
			if tt.prepare != nil {
				tt.prepare(service)
			}

			h := &BlobsHandler{
				service: service,
				log:     &logger,
			}
			got, err := h.BlobDel(tt.args.ctx, tt.args.req)
			if err != nil {
				if !errors.Is(tt.wantErr, err) {
					t.Errorf("BlobDel() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BlobDel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlobsHandler_BlobGet(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.BlobGetRequest
	}
	var (
		logger = zerolog.New(io.Discard)
	)

	tests := []struct {
		name     string
		args     args
		prepare  func(usecase *mocks.MockBlobUsecaseInf)
		wantResp *pb.BlobGetResponse
		wantErr  error
	}{
		{
			name: "BlobGet Test1: success",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "user_ID1"})),
				req: &pb.BlobGetRequest{
					CredID: "cred_ID1",
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {
				gomock.InOrder(
					usecase.EXPECT().
						GetBlob(gomock.Any(), "user_ID1", "cred_ID1").
						Return(&structs.CryptoBlob{
							ID:     "123123",
							UserID: "444444",
							Blob:   "000000",
						}, nil),
				)
			},
			wantResp: &pb.BlobGetResponse{
				Cred: &pb.CryptoBlob{
					ID:   "123123",
					Blob: "000000",
				},
			},
			wantErr: nil,
		},
		{
			name: "BlobGet Test2: error",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "user_ID1"})),
				req: &pb.BlobGetRequest{
					CredID: "cred_ID1",
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {
				gomock.InOrder(
					usecase.EXPECT().
						GetBlob(gomock.Any(), "user_ID1", "cred_ID1").
						Return(nil, myerrors.ErrBlobGet),
				)
			},
			wantResp: nil,
			wantErr:  myerrors.ErrBlobGet,
		},
		{
			name: "BlobAdd Test3: Empty token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
				req: &pb.BlobGetRequest{
					CredID: "cred_ID1",
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {

			},
			wantResp: nil,
			wantErr:  myerrors.ErrEmptyToken,
		},
		{
			name: "BlobAdd Test4: No token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), nil),
				req: &pb.BlobGetRequest{
					CredID: "cred_ID1",
				},
			},
			wantResp: nil,
			wantErr:  myerrors.ErrEmptyToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockBlobUsecaseInf(ctrl)
			if tt.prepare != nil {
				tt.prepare(service)
			}

			h := &BlobsHandler{
				service: service,
				log:     &logger,
			}
			gotResp, err := h.BlobGet(tt.args.ctx, tt.args.req)
			if err != nil {
				if !errors.Is(tt.wantErr, err) {
					t.Errorf("BlobGet() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("BlobGet() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestBlobsHandler_BlobList(t *testing.T) {
	type args struct {
		ctx context.Context
		req *empty.Empty
	}
	var (
		logger = zerolog.New(io.Discard)
	)

	tests := []struct {
		name     string
		args     args
		prepare  func(usecase *mocks.MockBlobUsecaseInf)
		wantResp *pb.BlobListResponse
		wantErr  error
	}{
		{
			name: "BlobList Test1: success",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "user_ID1"})),
				req: new(empty.Empty),
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {
				gomock.InOrder(
					usecase.EXPECT().
						ListBlobs(gomock.Any(), "user_ID1").
						Return([]*structs.CryptoBlob{
							{
								ID:     "123123",
								UserID: "user_ID1",
								Blob:   "000000",
							},
							{
								ID:     "12311111123",
								UserID: "user_ID1",
								Blob:   "000022222200",
							},
						}, nil),
				)
			},
			wantResp: &pb.BlobListResponse{
				Blobs: []*pb.CryptoBlob{
					{
						ID:   "123123",
						Blob: "000000",
					},
					{
						ID:   "12311111123",
						Blob: "000022222200",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "BlobList Test2: errors",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "user_ID1"})),
				req: new(empty.Empty),
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {
				gomock.InOrder(
					usecase.EXPECT().
						ListBlobs(gomock.Any(), "user_ID1").
						Return(nil, myerrors.ErrBlobList),
				)
			},
			wantResp: nil,
			wantErr:  myerrors.ErrBlobList,
		},
		{
			name: "BlobList Test3: Empty token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
				req: new(empty.Empty),
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {

			},
			wantResp: nil,
			wantErr:  myerrors.ErrEmptyToken,
		},
		{
			name: "BlobList Test4: No token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), nil),
				req: new(empty.Empty),
			},
			wantResp: nil,
			wantErr:  myerrors.ErrEmptyToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockBlobUsecaseInf(ctrl)
			if tt.prepare != nil {
				tt.prepare(service)
			}

			h := &BlobsHandler{
				service: service,
				log:     &logger,
			}
			gotResp, err := h.BlobList(tt.args.ctx, tt.args.req)
			if err != nil {
				if !errors.Is(tt.wantErr, err) {
					t.Errorf("BlobList() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("BlobList() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestBlobsHandler_BlobUpd(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.BlobUpdRequest
	}
	var (
		logger = zerolog.New(io.Discard)
	)

	tests := []struct {
		name    string
		args    args
		prepare func(usecase *mocks.MockBlobUsecaseInf)
		want    *empty.Empty
		wantErr error
	}{
		{
			name: "BlobEdit Test1: success",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "user_ID1"})),
				req: &pb.BlobUpdRequest{
					Blob: &pb.CryptoBlob{
						ID:   "11123123123",
						Blob: "33333333333",
					},
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {
				gomock.InOrder(
					usecase.EXPECT().UpdBlob(gomock.Any(), "user_ID1", &structs.CryptoBlob{
						ID:     "11123123123",
						UserID: "user_ID1",
						Blob:   "33333333333",
					}).Return(nil),
				)
			},
			want:    new(empty.Empty),
			wantErr: nil,
		},
		{
			name: "BlobEdit Test2: some error",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "user_1"})),
				req: &pb.BlobUpdRequest{
					Blob: &pb.CryptoBlob{
						ID:   "11123123123",
						Blob: "33333333333",
					},
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {
				gomock.InOrder(
					usecase.EXPECT().UpdBlob(gomock.Any(), gomock.Any(), &structs.CryptoBlob{
						ID:     "11123123123",
						UserID: "user_1",
						Blob:   "33333333333",
					}).Return(myerrors.ErrBlobAdd),
				)
			},
			want:    nil,
			wantErr: myerrors.ErrBlobUpd,
		},
		{
			name: "BlobEdit Test3: Empty token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
				req: &pb.BlobUpdRequest{
					Blob: &pb.CryptoBlob{
						ID:   "11123123123",
						Blob: "33333333333",
					},
				},
			},
			prepare: func(usecase *mocks.MockBlobUsecaseInf) {

			},
			want:    nil,
			wantErr: myerrors.ErrEmptyToken,
		},
		{
			name: "BlobEdit Test4: No token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), nil),
				req: &pb.BlobUpdRequest{
					Blob: &pb.CryptoBlob{
						ID:   "11123123123",
						Blob: "33333333333",
					},
				},
			},
			want:    nil,
			wantErr: myerrors.ErrEmptyToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockBlobUsecaseInf(ctrl)
			if tt.prepare != nil {
				tt.prepare(service)
			}

			h := &BlobsHandler{
				service: service,
				log:     &logger,
			}
			got, err := h.BlobUpd(tt.args.ctx, tt.args.req)
			if err != nil {
				if !errors.Is(tt.wantErr, err) {
					t.Errorf("BlobUpd() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BlobUpd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBlobsHandler(t *testing.T) {

	var (
		logger = zerolog.New(io.Discard)
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockBlobUsecaseInf(ctrl)

	tests := []struct {
		name string
		want *BlobsHandler
	}{
		{
			name: "TestNewBlobsHnd Test1",
			want: &BlobsHandler{
				service: service,
				log:     &logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := NewBlobsHandler(&logger, service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlobsHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
