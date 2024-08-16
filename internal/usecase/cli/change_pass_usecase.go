package cli

import (
	"context"

	pb "passkeeper/internal/transport/grpc/protofiles/proto"
)

func (c *ClientUC) ChangePass(ctx context.Context, curPass, newPass, repeat string) error {
	oldHash := Hash(curPass)
	newHash := Hash(newPass)
	repeatH := Hash(repeat)

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
