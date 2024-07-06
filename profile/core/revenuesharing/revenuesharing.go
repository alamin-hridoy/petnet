package revenuesharing

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
)

type RevenueSharingStore interface {
	CreateRevenueSharing(ctx context.Context, req storage.RevenueSharing) (*storage.RevenueSharing, error)
	UpdateRevenueSharing(ctx context.Context, req storage.RevenueSharing) (*storage.RevenueSharing, error)
	CreateRevenueSharingTier(ctx context.Context, req storage.RevenueSharingTier) (*storage.RevenueSharingTier, error)
	UpdateRevenueSharingTier(ctx context.Context, req storage.RevenueSharingTier) (*storage.RevenueSharingTier, error)
	GetRevenueSharingList(ctx context.Context, req storage.RevenueSharing) ([]storage.RevenueSharing, error)
	GetRevenueSharingTierList(ctx context.Context, req storage.RevenueSharingTier) ([]storage.RevenueSharingTier, error)
	DeleteRevenueSharing(ctx context.Context, req storage.RevenueSharing) error
	DeleteRevenueSharingTier(ctx context.Context, req storage.RevenueSharingTier) error
	DeleteRevenueSharingTierById(ctx context.Context, req storage.RevenueSharingTier) error
	GetPartnerTransactionType(ctx context.Context, req *storage.GetAllPartnerListReq) ([]storage.GetAllPartnerListReq, error)
}

type Svc struct {
	st RevenueSharingStore
}

func New(st *postgres.Storage) *Svc {
	return &Svc{st: st}
}
