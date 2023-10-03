package client

import (
	"context"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/pokt-foundation/portal-db/v2/types"
	"github.com/pokt-foundation/utils-go/environment"
	"github.com/stretchr/testify/assert"
)

// Load all env vars needed by test
var _ = godotenv.Load("../.env")

func Test_GetTotalRelaysForPortalAppIDs(t *testing.T) {
	tests := []struct {
		name          string
		config        Config
		from          time.Time
		to            time.Time
		portalAppIDs  []types.PortalAppID
		expected      []AnalyticsRelaysTotal
		expectedError error
	}{
		{
			name: "Should return total relays for given PortalAppIDs within the date range",
			config: Config{
				APIKey: environment.MustGetString("DWH_API_KEY"),
				URL:    environment.MustGetString("DWH_API_URL"),
			},
			from: time.Date(2023, 8, 4, 0, 0, 0, 0, time.UTC),
			to:   time.Date(2023, 9, 4, 0, 0, 0, 0, time.UTC),
			portalAppIDs: []types.PortalAppID{
				"44c0823fbdf0aed3fa2d6357",
				"3742b06f9e13c9ea22a8d599",
			},
			expected:      []AnalyticsRelaysTotal{},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dwhClient, err := NewDWHClient(test.config)
			if err != nil {
				t.Fatal(err)
			}

			actual, err := dwhClient.GetTotalRelaysForPortalAppIDs(context.Background(), GetTotalRelaysForPortalAppIDsParams{
				From:         test.from,
				To:           test.to,
				PortalAppIDs: test.portalAppIDs,
			})
			assert.Equal(t, test.expectedError, err)
			assert.Len(t, actual, len(test.expected))
			assert.NotEmpty(t, test.expected, actual)
		})
	}
}
