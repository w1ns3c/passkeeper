package handlers

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/w1nsec/passkeeper/internal/config"
	"github.com/w1nsec/passkeeper/internal/entities"
	pb "github.com/w1nsec/passkeeper/internal/transport/grpc/protofiles"
	"github.com/w1nsec/passkeeper/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrNoTokenMsg = "no token in context"
	ErrNoToken    = status.Error(codes.Unauthenticated, ErrRegisterMsg)

	ErrEmptyTokenMsg = "token is empty"
	ErrEmptyToken    = status.Error(codes.Unauthenticated, ErrEmptyTokenMsg)

	ErrCredNoAddMsg = "cred not add"
	ErrCredNoAdd    = status.Error(codes.Internal, ErrCredNoAddMsg)
)

type CredsHandler struct {
	pb.UnimplementedCredSvcServer
	service usecase.CredUsecaseInf
	log     *zerolog.Logger
}

func (h *CredsHandler) CredAdd(ctx context.Context, in *pb.CredAddRequest) (*pb.CredAddResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrNoToken
	}

	tokens := md.Get(config.TokenHeader)
	if len(tokens) == 0 {
		return nil, ErrEmptyToken
	}
	token := tokens[0]

	cred := &entities.Credential{
		Login:       in.Cred.Login,
		Password:    in.Cred.Password,
		Date:        time.Now(),
		Resource:    in.Cred.Resource,
		Description: in.Cred.Description,
	}

	err := h.service.AddCredential(ctx, token, cred)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredNoAddMsg)

		return nil, ErrCredNoAdd
	}

	return &pb.CredAddResponse{}, nil
}

func (h *CredsHandler) CredGet(ctx context.Context, in *pb.CredGetRequest) (*pb.CredGetResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrNoToken
	}

	tokens := md.Get(config.TokenHeader)
	if len(tokens) == 0 {
		return nil, ErrEmptyToken
	}
	token := tokens[0]

	h.service.GetCredential(ctx, token, in.CredID)

}

func (h CredsHandler) CredUpd(ctx context.Context, in *pb.CredUpdRequest) (*pb.CredUpdResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h CredsHandler) CredDel(ctx context.Context, in *pb.CredDelRequest) (*pb.CredDelResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h CredsHandler) CredList(ctx context.Context, in *pb.CredListRequest) (*pb.CredListResponse, error) {
	//TODO implement me
	panic("implement me")
}
