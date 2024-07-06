package hydra

import (
	"errors"
	"net/http"
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ory/hydra-client-go/client"
)

type transport struct {
	Transport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	rq := req.Clone(req.Context())
	rq.Header.Set("X-Forwarded-Proto", "https")
	return t.Transport.RoundTrip(rq)
}

var ErrInvalidToken = errors.New("invalid token")

type Service struct {
	// list of method to ignored
	ignoredMethods map[string]struct{}

	// authScopes is a mapping between OAuth2 scope and full method name.
	// If the map is not empty, check if scope is matching with value returned in claims.
	// Otherwise all methods are allowed.
	authScopes map[string]string

	// knownAudience is a whitelist defining the audiences (list of URLs) a token
	// is authorized. URLs MUST NOT contain whitespaces.
	// If the list is not empty, check if the claim's audience belongs to the known
	// audience. Otherwise, all audiences are allowed.
	//
	// TODO:
	// this field will probably be invalidated once OB-3148 rolled out;
	// adding it as quick fix for now to quickly deliver the business requirements
	// that tokens for sandbox and live shouldn't be interchangeable,
	// see OB-3540 and OB-3584 - deprecate as need be.
	knownAudience []string

	admin

	// when set to true, the metadata function will not return a fatal error
	optional bool
}

// Option is option to configure UnaryServerInterceptor method.
type Option func(*Service)

// WithIgnoredMethods configures Service to ignore given method and pass to next handler.
func WithIgnoredMethods(methods []string) Option {
	return func(s *Service) {
		ignoredMethods := make(map[string]struct{}, len(methods))
		for _, method := range methods {
			ignoredMethods[method] = struct{}{}
		}
		s.ignoredMethods = ignoredMethods
	}
}

// WithKnownAudience configures Service to check if the audience returned in claims matched
// the pre-configured known audiences.
func WithKnownAudience(aud []string) Option {
	return func(s *Service) {
		s.knownAudience = aud
	}
}

// WithAuthScopes configures Service to check if the scope returned in claims matched a pre configured map.
func WithAuthScopes(scopes map[string]string) Option {
	return func(s *Service) {
		s.authScopes = scopes
	}
}

// WithOptional being set will allow the metaloader function to not return any fatal error that breaks
// the metadata chain
func WithOptional() Option {
	return func(s *Service) {
		s.optional = true
	}
}

// NewService creates a new Service instance.
func NewService(httpClient *http.Client, hydraAdminURL string, opts ...Option) (*Service, error) {
	parsedURL, err := url.Parse(hydraAdminURL)
	if err != nil {
		return nil, err
	}
	if parsedURL.String() == "" {
		return nil, errors.New("empty hydra url given")
	}

	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}

	ht := httptransport.New(parsedURL.Host, parsedURL.Path, []string{parsedURL.Scheme})
	ht.Transport = &transport{
		Transport: http.DefaultTransport,
	}
	s.admin = admin{cl: client.New(ht, nil).Admin}
	return s, nil
}
