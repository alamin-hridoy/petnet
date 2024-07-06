package remittoaccount

import (
	"context"
	"encoding/json"
	"net/url"
	"path"
)

type iPerahubReqService interface {
	RtaPost(ctx context.Context, url string, body interface{}) (json.RawMessage, error)
}

// Client is client for Perahub Remit-to-Account API (1.0)
// perahub doc:https://perahub-rta-apidocs-bpsk2lcqiq-as.a.run.app/
type Client struct {
	baseUrl   *url.URL
	phService iPerahubReqService
}

// NewRTAClient ...
func NewRTAClient(phService iPerahubReqService, baseUrl *url.URL) *Client {
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
