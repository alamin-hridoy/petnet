package remittance

import (
	"strings"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
)

type Svc struct {
	ph *perahub.Svc
	st *postgres.Storage
}

func New(st *postgres.Storage, prh *perahub.Svc) *Svc {
	s := &Svc{
		ph: prh,
		st: st,
	}
	return s
}

func handlePerahubError(err error) *coreerror.Error {
	switch err.(type) {
	case *perahub.Error:
		pErr, ok := err.(*perahub.Error)
		if !ok || pErr == nil {
			return coreerror.ToCoreError(err)
		}

		if pErr.GRPCCode == codes.Internal {
			msg := strings.ToUpper(pErr.Msg)
			if strings.Contains(msg, "SQLSTATE") {
				return coreerror.NewCoreError(codes.Internal, coreerror.MsgPerahubDatabaseError)
			}

			if strings.Contains(msg, "NOT FOUND") ||
				strings.Contains(msg, "INVALID SEND VALIDATION REFERENCE NUMBER") ||
				strings.Contains(msg, "INVALID PERAHUB REFERENCE NUMBER") ||
				strings.Contains(msg, "INVALID PAYOUT VALIDATION REFERENCE NUMBER") {
				return coreerror.NewCoreError(codes.NotFound, pErr.Msg)
			}

			return coreerror.NewCoreError(codes.Internal, coreerror.MsgPerahubInternalError)
		}

		return coreerror.ToCoreError(err)

	default:
		return coreerror.ToCoreError(err)
	}
}
