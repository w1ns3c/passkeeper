package handlers

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
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
	ErrNoToken    = status.Error(codes.Unauthenticated, ErrNoTokenMsg)

	ErrEmptyTokenMsg = "token is empty"
	ErrEmptyToken    = status.Error(codes.Unauthenticated, ErrEmptyTokenMsg)

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
	service usecase.CredUsecaseInf
	log     *zerolog.Logger
}

func (h *CredsHandler) CredAdd(ctx context.Context, req *pb.CredAddRequest) (*empty.Empty, error) {
	token, err := ExtractUserToken(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	cred := &entities.Credential{
		Login:       req.Cred.Login,
		Password:    req.Cred.Password,
		Date:        time.Now(),
		Resource:    req.Cred.Resource,
		Description: req.Cred.Description,
	}

	err = h.service.AddCredential(ctx, token, cred)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredAddMsg)

		return nil, ErrCredAdd
	}

	return new(empty.Empty), nil
}

func (h *CredsHandler) CredGet(ctx context.Context, req *pb.CredGetRequest) (resp *pb.CredGetResponse, err error) {
	token, err := ExtractUserToken(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	cr, err := h.service.GetCredential(ctx, token, req.CredID)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredGetMsg)

		return nil, ErrCredGet
	}

	resp = &pb.CredGetResponse{
		Cred: &pb.Credential{
			Login:       cr.Login,
			Password:    cr.Password,
			Resource:    cr.Resource,
			Description: cr.Description,
			Date:        cr.Date.Format(time.DateTime),
		},
	}
	return resp, nil
}

func (h *CredsHandler) CredUpd(ctx context.Context, req *pb.CredUpdRequest) (*empty.Empty, error) {
	token, err := ExtractUserToken(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	t, err := time.Parse(time.DateTime, req.Cred.Date)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredUpdMsg)

		return nil, ErrCredUpd
	}

	cred := &entities.Credential{
		Login:       req.Cred.Login,
		Password:    req.Cred.Password,
		Date:        t,
		Resource:    req.Cred.Resource,
		Description: req.Cred.Description,
	}

	err = h.service.AddCredential(ctx, token, cred)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredUpdMsg)

		return nil, ErrCredUpd
	}

	return new(empty.Empty), nil
}

func (h *CredsHandler) CredDel(ctx context.Context, req *pb.CredDelRequest) (*empty.Empty, error) {
	token, err := ExtractUserToken(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	err = h.service.DeleteCredential(ctx, token, req.CredID)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredDelMsg)

		return nil, ErrCredDel
	}

	return new(empty.Empty), nil
}

func (h *CredsHandler) CredList(ctx context.Context, req *empty.Empty) (resp *pb.CredListResponse, err error) {
	token, err := ExtractUserToken(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	creds, err := h.service.ListCredentials(ctx, token)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrCredListMsg)

		return nil, ErrCredList
	}

	resp.Creds = make([]*pb.Credential, len(creds))
	for i := 0; i < len(creds); i++ {
		resp.Creds[i] = &pb.Credential{
			Login:       creds[i].Login,
			Password:    creds[i].Password,
			Resource:    creds[i].Resource,
			Description: creds[i].Description,
			Date:        creds[i].Date.Format(time.DateTime),
		}
	}

	return resp, nil
}

func ExtractUserToken(ctx context.Context) (token string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrNoToken
	}

	tokens := md.Get(config.TokenHeader)
	if len(tokens) == 0 {
		return "", ErrEmptyToken
	}
	token = tokens[0]

	return token, nil
}
