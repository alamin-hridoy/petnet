package revenue_commission

import (
	"context"
	"encoding/json"
	"net/url"
	"path"
)

//go:generate mockery --name=iPerahubService --structname=MockPerahubService --filename generated_mock_perahub_service_test.go --testonly --output . --outpkg revenue_commission
type iPerahubService interface {
	PostRevComm(ctx context.Context, subUrl string, body interface{}) (json.RawMessage, error)
	PutRevComm(ctx context.Context, url string, body interface{}) (json.RawMessage, error)
	GetRevComm(ctx context.Context, url string) (json.RawMessage, error)
	DeleteRevComm(ctx context.Context, url string) (json.RawMessage, error)
}

// Client is client for perahub revenue commission
// perahub doc:https://drp-maintenance-docs-dev-bpsk2lcqiq-as.a.run.app/
type Client struct {
	baseUrl   *url.URL
	phService iPerahubService
}

// NewRevCommClient ...
func NewRevCommClient(phService iPerahubService, baseUrl *url.URL) *Client {
	return &Client{
		baseUrl:   baseUrl,
		phService: phService,
	}
}

// CommissionType is commission type, can be either absolute or range
type CommissionType string

// TrxType is transaction type, can be either inbound or outbound
type TrxType string

const (
	CommissionTypeAbsolute CommissionType = "absolute"
	CommissionTypeRange    CommissionType = "range"
	CommissionTypePercent  CommissionType = "percentage"

	TrxTypeInbound  TrxType = "inbound"
	TrxTypeOutbound TrxType = "outbound"
)

func (c *Client) getUrl(subUrl string) string {
	u := *c.baseUrl
	u.Path = path.Join(c.baseUrl.Path, subUrl)
	return u.String()
}
