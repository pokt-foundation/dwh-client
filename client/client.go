package client

import (
	"context"
	"errors"
	"net/http"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pokt-foundation/portal-http-db/v2/types"
)

const dwhAPIKeyHeader = "Portal-DWH-Service-Api-Key"

var (
	ErrNoPortalAppIDS = errors.New("no portal app ids provided")
	ErrNoDateRange    = errors.New("both from and to dates must be provided")

	ErrJSON204           = errors.New("error 204: no content")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrNotFound          = errors.New("resource not found")
	ErrEmptybody         = errors.New("empty response body")
	ErrEmptyJSON200Data  = errors.New("empty json 200 data")
	ErrEmptyHTTPResponse = errors.New("empty http response")
	ErrJSONDefault       = errors.New("unknown error")
)

type (
	Config struct {
		APIKey, URL string
	}

	// DWHClient struct contains all the possible methods to interact with the Data Warehouse API.
	DWHClient struct {
		IDWHClient
		client *ClientWithResponses
	}

	IDWHClient interface {
		GetTotalRelaysForPortalAppIDs(ctx context.Context, params GetTotalRelaysForPortalAppIDsParams) ([]AnalyticsRelaysTotal, error)
	}

	GetTotalRelaysForPortalAppIDsParams struct {
		From, To     time.Time
		PortalAppIDs []types.PortalAppID
	}
)

func NewDWHClient(config Config) (*DWHClient, error) {
	dwhClient, err := NewClientWithResponses(config.URL,
		WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set(dwhAPIKeyHeader, config.APIKey)
			return nil
		}),
	)

	return &DWHClient{client: dwhClient}, err
}

// GetTotalRelaysForPortalAppIDs returns the total relays for the given portal app ids.
func (d *DWHClient) GetTotalRelaysForPortalAppIDs(ctx context.Context, params GetTotalRelaysForPortalAppIDsParams) ([]AnalyticsRelaysTotal, error) {
	if params.From.IsZero() || params.To.IsZero() {
		return nil, ErrNoDateRange
	}
	if len(params.PortalAppIDs) == 0 {
		return nil, ErrNoPortalAppIDS
	}

	portalAppIDStrs := []string{}
	for _, portalAppID := range params.PortalAppIDs {
		portalAppIDStrs = append(portalAppIDStrs, string(portalAppID))
	}

	response, err := d.client.GetAnalyticsRelaysTotalCategoryWithResponse(
		ctx,
		"application_id",
		&GetAnalyticsRelaysTotalCategoryParams{
			From:          openapi_types.Date{Time: params.From},
			To:            openapi_types.Date{Time: params.To},
			CategoryValue: portalAppIDStrs,
		},
	)
	if err != nil {
		return nil, err
	}

	// Handle non-200 error responses
	if response.JSON200 == nil || response.JSON200.Data == nil {
		switch {
		case response.JSON204 != nil:
			return nil, ErrJSON204
		case response.JSON401 != nil:
			return nil, ErrUnauthorized
		case response.JSON404 != nil:
			return nil, ErrNotFound
		case response.Body == nil:
			return nil, ErrEmptybody
		case response.HTTPResponse == nil:
			return nil, ErrEmptyHTTPResponse
		case response.JSON200.Data == nil:
			return nil, ErrEmptyJSON200Data
		case response.JSONDefault != nil:
			return nil, ErrJSONDefault
		}
	}

	responseData := []AnalyticsRelaysTotal{}
	for _, data := range *response.JSON200.Data {
		analyticsRelaysTotal, err := data.AsAnalyticsRelaysTotal()
		if err != nil {
			return nil, err
		}
		responseData = append(responseData, analyticsRelaysTotal)
	}

	return responseData, nil
}
