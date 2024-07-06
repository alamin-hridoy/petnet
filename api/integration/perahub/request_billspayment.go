package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"brank.as/petnet/api/core/static"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
)

type BCGetTokenRequest struct {
	GrantType string `json:"grant_type"`
	TpaID     string `json:"tpa_id"`
	Scope     string `json:"scope"`
}

type BCGetTokenResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Result  BCGetTokenResult `json:"result"`
	RemcoID int              `json:"remco_id"`
}

type BCGetTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (s *Svc) BCGetToken(ctx context.Context, req BCGetTokenRequest) (*BCGetTokenResponse, error) {
	res, err := s.BillsPost(ctx, s.getUrl("bayad/bayad-center/token"), req)
	if err != nil {
		return nil, err
	}
	bcgt := &BCGetTokenResponse{}
	if err := json.Unmarshal(res, bcgt); err != nil {
		return nil, err
	}
	return bcgt, nil
}

// BillsPost request to perahub bill pay
func (s *Svc) BillsPost(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
	log := logging.FromContext(ctx)
	ptnr := phmw.GetPartner(ctx)
	byc := static.BYCBP
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

	if ptnr == byc && url != s.getUrl("bayad/bayad-center/token") {
		res, err := s.BCGetToken(ctx, BCGetTokenRequest{
			GrantType: "client_credentials",
			TpaID:     "PP01",
			Scope:     "mecom-auth/all",
		})
		if err != nil {
			logging.WithError(err, log).Error("couldn't get token")
			return nil, err
		}
		bcToken := res.Result.AccessToken
		req.Header.Add("X-Bayad-Center-Token", bcToken)
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
			Type:       BillsError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("bills payment error")

		return nil, err
	}

	l.Debug("bills payment success")

	return b, nil
}

func (s *Svc) BillsGet(ctx context.Context, url string) (json.RawMessage, error) {
	log := logging.FromContext(ctx)
	ptnr := phmw.GetPartner(ctx)
	byc := static.BYCBP
	renewCache := false
renewAuthCache:
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)

		return nil, err
	}

	if ptnr == byc && url != s.getUrl("bayad/bayad-center/token") {
		res, err := s.BCGetToken(ctx, BCGetTokenRequest{
			GrantType: "client_credentials",
			TpaID:     "PP01",
			Scope:     "mecom-auth/all",
		})
		if err != nil {
			logging.WithError(err, log).Error("couldn't get token")
			return nil, err
		}
		bcToken := res.Result.AccessToken
		req.Header.Add("X-Bayad-Center-Token", bcToken)
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
			Type:       BillsError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("bills payment error")

		return nil, err
	}

	l.Debug("bilss payment success")

	return b, nil
}
