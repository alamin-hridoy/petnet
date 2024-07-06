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

	"google.golang.org/grpc/codes"

	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/serviceutil/logging"
)

type remitanceError struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Errors  interface{} `json:"errors"`
}

func (s *Svc) remitanceURL(name string) string {
	u := *s.phRemittanceUrl
	if u.Path == "" {
		u = *s.nonexUrl
		u.Path = strings.ReplaceAll(u.Path, "remit/nonex/", "remit/dmt/")
	}
	u.Path = path.Join(u.Path, name)
	return u.String()
}

// remitanceGet request to perahub remitance
func (s *Svc) remitanceGet(ctx context.Context, url string) (json.RawMessage, error) {
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
			Type:       RemitanceError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("remitance error")
		return nil, err
	} else {
		l.Debug("remitance success")
	}
	return b, nil
}

// remitancePost request to perahub remitance
func (s *Svc) remitancePost(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
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
			Type:       RemitanceError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("remittance error")
		return nil, err
	} else {
		l.Debug("remittance success")
	}
	return b, nil
}

// remitancePut request to perahub remitance
func (s *Svc) remitancePut(ctx context.Context, url string, body interface{}) (json.RawMessage, error) {
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
			Type:       RemitanceError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("remitance error")
		return nil, err
	} else {
		l.Debug("remitance success")
	}
	return b, nil
}

// remitanceDelete request to perahub remitance
func (s *Svc) remitanceDelete(ctx context.Context, url string) (json.RawMessage, error) {
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
			Type:       RemitanceError,
		}
	}

	l := log.WithField("http_status", resp.StatusCode).WithField("response_body", json.RawMessage(buf.String()))
	if err := handleRemittanceErr(ctx, b, resp.StatusCode); err != nil {
		l.Error("remitance error")
		return nil, err
	} else {
		l.Debug("remitance success")
	}
	return b, nil
}

func handleRemittanceErr(ctx context.Context, b []byte, sts int) error {
	log := logging.FromContext(ctx)
	if sts == http.StatusOK {
		return nil
	}

	remitErr := &remitanceError{}
	if err := json.Unmarshal(b, remitErr); err != nil {
		logging.WithError(err, log).Error("non standard error")
		return &Error{
			Code:       strconv.Itoa(sts),
			GRPCCode:   codes.Internal,
			Msg:        codes.Internal.String(),
			UnknownErr: strings.Join(strings.Fields(string(b)), ""),
			Type:       RemitanceError,
		}
	}

	var stscode codes.Code
	switch sts {
	case http.StatusConflict:
		stscode = codes.AlreadyExists
	case http.StatusUnprocessableEntity:
		stscode = codes.InvalidArgument
	case http.StatusNotFound:
		stscode = codes.NotFound
	case http.StatusBadRequest:
		stscode = codes.InvalidArgument
	case http.StatusPaymentRequired:
		stscode = codes.InvalidArgument
	case http.StatusCreated:
		stscode = codes.InvalidArgument
	default:
		stscode = codes.Internal
	}

	serr := &Error{}
	serr.Msg = remitErr.Message
	serr.Type = RemitanceError
	serr.GRPCCode = stscode
	serr.Msg, serr.Errors = parseErrorMessage(remitErr.Errors)

	switch v := remitErr.Code.(type) {
	case string:
		serr.Code = v
	case int:
		serr.Code = strconv.Itoa(v)
	case int32:
		serr.Code = strconv.Itoa(int(v))
	case float64:
		serr.Code = strconv.Itoa(int(v))
	default:
		log.Debugf("unexpected type %T", v)
	}

	if serr.Code == "" {
		serr.Code = strconv.Itoa(sts)
	}

	if remitErr.Error == nil {
		if serr.Msg == "" {
			serr.Msg = stscode.String()
			serr.UnknownErr = strings.Join(strings.Fields(string(b)), "")
		}

		return serr
	}

	msg, details := parseErrorMessage(remitErr.Error)
	if msg != "" {
		serr.Msg = msg
	}

	if remitErr.Errors != nil {
		msg, details = parseErrorMessage(remitErr.Errors)
		if msg != "" {
			serr.Msg = msg
		}
	}

	if len(details) > 0 {
		serr.Errors = details
	}

	if serr.Msg == "" {
		serr.Msg = stscode.String()
		serr.UnknownErr = strings.Join(strings.Fields(string(b)), "")
	}

	return serr
}

func parseErrorMessage(e interface{}) (string, map[string][]string) {
	switch v := e.(type) {
	case string:
		return v, nil

	case []string:
		if len(v) != 0 {
			return strings.Join(v, ", "), nil
		}

	case map[string][]string:
		fieldMsg := ""
		separator := ""
		for _, msgInterface := range v {
			msg, _ := parseErrorMessage(msgInterface)
			if msg == "" {
				continue
			}

			fieldMsg += separator + msg
			separator = ", "
		}

		return fieldMsg, v

	case map[string]interface{}:
		fieldMsg := ""
		separator := ""
		for key, msgInterface := range v {
			if key == "message" {
				return parseErrorMessage(msgInterface)
			}

			msg, _ := parseErrorMessage(msgInterface)
			if msg == "" {
				continue
			}

			fieldMsg += separator + msg
			separator = ", "
		}

		return fieldMsg, nil

	case []interface{}:
		for _, k := range v {
			return parseErrorMessage(k)
		}

	}

	return "", nil
}
