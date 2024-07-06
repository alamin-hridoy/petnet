package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) GetPartnerTransactionType(ctx context.Context, req *rc.GetPartnerTransactionTypeRequest) (*rc.GetPartnerTransactionTypeResponse, error) {
	res, err := s.st.GetPartnerTransactionType(ctx, &storage.GetAllPartnerListReq{
		Partner: req.GetPartners(),
	})
	if err != nil {
		return nil, err
	}
	var rcList []*rc.PartnerDetail
	ptnrTransactions := make(map[string][]string)
	for _, v := range res {
		ptnrTransactions[v.Partner] = append(ptnrTransactions[v.Partner], v.TransactionType)
	}
	for _, v := range res {
		rcList = append(rcList, &rc.PartnerDetail{
			Stype:            v.Partner,
			Name:             v.Partner,
			TransactionTypes: ptnrTransactions[v.Partner],
		})
	}
	return &rc.GetPartnerTransactionTypeResponse{
		PartnerDetails: rcList,
	}, nil
}
