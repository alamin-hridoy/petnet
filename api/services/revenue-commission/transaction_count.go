package revenue_commission

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	revcom "brank.as/petnet/gunk/drp/v1/revenue-commission"
)

// CreateDSATransactionCount Creates DSA Transaction Count record.
func (s *Svc) CreateDSATransactionCount(context.Context, *revcom.CreateDSATransactionCountRequest) (*revcom.DSATransactionCount, error) {
	return nil, errUnImplemented
}

// UpdateDSATransactionCount Updates DSA Transaction Count record.
func (s *Svc) UpdateDSATransactionCount(context.Context, *revcom.UpdateDSATransactionCountRequest) (*revcom.DSATransactionCount, error) {
	return nil, errUnImplemented
}

// GetDSATransactionCountByID Gets DSA Transaction Count record by transaction count ID.
func (s *Svc) GetDSATransactionCountByID(context.Context, *revcom.GetDSATransactionCountByIDRequest) (*revcom.DSATransactionCount, error) {
	return nil, errUnImplemented
}

// DeleteDSATransactionCount Deletes DSA Transaction Count record by transaction count ID.
func (s *Svc) DeleteDSATransactionCount(context.Context, *revcom.DeleteDSATransactionCountRequest) (*revcom.DSATransactionCount, error) {
	return nil, errUnImplemented
}

// ListDSATransactionCountAll List all DSA Transaction Count records.
func (s *Svc) ListDSATransactionCountAll(context.Context, *emptypb.Empty) (*revcom.ListDSATransactionCountResponse, error) {
	return nil, errUnImplemented
}

// ListDSATransactionCountByYearMonth List DSA Transaction Count records by year and month.
func (s *Svc) ListDSATransactionCountByYearMonth(context.Context, *revcom.ListDSATransactionCountByYearMonthRequest) (*revcom.ListDSATransactionCountResponse, error) {
	return nil, errUnImplemented
}
