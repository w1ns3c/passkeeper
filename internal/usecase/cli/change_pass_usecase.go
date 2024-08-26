package cli

import (
	"context"
	"passkeeper/internal/utils/hashes"

	pb "passkeeper/internal/transport/grpc/protofiles/proto"
)

func (c *ClientUC) ChangePass(ctx context.Context, curPass, newPass, repeat string) error {
	oldHash := hashes.Hash(curPass)
	newHash := hashes.Hash(newPass)
	repeatH := hashes.Hash(repeat)

	if newHash != repeatH {
		return ErrPassNotSame
	}

	req := &pb.UserChangePassReq{
		OldPass: oldHash,
		NewPass: newPass,
		Repeat:  repeatH,
	}

	_, err := c.passSvc.ChangePass(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
