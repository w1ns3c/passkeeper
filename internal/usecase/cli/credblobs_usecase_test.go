package cli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"passkeeper/internal/entities/structs"
)

func TestSortCredsByDate(t *testing.T) {
	var (
		t1     = time.Now()
		creds1 = []*structs.Credential{
			{
				ID:          "1111",
				Resource:    "res_1111",
				Password:    "pass_1111",
				Description: "desc_1111",
				Date:        t1,
			},
			{
				ID:          "2222",
				Resource:    "res_2222",
				Password:    "pass_2222",
				Description: "desc_2222",
				Date:        t1.Add(time.Hour * 4), // t1 + 4 hour
			},
			{
				ID:          "3333",
				Resource:    "res_3333",
				Password:    "pass_3333",
				Description: "desc_3333",
				Date:        t1.Add(time.Hour * -2), // t1 - 2 hour
			},
		}
		order1 = []string{"2222", "1111", "3333"}
	)
	tests := []struct {
		name    string
		creds   []*structs.Credential
		orderID []string
	}{
		{
			name:    "Check 1",
			creds:   creds1,
			orderID: order1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.orderID) != len(tt.creds) {
				t.Errorf("wrong test, len of order: %d and len of creds: %d, should be the same",
					len(tt.orderID), len(tt.creds))
			}
			SortCredsByDate(tt.creds)

			for ind, id := range tt.orderID {
				require.Equal(t, tt.creds[ind].ID, id, "id not the same")
			}

		})
	}
}
