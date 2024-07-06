package bpi

import (
	"context"

	"google.golang.org/grpc/codes"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
)

func (s *Svc) Kind() string {
	return static.BPICode
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

func (s *Svc) StageCreateRemit(ctx context.Context, r core.Remittance) (*core.RemitResponse, error) {
	return nil, coreerror.NewCoreError(codes.Unavailable, "service not available for bpi")
}

func handleBPIError(err error) *coreerror.Error {
	switch err.(type) {
	case *perahub.Error:
		pErr, ok := err.(*perahub.Error)
		if !ok || pErr == nil {
			return coreerror.ToCoreError(err)
		}

		// status 400, code "99" - Control Number not found
		if pErr.GRPCCode == codes.InvalidArgument && pErr.Code == "99" {
			return coreerror.NewCoreError(codes.NotFound, coreerror.MsgControlNumberNotFound)
		}

		return coreerror.ToCoreError(err)

	default:
		return coreerror.ToCoreError(err)
	}
}
