package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pokt-foundation/portal-http-db/v2/types"
)

const dwhAPIKeyHeader = "Portal-DWH-Service-Api-Key"

var (
	ErrNoPortalAppIDs = errors.New("no portal app ids provided")
	ErrNoDateRange    = errors.New("no date range provided")
	ErrNoIDs          = errors.New("no ids provided for category type '%s")

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
		GetTotalRelaysForPortalAppIDs(ctx context.Context, params GetTotalRelaysForPortalAppIDsParams) ([]PortalAppRelaysTotal, error)
		GetTotalRelaysForAccountIDs(ctx context.Context, params GetTotalRelaysForAccountIDsParams) ([]AccountRelaysTotal, error)
	}

	GetTotalRelaysForPortalAppIDsParams struct {
		From, To     time.Time
		PortalAppIDs []types.PortalAppID
	}
	GetTotalRelaysForAccountIDsParams struct {
		From, To   time.Time
		AccountIDs []types.AccountID
	}

	PortalAppRelaysTotal struct {
		PortalAppID    types.PortalAppID
		Count          int
		AverageLatency float32
		RateError      float32
		RateSuccess    float32
		From, To       time.Time
	}

	AccountRelaysTotal struct {
		AccountID      types.AccountID
		Count          int
		AverageLatency float32
		RateError      float32
		RateSuccess    float32
		From, To       time.Time
	}
)

/* ---------- Factory Func ---------- */

func NewDWHClient(config Config) (*DWHClient, error) {
	dwhClient, err := NewClientWithResponses(config.URL,
		WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set(dwhAPIKeyHeader, config.APIKey)
			return nil
		}),
	)

	return &DWHClient{client: dwhClient}, err
}

/* ---------- Public Methods ---------- */

// GetTotalRelaysForPortalAppIDs returns the total relays for the given portal app IDs.
func (d *DWHClient) GetTotalRelaysForPortalAppIDs(ctx context.Context, params GetTotalRelaysForPortalAppIDsParams) ([]PortalAppRelaysTotal, error) {
	portalAppIDStrings := make([]string, len(params.PortalAppIDs))
	for i, portalAppID := range params.PortalAppIDs {
		portalAppIDStrings[i] = string(portalAppID)
	}

	responseData, err := d.getTotalRelays(ctx, params.From, params.To, portalAppIDStrings, GetAnalyticsRelaysTotalCategoryParamsCategoryApplicationId)
	if err != nil {
		return nil, err
	}

	portalAppRelaysTotals := make([]PortalAppRelaysTotal, len(responseData))
	for i, data := range responseData {
		portalAppRelaysTotals[i] = parsePortalAppIDAnalyticsRelaysTotal(data)
	}

	return portalAppRelaysTotals, nil
}

// GetTotalRelaysForAccountIDs returns the total relays for the given account IDs.
func (d *DWHClient) GetTotalRelaysForAccountIDs(ctx context.Context, params GetTotalRelaysForAccountIDsParams) ([]AccountRelaysTotal, error) {
	accountIDStrings := make([]string, len(params.AccountIDs))
	for i, accountID := range params.AccountIDs {
		accountIDStrings[i] = string(accountID)
	}

	responseData, err := d.getTotalRelays(ctx, params.From, params.To, accountIDStrings, GetAnalyticsRelaysTotalCategoryParamsCategoryAccountId)
	if err != nil {
		return nil, err
	}

	accountRelaysTotals := make([]AccountRelaysTotal, len(responseData))
	for i, data := range responseData {
		accountRelaysTotals[i] = parseAccountIDAnalyticsRelaysTotal(data)
	}

	return accountRelaysTotals, nil
}

/* ---------- Private Methods/funcs ---------- */

// getTotalRelays fetches total relays from the data warehouse client.
func (d *DWHClient) getTotalRelays(ctx context.Context, from, to time.Time, ids []string, categoryType GetAnalyticsRelaysTotalCategoryParamsCategory) ([]AnalyticsRelaysTotal, error) {
	if from.IsZero() || to.IsZero() {
		return nil, ErrNoDateRange
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf(ErrNoIDs.Error(), categoryType)
	}

	response, err := d.client.GetAnalyticsRelaysTotalCategoryWithResponse(
		ctx,
		categoryType,
		&GetAnalyticsRelaysTotalCategoryParams{
			From:          openapi_types.Date{Time: from},
			To:            openapi_types.Date{Time: to},
			CategoryValue: ids,
		},
	)
	if err != nil {
		return nil, err
	}

	return d.parseResponse(response)
}

// parseResponse parses the response from the data warehouse client.
func (d *DWHClient) parseResponse(response *GetAnalyticsRelaysTotalCategoryResponse) ([]AnalyticsRelaysTotal, error) {
	// Handle non-200 error responses
	if response.JSON200 == nil || response.JSON200.Data == nil {
		return nil, d.handleNon200Response(response)
	}

	var responseData []AnalyticsRelaysTotal
	for _, data := range *response.JSON200.Data {
		analyticsRelaysTotal, err := data.AsAnalyticsRelaysTotal()
		if err != nil {
			return nil, err
		}
		responseData = append(responseData, analyticsRelaysTotal)
	}
	return responseData, nil
}

// handleNon200Response handles non-200 error responses.
func (d *DWHClient) handleNon200Response(response *GetAnalyticsRelaysTotalCategoryResponse) error {
	switch {
	case response.JSON204 != nil:
		return ErrJSON204
	case response.JSON401 != nil:
		return ErrUnauthorized
	case response.JSON404 != nil:
		return ErrNotFound
	case response.Body == nil:
		return ErrEmptybody
	case response.HTTPResponse == nil:
		return ErrEmptyHTTPResponse
	case response.JSON200.Data == nil:
		return ErrEmptyJSON200Data
	case response.JSONDefault != nil:
		return ErrJSONDefault
	}
	return nil
}

func parsePortalAppIDAnalyticsRelaysTotal(analyticsRelaysTotal AnalyticsRelaysTotal) PortalAppRelaysTotal {
	parsedData := PortalAppRelaysTotal{
		PortalAppID:    types.PortalAppID(derefString(analyticsRelaysTotal.CategoryValue)),
		Count:          derefInt(analyticsRelaysTotal.CountTotal),
		AverageLatency: derefFloat32(analyticsRelaysTotal.AvgLatency),
		RateError:      derefFloat32(analyticsRelaysTotal.RateError),
		RateSuccess:    derefFloat32(analyticsRelaysTotal.RateSuccess),
		From:           convertDate(analyticsRelaysTotal.From),
		To:             convertDate(analyticsRelaysTotal.To),
	}
	return parsedData
}

func parseAccountIDAnalyticsRelaysTotal(analyticsRelaysTotal AnalyticsRelaysTotal) AccountRelaysTotal {
	parsedData := AccountRelaysTotal{
		AccountID:      types.AccountID(derefString(analyticsRelaysTotal.CategoryValue)),
		Count:          derefInt(analyticsRelaysTotal.CountTotal),
		AverageLatency: derefFloat32(analyticsRelaysTotal.AvgLatency),
		RateError:      derefFloat32(analyticsRelaysTotal.RateError),
		RateSuccess:    derefFloat32(analyticsRelaysTotal.RateSuccess),
		From:           convertDate(analyticsRelaysTotal.From),
		To:             convertDate(analyticsRelaysTotal.To),
	}
	return parsedData
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

func derefFloat32(f *float32) float32 {
	if f != nil {
		return *f
	}
	return 0.0
}

func convertDate(d *openapi_types.Date) time.Time {
	if d != nil {
		return d.Time
	}
	return time.Time{}
}
