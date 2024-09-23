package handlers

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"

	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
	"passkeeper/internal/entities/structs"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/srv/blobsUC"
)

// BlobsHandler handle blobs requests
type BlobsHandler struct {
	pb.UnimplementedBlobSvcServer
	service blobsUC.BlobUsecaseInf
	log     *zerolog.Logger
}

// NewBlobsHandler is a constructor for BlobsHandler
func NewBlobsHandler(logger *zerolog.Logger, service blobsUC.BlobUsecaseInf) *BlobsHandler {
	return &BlobsHandler{
		UnimplementedBlobSvcServer: pb.UnimplementedBlobSvcServer{},
		service:                    service,
		log:                        logger,
	}
}

// BlobAdd handle blob add request
func (h *BlobsHandler) BlobAdd(ctx context.Context, req *pb.BlobAddRequest) (*empty.Empty, error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	blob := &structs.CryptoBlob{
		ID:     req.Cred.ID,
		UserID: userID,
		Blob:   req.Cred.Blob,
	}

	err = h.service.AddBlob(ctx, userID, blob)
	if err != nil {
		h.log.Error().
			Err(err).Msg(myerrors.ErrCredAddMsg)

		return nil, myerrors.ErrCredAdd
	}

	return new(empty.Empty), nil
}

// BlobGet handle blob get request
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
			Err(err).Msg(myerrors.ErrCredGetMsg)

		return nil, myerrors.ErrCredGet
	}

	resp = &pb.BlobGetResponse{
		Cred: &pb.CryptoBlob{
			ID:   cr.ID,
			Blob: cr.Blob,
		},
	}
	return resp, nil
}

// BlobUpd handle blob update request
func (h *BlobsHandler) BlobUpd(ctx context.Context, req *pb.BlobUpdRequest) (*empty.Empty, error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	blob := &structs.CryptoBlob{
		ID:     req.Blob.ID,
		UserID: userID,
		Blob:   req.Blob.Blob,
	}

	err = h.service.UpdBlob(ctx, userID, blob)
	if err != nil {
		h.log.Error().
			Err(err).Msg(myerrors.ErrCredUpdMsg)

		return nil, myerrors.ErrCredUpd
	}

	return new(empty.Empty), nil
}

// BlobDel handle blob delete request
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
			Err(err).Msg(myerrors.ErrCredDelMsg)

		return nil, myerrors.ErrCredDel
	}

	return new(empty.Empty), nil
}

// BlobList handle request that ask to return all blobs
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
			Err(err).Msg(myerrors.ErrCredListMsg)

		return nil, myerrors.ErrCredList
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
