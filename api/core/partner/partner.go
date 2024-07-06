package partner

import (
	"context"
	"log"

	"brank.as/petnet/api/core"
	aya "brank.as/petnet/api/core/partner/ayannah"
	bpi "brank.as/petnet/api/core/partner/bpi"
	cebint "brank.as/petnet/api/core/partner/cebint"
	ceb "brank.as/petnet/api/core/partner/cebuana"
	ic "brank.as/petnet/api/core/partner/instacash"
	ie "brank.as/petnet/api/core/partner/intelexpress"
	ir "brank.as/petnet/api/core/partner/iremit"
	jpr "brank.as/petnet/api/core/partner/japanremit"
	mb "brank.as/petnet/api/core/partner/metrobank"
	pr "brank.as/petnet/api/core/partner/perahubremit"
	rmg "brank.as/petnet/api/core/partner/remitly"
	ria "brank.as/petnet/api/core/partner/ria"
	tfg "brank.as/petnet/api/core/partner/transfast"
	unt "brank.as/petnet/api/core/partner/uniteller"
	ussc "brank.as/petnet/api/core/partner/ussc"
	wise "brank.as/petnet/api/core/partner/wise"
	wug "brank.as/petnet/api/core/partner/wu"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
)

const (
	terminalID = "PH259AMT001A" // petnet will provide this operator id
	operatorID = "drp"          // petnet will provide this terminal id
)

type Guider interface {
	InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error)
	Kind() string
}

type Svc struct {
	guiders map[string]Guider
	st      *postgres.Storage
	ph      *perahub.Svc
}

func New(st *postgres.Storage, ph *perahub.Svc) *Svc {
	gs := []Guider{wug.New(st, ph), rmg.New(st, ph), tfg.New(st, ph), wise.New(st, ph), unt.New(st, ph), ceb.New(st, ph), ussc.New(), ir.New(), ria.New(), mb.New(), bpi.New(), ic.New(), jpr.New(), aya.New(), cebint.New(), ie.New(), pr.New(st, ph)}
	s := &Svc{
		guiders: make(map[string]Guider, len(gs)),
		st:      st,
		ph:      ph,
	}
	for i, r := range gs {
		switch {
		case r == nil:
			log.Fatalf("guider %d nil", i)
		case r.Kind() == "":
			log.Fatalf("guider %d missing partner type", i)
		}
		s.guiders[r.Kind()] = r
	}
	return s
}
