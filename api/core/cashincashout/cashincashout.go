package cashincashout

import (
	"strings"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
	"google.golang.org/grpc/codes"
)

type Svc struct {
	ph *perahub.Svc
	st *postgres.Storage
}

func New(prh *perahub.Svc, st *postgres.Storage) *Svc {
	s := &Svc{
		ph: prh,
		st: st,
	}
	return s
}

func handleCiCoError(err error) *coreerror.Error {
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

			if strings.Contains(msg, "NOT FOUND") {
				return coreerror.NewCoreError(codes.NotFound, pErr.Msg)
			}

			return coreerror.NewCoreError(codes.Internal, coreerror.MsgPerahubInternalError)
		}

		return coreerror.ToCoreError(err)

	default:
		return coreerror.ToCoreError(err)
	}
}
