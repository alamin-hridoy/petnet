package bills_payment

import (
	"context"
	"encoding/json"
	"time"

	"brank.as/petnet/api/storage"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) BillsPaymentTransactList(ctx context.Context, r *bp.BillsPaymentTransactListRequest) (*bp.BillsPaymentTransactListResponse, error) {
	dateFrom, err := time.Parse("2006-01-02", r.From)
	if err != nil {
		dateFrom = time.Time{}
	}
	dateUntil, err := time.Parse("2006-01-02", r.Until)
	if err != nil {
		dateUntil = time.Now()
	}

	res, err := s.st.ListBillPayment(ctx, storage.BillPaymentFilter{
		BillPaymentStatus: string(storage.SuccessStatus),
		SortByColumn:      storage.BillPayColumn(r.GetSortByColumn().String()),
		SortOrder:         storage.SortOrder(r.GetSortOrder().String()),
		Limit:             int(r.GetLimit()),
		Offset:            int(r.GetOffset()),
		From:              dateFrom,
		Until:             dateUntil,
		OrgID:             r.GetOrgID(),
		ReferenceNumber:   r.GetReferenceNumber(),
		ExcludePartners:   r.GetExcludePartners(),
	})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, storage.ErrNotFound
	}
	bills := []*bp.BillsPayment{}
	for _, v := range res {
		ttlAmnt, svChrg := &bp.Amount{}, &bp.Amount{}
		json.Unmarshal(v.TotalAmount, ttlAmnt)
		json.Unmarshal(v.ServiceCharge, svChrg)
		bills = append(bills, &bp.BillsPayment{
			ReferenceNumber: v.ReferenceNumber,
			Partner:         v.PartnerID,
			TotalAmount:     ttlAmnt,
			TransactFee:     svChrg,
			TransactCommission: &bp.Amount{
				Amount:   "0",
				Currency: "PHP",
			},
			TransactionCompletedTime: timestamppb.New(v.TrxDate),
		})
	}

	return &bp.BillsPaymentTransactListResponse{
		Next:          r.Offset + int32(len(bills)),
		BillsPayments: bills,
		Total:         int32(res[0].Total),
	}, nil
}
