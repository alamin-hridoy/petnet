package userremit

import (
	"context"
	"time"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/remit"
	"brank.as/petnet/api/storage/postgres"
)

type Svc struct {
	rmt *remit.Svc
	st  *postgres.Storage
}

func New(rmt *remit.Svc, st *postgres.Storage) *Svc {
	return &Svc{rmt: rmt, st: st}
}

func toDate(dt string) (*core.Date, error) {
	tm, err := time.Parse("01/02/2006", dt)
	if err != nil {
		return nil, err
	}
	return &core.Date{
		Year:  tm.Format("2006"),
		Month: tm.Format("01"),
		Day:   tm.Format("02"),
	}, nil
}

func (s *Svc) ProcessRemit(ctx context.Context, r core.ProcessRemit, partner string) (*core.ProcessRemit, error) {
	return s.rmt.ProcessRemit(ctx, r, partner)
}

func (s *Svc) SearchRemit(ctx context.Context, r core.SearchRemit, partner string) (*core.SearchRemit, error) {
	return s.rmt.SearchRemit(ctx, r, partner)
}

func (s *Svc) StageDisburseRemit(ctx context.Context, r core.Remittance, partner string) (*core.Remittance, error) {
	return s.rmt.StageDisburseRemit(ctx, r, partner)
}

func (s *Svc) ListRemit(ctx context.Context, r core.FilterList) (*core.SearchRemitResponse, error) {
	return s.rmt.ListRemit(ctx, r)
}
