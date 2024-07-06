package ria

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
	return static.RIACode
}

type remitStore interface {
	CreateRemitCache(context.Context, storage.RemitCache) (*storage.RemitCache, error)
	CreateRemitHistory(context.Context, storage.RemitHistory) (*storage.RemitHistory, error)
	GetRemitCache(context.Context, string) (*storage.RemitCache, error)
	GetRemitHistory(context.Context, string) (*storage.RemitHistory, error)
	ListRemitHistory(context.Context, storage.LRHFilter) ([]storage.RemitHistory, error)
	UpdateRemitCache(context.Context, storage.RemitCache) (*storage.RemitCache, error)
	UpdateRemitHistory(context.Context, storage.RemitHistory) (*storage.RemitHistory, error)
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
	return nil, coreerror.NewCoreError(codes.Unavailable, "service not available for Ria")
}

func handleRiaError(err error) *coreerror.Error {
	switch err.(type) {
	case *perahub.Error:
		pErr, ok := err.(*perahub.Error)
		if !ok || pErr == nil {
			return coreerror.ToCoreError(err)
		}

		msg := pErr.Msg
		if pErr.UnknownErr != "" {
			msg = pErr.UnknownErr
		}

		if pErr.GRPCCode == codes.InvalidArgument {
			if pErr.Code == "2001" {
				return coreerror.NewCoreError(codes.NotFound, coreerror.MsgControlNumberNotFound)
			}

			if pErr.Code == "2002" {
				return coreerror.NewCoreError(codes.AlreadyExists, coreerror.MsgTransactionAlreadyClaimed)
			}

			if pErr.Code == "2005" {
				// Order must be successfully verified using the OP_VerifyOrderForPayout method prior to confirming payment.
				return coreerror.NewCoreError(codes.FailedPrecondition, coreerror.MsgNeedToConfirm)
			}
		}

		return coreerror.NewCoreError(pErr.GRPCCode, msg)

	default:
		return coreerror.ToCoreError(err)
	}
}
