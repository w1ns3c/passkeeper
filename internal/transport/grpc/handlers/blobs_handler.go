package handlers

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/srv/blobsUC"
)

var (
	ErrCredAddMsg = "cred not added"
	ErrCredAdd    = status.Error(codes.Internal, ErrCredAddMsg)

	ErrCredGetMsg = "cred can't get"
	ErrCredGet    = status.Error(codes.Internal, ErrCredGetMsg)

	ErrCredUpdMsg = "cred not updated"
	ErrCredUpd    = status.Error(codes.Internal, ErrCredUpdMsg)

	ErrCredDelMsg = "cred not deleted"
	ErrCredDel    = status.Error(codes.Internal, ErrCredDelMsg)

	ErrCredListMsg = "creds not listed"
	ErrCredList    = status.Error(codes.Internal, ErrCredListMsg)
)

type BlobsHandler struct {
	pb.UnimplementedBlobSvcServer
	service blobsUC.BlobUsecaseInf
	log     *zerolog.Logger
}

func NewBlobsHandler(logger *zerolog.Logger, service blobsUC.BlobUsecaseInf) *BlobsHandler {
	return &BlobsHandler{
		UnimplementedBlobSvcServer: pb.UnimplementedBlobSvcServer{},
		service:                    service,
		log:                        logger,
	}
}

func (h *BlobsHandler) BlobAdd(ctx context.Context, req *pb.BlobAddRequest) (*empty.Empty, error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	blob := &entities.CryptoBlob{
		ID:     req.Cred.ID,
		UserID: userID,
		Blob:   req.Cred.Blob,
	}

	err = h.service.AddBlob(ctx, userID, blob)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredAddMsg)

		return nil, ErrCredAdd
	}

	return new(empty.Empty), nil
}

func (h *BlobsHandler) BlobGet(ctx context.Context, req *pb.BlobGetRequest) (resp *pb.BlobGetResponse, err error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	cr, err := h.service.GetBlob(ctx, userID, req.CredID)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredGetMsg)

		return nil, ErrCredGet
	}

	resp = &pb.BlobGetResponse{
		Cred: &pb.CryptoBlob{
			ID:   cr.ID,
			Blob: cr.Blob,
		},
	}
	return resp, nil
}

func (h *BlobsHandler) BlobUpd(ctx context.Context, req *pb.BlobUpdRequest) (*empty.Empty, error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	blob := &entities.CryptoBlob{
		ID:     req.Blob.ID,
		UserID: userID,
		Blob:   req.Blob.Blob,
	}

	err = h.service.UpdBlob(ctx, userID, blob)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredUpdMsg)

		return nil, ErrCredUpd
	}

	return new(empty.Empty), nil
}

func (h *BlobsHandler) BlobDel(ctx context.Context, req *pb.BlobDelRequest) (*empty.Empty, error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	err = h.service.DelBlob(ctx, userID, req.CredID)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredDelMsg)

		return nil, ErrCredDel
	}

	return new(empty.Empty), nil
}

func (h *BlobsHandler) BlobList(ctx context.Context, req *empty.Empty) (resp *pb.BlobListResponse, err error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	h.log.Info().
		Msgf("User \"%s\" request blobs list", userID)

	blobs, err := h.service.ListBlobs(ctx, userID)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredListMsg)

		return nil, ErrCredList
	}

	h.log.Info().
		Msgf("User \"%s\" have: %d blobs", userID, len(blobs))

	resp = &pb.BlobListResponse{
		Blobs: make([]*pb.CryptoBlob, len(blobs)),
	}
	for i := 0; i < len(blobs); i++ {
		resp.Blobs[i] = &pb.CryptoBlob{
			ID:   blobs[i].ID,
			Blob: blobs[i].Blob,
		}
	}

	return resp, nil
}
