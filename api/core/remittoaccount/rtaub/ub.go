package rtaub

import (
	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	rtai "brank.as/petnet/api/integration/remittoaccount"
	"brank.as/petnet/api/storage/postgres"
)

func (s *Svc) Kind() string {
	return static.UBRTA
}

type Svc struct {
	remitAcc *rtai.Client
	st       *postgres.Storage
}

func New(store *postgres.Storage, RemitAcc *rtai.Client) *Svc {
	return &Svc{
		remitAcc: RemitAcc,
		st:       store,
	}
}

func handleUBError(err error) *coreerror.Error {
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
