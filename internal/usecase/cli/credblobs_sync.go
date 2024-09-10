package cli

import (
	"context"
	"google.golang.org/grpc/metadata"
	"passkeeper/internal/config"
	"time"
)

func (c *ClientUC) SyncCreds(ctx context.Context) {

	ticker := time.NewTicker(c.SyncTime)
	for range ticker.C {
		// skip if user not authed
		// or edit/add new cred
		c.m.RLock()
		if (c.User == nil && c.Creds == nil) || !c.Authed {
			c.m.RUnlock()
			continue
		}
		c.m.RUnlock()

		md := metadata.New(map[string]string{config.TokenHeader: c.GetToken()})
		ctx = metadata.NewOutgoingContext(ctx, md)

		err := c.ListCreds(ctx)
		if err != nil {
			if c.log != nil {
				c.log.Error().
					Err(err).Msg("unable to sync credentials")
			}
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
