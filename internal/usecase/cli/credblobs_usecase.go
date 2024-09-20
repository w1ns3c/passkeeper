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

// GetBlobs func is client logic for tui app
// to get credential blobs from server and decrypt them to Credential/Card/Note entities
func (c *ClientUC) GetBlobs(ctx context.Context) error {

	resp, err := c.credsSvc.BlobList(ctx, new(empty.Empty))
	if err != nil {
		return err
	}

	creds := make([]*entities.Credential, 0)
	cards := make([]*entities.Card, 0)
	notes := make([]*entities.Note, 0)

	for i := 0; i < len(resp.Blobs); i++ {
		blob := &entities.CryptoBlob{
			ID:     resp.Blobs[i].ID,
			UserID: c.User.ID,
			Blob:   resp.Blobs[i].Blob,
		}

		tmpCred, err := hashes.DecryptBlob(blob, c.User.Secret)
		if err != nil {
			// TODO handle ERRORS!!!
			c.log.Error().
				Err(err).Msg("can't decrypt cipher blob")

			continue
		}

		switch tmpCred.(type) {
		case *entities.Card:
			cards = append(cards, tmpCred.(*entities.Card))
		case *entities.Note:
			notes = append(notes, tmpCred.(*entities.Note))
		case *entities.Credential:
			creds = append(creds, tmpCred.(*entities.Credential))
		default:
			c.log.Error().
				Err(fmt.Errorf("unknown type of blob")).Send()
		}

	}

	SortCredsByDate(creds)
	SortNotesByDate(notes)

	c.m.Lock()
	c.Creds = creds
	c.Cards = cards
	c.Notes = notes
	c.m.Unlock()

	return nil
}

func (c *ClientUC) EditBlob(ctx context.Context, cred entities.CredInf, ind int) (err error) {
	// check that blob with ind exist
	switch cred.(type) {
	case *entities.Credential:
		if ind < 0 && ind >= len(c.Creds) {
			return fmt.Errorf("invalid index of creds")
		}
	case *entities.Card:
		if ind < 0 && ind >= len(c.Cards) {
			return fmt.Errorf("invalid index of cards")
		}
	case *entities.Note:
		if ind < 0 && ind >= len(c.Notes) {
			return fmt.Errorf("invalid index of notes")
		}
	default:
		return fmt.Errorf("unknown type of blob")
	}

	blob, err := hashes.EncryptBlob(cred, c.User.Secret)
	if err != nil {
		return err
	}

	if blob == nil {
		return fmt.Errorf("something go wrong, blob is nil")
	}

	req := &pb.BlobUpdRequest{
		Blob: &pb.CryptoBlob{
			ID:   blob.ID,
			Blob: blob.Blob,
		},
	}

	_, err = c.credsSvc.BlobUpd(ctx, req)
	if err != nil {
		return err
	}

	c.m.Lock()
	defer c.m.Unlock()
	// save creds in local app
	switch cred.(type) {
	case *entities.Credential:
		if err = entities.SaveCred(c.Creds, ind, cred.(*entities.Credential)); err != nil {
			return err
		}

	case *entities.Card:
		if err = entities.SaveCard(c.Cards, ind, cred.(*entities.Card)); err != nil {
			return err
		}

	case *entities.Note:
		if err = entities.SaveNote(c.Notes, ind, cred.(*entities.Note)); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown credential type")
	}

	return err
}

