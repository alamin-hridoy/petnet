package remit

import (
	"context"
	"fmt"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/serviceutil/logging"
)

type remitStore interface {
	GetRemitCache(context.Context, string) (*storage.RemitCache, error)
	ListRemitHistory(context.Context, storage.LRHFilter) ([]storage.RemitHistory, error)
}

type Remitter interface {
	StageCreateRemit(context.Context, core.Remittance) (*core.RemitResponse, error)
	StageDisburseRemit(context.Context, core.Remittance) (*core.Remittance, error)
	ProcessRemit(context.Context, core.ProcessRemit) (*core.ProcessRemit, error)
	Search(context.Context, core.SearchRemit) (*core.SearchRemit, error)
	Kind() string
}

type Svc struct {
	remitters map[string]Remitter
	st        remitStore
	ph        *perahub.Svc
	rc        *static.Svc
}

func New(store remitStore, wu *perahub.Svc, rc *static.Svc, rs []Remitter) (*Svc, error) {
	s := &Svc{
		remitters: make(map[string]Remitter, len(rs)),
		st:        store,
		ph:        wu,
		rc:        rc,
	}
	for i, r := range rs {
		switch {
		case r == nil:
			return nil, fmt.Errorf("remitter %d nil", i)
		case r.Kind() == "":
			return nil, fmt.Errorf("remitter %d missing partner type", i)
		}
		s.remitters[r.Kind()] = r
	}
	return s, nil
}

func (s *Svc) GetPartnerByTxnID(ctx context.Context, txnID string) (string, error) {
	log := logging.FromContext(ctx)
	rm, err := s.st.GetRemitCache(ctx, txnID)
	if err != nil || rm.RemcoID == "" {
		logging.WithError(err, log).Error("partner not found")
		return "", err
	}
	return rm.RemcoID, nil
}
