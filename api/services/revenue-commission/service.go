package revenue_commission

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	revcom_int "brank.as/petnet/api/integration/revenue-commission"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/gunk/drp/v1/dsa"
	revcom "brank.as/petnet/gunk/drp/v1/revenue-commission"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	rsr "brank.as/petnet/gunk/dsa/v2/revenuesharingreport"
)

// iDSAStore contract for DSA store
type iDSAStore interface {
	GetDSAByID(ctx context.Context, dsaID uint32) (*revcom_int.DSA, error)
	ListDSA(ctx context.Context) ([]revcom_int.DSA, error)
	CreateDSA(ctx context.Context, req *revcom_int.SaveDSARequest) (*revcom_int.DSA, error)
	UpdateDSA(ctx context.Context, req *revcom_int.SaveDSARequest) (*revcom_int.DSA, error)
	DeleteDSA(ctx context.Context, dsaID uint32) (*revcom_int.DSA, error)
}

// iCommissionFeeStore contract for Remco CommissionFee store
type iCommissionFeeStore interface {
	ListRemcoCommissionFee(ctx context.Context) ([]revcom_int.RemcoCommissionFee, error)
	GetRemcoCommissionFeeByID(ctx context.Context, feeID uint32) (*revcom_int.RemcoCommissionFee, error)
	CreateRemcoCommissionFee(ctx context.Context, request *revcom_int.SaveRemcoCommissionFeeRequest) (*revcom_int.RemcoCommissionFee, error)
	UpdateRemcoCommissionFee(ctx context.Context, request *revcom_int.SaveRemcoCommissionFeeRequest) (*revcom_int.RemcoCommissionFee, error)
	DeleteRemcoCommissionFee(ctx context.Context, feeID uint32) (*revcom_int.RemcoCommissionFee, error)
	CreateTransactionCount(ctx context.Context, req *revcom_int.SaveTransactionCountRequest) (*revcom_int.TransactionCount, error)
}

// iDSACommissionStore contract for DSA Commission Store
type iDSACommissionStore interface {
	GetDSACommissionByID(ctx context.Context, commissionID uint32) (*revcom_int.DSACommission, error)
	ListDSACommission(ctx context.Context) ([]revcom_int.DSACommission, error)
	CreateDSACommission(ctx context.Context, req *revcom_int.SaveDSACommissionRequest) (*revcom_int.DSACommission, error)
	UpdateDSACommission(ctx context.Context, req *revcom_int.SaveDSACommissionRequest) (*revcom_int.DSACommission, error)
	DeleteDSACommission(ctx context.Context, commissionID uint32) (*revcom_int.DSACommission, error)

	GetCommissionTierByID(ctx context.Context, tierID uint32) (*revcom_int.CommissionTier, error)
	ListCommissionTier(ctx context.Context) ([]revcom_int.CommissionTier, error)
	CreateCommissionTier(ctx context.Context, req *revcom_int.SaveCommissionTierRequest) (*revcom_int.CommissionTier, error)
	UpdateCommissionTier(ctx context.Context, req *revcom_int.SaveCommissionTierRequest) (*revcom_int.CommissionTier, error)
	DeleteCommissionTier(ctx context.Context, tierID uint32) (*revcom_int.CommissionTier, error)
}

// iStore contract for DSA store
type iStore interface {
	GetTransactionReport(ctx context.Context, pf *storage.LRHFilter) (*storage.RemitHistory, error)
}

// iRevenueReportClient contract for saving revenue sharing report
type iRevenueReportClient interface {
	// Create Revenue Sharing Report
	CreateRevenueSharingReport(ctx context.Context, in *rsr.CreateRevenueSharingReportRequest, opts ...grpc.CallOption) (*rsr.CreateRevenueSharingReportResponse, error)
}

// Svc ...
// TODO(vitthal): Should move service and integration under profile?
type Svc struct {
	dsaStore           iDSAStore
	store              iStore
	commissionFeeStore iCommissionFeeStore
	revenueReport      iRevenueReportClient
	dsaCommissionStore iDSACommissionStore
	revcom.UnimplementedRevenueCommissionServiceServer
	dsa.UnimplementedDSAServiceServer
	profileService ppb.OrgProfileServiceClient
}

// Option is type for creating service Svc with options
type Option func(service *Svc)

// NewRevenueCommissionService ...
func NewRevenueCommissionService(opts ...Option) *Svc {
	s := &Svc{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithDSAStore ...
func WithDSAStore(dsaStore iDSAStore) Option {
	return func(s *Svc) {
		s.dsaStore = dsaStore
	}
}

// WithCommissionFeeStore ...
func WithCommissionFeeStore(commissionFeeStore iCommissionFeeStore) Option {
	return func(s *Svc) {
		s.commissionFeeStore = commissionFeeStore
	}
}

// WithStorage ...
func WithStorage(store *postgres.Storage) Option {
	return func(s *Svc) {
		s.store = store
	}
}

// WithRevenueSharingReport ...
func WithRevenueSharingReport(revenueReport rsr.RevenueSharingReportServiceClient) Option {
	return func(s *Svc) {
		s.revenueReport = revenueReport
	}
}

// With ProfileDSA ...
func WithProfileDSA(profileClient ppb.OrgProfileServiceClient) Option {
	return func(s *Svc) {
		s.profileService = profileClient
	}
}

// WithDSACommissionStore ...
func WithDSACommissionStore(dsaCommStore iDSACommissionStore) Option {
	return func(s *Svc) {
		s.dsaCommissionStore = dsaCommStore
	}
}

// RegisterSvc register the remit service.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	revcom.RegisterRevenueCommissionServiceServer(srv, s)
	dsa.RegisterDSAServiceServer(srv, s)
	return nil
}

// RegisterGateway ...
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	err = revcom.RegisterRevenueCommissionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return err
	}

	return dsa.RegisterDSAServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
