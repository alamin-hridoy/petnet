package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
)

// RtaPost request to perahub rta
func (s *Svc) RtaPost(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
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
			Type:       RTAError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("rta error")

		return nil, err
	}

	l.Debug("rta success")

	return b, nil
}
