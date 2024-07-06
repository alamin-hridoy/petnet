package bayadcenter

import (
	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	bpi "brank.as/petnet/api/integration/bills-payment"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
)

func (s *Svc) Kind() string {
	return static.BYCBP
}

type Svc struct {
	billAcc *bpi.Client
	st      *postgres.Storage
}

func New(store *postgres.Storage, BillAcc *bpi.Client) *Svc {
	return &Svc{
		billAcc: BillAcc,
		st:      store,
	}
}

func handleBayadcenterError(err error) *coreerror.Error {
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
