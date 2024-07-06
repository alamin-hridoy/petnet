package perahub

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetMicroInsurance send get request to perahub nonex
func (s *Svc) GetMicroInsurance(ctx context.Context, url string) (json.RawMessage, error) {
	return s.getOrDeletePerahub(ctx, http.MethodGet, url, s.nonexAPIKey)
}

// PostMicroInsurance send post request to perahub nonex
func (s *Svc) PostMicroInsurance(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
	return s.postOrPutPerahub(ctx, http.MethodPost, url, body, s.nonexAPIKey)
}
