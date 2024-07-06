package microinsurance

import (
	"context"
	"encoding/json"
	"net/url"
	"path"
)

//go:generate mockery --name=iPerahubService --structname=MockPerahubService --filename generated_mock_perahub_service_test.go --testonly --output . --outpkg microinsurance
type iPerahubService interface {
	PostMicroInsurance(ctx context.Context, subUrl string, body interface{}) (json.RawMessage, error)
	GetMicroInsurance(ctx context.Context, url string) (json.RawMessage, error)
}

// Client is client for micro insurance
// perahub doc:https://insurance-dev-apidocs-bpsk2lcqiq-as.a.run.app/
type Client struct {
	baseUrl   *url.URL
	phService iPerahubService
}

// NewMicroInsuranceClient ...
func NewMicroInsuranceClient(phService iPerahubService, baseUrl *url.URL) *Client {
	return &Client{
		baseUrl:   baseUrl,
		phService: phService,
	}
}

func (c *Client) getUrl(subUrl string) string {
	u := *c.baseUrl
	u.Path = path.Join(c.baseUrl.Path, subUrl)
	return u.String()
}
