package handlers

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
	pb "github.com/w1nsec/passkeeper/internal/transport/grpc/protofiles"
	"github.com/w1nsec/passkeeper/internal/usecase"
)

type UserPassHandler struct {
	pb.UnimplementedUserPassSvcServer
	service usecase.UserUsecaseInf
	log     *zerolog.Logger
}

func (h *UserPassHandler) ChangePass(ctx context.Context, req *pb.UserChangePassReq) (*empty.Empty, error) {
	userID, err := ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	err = h.service.ChangePassword(ctx, userID,
		req.OldPass, req.NewPass, req.Repeat)

	return nil, err
}
