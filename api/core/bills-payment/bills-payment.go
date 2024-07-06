package bills_payment

import (
	"context"
	"fmt"

	bpi "brank.as/petnet/api/integration/bills-payment"
	"brank.as/petnet/api/storage"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
)

type Biller interface {
	Transact(ctx context.Context, r *bp.BPTransactRequest) (*bp.BPTransactResponse, error)
	Validate(ctx context.Context, r *bp.BPValidateRequest) (*bp.BPValidateResponse, error)
	TransactInquire(ctx context.Context, r *bp.BPTransactInquireRequest) (*bp.BPTransactInquireResponse, error)
	BillerList(ctx context.Context, r *bp.BPBillerListRequest) (*bp.BPBillerListResponse, error)
	Kind() string
}

type BillerStore interface {
	ListBillPayment(ctx context.Context, f storage.BillPaymentFilter) ([]storage.BillPayment, error)
}

type Svc struct {
	billers map[string]Biller
	billAcc *bpi.Client
	st      BillerStore
}

func New(rs []Biller, BillAcc *bpi.Client, st BillerStore) (*Svc, error) {
	s := &Svc{
		billers: make(map[string]Biller, len(rs)),
		billAcc: BillAcc,
		st:      st,
	}
	for i, r := range rs {
		switch {
		case r == nil:
			return nil, fmt.Errorf("bills payment %d nil", i)
		case r.Kind() == "":
			return nil, fmt.Errorf("bills payment %d missing partner", i)
		}
		s.billers[r.Kind()] = r
	}
	return s, nil
}
