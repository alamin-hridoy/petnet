package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	// hydra
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"

	// github
	"github.com/pkg/errors"
)

type AdminClient struct {
	host        string
	tr          *transport
	mux         *sync.Mutex
	adminClient admin.ClientService
}

// Creates an hydra admin client for interacting with OAuth2 clients
func NewAdminClient(host string) (*AdminClient, error) {
	s := &AdminClient{
		tr: &transport{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          10,
				MaxIdleConnsPerHost:   1,
				MaxConnsPerHost:       5,
				IdleConnTimeout:       5 * time.Second,
				TLSHandshakeTimeout:   2 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
		mux:  &sync.Mutex{},
		host: host,
	}

	adminClient, err := s.newAdminClient()
	if err != nil {
		return nil, Error{
			Message: `hydra admin client initialization failed`,
			Err:     err,
			Code:    HydraError,
		}
	}
	s.adminClient = adminClient

	return s, nil
}

const (
	MethodPublic  = "none"
	MethodPrivate = "client_secret_basic"
)

const AuthBackend = "authenticator"

type AuthClient struct {
	OwnerID                           string
	ClientID                          string
	ClientName                        string
	RedirectURIs                      []string
	CORS                              []string
	PostLogoutRedirectURIs            []string
	FrontChannelLogoutURI             string
	FrontChannelLogoutSessionRequired bool
	BackChannelLogoutURI              string
	BackChannelLogoutSessionRequired  bool
	LogoURL                           string
	GrantTypes                        []string
	ResponseTypes                     []string
	Scopes                            []string
	Audience                          []string
	Secret                            string
	SubjectType                       string
	AuthMethod                        string
	AuthBackend                       string
}

// CreateClient within the hydra service
func (s *AdminClient) CreateClient(ctx context.Context, client AuthClient) (*AuthClient, error) {
	for _, aud := range client.Audience {
		if _, err := url.Parse(aud); err != nil {
			return nil, Error{
				Message: `invalid audience URL`,
				Err:     err,
				Code:    InvalidClientParam,
			}
		}
	}

	var md map[string]string
	if client.AuthBackend != "" {
		md = map[string]string{AuthBackend: client.AuthBackend}
	}

	cl := &models.OAuth2Client{
		AllowedCorsOrigins:                client.CORS,
		Audience:                          client.Audience,
		BackchannelLogoutSessionRequired:  client.BackChannelLogoutSessionRequired,
		BackchannelLogoutURI:              client.BackChannelLogoutURI,
		ClientID:                          client.ClientID,
		ClientName:                        client.ClientName,
		ClientSecret:                      client.Secret,
		FrontchannelLogoutSessionRequired: client.FrontChannelLogoutSessionRequired,
		FrontchannelLogoutURI:             client.FrontChannelLogoutURI,
		GrantTypes:                        client.GrantTypes,
		LogoURI:                           client.LogoURL,
		Metadata:                          md,
		Owner:                             client.OwnerID,
		PostLogoutRedirectUris:            client.PostLogoutRedirectURIs,
		RedirectUris:                      client.RedirectURIs,
		ResponseTypes:                     client.ResponseTypes,
		Scope:                             strings.Join(client.Scopes, " "),
		SubjectType:                       client.SubjectType,
		TokenEndpointAuthMethod:           client.AuthMethod,
	}
	m := admin.NewCreateOAuth2ClientParams().
		WithBody(cl).
		WithContext(ctx)

	res, err := s.adminClient.CreateOAuth2Client(m)
	if err != nil {
		return nil, Error{
			Message: "failed to create oauth2 client",
			Err:     err,
			Code:    HydraError,
		}
	}
	if res.Payload == nil {
		return nil, Error{
			Message: "failed to create oauth2 client",
			Code:    HydraError,
		}
	}

	if client.ClientID != "" && res.Payload.ClientID != client.ClientID {
		return nil, Error{
			Message: fmt.Sprintf("failed to create client with id: %s", client.ClientID),
			Code:    HydraError,
		}
	}

	cl = res.Payload
	return &AuthClient{
		OwnerID:                           cl.Owner,
		ClientID:                          cl.ClientID,
		ClientName:                        cl.ClientName,
		RedirectURIs:                      cl.RedirectUris,
		PostLogoutRedirectURIs:            cl.PostLogoutRedirectUris,
		FrontChannelLogoutURI:             cl.FrontchannelLogoutURI,
		FrontChannelLogoutSessionRequired: cl.FrontchannelLogoutSessionRequired,
		BackChannelLogoutURI:              cl.BackchannelLogoutURI,
		BackChannelLogoutSessionRequired:  cl.BackchannelLogoutSessionRequired,
		LogoURL:                           cl.LogoURI,
		GrantTypes:                        cl.GrantTypes,
		ResponseTypes:                     cl.ResponseTypes,
		Scopes:                            strings.Split(cl.Scope, " "),
		Audience:                          cl.Audience,
		Secret:                            cl.ClientSecret,
		SubjectType:                       cl.SubjectType,
		AuthMethod:                        cl.TokenEndpointAuthMethod,
	}, nil
}

// GetClient within the hydra service
func (s *AdminClient) GetClient(ctx context.Context, clientID string) (*AuthClient, error) {
	if clientID == "" {
		return nil, missingClientID
	}

	m := admin.NewGetOAuth2ClientParams().
		WithID(clientID).
		WithContext(ctx)
	res, err := s.adminClient.GetOAuth2Client(m)
	if err != nil {
		return nil, Error{
			Message: fmt.Sprintf("failed to fetch client with id: %s", clientID),
			Code:    HydraError,
			Err:     err,
		}
	}
	cl := res.Payload
	return &AuthClient{
		OwnerID:                           cl.Owner,
		CORS:                              cl.AllowedCorsOrigins,
		ClientID:                          cl.ClientID,
		ClientName:                        cl.ClientName,
		RedirectURIs:                      cl.RedirectUris,
		LogoURL:                           cl.LogoURI,
		PostLogoutRedirectURIs:            cl.PostLogoutRedirectUris,
		FrontChannelLogoutURI:             cl.FrontchannelLogoutURI,
		FrontChannelLogoutSessionRequired: cl.FrontchannelLogoutSessionRequired,
		BackChannelLogoutURI:              cl.BackchannelLogoutURI,
		BackChannelLogoutSessionRequired:  cl.BackchannelLogoutSessionRequired,
		GrantTypes:                        cl.GrantTypes,
		ResponseTypes:                     cl.ResponseTypes,
		Scopes:                            strings.Split(cl.Scope, " "),
		Audience:                          cl.Audience,
		SubjectType:                       cl.SubjectType,
		AuthMethod:                        cl.TokenEndpointAuthMethod,
	}, nil
}

// UpdateClient within the hydra service
func (s *AdminClient) UpdateClient(ctx context.Context, client AuthClient) error {
	if client.ClientID == "" {
		return missingClientID
	}

	_, err := s.GetClient(ctx, client.ClientID)
	if err != nil {
		return err
	}

	s.mux.Lock()
	defer s.mux.Unlock()
	return s.updateClientNoLock(ctx, client)
}

func (s *AdminClient) updateClientNoLock(ctx context.Context, client AuthClient) error {
	cl := &models.OAuth2Client{
		Owner:                             client.OwnerID,
		ClientID:                          client.ClientID,
		ClientName:                        client.ClientName,
		RedirectUris:                      client.RedirectURIs,
		AllowedCorsOrigins:                client.CORS,
		PostLogoutRedirectUris:            client.PostLogoutRedirectURIs,
		FrontchannelLogoutURI:             client.FrontChannelLogoutURI,
		FrontchannelLogoutSessionRequired: client.FrontChannelLogoutSessionRequired,
		BackchannelLogoutURI:              client.BackChannelLogoutURI,
		BackchannelLogoutSessionRequired:  client.BackChannelLogoutSessionRequired,
		GrantTypes:                        client.GrantTypes,
		ResponseTypes:                     client.ResponseTypes,
		Scope:                             strings.Join(client.Scopes, " "),
		Audience:                          client.Audience,
		ClientSecret:                      client.Secret,
		SubjectType:                       client.SubjectType,
		TokenEndpointAuthMethod:           client.AuthMethod,
	}
	m := admin.NewUpdateOAuth2ClientParams().
		WithID(client.ClientID).
		WithBody(cl).
		WithContext(ctx)
	if _, err := s.adminClient.UpdateOAuth2Client(m); err != nil {
		return Error{
			Message: fmt.Sprintf("failed to update client with id: %s", client.ClientID),
			Err:     err,
			Code:    HydraError,
		}
	}
	return nil
}

// DeleteClient within the hydra service
func (s *AdminClient) DeleteClient(ctx context.Context, clientID string) error {
	if clientID == "" {
		return missingClientID
	}

	m := admin.NewDeleteOAuth2ClientParams().
		WithID(clientID).
		WithContext(ctx)
	if _, err := s.adminClient.DeleteOAuth2Client(m); err != nil {
		return Error{
			Message: fmt.Sprintf("failed to delete client with id: %s", clientID),
			Err:     err,
			Code:    HydraError,
		}
	}
	return nil
}

// AddRedirectURI appends a new redirect_uri to specified client
func (s *AdminClient) AddRedirectURI(ctx context.Context, clientID, uri string) error {
	switch "" {
	case clientID:
		return missingClientID
	case uri:
		return missingURI
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	client, err := s.GetClient(ctx, clientID)
	if err != nil {
		return err
	}
	client.RedirectURIs = append(client.RedirectURIs, uri)
	return s.updateClientNoLock(ctx, *client)
}

// AddPostLogoutRedirectURI appends a new post_logout_redirect_uris to specified client
func (s *AdminClient) AddPostLogoutRedirectURI(ctx context.Context, clientID, uri string) error {
	switch "" {
	case clientID:
		return missingClientID
	case uri:
		return missingURI
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	client, err := s.GetClient(ctx, clientID)
	if err != nil {
		return err
	}
	client.PostLogoutRedirectURIs = append(client.PostLogoutRedirectURIs, uri)
	return s.updateClientNoLock(ctx, *client)
}

// newAdminClient creates a hydra admin client.
func (s *AdminClient) newAdminClient() (admin.ClientService, error) {
	hydraURL, err := url.Parse(s.host)
	if err != nil {
		return nil, errors.Wrap(err, "parsing hydra service url failed")
	}
	ht := httptransport.New(hydraURL.Host, hydraURL.Path, []string{hydraURL.Scheme})
	ht.Transport = s.tr

	client := admin.New(ht, nil)
	isReadyCh := waitHydraReady(client)
	if !<-isReadyCh {
		return nil, errors.New("hydra instance was not ready")
	}
	return client, nil
}

func waitHydraReady(client admin.ClientService) chan bool {
	isReadyCh := make(chan bool)
	isReadyParams := admin.NewIsInstanceAliveParams()

	go func() {
		const (
			maxAttempts = 20
			backoff     = 500 * time.Millisecond
		)

		for numAttempts := 1; numAttempts <= maxAttempts; numAttempts++ {
			res, err := client.IsInstanceAlive(isReadyParams)
			if err == nil && res.Payload.Status == "ok" {
				isReadyCh <- true
				return
			}

			if numAttempts > maxAttempts {
				break
			}

			time.Sleep(time.Duration(numAttempts) * backoff)
		}

		isReadyCh <- false
	}()

	return isReadyCh
}

type transport struct {
	Transport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	rq := req.Clone(req.Context())
	rq.Header.Set("X-Forwarded-Proto", "https")
	return t.Transport.RoundTrip(rq)
}
