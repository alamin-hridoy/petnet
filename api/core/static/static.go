package static

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage/postgres"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	WUCode       = "WU"
	IRCode       = "IR"
	TFCode       = "TF"
	RMCode       = "RM"
	RIACode      = "RIA"
	MBCode       = "MB"
	BPICode      = "BPI"
	USSCCode     = "USSC"
	WISECode     = "WISE"
	ICCode       = "IC"
	JPRCode      = "JPR"
	UNTCode      = "UNT"
	CEBCode      = "CEB"
	AYACode      = "AYA"
	CEBINTCode   = "CEBINT"
	IECode       = "IE"
	MBRTA        = "MB"
	UBRTA        = "UB"
	BPIRTA       = "BPI"
	PerahubRemit = "PerahubRemit"
	ECPBP        = "ECP"
	BYCBP        = "BYC"
	MLPBP        = "MLP"
)

type WUTxType = string

const (
	WUSendMoney      WUTxType = "SO"
	WUDirectBank     WUTxType = "D2B"
	WUMobileTransfer WUTxType = "MMT"
	WUQuickPay       WUTxType = "QP"
)

// if partner doesn't have a name for the codes we can use these general codes
const (
	Payout  = "PO"
	Sendout = "SO"
)

var Partners = map[string][]core.Remco{
	"PH": {
		{
			Code: WUCode,
			Name: "Western Union",
			SendTypes: map[string]core.SendRemitType{
				"Send": {
					Code:        WUSendMoney,
					Description: "Send money to an individual for pick up at any Western Union location.",
					Receiver:    true,
				},
				"Direct": {
					Code:        WUDirectBank,
					Description: "Send money directly to a bank account.",
					Receiver:    true,
					BankAccount: true,
				},
				"Mobile": {
					Code:        WUMobileTransfer,
					Description: "Send money to a mobile number.",
					Receiver:    true,
				},
				"QuickPay": {
					Code:        WUQuickPay,
					Description: "Send money to a buiness that accepts Western Union QuickPay.",
					Business:    true,
				},
			},
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: IRCode,
			Name: "iRemit",
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: TFCode,
			Name: "Transfast",
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: RMCode,
			Name: "Remitly",
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: RIACode,
			Name: "Ria",
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: MBCode,
			Name: "Metrobank",
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: BPICode,
			Name: "BPI",
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: USSCCode,
			Name: "USSC",
			SendTypes: map[string]core.SendRemitType{
				"Send": {
					Code:        Sendout,
					Description: "Send money transaction.",
					Receiver:    true,
				},
			},
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: WISECode,
			Name: "TransferWise",
			SendTypes: map[string]core.SendRemitType{
				"Send": {
					Code:        Sendout,
					Description: "Send money transaction.",
					Receiver:    true,
				},
			},
		},
		{
			Code: ICCode,
			Name: "InstaCash",
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: JPRCode,
			Name: "JapanRemit",
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: UNTCode,
			Name: "Uniteller",
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: CEBINTCode,
			Name: "Cebuana intl",
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: CEBCode,
			Name: "Cebuana",
			SendTypes: map[string]core.SendRemitType{
				"Send": {
					Code:        Sendout,
					Description: "Send money transaction.",
					Receiver:    true,
				},
			},
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: AYACode,
			Name: "Ayannah",
			SendTypes: map[string]core.SendRemitType{
				"Send": {
					Code:        Sendout,
					Description: "Send money transaction.",
					Receiver:    true,
				},
			},
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: IECode,
			Name: "IntelExpress",
			SendTypes: map[string]core.SendRemitType{
				"Send": {
					Code:        Sendout,
					Description: "Send money transaction.",
					Receiver:    true,
				},
			},
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
		{
			Code: PerahubRemit,
			Name: "PerahubRemit",
			SendTypes: map[string]core.SendRemitType{
				"Send": {
					Code:        Sendout,
					Description: "Send money transaction.",
					Receiver:    true,
				},
			},
			DisburseTypes: map[string]core.DisburseRemitType{
				"Payout": {
					Code:        Payout,
					Description: "Payout a transaction.",
				},
			},
		},
	},
}

type CountrySource interface {
	CurrencyCodes(ctx context.Context, cty, cur string) (map[string]string, error)
}

type Svc struct {
	all []core.Remco
	cty CountrySource
	st  *postgres.Storage
}

func New(cty CountrySource, st *postgres.Storage) *Svc {
	a := []core.Remco{}
	for k := range Partners {
		a = append(a, Partners[k]...)
	}
	return &Svc{all: a, cty: cty, st: st}
}

func (s *Svc) ListPartners(ctx context.Context, cty string) ([]core.Remco, error) {
	if cty == "" {
		return s.all, nil
	}
	r := Partners[cty]
	if len(r) == 0 {
		return nil, status.Errorf(codes.NotFound, "no partners available in %q", cty)
	}
	return r, nil
}

// PartnerCode validates that the remittance partner code exists.
func (s *Svc) PartnerCode(ctx context.Context, partner string) error {
	for _, p := range s.all {
		if p.Code != partner {
			continue
		}
		return nil
	}
	return status.Error(codes.InvalidArgument, "invalid remittance partner")
}

func (s *Svc) SendRemitType(ctx context.Context, partner, remTyp string, code bool) (*core.SendRemitType, error) {
	for _, p := range s.all {
		if p.Code != partner {
			continue
		}
		for k, t := range p.SendTypes {
			if code {
				if t.Code != remTyp {
					continue
				}
				return &t, nil
			} else {
				if k != remTyp {
					continue
				}
				return &t, nil
			}
		}
	}
	return nil, status.Error(codes.InvalidArgument, "invalid remittance type")
}

func (s *Svc) DisburseRemitType(ctx context.Context, partner, remTyp string, code bool) (*core.DisburseRemitType, error) {
	for _, p := range s.all {
		if p.Code != partner {
			continue
		}
		for k, t := range p.DisburseTypes {
			if code {
				if t.Code != remTyp {
					continue
				}
				return &t, nil
			} else {
				if k != remTyp {
					continue
				}
				return &t, nil
			}
		}
	}
	return nil, status.Error(codes.InvalidArgument, "invalid remittance type")
}

func PartnerExists(partner string, country string) bool {
	pns := Partners[country]
	for _, pn := range pns {
		if pn.Code == partner {
			return true
		}
	}
	return false
}
