package cli

import (
	"context"
	"fmt"
	"sort"

	"github.com/golang/protobuf/ptypes/empty"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
)

// ListCreds func is client logic for tui app
// to get credential blobs from server and decrypt them to Credential entities
func (c *ClientUC) ListCreds(ctx context.Context) error {

	resp, err := c.credsSvc.CredList(ctx, new(empty.Empty))
	if err != nil {
		return err
	}

	creds := make([]*entities.Credential, len(resp.Creds))

	for i := 0; i < len(resp.Creds); i++ {
		blob := &entities.CredBlob{
			ID:     resp.Creds[i].ID,
			UserID: c.User.ID,
			Blob:   resp.Creds[i].Blob,
		}

		cred, err := hashes.DecryptBlob(blob, c.User.Secret)
		if err != nil {
			// TODO handle ERRORS!!!
		}

		creds[i] = cred
	}

	SortCredsByDate(creds)

	c.m.Lock()
	c.Creds = creds
	c.m.Unlock()

	return nil
}

func (c *ClientUC) EditCred(ctx context.Context, cred *entities.Credential, ind int) (err error) {

	//ID:          form.tuiApp.Creds[ind].ID,
	// update credID
	cred.ID = c.Creds[ind].ID

	blob, err := hashes.EncryptBlob(cred, c.User.Secret)
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
	if err != nil {
		return err
	}

	c.m.Lock()
	defer c.m.Unlock()
	// save creds in local app
	if err = entities.SaveCred(c.Creds, ind, cred); err != nil {
		return err
	}

	return err
}

func (c *ClientUC) AddCred(ctx context.Context, cred *entities.Credential) (err error) {

	blob, err := hashes.EncryptBlob(cred, c.User.Secret)
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
	if err != nil {
		// can't save creds on server
		return err
	}

	c.m.Lock()
	defer c.m.Unlock()
	tmpCreds, err := entities.AddCred(c.Creds, cred)
	if err != nil {
		// can't save creds localy
		return err
	}
	c.Creds = tmpCreds

	return nil
}

func (c *ClientUC) DelCred(ctx context.Context, ind int) (err error) {
	if ind < 0 && ind >= len(c.Creds) {
		return fmt.Errorf("invalid index")
	}

	credID := c.Creds[ind].ID

	req := &pb.CredDelRequest{CredID: credID}
	_, err = c.credsSvc.CredDel(ctx, req)
	if err != nil {

		return err
	}

	c.m.Lock()
	defer c.m.Unlock()
	newCreds, err := entities.Delete(c.Creds, ind)
	if err != nil {
		return err
	}
	c.Creds = newCreds

	return err
}

func (c *ClientUC) GetCreds() []*entities.Credential {
	c.m.Lock()

	tmpCreds := make([]*entities.Credential, len(c.Creds))
	copy(tmpCreds, c.Creds)

	c.m.Unlock()

	return tmpCreds
}

func (c *ClientUC) GetCredByIND(ind int) (cred *entities.Credential, err error) {
	if ind < 0 || ind >= len(c.Creds) {
		return nil, fmt.Errorf("invalid index")
	}

	c.m.Lock()
	defer c.m.Unlock()

	return c.Creds[ind], nil
}

// SortCredsByDate sort creds, now the first cred is the latest added
func SortCredsByDate(creds []*entities.Credential) {
	sort.Slice(creds, func(i, j int) bool {
		if creds[i].Date.After(creds[j].Date) {
			return true
		}
		return false
	})
}

func (c *ClientUC) GetCards() []*entities.Card {
	return c.Cards
}

func (c *ClientUC) GetNotes() []*entities.Note {
	return c.Notes
}
