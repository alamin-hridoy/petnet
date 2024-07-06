package keto

import (
	"github.com/go-openapi/strfmt"
	"github.com/ory/keto-client-go/client"
)

type Svc struct {
	cl *client.OryKeto
}

func New(baseURL string) *Svc {
	cfg := client.DefaultTransportConfig()
	if baseURL != "" {
		cfg.WithHost(baseURL)
		cfg.WithSchemes([]string{"http"})
	}

	c := client.NewHTTPClientWithConfig(strfmt.Default, cfg)
	return &Svc{cl: c}
}
