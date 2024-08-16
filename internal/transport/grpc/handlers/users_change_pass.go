package handlers

import (
	"context"

	"passkeeper/internal/usecase/srv/usersUC"

	pb "passkeeper/internal/transport/grpc/protofiles/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
)

type UserPassHandler struct {
	pb.UnimplementedUserPassSvcServer
	service usersUC.UserUsecaseInf
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