func (c *ClientUC) AddBlob(ctx context.Context, cred entities.CredInf) (err error) {

	blob, err := hashes.EncryptBlob(cred, c.User.Secret)
	if err != nil {
		return err
	}

	if blob == nil {
		return fmt.Errorf("something go wrong, blob is nil")
	}

	req := &pb.BlobAddRequest{
		Cred: &pb.CryptoBlob{
			ID:   blob.ID,
			Blob: blob.Blob,
		},
	}

	_, err = c.credsSvc.BlobAdd(ctx, req)
	if err != nil {
		// can't save creds on server
		// can't save creds localy
		c.log.Error().
			Err(err).Msg("can't save new cred remotely")
		return err
	}

	c.m.Lock()
	defer c.m.Unlock()

	switch cred.(type) {
	case *entities.Credential:
		tmpCreds, err := entities.AddCred(c.Creds, cred.(*entities.Credential))
		if err != nil {
			// can't save creds localy
			c.log.Error().
				Err(err).Msg("can't save new cred locally")
			return err
		}
		c.Creds = tmpCreds

	case *entities.Card:
		tmpCards, err := entities.AddCard(c.Cards, cred.(*entities.Card))
		if err != nil {
			// can't save creds localy
			c.log.Error().
				Err(err).Msg("can't save new cred locally")
			return err
		}
		c.Cards = tmpCards

	case *entities.Note:
		tmpNotes, err := entities.AddNote(c.Notes, cred.(*entities.Note))
		if err != nil {
			// can't save creds localy
			c.log.Error().
				Err(err).Msg("can't save new cred locally")
			return err
		}
		c.Notes = tmpNotes

	default:
		return fmt.Errorf("unknown credential type")
	}

	return nil
}

func (c *ClientUC) DelBlob(ctx context.Context, ind int, blobType entities.BlobType) (err error) {
	var delID string

	// check that blob with ind exist
	switch blobType {
	case entities.UserCred:
		tmp, err := c.GetCredByIND(ind)
		if err != nil {
			return fmt.Errorf("invalid index of creds")
		}
		delID = tmp.ID

	case entities.UserCard:
		tmp, err := c.GetCardByIND(ind)
		if err != nil {
			return fmt.Errorf("invalid index of creds")
		}
		delID = tmp.ID

	case entities.UserNote:
		tmp, err := c.GetNoteByIND(ind)
		if err != nil {
			return fmt.Errorf("invalid index of notes")
		}
		delID = tmp.ID

	default:
		return fmt.Errorf("unknown type of blob")
	}

	req := &pb.BlobDelRequest{CredID: delID}
	_, err = c.credsSvc.BlobDel(ctx, req)
	if err != nil {

		return err
	}

	c.m.Lock()
	defer c.m.Unlock()

	// update blobs values
	switch blobType {
	case entities.UserCred:
		newCreds, err := entities.DeleteCred(c.Creds, ind)
		if err != nil {
			return err
		}
		c.Creds = newCreds

	case entities.UserCard:
		newCards, err := entities.DeleteCard(c.Cards, ind)
		if err != nil {
			return err
		}
		c.Cards = newCards

	case entities.UserNote:
		newNotes, err := entities.DeleteNote(c.Notes, ind)
		if err != nil {
			return err
		}
		c.Notes = newNotes

	default:
		return fmt.Errorf("unknown type of blob")
	}

	return err
}

func (c *ClientUC) GetCredByIND(ind int) (cred *entities.Credential, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	if ind < 0 || ind >= len(c.Creds) {
		return nil, fmt.Errorf("invalid index")
	}

	return c.Creds[ind], nil
}

func (c *ClientUC) GetCardByIND(ind int) (card *entities.Card, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	if ind < 0 || ind >= len(c.Cards) {
		return nil, fmt.Errorf("invalid index")
	}

	return c.Cards[ind], nil
}

func (c *ClientUC) GetNoteByIND(ind int) (note *entities.Note, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	if ind < 0 || ind >= len(c.Notes) {
		return nil, fmt.Errorf("invalid index")
	}

	return c.Notes[ind], nil
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

// SortNotesByDate sort cards, now the first cred is the latest added
func SortNotesByDate(notes []*entities.Note) {
	sort.Slice(notes, func(i, j int) bool {
		if notes[i].Date.After(notes[j].Date) {
			return true
		}
		return false
	})
}

func (c *ClientUC) GetCreds() []*entities.Credential {
	c.m.Lock()

	tmpCreds := make([]*entities.Credential, len(c.Creds))
	copy(tmpCreds, c.Creds)

	c.m.Unlock()

	return tmpCreds
}

func (c *ClientUC) GetCards() []*entities.Card {
	return c.Cards
}

func (c *ClientUC) GetNotes() []*entities.Note {
	return c.Notes
}
