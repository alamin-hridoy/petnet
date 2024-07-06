package ayannah

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
	return static.AYACode
}

type remitStore interface {
	CreateRemitCache(context.Context, storage.RemitCache) (*storage.RemitCache, error)
	GetRemitCache(context.Context, string) (*storage.RemitCache, error)
	UpdateRemitCache(context.Context, storage.RemitCache) (*storage.RemitCache, error)
	CreateRemitHistory(context.Context, storage.RemitHistory) (*storage.RemitHistory, error)
	GetRemitHistory(context.Context, string) (*storage.RemitHistory, error)
	ListRemitHistory(context.Context, storage.LRHFilter) ([]storage.RemitHistory, error)
	UpdateRemitHistory(context.Context, storage.RemitHistory) (*storage.RemitHistory, error)
	OrderIDExists(context.Context, string) bool
}

type Svc struct {
	st  remitStore
	ph  *perahub.Svc
	stc *static.Svc
}

func New(st remitStore, ph *perahub.Svc, stc *static.Svc) *Svc {
	return &Svc{
		st:  st,
		ph:  ph,
		stc: stc,
	}
}

func handleAyannahError(err error) *coreerror.Error {
	switch err.(type) {
	case *perahub.Error:
		pErr, ok := err.(*perahub.Error)
		if !ok || pErr == nil {
			return coreerror.ToCoreError(err)
		}

		// status 400, code "90" - Payout Conflict
		if pErr.GRPCCode == codes.InvalidArgument {
			if pErr.Code == "90" || strings.Contains(pErr.Msg, "Transaction Paid Out") {
				return coreerror.NewCoreError(codes.AlreadyExists, coreerror.MsgTransactionAlreadyClaimed)
			}
		}

		return coreerror.ToCoreError(err)

	default:
		return coreerror.ToCoreError(err)
	}
}
