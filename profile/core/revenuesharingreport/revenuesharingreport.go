package revenuesharingreport

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
)

type RevenueSharingReportStore interface {
	CreateRevenueSharingReport(ctx context.Context, req storage.RevenueSharingReport) (*storage.RevenueSharingReport, error)
	GetRevenueSharingReportList(ctx context.Context, req storage.RevenueSharingReport) ([]storage.RevenueSharingReport, error)
	UpdateRevenueSharingReport(ctx context.Context, req storage.RevenueSharingReport) (*storage.RevenueSharingReport, error)
}

type Svc struct {
	st RevenueSharingReportStore
}

func New(st *postgres.Storage) *Svc {
	return &Svc{st: st}
}
