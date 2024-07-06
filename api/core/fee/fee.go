package fee

import (
	"context"
	"log"

	"brank.as/petnet/api/core"
	usscf "brank.as/petnet/api/core/fee/ussc"
	wuf "brank.as/petnet/api/core/fee/wu"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
)

type Fee interface {
	FeeInquiry(ctx context.Context, r core.FeeInquiryReq) (map[string]string, error)
	Kind() string
}

type Svc struct {
	fee map[string]Fee
	st  *postgres.Storage
	ph  *perahub.Svc
}

func New(st *postgres.Storage, ph *perahub.Svc) *Svc {
	fs := []Fee{wuf.New(ph), usscf.New(ph)}
	s := &Svc{
		fee: make(map[string]Fee, len(fs)),
		ph:  ph,
	}
	for i, r := range fs {
		switch {
		case r == nil:
			log.Fatalf("fee %d nil", i)
		case r.Kind() == "":
			log.Fatalf("fee %d missing partner type", i)
		}
		s.fee[r.Kind()] = r
	}
	return s
}
