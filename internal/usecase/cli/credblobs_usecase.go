package cli

import (
	"context"
	"fmt"
	"sort"

	"github.com/golang/protobuf/ptypes/empty"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/structs"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
)

// GetBlobs func is client logic for tui app
// to get credential blobs from server and decrypt them to Credential/Card/Note entities
func (c *ClientUC) GetBlobs(ctx context.Context) error {

	resp, err := c.credsSvc.BlobList(ctx, new(empty.Empty))
	if err != nil {
		return err
	}

	creds := make([]*structs.Credential, 0)
	cards := make([]*structs.Card, 0)
	notes := make([]*structs.Note, 0)
	files := make([]*structs.File, 0)

	for i := 0; i < len(resp.Blobs); i++ {
		blob := &structs.CryptoBlob{
			ID:     resp.Blobs[i].ID,
			UserID: c.User.ID,
			Blob:   resp.Blobs[i].Blob,
		}

		tmp, err := hashes.DecryptBlob(blob, c.User.Secret)
		if err != nil {
			c.log.Error().
				Err(err).Msg("can't decrypt cipher blob")

			continue
		}

		switch tmp.(type) {
		case *structs.Card:
			cards = append(cards, tmp.(*structs.Card))
		case *structs.Note:
			notes = append(notes, tmp.(*structs.Note))
		case *structs.Credential:
			creds = append(creds, tmp.(*structs.Credential))
		case *structs.File:
			files = append(files, tmp.(*structs.File))
		default:
			c.log.Error().
				Err(fmt.Errorf("unknown type of blob")).Send()
		}

	}

	SortCredsByDate(creds)
	SortNotesByDate(notes)

	c.m.Lock()
	defer c.m.Unlock()

	c.Creds = creds
	c.Cards = cards
	c.Notes = notes
	c.Files = files

	c.log.Info().Msgf("sync sum blobs:  %d", len(resp.Blobs))
	c.log.Info().Msgf("decrypted blobs: %d", len(c.Creds)+len(c.Cards)+len(c.Notes))
	c.log.Info().
		Int("creds", len(c.Creds)).
		Int("cards", len(c.Cards)).
		Int("notes", len(c.Notes)).
		Int("files", len(c.Files)).Msgf("blobs by type:")

	return nil
}

