package mw

import "net/http"

var _ http.RoundTripper = RoundTripFunc(nil)

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (r RoundTripFunc) RoundTrip(rq *http.Request) (*http.Response, error) { return r(rq) }

func JSONType(r http.RoundTripper) http.RoundTripper {
	if r == nil {
		r = http.DefaultTransport
	}
	return RoundTripFunc(func(rq *http.Request) (*http.Response, error) {
		rq.Header.Set("Content-Type", "application/json")
		rq.Header.Set("Accept", "application/json")
		return r.RoundTrip(rq)
	})
}
