package perahub

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"brank.as/petnet/api/core/static"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
)

// TODO: deprecate to use post helper
type ResponseHeader struct {
	ErrorCode string `json:"errorcode"`
	Message   string `json:"message"`
}

type response struct {
	Header struct {
		ErrorCode string `json:"errorcode"`
		Message   string `json:"message"`
	} `json:"header"`
	Body json.RawMessage `json:"body"`
}

type PerahubRequest struct {
	WU InnerRequest `json:"uspwuapi"`
}

type InnerRequest struct {
	Header    RequestHeader `json:"header"`
	Body      RequestBody   `json:"body"`
	Signature string        `json:"signature"`
}

type RequestHeader struct {
	Coy          string      `json:"coy"`
	Token        string      `json:"token"`
	LocationCode string      `json:"location_code"`
	UserCode     json.Number `json:"user_code"`
	ClientIP     string      `json:"clientip"`
	IsWeb        string      `json:"isweb"`
}

type RequestBody struct {
	Module  string          `json:"module"`
	Request string          `json:"request"`
	Param   json.RawMessage `json:"param"`
}

type oauth struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

var (
	ptnrAuthMu = sync.Mutex{}
	ptnrAuth   = map[string]oauth{}

	putAndPostMethods   = map[string]bool{http.MethodPost: true, http.MethodPut: true}
	getAndDeleteMethods = map[string]bool{http.MethodGet: true, http.MethodDelete: true}
)

func setPartnerAuth(ptnr string, a oauth) {
	ptnrAuthMu.Lock()
	defer ptnrAuthMu.Unlock()
	ptnrAuth[ptnr] = a
}

func getPartnerAuth(ptnr string) (*oauth, error) {
	ptnrAuthMu.Lock()
	defer ptnrAuthMu.Unlock()
	a, ok := ptnrAuth[ptnr]
	if !ok {
		return nil, fmt.Errorf("auth doesn't exist for partner: %s", ptnr)
	}
	return &a, nil
}

func (s *Svc) setAuthBearer(ctx context.Context, req *http.Request, renewCache bool) (bool, error) {
	ptnr := phmw.GetPartner(ctx)
	var token string
	switch ptnr {
	case static.WISECode:
		if req.URL.String() == s.nonexUrl.String()+"transferwise/oauth/token" {
			return false, nil
		}
		a, err := getPartnerAuth(ptnr)
		if err != nil || time.Now().After(a.Expiry.Add(-time.Minute)) || renewCache {
			creds := s.ptnrAuthCreds[ptnr]
			res, err := s.WISEGetTokens(ctx, WISEGetTokensReq{
				ClientID:     creds.ClientID,
				ClientSecret: creds.ClientSecret,
			})
			if err != nil {
				return false, err
			}
			exp, err := res.ExpiresIn.Int64()
			if err != nil {
				return false, err
			}
			setPartnerAuth(ptnr, oauth{
				AccessToken:  res.AccessToken,
				RefreshToken: res.RefreshToken,
				Expiry:       time.Now().Add(time.Duration(exp) * time.Second),
			})
			a, _ = getPartnerAuth(ptnr)
		}
		token = a.AccessToken
	default:
		return false, nil
	}
	req.Header.Add("Authorization", "Bearer "+token)
	return true, nil
}

