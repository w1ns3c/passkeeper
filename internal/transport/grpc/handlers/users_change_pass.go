package handlers

import (
	"context"
	"github.com/w1ns3c/passkeeper/internal/usecase/srv"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
	pb "github.com/w1ns3c/passkeeper/internal/transport/grpc/protofiles/proto"
)

type UserPassHandler struct {
	pb.UnimplementedUserPassSvcServer
	service srv.UserUsecaseInf
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
