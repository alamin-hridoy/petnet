package perahub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/logging"
)

const (
	TFErr           = "Transfast Error"
	IRErr           = "IRemit Error"
	RIAErr          = "Ria Error"
	MBErr           = "Metrobank Error"
	RMErr           = "Remitly Error"
	WUErr           = "WU"
	BPErr           = "BPI Error"
	USSCErr         = "USSC Error"
	ICErr           = "InstaCash Error"
	JPRErr          = "JapanRemit Error"
	UNTErr          = "Uniteller Error"
	WISEErr         = "WISE Error"
	CEBErr          = "Cebuana Error"
	CEBINTErr       = "Cebuana Intl Error"
	AYANNAHErr      = "Ayannah Error"
	IEErr           = "IntelExpress Error"
	PerahubRemitErr = "PerahubRemit Error"
)

var ptnrErr = map[string]struct{}{
	TFErr:  {},
	IRErr:  {},
	RIAErr: {},
	MBErr:  {},
	// todo: remitly doesn't have this structure at the moment, it returns
	// the error in the format
	//{
	//   "code": 404,
	//   "message": "Transfer not found.",
	//   "remco_id": 21
	//}
	// let petned know about this
	RMErr:      {},
	WUErr:      {},
	BPErr:      {},
	USSCErr:    {},
	ICErr:      {},
	JPRErr:     {},
	UNTErr:     {},
	CEBErr:     {},
	CEBINTErr:  {},
	AYANNAHErr: {},
	IEErr:      {},
}

type nonexError struct {
	Code    interface{}         `json:"code"`
	Msg     string              `json:"message"`
	Error   interface{}         `json:"error"`
	Errors  map[string][]string `json:"errors"`
	Details []Detail            `json:"details"`
}

type billerError struct {
	Code    interface{}         `json:"code"`
	Msg     string              `json:"message"`
	Error   interface{}         `json:"error"`
	Result  interface{}         `json:"result"`
	RemcoID int                 `json:"remco_id"`
	Errors  map[string][]string `json:"errors"`
	Details []Detail            `json:"details"`
}

type ErrorType struct {
	Type string `json:"type"`
	Msg  string `json:"message"`
}

type Detail struct {
	Msg string `json:"message"`
}

type Error struct {
	Code       string
	GRPCCode   codes.Code
	Msg        string
	UnknownErr string
	Type       ErrType
	Errors     map[string][]string
}

func (r *Error) Error() string {
	return r.Msg
}

type (
	ErrMsg  string
	ErrCode string
	ErrType string
)

const BillPayError = "BILLER"

const (
	NonexError          ErrType = "NONEX"
	NonexConError       ErrType = "NONEXCON"
	BillerError         ErrType = "BILLER"
	RemitanceError      ErrType = "REMITANCE"
	CicoError           ErrType = "CICO"
	DRPError            ErrType = "DRP"
	RTAError            ErrType = "RTA"
	BillsError          ErrType = "BILLSPAYMENT"
	PartnerError        ErrType = "PARTNER"
	MicroInsuranceError ErrType = "MICROINSURANCE"
)

var errMethodNotAllowed error = &Error{
	GRPCCode: codes.Unimplemented,
	Code:     codes.Unimplemented.String(),
	Msg:      "Method not allowed",
	Type:     DRPError,
}

