package handlers

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"

	"passkeeper/internal/entities/hashes"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/srv/usersUC"
)

// UserChangePassHandler handle user's request to change password
//
//go:generate mockgen -source=../protofiles/proto/users_change_pwd_grpc.pb.go -destination=../../../../mocks/gservice/user_change_pass.go -package=mocks
type UserChangePassHandler struct {
	pb.UnimplementedUserChangePassSvcServer
	service usersUC.UserUsecaseInf
	log     *zerolog.Logger
}

// NewUserChangePassHandler is a constructor for UserChangePassHandler
func NewUserChangePassHandler(logger *zerolog.Logger, service usersUC.UserUsecaseInf) *UserChangePassHandler {
	return &UserChangePassHandler{
		UnimplementedUserChangePassSvcServer: pb.UnimplementedUserChangePassSvcServer{},
		service:                              service,
		log:                                  logger,
	}
}

// ChangePass handle user request to change password
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
