package cashincashout

import (
	"context"
	"strconv"
	"time"

	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	cbp "brank.as/petnet/gunk/drp/v1/cashincashout"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ListRemit historical transactions.
func (s *Svc) CICOTransactList(ctx context.Context, f *cbp.CICOTransactListRequest) (*cbp.CICOTransactListResponse, error) {
	log := logging.FromContext(ctx)
	var err error
	dateFrom, err := time.Parse("2006-01-02", f.GetFrom())
	if err != nil {
		dateFrom = time.Time{}
	}
	dateUntil, err := time.Parse("2006-01-02", f.GetUntil())
	if err != nil {
		dateUntil = time.Now()
	}

	fl := storage.CashInCashOutTrxListFilter{
		From:             dateFrom,
		Until:            dateUntil,
		Limit:            int(f.GetLimit()),
		Offset:           int(f.GetOffset()),
		SortOrder:        storage.SortOrder(f.GetSortOrder().String()),
		SortByColumn:     storage.CICOHistoryColumn(f.GetSortByColumn().String()),
		ReferenceNumber:  f.GetReferenceNumber(),
		ExcludeProviders: f.GetExcludeProviders(),
		OrgID:            f.GetOrgID(),
	}

	rhs, err := s.st.ListCICOTrx(ctx, fl)
	if err != nil {
		log.Error(err)
	}
	if len(rhs) == 0 {
		return nil, storage.ErrNotFound
	}

	cicotrx := []*cbp.CICOTransact{}
	for _, v := range rhs {
		ttlAmnt, svChrg := &cbp.Amount{
			Amount:   strconv.Itoa(v.TotalAmount),
			Currency: "PHP",
		}, &cbp.Amount{
			Amount:   strconv.Itoa(v.Charges),
			Currency: "PHP",
		}
		cicotrx = append(cicotrx, &cbp.CICOTransact{
			ReferenceNumber:          v.ReferenceNumber,
			Provider:                 v.Provider,
			TotalAmount:              ttlAmnt,
			TransactFee:              svChrg,
			TransactCommission:       &cbp.Amount{Amount: "0", Currency: "PHP"},
			TransactionCompletedTime: timestamppb.New(v.TrxDate),
		})
	}

	tot := 0
	if len(rhs) > 0 {
		tot = rhs[0].Total
	}

	return &cbp.CICOTransactListResponse{
		Next:          f.Offset + int32(len(cicotrx)),
		CICOTransacts: cicotrx,
		Total:         int32(tot),
	}, nil
}

func getDSAOrgID(ctx context.Context) string {
	ot := phmw.GetOrgType(ctx)
	switch ot {
	// this happens when API is used internally and means that either the dsa
	// or admin platform is used to get transactions
	case ppb.OrgType_PetNet.String(), ppb.OrgType_DSA.String():
		return phmw.GetDSAOrgID(ctx)
	}
	// this happens when API is used externally and means the user authenticated
	// with api client credentials to get token
	return hydra.OrgID(ctx)
}