func handleNonexErr(ctx context.Context, b []byte, url string, sts int) error {
	log := logging.FromContext(ctx)
	if sts == http.StatusOK || sts == http.StatusCreated || sts == http.StatusAccepted {
		return nil
	}

	nErr := &nonexError{}

	if err := json.Unmarshal(b, nErr); err != nil {
		logging.WithError(err, log).Error("unmarshal nonex error")
		// Should be perahub internal error, because invalid error message from perahub
		return &Error{
			Code:       strconv.Itoa(int(codes.Internal)),
			GRPCCode:   codes.Internal,
			Msg:        "Perahub internal error",
			UnknownErr: strings.Join(strings.Fields(string(b)), ""),
			Type:       NonexError,
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
	default:
		stscode = codes.Internal
	}

	serr := &Error{}
	serr.Msg = nErr.Msg
	serr.Type = PartnerError
	serr.GRPCCode = stscode
	switch v := nErr.Error.(type) {
	case string:
		if v != "" {
			serr.Msg = v
		}
	case map[string]interface{}:
		msg, ok := v["message"]
		if !ok {
			break
		}

		switch innerMsg := msg.(type) {
		case string:
			if innerMsg != "" {
				serr.Msg = innerMsg
			}

		case []interface{}:
			if len(innerMsg) > 0 {
				if inMap, ok := innerMsg[0].(map[string]interface{}); ok {
					if inMsg, inMapOk := inMap["Message"]; inMapOk && inMsg != "" {
						serr.Msg = inMsg.(string)
					}
				}
			}
		}
	}
	for _, d := range nErr.Details {
		if d.Msg != "" {
			serr.Msg = d.Msg
		}
	}
	switch v := nErr.Code.(type) {
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

	// TODO(vitthal): move to specific parther
	if sts == http.StatusBadRequest && nErr.Code == "99" {
		return status.Error(codes.NotFound, "404 Not Found")
	}
	if sts == http.StatusBadRequest && nErr.Code == "10202001" {
		return status.Error(codes.AlreadyExists, "Transaction Already Claimed")
	}
	if serr.Msg == "" {
		serr.Msg = stscode.String()
		serr.UnknownErr = strings.Join(strings.Fields(string(b)), "")
	}
	if nErr.Errors != nil {
		serr.Errors = nErr.Errors
	}
	return serr
}

func handleWUErr(ctx context.Context, r response, sts int) error {
	s := strings.SplitN(r.Header.Message, " ", 2)
	var code, msg string
	if len(s) == 1 {
		msg = r.Header.Message
	} else {
		code, msg = s[0], strings.TrimSpace(strings.TrimPrefix(s[1], "-"))
	}

	if _, ok := ptnrErr[r.Header.ErrorCode]; ok {
		return &Error{
			GRPCCode: codes.Internal,
			Code:     code,
			Msg:      msg,
			Type:     PartnerError,
		}
	}
	if r.Header.ErrorCode != "1" && !strings.Contains(r.Header.Message, "Success") {
		return &Error{
			GRPCCode:   codes.Internal,
			Code:       code,
			Msg:        "internal error occurred",
			Type:       NonexError,
			UnknownErr: string(r.Body) + r.Header.ErrorCode + r.Header.Message,
		}
	}
	return nil
}

// GRPCError creates a nonex GRPC error with details.
func GRPCError(c codes.Code, msg string, details *ppb.Error) error {
	if c == codes.OK {
		c = codes.Internal
	}
	p := status.New(c, msg)
	if e, err := p.WithDetails(details); err == nil {
		return e.Err()
	}
	return p.Err()
}

func FromError(err error) *ppb.Error {
	fmt.Printf("from %+v\\n", err)
	s, ok := status.FromError(err)
	if !ok {
		fmt.Printf("%+v %T", err, err)
		return nil
	}
	fmt.Printf("status %+v\\n", s)
	if len(s.Proto().GetDetails()) == 0 {
		fmt.Println(s.Proto(), s.Proto().GetDetails())
		return nil
	}
	fmt.Printf("details %+v\\n", s.Proto())
	det := &ppb.Error{}
	if err := s.Proto().Details[0].UnmarshalTo(det); err != nil {
		fmt.Println("unmarshal", err)
		return nil
	}
	fmt.Printf("partner error %+v\\n", det)
	return det
}

func handleBillsPaymentErr(ctx context.Context, b []byte, url string, sts int) error {
	log := logging.FromContext(ctx)
	if sts == http.StatusOK {
		return nil
	}

	billErr := &billerError{}
	if err := json.Unmarshal(b, billErr); err != nil {
		logging.WithError(err, log).Error("non standard error")
		return &Error{
			Code:       strconv.Itoa(sts),
			GRPCCode:   codes.Internal,
			Msg:        codes.Internal.String(),
			UnknownErr: strings.Join(strings.Fields(string(b)), ""),
			Type:       BillerError,
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
	serr.Msg = billErr.Msg
	serr.Type = BillerError
	serr.GRPCCode = stscode
	switch v := billErr.Error.(type) {
	case string:
		if v != "" {
			serr.Msg = v
		}
	case map[string]interface{}:
		msg, ok := v["message"].(string)
		if ok {
			serr.Msg = msg
		}
	}
	for _, d := range billErr.Details {
		if d.Msg != "" {
			serr.Msg = d.Msg
		}
	}
	switch v := billErr.Code.(type) {
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

	if serr.Msg == "" {
		serr.Msg = stscode.String()
		serr.UnknownErr = strings.Join(strings.Fields(string(b)), "")
	}
	if billErr.Errors != nil {
		serr.Errors = billErr.Errors
	}
	return serr
}

func ConvertErr(err error, errType ErrType) error {
	if err == nil {
		return nil
	}
	er := status.Convert(err)
	cd := er.Code()
	if errType == NonexConError {
		cd = codes.OutOfRange
	}
	return status.Error(cd, er.Message())
}
