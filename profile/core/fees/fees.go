package fees

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
)

type FeeStore interface {
	CreateOrgFees(ctx context.Context, fee storage.FeeCommission) (*storage.FeeCommission, error)
	UpsertOrgFees(ctx context.Context, fee storage.FeeCommission) (*storage.FeeCommission, error)
	ListOrgFees(ctx context.Context, oid string, f storage.LimitOffsetFilter) ([]storage.FeeCommission, error)
	UpsertRate(ctx context.Context, fee storage.Rate) (*storage.Rate, error)
	CreateFeeCommissionRate(ctx context.Context, rate storage.Rate) (*storage.Rate, error)
	ListFeesCommissionRate(ctx context.Context, fcid string) ([]storage.Rate, error)
}

type Svc struct {
	st FeeStore
}

func New(st *postgres.Storage) *Svc {
	return &Svc{st: st}
}
