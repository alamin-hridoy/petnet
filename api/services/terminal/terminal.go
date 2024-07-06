package terminal

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/storage"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pariz/gountries"

	tpb "brank.as/petnet/gunk/drp/v1/terminal"
)

type (
	WUVal struct{}
	IRVal struct {
		q *gountries.Query
	}
	TFVal struct {
		q *gountries.Query
	}
	RMVal struct {
		q *gountries.Query
	}
	RIAVal struct {
		q *gountries.Query
	}
	MBVal struct {
		q *gountries.Query
	}
	BPIVal struct {
		q *gountries.Query
	}
	USSCVal struct {
		q *gountries.Query
	}
	ICVal struct {
		q *gountries.Query
	}
	JPRVal struct {
		q *gountries.Query
	}
	WISEVal struct {
		q *gountries.Query
	}
	UNTVal struct {
		q *gountries.Query
	}
	CEBVal struct {
		q *gountries.Query
	}
	CEBIVal struct {
		q *gountries.Query
	}
	AYAVal struct {
		q *gountries.Query
	}
	IEVal struct {
		q *gountries.Query
	}
	PHUBVal struct {
		q *gountries.Query
	}
)

func (*WUVal) Kind() string {
	return static.WUCode
}

func (*IRVal) Kind() string {
	return static.IRCode
}

func (*TFVal) Kind() string {
	return static.TFCode
}

func (*RMVal) Kind() string {
	return static.RMCode
}

func (*RIAVal) Kind() string {
	return static.RIACode
}

func (*MBVal) Kind() string {
	return static.MBCode
}

func (*BPIVal) Kind() string {
	return static.BPICode
}

func (*USSCVal) Kind() string {
	return static.USSCCode
}

func (*ICVal) Kind() string {
	return static.ICCode
}

func (*JPRVal) Kind() string {
	return static.JPRCode
}

func (*WISEVal) Kind() string {
	return static.WISECode
}

func (*UNTVal) Kind() string {
	return static.UNTCode
}

func (*CEBVal) Kind() string {
	return static.CEBCode
}

func (*CEBIVal) Kind() string {
	return static.CEBINTCode
}

func (*AYAVal) Kind() string {
	return static.AYACode
}

func (*IEVal) Kind() string {
	return static.IECode
}

func (*PHUBVal) Kind() string {
	return static.PerahubRemit
}

type Validator interface {
	CreateRemitValidate(context.Context, *tpb.CreateRemitRequest, *core.SendRemitType) (*core.Remittance, error)
	ConfirmRemitValidate(context.Context, *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error)
	LookupRemitValidate(context.Context, *tpb.LookupRemitRequest) (*core.SearchRemit, error)
	DisburseRemitValidate(context.Context, *tpb.DisburseRemitRequest) (*core.Remittance, error)
	Kind() string
}

func NewValidators(q *gountries.Query) []Validator {
	return []Validator{&WUVal{}, &IRVal{q}, &TFVal{q}, &RMVal{q}, &RIAVal{q}, &MBVal{q}, &BPIVal{q}, &USSCVal{q}, &ICVal{q}, &JPRVal{q}, &WISEVal{q}, &UNTVal{q}, &CEBVal{q}, &CEBIVal{q}, &AYAVal{q}, &IEVal{q}, &PHUBVal{q}}
}

type RemitStore interface {
	StageCreateRemit(context.Context, core.Remittance, string) (*core.RemitResponse, error)
	StageDisburseRemit(context.Context, core.Remittance, string) (*core.Remittance, error)
	ProcessRemit(context.Context, core.ProcessRemit, string) (*core.ProcessRemit, error)
	SearchRemit(context.Context, core.SearchRemit, string) (*core.SearchRemit, error)
	ListRemit(ctx context.Context, f core.FilterList) (*core.SearchRemitResponse, error)
	GetPartnerByTxnID(context.Context, string) (string, error)
}

type RemitLookup interface {
	SendRemitType(ctx context.Context, partner, remTyp string, code bool) (*core.SendRemitType, error)
	DisburseRemitType(ctx context.Context, partner, remTyp string, code bool) (*core.DisburseRemitType, error)
	GetISO(ctx context.Context, c string) (*storage.ISOCty, error)
}

type Svc struct {
	tpb.UnimplementedTerminalServiceServer
	remit      RemitStore
	validators map[string]Validator
	lk         RemitLookup
}

// New Remit service.
func New(remit RemitStore, lk RemitLookup, vs []Validator) (*Svc, error) {
	s := &Svc{
		validators: make(map[string]Validator, len(vs)),
		remit:      remit,
		lk:         lk,
	}
	for i, v := range vs {
		switch {
		case v == nil:
			return nil, fmt.Errorf("validator %d nil", i)
		case v.Kind() == "":
			return nil, fmt.Errorf("validator %d missing partner type", i)
		}
		s.validators[v.Kind()] = v
	}
	return s, nil
}

// Register the remit service.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	tpb.RegisterTerminalServiceServer(srv, s)
	return nil
}

// RegisterGateway ...
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return tpb.RegisterTerminalServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
