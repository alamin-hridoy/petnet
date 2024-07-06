package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"

	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
)

func (s *Svc) cicoURL(name string) string {
	u := *s.cicoUrl
	if u.Path == "" {
		u = *s.nonexUrl
		u.Path = strings.ReplaceAll(u.Path, "remit/nonex/", "cico/wrapper")
	}
	u.Path = path.Join(u.Path, name)
	return u.String()
}

// cicoGet request to perahub cico
func (s *Svc) cicoGet(ctx context.Context, url string) (json.RawMessage, error) {
	log := logging.FromContext(ctx)
	renewCache := false
renewAuthCache:
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)

		return nil, err
	}

	bearerSet, err := s.setAuthBearer(ctx, req, renewCache)
	if err != nil {
		logging.WithError(err, log).Error("setting bearer token for: ", phmw.GetPartner(ctx))

		return nil, err
	}

	setPerahubHeaders(req, s.nonexAPIKey)
	resp, err := s.cl.Do(req.WithContext(ctx))
	if err != nil {
		log.Error(err)

		return nil, err
	}

	if bearerSet && resp.StatusCode == http.StatusUnauthorized && !renewCache {
		renewCache = true
		resp.Body.Close()
		goto renewAuthCache
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)
	resp.Body = io.NopCloser(io.TeeReader(resp.Body, buf))

	var b json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		log.WithField("http_status", resp.StatusCode).WithField("response_body", buf.String()).Debug("invalid json")

		return nil, &Error{
			GRPCCode:   codes.Internal,
			Code:       strconv.Itoa(resp.StatusCode),
			Msg:        codes.Internal.String(),
			UnknownErr: buf.String(),
			Type:       CicoError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("cico error")

		return nil, err
	}

	l.Debug("cico success")

	return b, nil
}

// cicoPost request to perahub cico
func (s *Svc) cicoPost(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
	log := logging.FromContext(ctx)
	reqBody, err := json.Marshal(body)
	if err != nil {
		log.Error(err)

		return nil, err
	}

	logReqBody(ctx, url, json.RawMessage(reqBody))

	renewCache := false
renewAuthCache:
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Error(err)

		return nil, err
	}

	bearerSet, err := s.setAuthBearer(ctx, req, renewCache)
	if err != nil {
		logging.WithError(err, log).Error("setting bearer token for: ", phmw.GetPartner(ctx))

		return nil, err
	}

	setPerahubHeaders(req, s.nonexAPIKey)
	resp, err := s.cl.Do(req.WithContext(ctx))
	if err != nil {
		log.Error(err)

		return nil, err
	}

	if bearerSet && resp.StatusCode == http.StatusUnauthorized && !renewCache {
		renewCache = true
		resp.Body.Close()
		goto renewAuthCache
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)
	resp.Body = io.NopCloser(io.TeeReader(resp.Body, buf))

	var b json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		log.WithField("http_status", resp.StatusCode).WithField("response_body", buf.String()).Debug("invalid json")

		return nil, &Error{
			GRPCCode:   codes.Internal,
			Code:       strconv.Itoa(resp.StatusCode),
			Msg:        codes.Internal.String(),
			UnknownErr: buf.String(),
			Type:       CicoError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("cico error")

		return nil, err
	}

	l.Debug("cico success")

	return b, nil
}

// cicoPut request to perahub cico
func (s *Svc) cicoPut(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
	log := logging.FromContext(ctx)
	reqBody, err := json.Marshal(body)
	if err != nil {
		log.Error(err)

		return nil, err
	}

	logReqBody(ctx, url, json.RawMessage(reqBody))

	renewCache := false
renewAuthCache:
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Error(err)

		return nil, err
	}

	bearerSet, err := s.setAuthBearer(ctx, req, renewCache)
	if err != nil {
		logging.WithError(err, log).Error("setting bearer token for: ", phmw.GetPartner(ctx))

		return nil, err
	}

	setPerahubHeaders(req, s.nonexAPIKey)
	resp, err := s.cl.Do(req.WithContext(ctx))
	if err != nil {
		log.Error(err)

		return nil, err
	}

	if bearerSet && resp.StatusCode == http.StatusUnauthorized && !renewCache {
		renewCache = true
		resp.Body.Close()
		goto renewAuthCache
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)
	resp.Body = io.NopCloser(io.TeeReader(resp.Body, buf))

	var b json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		log.WithField("http_status", resp.StatusCode).WithField("response_body", buf.String()).Debug("invalid json")

		return nil, &Error{
			GRPCCode:   codes.Internal,
			Code:       strconv.Itoa(resp.StatusCode),
			Msg:        codes.Internal.String(),
			UnknownErr: buf.String(),
			Type:       CicoError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("cico error")

		return nil, err
	}

	l.Debug("cico success")

	return b, nil
}

// cicoDelete request to perahub cico
func (s *Svc) cicoDelete(ctx context.Context, url string) (json.RawMessage, error) {
	log := logging.FromContext(ctx)
	renewCache := false
renewAuthCache:
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Error(err)

		return nil, err
	}

	bearerSet, err := s.setAuthBearer(ctx, req, renewCache)
	if err != nil {
		logging.WithError(err, log).Error("setting bearer token for: ", phmw.GetPartner(ctx))

		return nil, err
	}

	setPerahubHeaders(req, s.nonexAPIKey)
	resp, err := s.cl.Do(req.WithContext(ctx))
	if err != nil {
		log.Error(err)

		return nil, err
	}

	if bearerSet && resp.StatusCode == http.StatusUnauthorized && !renewCache {
		renewCache = true
		resp.Body.Close()
		goto renewAuthCache
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)
	resp.Body = io.NopCloser(io.TeeReader(resp.Body, buf))

	var b json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		log.WithField("http_status", resp.StatusCode).WithField("response_body", buf.String()).Debug("invalid json")

		return nil, &Error{
			GRPCCode:   codes.Internal,
			Code:       strconv.Itoa(resp.StatusCode),
			Msg:        codes.Internal.String(),
			UnknownErr: buf.String(),
			Type:       CicoError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("cico error")

		return nil, err
	}

	l.Debug("cico success")

	return b, nil
}
