package syncUC

import (
	"context"
	"time"
)

type SyncUsecase interface {
	SyncBlobs(ctx context.Context)
	StopSync()
	ContinueSync()
	CheckSync() bool
	GetSyncTime() time.Duration
}
