package perahub

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetRevComm send get request to perahub nonex
func (s *Svc) GetRevComm(ctx context.Context, url string) (json.RawMessage, error) {
	return s.getOrDeletePerahub(ctx, http.MethodGet, url, s.nonexAPIKey)
}

// PostRevComm send post request to perahub nonex
func (s *Svc) PostRevComm(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
	return s.postOrPutPerahub(ctx, http.MethodPost, url, body, s.nonexAPIKey)
}

// PutRevComm send put request to perahub nonex
func (s *Svc) PutRevComm(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
	return s.postOrPutPerahub(ctx, http.MethodPut, url, body, s.nonexAPIKey)
}

// DeleteRevComm send delete request to perahub nonex
func (s *Svc) DeleteRevComm(ctx context.Context, url string) (json.RawMessage, error) {
	return s.getOrDeletePerahub(ctx, http.MethodDelete, url, s.nonexAPIKey)
}
