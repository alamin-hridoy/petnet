package wu

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
)

func (s *Svc) Kind() string {
	return static.WUCode
}

type remitStore interface {
	CreateRemitCache(context.Context, storage.RemitCache) (*storage.RemitCache, error)
	CreateRemitHistory(context.Context, storage.RemitHistory) (*storage.RemitHistory, error)
	GetRemitCache(context.Context, string) (*storage.RemitCache, error)
	GetRemitHistory(context.Context, string) (*storage.RemitHistory, error)
	ListRemitHistory(context.Context, storage.LRHFilter) ([]storage.RemitHistory, error)
	UpdateRemitCache(context.Context, storage.RemitCache) (*storage.RemitCache, error)
	UpdateRemitHistory(context.Context, storage.RemitHistory) (*storage.RemitHistory, error)
	OrderIDExists(context.Context, string) bool
}

type Svc struct {
	st          remitStore
	ph          *perahub.Svc
	stc         *static.Svc
	perahubMock bool
}

func New(st remitStore, ph *perahub.Svc, stc *static.Svc, phm bool) *Svc {
	return &Svc{
		st:          st,
		ph:          ph,
		stc:         stc,
		perahubMock: phm,
	}
}

func handleWUError(err error) *coreerror.Error {
	switch err.(type) {
	case *coreerror.Error:
		return err.(*coreerror.Error)

	case *perahub.Error:
		pErr, ok := err.(*perahub.Error)
		if !ok || pErr == nil {
			return coreerror.ToCoreError(err)
		}

		// TODO(vitthal): handle other error codes
		if strings.Contains(pErr.Msg, "NO RECORD FOUND") || strings.Contains(pErr.Msg, "NO TRANSACTION FOUND") {
			return coreerror.NewCoreError(codes.NotFound, "not found")
		}

		if strings.Contains(pErr.Msg, "SQL statement") {
			return coreerror.NewCoreError(codes.Internal, coreerror.MsgPerahubDatabaseError)
		}

		if strings.Contains(pErr.Msg, "is not a valid instance") {
			msg := getStringInBetween(pErr.Msg, "element ", " instance of type")
			if msg == "" {
				msg = pErr.Msg
			}

			return coreerror.NewCoreError(codes.InvalidArgument, msg)
		}

		if strings.Contains(pErr.Msg, "USPError in parameter") {
			msg := getStringInBetween(pErr.Msg, "[\"", "\"]")
			if msg == "" {
				msg = pErr.Msg
			}

			return coreerror.NewCoreError(codes.InvalidArgument, msg)
		}

		code := getGrpcCodeFromWUCode(pErr.Code)
		if code == codes.Internal {
			return coreerror.NewCoreError(codes.Internal, coreerror.MsgPerahubInternalError)
		}

		if code == codes.Unavailable {
			return coreerror.NewCoreError(codes.Unavailable, coreerror.MsgConnectionError)
		}

		return coreerror.NewCoreError(code, pErr.Msg)

	default:
		return coreerror.ToCoreError(err)
	}
}

func getGrpcCodeFromWUCode(wuCode string) codes.Code {
	switch wuCode {
	case "T5803", "T6034", "T6082", "T6081", "T6006", "T5705", "Missing":
		return codes.InvalidArgument
	case "U9035":
		return codes.AlreadyExists
	case "T0851":
		return codes.DeadlineExceeded
	case "502", "Failed":
		return codes.Unavailable
	default:
		return codes.Internal
	}
}

func getStringInBetween(str string, start string, end string) (result string) {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return str
	}

	s := strings.Index(str, start)
	if s == -1 {
		s = 0
	} else {
		s += len(start)
	}

	e := strings.Index(str, end)
	if e == -1 {
		e = len(str)
	}

	return str[s:e]
}
