package mainpkg

import (
	"net/http"

	"github.com/gorilla/mux"
)

// WebHandler identifies and handles HTTP requests.
// Allows HTTP and gRPC handlers to listen on the same port.
type WebHandler interface {
	IsWebRequest(*http.Request) bool
	http.Handler
}

type Middleware interface {
	WrapHandler(next http.Handler) http.Handler
}

type exclusiveHTTP struct{ http.Handler }

// HTTPOnly serves all requests using the given HTTP handler.
func HTTPOnly(h http.Handler) WebHandler              { return exclusiveHTTP{Handler: h} }
func (exclusiveHTTP) IsWebRequest(*http.Request) bool { return true }

type gorilla struct{ *mux.Router }

// WrapGorilla conforms a gorilla mux to a WebHandler.
func WrapGorilla(m *mux.Router) WebHandler          { return gorilla{Router: m} }
func (w gorilla) IsWebRequest(r *http.Request) bool { return w.Router.Match(r, &mux.RouteMatch{}) }
