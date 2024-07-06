package intelexpress

import (
	"context"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
)

func (s *Svc) Kind() string {
	return static.IECode
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

func handleIntelExpressError(err error) *coreerror.Error {
	switch err.(type) {
	case *perahub.Error:
		pErr, ok := err.(*perahub.Error)
		if !ok || pErr == nil {
			return coreerror.ToCoreError(err)
		}

		// TODO(vitthal): handle other error codes

		return coreerror.ToCoreError(err)

	default:
		return coreerror.ToCoreError(err)
	}
}
