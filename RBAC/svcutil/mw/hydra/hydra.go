package hydra

import (
	"errors"
	"net/http"
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ory/hydra-client-go/client/admin"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	errMissingClientID = errors.New("token is missing client ID")
	errMissingToken    = errors.New("token is missing")
)

type Service struct {
	// list of method to ignored
	ignoredMethods map[string]struct{}

	// scopes is a mapping between OAuth2 scope and full method name.
	// If the map is not empty, check if scope is matching with value returned in claims.
	// Otherwise all methods are allowed.
	scopes map[string]string

	// knownAudience is a whitelist defining the audiences (list of URLs) a token
	// is authorized. URLs MUST NOT contain whitespaces.
	// If the list is not empty, check if the claim's audience belongs to the known
	// audience. Otherwise, all audiences are allowed.
	knownAudience []string

	cl admin.ClientService

	// when set to true, the metadata function will not return a fatal error
	optional bool
}

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
func WithKnownAudience(aud []string) Option { return func(s *Service) { s.knownAudience = aud } }

// WithAuthScopes configures Service to check if the scope returned in claims matched a pre configured map.
func WithAuthScopes(scopes map[string]string) Option { return func(s *Service) { s.scopes = scopes } }

// WithOptional being set will allow the metaloader function to not return any fatal error that breaks
// the metadata chain
func WithOptional() Option { return func(s *Service) { s.optional = true } }

// Option is option to configure UnaryServerInterceptor method.
type Option func(*Service)

// NewService creates a new Service instance.
func NewService(adminURL string, opts ...Option) (*Service, error) {
	parsedURL, err := url.Parse(adminURL)
	if err != nil {
		return nil, err
	}
	if parsedURL.String() == "" {
		return nil, errors.New("empty hydra url given")
	}

	ht := httptransport.New(parsedURL.Host, parsedURL.Path, []string{parsedURL.Scheme})
	ht.Transport = fakeTLS(http.DefaultTransport)
	s := &Service{
		cl: admin.New(ht, nil),
	}
	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (r RoundTripFunc) RoundTrip(rq *http.Request) (*http.Response, error) { return r(rq) }

func fakeTLS(rt http.RoundTripper) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		rq := req.Clone(req.Context())
		rq.Header.Set("X-Forwarded-Proto", "https")
		return rt.RoundTrip(rq)
	})
}
