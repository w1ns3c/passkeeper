package blobsUC

import (
	"context"

	"passkeeper/internal/entities/structs"
)

type BlobsActionsUsecase interface {
	GetBlobs(ctx context.Context) (err error)
	EditBlob(ctx context.Context, cred structs.CredInf, ind int) (err error)
	AddBlob(ctx context.Context, cred structs.CredInf) (err error)
	DelBlob(ctx context.Context, ind int, blobType structs.BlobType) (err error)
}

type GetBlobsUsecase interface {
	GetCredByIND(credIND int) (cred *structs.Credential, err error)
	GetCardByIND(cardIND int) (cred *structs.Card, err error)
	GetNoteByIND(noteIND int) (cred *structs.Note, err error)
	GetFileByIND(ind int) (file *structs.File, err error)

	CredsLen() int
	CredsNotNil() bool

	GetCreds() []*structs.Credential
	GetCards() []*structs.Card
	GetNotes() []*structs.Note
	GetFiles() []*structs.File
}