// EditBlob change user blob, encrypt it and save changes in storage on server side
func (c *ClientUC) EditBlob(ctx context.Context, cred structs.CredInf, ind int) (err error) {
	// check that blob with ind exist
	switch cred.(type) {
	case *structs.Credential:
		if ind < 0 && ind >= len(c.Creds) {
			return fmt.Errorf("invalid index of creds")
		}
	case *structs.Card:
		if ind < 0 && ind >= len(c.Cards) {
			return fmt.Errorf("invalid index of cards")
		}
	case *structs.Note:
		if ind < 0 && ind >= len(c.Notes) {
			return fmt.Errorf("invalid index of notes")
		}
	case *structs.File:
		if ind < 0 && ind >= len(c.Files) {
			return fmt.Errorf("invalid index of files")
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

	var blobT string

	// save creds in local app
	switch cred.(type) {
	case *structs.Credential:
		if err = entities.SaveCred(c.Creds, ind, cred.(*structs.Credential)); err != nil {
			return err
		}
		blobT = "credential"

	case *structs.Card:
		if err = entities.SaveCard(c.Cards, ind, cred.(*structs.Card)); err != nil {
			return err
		}
		blobT = "card"

	case *structs.Note:
		if err = entities.SaveNote(c.Notes, ind, cred.(*structs.Note)); err != nil {
			return err
		}
		blobT = "note"

	case *structs.File:
		if err = entities.SaveFile(c.Files, ind, cred.(*structs.File)); err != nil {
			return err
		}
		blobT = "file"

	default:
		return fmt.Errorf("unknown credential type")
	}

	c.log.Info().
		Str("id", cred.GetID()).Msgf("blob (%s) updated", blobT)

	return err
}

// AddBlob encrypt new blob and save in storage on server side
func (c *ClientUC) AddBlob(ctx context.Context, cred structs.CredInf) (err error) {

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

	var blobT string
	switch cred.(type) {
	case *structs.Credential:
		tmpCreds, err := entities.AddCred(c.Creds, cred.(*structs.Credential))
		if err != nil {
			// can't save creds localy
			c.log.Error().
				Err(err).Msg("can't save new cred locally")
			return err
		}
		c.Creds = tmpCreds
		blobT = "credential"

	case *structs.Card:
		tmpCards, err := entities.AddCard(c.Cards, cred.(*structs.Card))
		if err != nil {
			// can't save creds localy
			c.log.Error().
				Err(err).Msg("can't save new cred locally")
			return err
		}
		c.Cards = tmpCards
		blobT = "card"

	case *structs.Note:
		tmpNotes, err := entities.AddNote(c.Notes, cred.(*structs.Note))
		if err != nil {
			// can't save creds localy
			c.log.Error().
				Err(err).Msg("can't save new cred locally")
			return err
		}
		c.Notes = tmpNotes
		blobT = "note"

	case *structs.File:
		tmpFiles, err := entities.AddFile(c.Files, cred.(*structs.File))
		if err != nil {
			// can't save creds localy
			c.log.Error().
				Err(err).Msg("can't save new cred locally")
			return err
		}
		c.Files = tmpFiles
		blobT = "file"

	default:
		return fmt.Errorf("unknown credential type")
	}

	c.log.Info().
		Str("id", cred.GetID()).Msgf("blob (%s) added", blobT)

	return nil
}

// DelBlob search blobID by ind and blobType on client side,
// then delete crypto blob from storage by blobID on server side
func (c *ClientUC) DelBlob(ctx context.Context, ind int, blobType structs.BlobType) (err error) {
	var delID string

	// check that blob with ind exist
	switch blobType {
	case structs.BlobCred:
		tmp, err := c.GetCredByIND(ind)
		if err != nil {
			return fmt.Errorf("invalid index of creds")
		}
		delID = tmp.ID

	case structs.BlobCard:
		tmp, err := c.GetCardByIND(ind)
		if err != nil {
			return fmt.Errorf("invalid index of creds")
		}
		delID = tmp.ID

	case structs.BlobNote:
		tmp, err := c.GetNoteByIND(ind)
		if err != nil {
			return fmt.Errorf("invalid index of notes")
		}
		delID = tmp.ID

	case structs.BlobFile:
		tmp, err := c.GetFileByIND(ind)
		if err != nil {
			return fmt.Errorf("invalid index of files")
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

	var blobT string

	// update blobs values
	switch blobType {
	case structs.BlobCred:
		newCreds, err := entities.DeleteCred(c.Creds, ind)
		if err != nil {
			return err
		}
		c.Creds = newCreds
		blobT = "credential"

	case structs.BlobCard:
		newCards, err := entities.DeleteCard(c.Cards, ind)
		if err != nil {
			return err
		}
		c.Cards = newCards
		blobT = "card"

	case structs.BlobNote:
		newNotes, err := entities.DeleteNote(c.Notes, ind)
		if err != nil {
			return err
		}
		c.Notes = newNotes
		blobT = "note"

	case structs.BlobFile:
		newFiles, err := entities.DeleteFile(c.Files, ind)
		if err != nil {
			return err
		}
		c.Files = newFiles
		blobT = "note"

	default:
		return fmt.Errorf("unknown type of blob")
	}

	c.log.Info().
		Str("id", delID).Msgf("blob (%s) deleted", blobT)

	return err
}

// GetCredByIND return cred by it's ind in slice (safety)
func (c *ClientUC) GetCredByIND(ind int) (cred *structs.Credential, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	if ind < 0 || ind >= len(c.Creds) {
		return nil, fmt.Errorf("invalid index")
	}

	return c.Creds[ind], nil
}

// GetCardByIND return cred by it's ind in slice (safety)
func (c *ClientUC) GetCardByIND(ind int) (card *structs.Card, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	if ind < 0 || ind >= len(c.Cards) {
		return nil, fmt.Errorf("invalid index")
	}

	return c.Cards[ind], nil
}

// GetNoteByIND return note by it's ind in slice (safety)
func (c *ClientUC) GetNoteByIND(ind int) (note *structs.Note, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	if ind < 0 || ind >= len(c.Notes) {
		return nil, fmt.Errorf("invalid index")
	}

	return c.Notes[ind], nil
}

// GetFileByIND return cred by it's ind in slice (safety)
func (c *ClientUC) GetFileByIND(ind int) (file *structs.File, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	if ind < 0 || ind >= len(c.Files) {
		return nil, fmt.Errorf("invalid index")
	}

	return c.Files[ind], nil
}

// SortCredsByDate sort creds, now the first cred is the latest added
func SortCredsByDate(creds []*structs.Credential) {
	sort.Slice(creds, func(i, j int) bool {
		if creds[i].Date.After(creds[j].Date) {
			return true
		}
		return false
	})
}

// SortNotesByDate sort cards, now the first cred is the latest added
func SortNotesByDate(notes []*structs.Note) {
	sort.Slice(notes, func(i, j int) bool {
		if notes[i].Date.After(notes[j].Date) {
			return true
		}
		return false
	})
}

// GetCreds return a copy of creds to view
func (c *ClientUC) GetCreds() []*structs.Credential {
	c.m.Lock()

	tmpCreds := make([]*structs.Credential, len(c.Creds))
	copy(tmpCreds, c.Creds)

	c.m.Unlock()

	return tmpCreds
}

// GetCards return a copy of cards to view
func (c *ClientUC) GetCards() []*structs.Card {
	return c.Cards
}

// GetNotes return a copy of notes to view
func (c *ClientUC) GetNotes() []*structs.Note {
	return c.Notes
}

// GetFiles return a copy of files to view
func (c *ClientUC) GetFiles() []*structs.File {
	return c.Files
}
