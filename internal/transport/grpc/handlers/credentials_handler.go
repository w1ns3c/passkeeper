package handlers

import (
	"context"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/usecase/srv/credentialsUC"

	pb "passkeeper/internal/transport/grpc/protofiles/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

type CredsHandler struct {
	pb.UnimplementedCredSvcServer
	service credentialsUC.CredUsecaseInf
	log     *zerolog.Logger
}

func NewCredsHandler(logger *zerolog.Logger, service credentialsUC.CredUsecaseInf) *CredsHandler {
	return &CredsHandler{
		UnimplementedCredSvcServer: pb.UnimplementedCredSvcServer{},
		service:                    service,
		log:                        logger,
	}
}

func (h *CredsHandler) CredAdd(ctx context.Context, req *pb.CredAddRequest) (*empty.Empty, error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	cred := &entities.CredBlob{
		ID:     req.Cred.ID,
		UserID: userID,
		Blob:   req.Cred.Blob,
	}

	err = h.service.AddCredential(ctx, userID, cred)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredAddMsg)

		return nil, ErrCredAdd
	}

	return new(empty.Empty), nil
}

func (h *CredsHandler) CredGet(ctx context.Context, req *pb.CredGetRequest) (resp *pb.CredGetResponse, err error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	cr, err := h.service.GetCredential(ctx, userID, req.CredID)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredGetMsg)

		return nil, ErrCredGet
	}

	resp = &pb.CredGetResponse{
		Cred: &pb.CredBlob{
			ID:   cr.ID,
			Blob: cr.Blob,
		},
	}
	return resp, nil
}

func (h *CredsHandler) CredUpd(ctx context.Context, req *pb.CredUpdRequest) (*empty.Empty, error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	cred := &entities.CredBlob{
		ID:     req.Cred.ID,
		UserID: userID,
		Blob:   req.Cred.Blob,
	}

	err = h.service.UpdateCredential(ctx, userID, cred)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredUpdMsg)

		return nil, ErrCredUpd
	}

	return new(empty.Empty), nil
}

func (h *CredsHandler) CredDel(ctx context.Context, req *pb.CredDelRequest) (*empty.Empty, error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	err = h.service.DeleteCredential(ctx, userID, req.CredID)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredDelMsg)

		return nil, ErrCredDel
	}

	return new(empty.Empty), nil
}

func (h *CredsHandler) CredList(ctx context.Context, req *empty.Empty) (resp *pb.CredListResponse, err error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	h.log.Info().
		Msgf("User \"%s\" request creds list", userID)

	creds, err := h.service.ListCredentials(ctx, userID)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredListMsg)

		return nil, ErrCredList
	}

	h.log.Info().
		Msgf("User \"%s\" have: %d creds", userID, len(creds))

	resp = &pb.CredListResponse{
		Blobs: make([]*pb.CredBlob, len(creds)),
	}
	for i := 0; i < len(creds); i++ {
		resp.Blobs[i] = &pb.CredBlob{
			ID:   creds[i].ID,
			Blob: creds[i].Blob,
		}
	}

	return resp, nil
}