// post request to perahub gateway
func (s *Svc) post(ctx context.Context, url string, body PerahubRequest) (json.RawMessage, error) {
	log := logging.FromContext(ctx)
	reqBody, err := json.Marshal(body)
	if err != nil {
		if err := ConvertErr(err, NonexError); err != nil {
			return nil, err
		}
		return nil, err
	}

	log.WithField("url", url).WithField("request body", json.RawMessage(reqBody)).Debug("sending")

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		if err := ConvertErr(err, NonexConError); err != nil {
			return nil, err
		}
		return nil, err
	}

	resp, err := s.cl.Do(req.WithContext(ctx))
	if err != nil {
		if err := ConvertErr(err, NonexConError); err != nil {
			return nil, err
		}
		return nil, err
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)
	resp.Body = io.NopCloser(io.TeeReader(resp.Body, buf))

	res := &response{}
	if err := json.NewDecoder(resp.Body).Decode(&struct {
		WU *response `json:"uspwuapi"`
	}{WU: res}); err != nil {
		log.WithField("http_status", resp.StatusCode).WithField("response_body", buf.String()).Debug("invalid json")
		return nil, &Error{
			GRPCCode:   codes.Internal,
			Code:       strconv.Itoa(resp.StatusCode),
			Msg:        codes.Internal.String(),
			UnknownErr: buf.String(),
			Type:       NonexError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleWUErr(ctx, *res, resp.StatusCode); err != nil {
		l.Error("wu error")
		return nil, err
	} else {
		l.Debug("wu success")
	}
	return res.Body, nil
}

// billsPaymentGet request to perahub biller
func (s *Svc) billsPaymentGet(ctx context.Context, url string) (json.RawMessage, error) {
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
			Type:       BillerError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleBillsPaymentErr(ctx, b, url, resp.StatusCode); err != nil {
		l.Error("Biller error")
		return nil, err
	} else {
		l.Debug("Biller success")
	}
	return b, nil
}

// getNonex request to perahub nonex
func (s *Svc) getNonex(ctx context.Context, url string) (json.RawMessage, error) {
	return s.getOrDeletePerahub(ctx, http.MethodGet, url, s.nonexAPIKey)
}

// postNonex request to perahub nonex
func (s *Svc) postNonex(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
	return s.postOrPutPerahub(ctx, http.MethodPost, url, body, s.nonexAPIKey)
}

func (s *Svc) postOrPutPerahub(ctx context.Context, method, url string, body interface{}, apiKey string) (json.RawMessage, error) {
	log := logging.FromContext(ctx)
	if _, ok := putAndPostMethods[method]; !ok {
		log.Error("invalid method for put or post: " + method)
		return nil, errMethodNotAllowed
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
		if err := ConvertErr(err, NonexError); err != nil {
			return nil, err
		}
		return nil, err
	}

	logReqBody(ctx, url, json.RawMessage(reqBody))

	renewCache := false
renewAuthCache:
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Error(err)
		if err := ConvertErr(err, NonexError); err != nil {
			return nil, err
		}
		return nil, err
	}

	bearerSet, err := s.setAuthBearer(ctx, req, renewCache)
	if err != nil {
		logging.WithError(err, log).Error("setting bearer token for: ", phmw.GetPartner(ctx))
		if err := ConvertErr(err, NonexError); err != nil {
			return nil, err
		}
		return nil, err
	}

	setPerahubHeaders(req, apiKey)
	resp, err := s.cl.Do(req.WithContext(ctx))
	if err != nil {
		log.Error(err)
		if err := ConvertErr(err, NonexConError); err != nil {
			return nil, err
		}
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
			Type:       NonexError,
		}
	}
	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleNonexErr(ctx, b, url, resp.StatusCode); err != nil {
		l.Error("nonex error")
		return nil, err
	} else {
		l.Debug("nonex success")
	}
	return b, nil
}

func (s *Svc) getOrDeletePerahub(ctx context.Context, method, url, apiKey string) (json.RawMessage, error) {
	log := logging.FromContext(ctx)
	if _, ok := getAndDeleteMethods[method]; !ok {
		log.Error("invalid method for get or delete: " + method)
		return nil, errMethodNotAllowed
	}

	renewCache := false
renewAuthCache:
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Error(err)
		if err := ConvertErr(err, NonexError); err != nil {
			return nil, err
		}
		return nil, err
	}

	bearerSet, err := s.setAuthBearer(ctx, req, renewCache)
	if err != nil {
		logging.WithError(err, log).Error("setting bearer token for: ", phmw.GetPartner(ctx))
		if err := ConvertErr(err, NonexConError); err != nil {
			return nil, err
		}
		return nil, err
	}

	setPerahubHeaders(req, apiKey)
	resp, err := s.cl.Do(req.WithContext(ctx))
	if err != nil {
		log.Error(err)
		if err := ConvertErr(err, NonexConError); err != nil {
			return nil, err
		}
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
			Type:       NonexError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleNonexErr(ctx, b, url, resp.StatusCode); err != nil {
		l.Error("nonex error")
		return nil, err
	} else {
		l.Debug("nonex success")
	}
	return b, nil
}

// billsPaymentPost request to perahub nonex
func (s *Svc) billsPaymentPost(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
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
			Type:       BillerError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleBillsPaymentErr(ctx, b, url, resp.StatusCode); err != nil {
		l.Error("nonex error")
		return nil, err
	} else {
		l.Debug("nonex success")
	}
	return b, nil
}

