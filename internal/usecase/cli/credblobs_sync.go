package cli

import (
	"context"
	"time"

	"google.golang.org/grpc/metadata"

	"passkeeper/internal/entities/config"
)

func (c *ClientUC) SyncBlobs(ctx context.Context) {

	ticker := time.NewTicker(c.SyncTime)
	for {
		select {
		case <-ticker.C:
			// skip if user not authed
			// or edit/add new cred
			c.m.RLock()
			if (c.User == nil && c.Creds == nil) || !c.Authed || c.viewPageFocus {
				c.m.RUnlock()
				continue
			}
			c.m.RUnlock()

			md := metadata.New(map[string]string{config.TokenHeader: c.GetToken()})
			newCtx := metadata.NewOutgoingContext(ctx, md)

			err := c.GetBlobs(newCtx)
			if err != nil {
				if c.log != nil {
					c.log.Error().
						Err(err).Msg("unable to sync credentials")
				}
			}

		case <-ctx.Done():
			return
		}
	}

}

func (c *ClientUC) StopSync() {
	c.m.Lock()
	c.viewPageFocus = true
	c.m.Unlock()
}

func (c *ClientUC) ContinueSync() {
	c.m.Lock()
	c.viewPageFocus = false
	c.m.Unlock()
}

func (c *ClientUC) CheckSync() bool {
	c.m.Lock()
	defer c.m.Unlock()

	return c.viewPageFocus
}
