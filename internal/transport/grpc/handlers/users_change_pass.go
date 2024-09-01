package handlers

import (
	"context"
	"passkeeper/internal/entities/hashes"

	"passkeeper/internal/usecase/srv/usersUC"

	pb "passkeeper/internal/transport/grpc/protofiles/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
)

type UserChangePassHandler struct {
	pb.UnimplementedUserChangePassSvcServer
	service usersUC.UserUsecaseInf
	log     *zerolog.Logger
}

func NewUserChangePassHandler(logger *zerolog.Logger, service usersUC.UserUsecaseInf) *UserChangePassHandler {
	return &UserChangePassHandler{
		UnimplementedUserChangePassSvcServer: pb.UnimplementedUserChangePassSvcServer{},
		service:                              service,
		log:                                  logger,
	}
}

func (h *UserChangePassHandler) ChangePass(ctx context.Context, req *pb.UserChangePassReq) (*empty.Empty, error) {
	userID, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		h.log.Error().
			Err(err).Send()

		return nil, err
	}

	err = h.service.ChangePassword(ctx, userID,
		req.OldPass, req.NewPass, req.Repeat)

	return nil, err
}
