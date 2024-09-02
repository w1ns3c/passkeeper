package cli

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
)

// GetCreds func is client logic for tui app
// to get credential blobs from server and decrypt them to Credential entities
func (c *ClientUC) GetCreds(ctx context.Context) (creds []*entities.Credential, err error) {

	resp, err := c.credsSvc.CredList(ctx, new(empty.Empty))
	if err != nil {
		return nil, err
	}

	creds = make([]*entities.Credential, len(resp.Creds))

	for i := 0; i < len(resp.Creds); i++ {
		blob := &entities.CredBlob{
			ID:     resp.Creds[i].ID,
			UserID: c.UserID,
			Blob:   resp.Creds[i].Blob,
		}

		cred, err := hashes.DecryptBlob(blob, c.CredsSecret)
		if err != nil {
			// TODO handle ERRORS!!!
		}

		creds[i] = cred
	}

	return creds, nil
}

func (c *ClientUC) EditCred(ctx context.Context, cred *entities.Credential) (err error) {

	blob, err := hashes.EncryptBlob(cred, c.CredsSecret)
	if err != nil {
		return err
	}

	if blob == nil {
		return fmt.Errorf("something go wrong, blob is nil")
	}

	req := &pb.CredUpdRequest{
		Cred: &pb.CredBlob{
			ID:   blob.ID,
			Blob: blob.Blob,
		},
	}

	_, err = c.credsSvc.CredUpd(ctx, req)
	return err
}

func (c *ClientUC) AddCred(ctx context.Context, cred *entities.Credential) (err error) {

	blob, err := hashes.EncryptBlob(cred, c.CredsSecret)
	if err != nil {
		return err
	}

	if blob == nil {
		return fmt.Errorf("something go wrong, blob is nil")
	}

	req := &pb.CredAddRequest{
		Cred: &pb.CredBlob{
			ID:   blob.ID,
			Blob: blob.Blob,
		},
	}

	_, err = c.credsSvc.CredAdd(ctx, req)
	return err
}

func (c *ClientUC) DelCred(ctx context.Context, credID string) (err error) {

	req := &pb.CredDelRequest{CredID: credID}
	_, err = c.credsSvc.CredDel(ctx, req)

	return err
}
