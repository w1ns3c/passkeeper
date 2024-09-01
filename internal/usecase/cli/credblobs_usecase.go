package cli

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
)

// GetCreds func is client logic for tui app
// to get credential blobs from server and decrypt them to Credential entities
func (c *ClientUC) GetCreds(ctx context.Context) (creds []*entities.Credential, err error) {

	resp, err := c.credsSvc.CredList(ctx, new(empty.Empty))
	if err != nil {
		return nil, err
	}

	creds = make([]*entities.Credential, len(resp.Creds))

	// client side, so there isn't interceptor that auto extract userID from token
	token, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := hashes.ExtractUserID(token)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(resp.Creds); i++ {
		blob := &entities.CredBlob{
			ID:     resp.Creds[i].ID,
			UserID: userID,
			Blob:   resp.Creds[i].Blob,
		}

		cred, err := hashes.DecryptBlob(blob, c.Secret)
		if err != nil {
			// TODO handle ERRORS!!!
		}

		creds[i] = cred
	}

	return creds, nil
}
