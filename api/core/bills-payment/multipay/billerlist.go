package multipay

import (
	"context"
	"strconv"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) BillerList(ctx context.Context, req *bp.BPBillerListRequest) (res *bp.BPBillerListResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.billAcc.BillsPaymentMultiPayBillerlist(ctx)
	if err != nil {
		logging.WithError(err, log).Error("Bills payment multipay biller list failed.")
		return nil, handleMultipayError(err)
	}

	var bl []*bp.BPBillerListResult
	for _, v := range rs.Result {
		blr := &bp.BPBillerListResult{
			BillerTag:     v.BillerTag,
			Description:   v.Description,
			ServiceCharge: int32(v.ServiceCharge),
			Category:      strconv.Itoa(v.Category),
			PartnerID:     int32(v.PartnerID),
		}
		if len(v.FieldList) > 0 {
			fl := []*bp.BPBillerListFieldList{}
			for _, flr := range v.FieldList {
				flre := &bp.BPBillerListFieldList{
					ID:          flr.ID,
					Type:        flr.Type,
					Label:       flr.Label,
					Order:       int32(flr.Order),
					Description: flr.Description,
					Placeholder: flr.Placeholder,
				}
				if len(flr.Rules) > 0 {
					rules := []*bp.BPBillerListRules{}
					for _, flrre := range flr.Rules {
						rules = append(rules, &bp.BPBillerListRules{
							Code:    int32(flrre.Code),
							Type:    flrre.Type,
							Value:   flrre.Value,
							Format:  flrre.Format,
							Message: flrre.Message,
							Options: flrre.Options,
						})
					}
					flre.Rules = rules
				}
			}
			blr.FieldList = fl
		}
		bl = append(bl, blr)
	}
	res = &bp.BPBillerListResponse{
		Code:    int32(rs.Code),
		Message: rs.Message,
		Result:  bl,
	}

	return res, nil
}
