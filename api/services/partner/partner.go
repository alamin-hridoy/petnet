package partner

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/serviceutil/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	ppb "brank.as/petnet/gunk/drp/v1/partner"
	pfppb "brank.as/petnet/gunk/dsa/v2/partner"
	pl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	pfsvc "brank.as/petnet/gunk/dsa/v2/service"
)

type (
	WUVal           struct{}
	TFVal           struct{}
	RMVal           struct{}
	WISEVal         struct{}
	UNTVal          struct{}
	CEBVal          struct{}
	USSCVal         struct{}
	IRVal           struct{}
	RIAVal          struct{}
	MBVal           struct{}
	BPIVal          struct{}
	ICVal           struct{}
	JPRVal          struct{}
	AYAVal          struct{}
	CEBINTVal       struct{}
	IEVal           struct{}
	PerahubRemitVal struct{}
)

func (*WUVal) Kind() string {
	return static.WUCode
}

func (*TFVal) Kind() string {
	return static.TFCode
}

func (*RMVal) Kind() string {
	return static.RMCode
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

func (*USSCVal) Kind() string {
	return static.USSCCode
}

func (*IRVal) Kind() string {
	return static.IRCode
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

func (*ICVal) Kind() string {
	return static.ICCode
}

func (*JPRVal) Kind() string {
	return static.JPRCode
}

func (*AYAVal) Kind() string {
	return static.AYACode
}

func (*CEBINTVal) Kind() string {
	return static.CEBINTCode
}

func (*IEVal) Kind() string {
	return static.IECode
}

func (*PerahubRemitVal) Kind() string {
	return static.PerahubRemit
}

type Validator interface {
	InputGuideValidate(context.Context, *ppb.InputGuideRequest) (*core.InputGuideRequest, error)
	Kind() string
}

func NewValidators() []Validator {
	return []Validator{&WUVal{}, &TFVal{}, &RMVal{}, &WISEVal{}, &UNTVal{}, &CEBVal{}, &USSCVal{}, &IRVal{}, &RIAVal{}, &MBVal{}, &BPIVal{}, &ICVal{}, &JPRVal{}, &AYAVal{}, &CEBINTVal{}, &IEVal{}, &PerahubRemitVal{}}
}

type RemcoStore interface {
	ListPartners(context.Context, string) ([]core.Remco, error)
	SendRemitType(ctx context.Context, partner, remTyp string, code bool) (*core.SendRemitType, error)
}

type PartnerStore interface {
	InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error)
	GetPartnersRemco(ctx context.Context) (*ppb.GetPartnersRemcoResponse, error)
}

type Svc struct {
	ppb.UnimplementedRemitPartnerServiceServer
	remit      RemcoStore
	partner    PartnerStore
	cl         pfppb.PartnerServiceClient
	scl        pfsvc.ServiceServiceClient
	plcl       pl.PartnerListServiceClient
	validators map[string]Validator
}

func New(remit RemcoStore, st PartnerStore, cl pfppb.PartnerServiceClient, scl pfsvc.ServiceServiceClient, plcl pl.PartnerListServiceClient, vs []Validator) *Svc {
	s := &Svc{
		validators: make(map[string]Validator, len(vs)),
		remit:      remit,
		partner:    st,
		cl:         cl,
		scl:        scl,
		plcl:       plcl,
	}
	for i, v := range vs {
		switch {
		case v == nil:
			log.Fatalf("validator %d nil", i)
		case v.Kind() == "":
			log.Fatalf("validator %d missing partner type", i)
		}
		s.validators[v.Kind()] = v
	}
	return s
}

// Register the remit service.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	ppb.RegisterRemitPartnerServiceServer(srv, s)
	return nil
}

// RegisterGateway parter endpoints.
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return ppb.RegisterRemitPartnerServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}

// UpdateRemcoId ...
func (s *Svc) UpdateRemcoId(ctx context.Context) error {
	log := logging.FromContext(ctx)
	_, err := s.GetPartnersRemco(ctx, &emptypb.Empty{})
	if err != nil {
		log.WithError(err).Error("failed to get remco list")
		return err
	}
	return nil
}
