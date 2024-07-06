package remittoaccount

import (
	"context"
	"fmt"

	rtai "brank.as/petnet/api/integration/remittoaccount"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
)

type Remitter interface {
	Inquire(ctx context.Context, r *rta.RTAInquireRequest) (*rta.RTAInquireResponse, error)
	Payment(ctx context.Context, r *rta.RTAPaymentRequest) (*rta.RTAPaymentResponse, error)
	Retry(ctx context.Context, r *rta.RTARetryRequest) (*rta.RTARetryResponse, error)
	Kind() string
}

type Svc struct {
	remitters map[string]Remitter
	remitAcc  *rtai.Client
}

func New(rs []Remitter, RemitAcc *rtai.Client) (*Svc, error) {
	s := &Svc{
		remitters: make(map[string]Remitter, len(rs)),
		remitAcc:  RemitAcc,
	}
	for i, r := range rs {
		switch {
		case r == nil:
			return nil, fmt.Errorf("remit to account %d nil", i)
		case r.Kind() == "":
			return nil, fmt.Errorf("remit to account %d missing partner", i)
		}
		s.remitters[r.Kind()] = r
	}
	return s, nil
}
