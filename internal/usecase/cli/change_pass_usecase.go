package cli

import (
	"context"

	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
)

// ChangePass provide user changing password functionality
func (c *ClientUC) ChangePass(ctx context.Context, curPass, newPass, repeat string) error {
	if newPass != repeat {
		return myerrors.ErrPassNotSame
	}

	oldHash := hashes.Hash(curPass)
	newHash := hashes.Hash(newPass)
	repeatH := hashes.Hash(repeat)

	req := &pb.UserChangePassReq{
		OldPass: oldHash,
		NewPass: newHash,
		Repeat:  repeatH,
	}

	_, err := c.passSvc.ChangePass(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