func (s *Svc) moduleURL(name, req string) string {
	u := *s.baseUrl
	switch name {
	case "wuso", "wupo", "wuqp", "wu":
		u.Path = path.Join(u.Path, "1.1", name+"-"+req)
	case "sdq":
		u.Path = path.Join(u.Path, "1.1", name)
	default:
		u.Path = path.Join(u.Path, name)
	}
	return u.String()
}

func (s *Svc) nonexURL(name string) string {
	u := *s.nonexUrl
	u.Path = path.Join(u.Path, name)
	return u.String()
}

func (s *Svc) phTransactURL(name string) string {
	u := *s.phTransactUrl
	u.Path = path.Join(u.Path, name)
	return u.String()
}

func (s *Svc) getUrl(name string) string {
	u := *s.billsUrl
	if u.Path == "" {
		u = *s.nonexUrl
		u.Path = strings.ReplaceAll(u.Path, "remit/nonex/", "billspay")
	}
	u.Path = path.Join(u.Path, name)
	return u.String()
}

func (s *Svc) billerURL(name string) string {
	u := *s.billerUrl
	if u.Path == "" {
		u = *s.nonexUrl
		u.Path = strings.ReplaceAll(u.Path, "remit/nonex/", "billspay/wrapper/api/")
	}
	u.Path = path.Join(u.Path, name)
	return u.String()
}

type RequestOptions struct {
	LocationCode string
	UserCode     json.Number
}

type RequestOptionFunc func(opt *RequestOptions)

func WithLocationCode(locationCode string) RequestOptionFunc {
	return func(opt *RequestOptions) {
		if locationCode != "" {
			opt.LocationCode = locationCode
		}
	}
}

func WithUserCode(userCode json.Number) RequestOptionFunc {
	return func(opt *RequestOptions) {
		if userCode != "" {
			opt.UserCode = userCode
		}
	}
}

func (s *Svc) newParahubRequest(ctx context.Context, module, request string, param interface{}, opts ...RequestOptionFunc) (*PerahubRequest, error) {
	prm, err := json.Marshal(param)
	if err != nil {
		if err := ConvertErr(err, NonexError); err != nil {
			return nil, err
		}
		return nil, err
	}
	sgn := "placeholderSignature"
	if s.signKey != nil {
		h := crypto.SHA256.New()
		h.Write(prm)
		s, err := rsa.SignPKCS1v15(rand.Reader, s.signKey, crypto.SHA256, h.Sum(nil))
		if err != nil {
			if err := ConvertErr(err, NonexError); err != nil {
				return nil, err
			}
			return nil, err
		}
		sgn = string(s)
	}

	var coy string
	lc := "1"
	if s.phEnv == "live" {
		dsa := phmw.GetDsaCode(ctx)
		coy = phmw.GetCoy(ctx)
		lc = "drp"
		if dsa == "DKT" {
			coy = "DKT"
			lc = "DKT"
		}
	} else {
		switch module {
		case "SignOn", "UpdateInfo":
			coy = "yondu"
		default:
			coy = "usp"
		}
	}

	if coy == "" {
		coy = "usp"
	}

	reqOpts := &RequestOptions{}
	for _, opt := range opts {
		opt(reqOpts)
	}

	userCode := json.Number("1")
	if reqOpts.UserCode != "" {
		userCode = reqOpts.UserCode
	}

	if reqOpts.LocationCode != "" {
		lc = reqOpts.LocationCode
	}

	return &PerahubRequest{
		InnerRequest{
			Header: RequestHeader{
				Coy:          coy,
				Token:        s.token,
				LocationCode: lc,
				UserCode:     userCode,
				ClientIP:     "127.0.0.1",
				IsWeb:        "1",
			},
			Body: RequestBody{
				Module:  module,
				Request: request,
				Param:   prm,
			},
			Signature: sgn,
		},
	}, nil
}

func logReqBody(ctx context.Context, url string, reqBody json.RawMessage) {
	log := logging.FromContext(ctx)
	switch phmw.GetPartner(ctx) {
	// obfuscate for wise client id and secret
	case static.WISECode:
		if strings.Contains(url, "transferwise/oauth/token") {
			log.WithField("url", url).WithField("request body", `{"client_id": "****", "client_secret": "****"}`).Debug("sending")
			return
		}
	}
	log.WithField("url", url).WithField("request body", json.RawMessage(reqBody)).Debug("sending")
}

func setPerahubHeaders(req *http.Request, apiKey string) {
	req.Header.Set("X-Perahub-Gateway-Token", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}
