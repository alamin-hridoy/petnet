package main

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	hydra "github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/spf13/viper"

	"brank.as/rbac/serviceutil/leaderelex/leader/client"
)

type cleanup struct {
	ival time.Duration
	keep time.Duration
	le   *client.Leader
	hy   admin.ClientService
}

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (rt RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return rt(req) }

func FakeTLSTerm(rt http.RoundTripper) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		rq := req.Clone(req.Context())
		rq.Header.Set("X-Forwarded-Proto", "https")
		return rt.RoundTrip(rq)
	})
}

// Create hydra cleanup
func HydraCleanup(conf *viper.Viper) (*cleanup, error) {
	parsedURL, err := url.Parse(conf.GetString("hydra.adminurl"))
	if err != nil {
		return nil, err
	}
	if parsedURL.String() == "" {
		return nil, errors.New("empty hydra url given")
	}

	ht := httptransport.New(parsedURL.Host, parsedURL.Path, []string{parsedURL.Scheme})
	ht.Transport = FakeTLSTerm(http.DefaultTransport)

	// keep duration for session tokens to protect.
	keep := conf.GetDuration("hydra.cleanup_keep")
	switch {
	case keep == 0:
		keep = -24 * time.Hour
	case keep > 0:
		keep = -keep
	}
	// worker interval
	ival := conf.GetDuration("hydra.cleanup_interval")
	switch {
	case ival == 0:
		ival = time.Hour
	case ival < 0:
		ival = -ival
	}
	return &cleanup{ival: ival, keep: keep, hy: hydra.New(ht, strfmt.Default).Admin}, nil
}

// Clean inactive tokens.
func (c *cleanup) Clean(ctx context.Context) error {
	f, err := c.hy.FlushInactiveOAuth2Tokens(
		admin.NewFlushInactiveOAuth2TokensParamsWithContext(ctx).
			WithBody(&models.FlushInactiveOAuth2TokensRequest{
				NotAfter: strfmt.DateTime(time.Now().Add(c.keep)),
			}),
	)
	if err != nil {
		return err
	}
	return f
}
