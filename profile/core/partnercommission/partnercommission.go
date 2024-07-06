package partnercommission

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
)

type PartnerCommissionStore interface {
	CreatePartnerCommission(ctx context.Context, req storage.PartnerCommission) (*storage.PartnerCommission, error)
	UpdatePartnerCommission(ctx context.Context, req storage.PartnerCommission) (*storage.PartnerCommission, error)
	CreatePartnerCommissionTier(ctx context.Context, req storage.PartnerCommissionTier) (*storage.PartnerCommissionTier, error)
	UpdatePartnerCommissionTier(ctx context.Context, req storage.PartnerCommissionTier) (*storage.PartnerCommissionTier, error)
	GetPartnerCommissionsList(ctx context.Context, req storage.PartnerCommission) ([]storage.PartnerCommission, error)
	GetPartnerCommissionsTierList(ctx context.Context, req storage.PartnerCommissionTier) ([]storage.PartnerCommissionTier, error)
	DeletePartnerCommission(ctx context.Context, req storage.PartnerCommission) error
	DeletePartnerCommissionTier(ctx context.Context, req storage.PartnerCommissionTier) error
	DeletePartnerCommissionTierById(ctx context.Context, req storage.PartnerCommissionTier) error
}

type Svc struct {
	st PartnerCommissionStore
}

func New(st *postgres.Storage) *Svc {
	return &Svc{st: st}
}
