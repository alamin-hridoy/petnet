package transfast

import (
	"context"

	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
)

func (s *Svc) Kind() string {
	return static.TFCode
}

type inputGuideStore interface {
	GetInputGuide(ctx context.Context, partner string) (*storage.InputGuide, error)
	CreateInputGuide(ctx context.Context, ig storage.InputGuide) (*storage.InputGuide, error)
}

type Svc struct {
	st inputGuideStore
	ph *perahub.Svc
}

func New(st inputGuideStore, ph *perahub.Svc) *Svc {
	return &Svc{
		st: st,
		ph: ph,
	}
}
