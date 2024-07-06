package bayadcenter

import (
	"context"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) BillerList(ctx context.Context, req *bp.BPBillerListRequest) (res *bp.BPBillerListResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.billAcc.BCBillerList(ctx)
	if err != nil {
		logging.WithError(err, log).Error("Bills payment bayadcenter biller list failed.")
		return nil, handleBayadcenterError(err)
	}

	var bl []*bp.BPBillerListResult
	for _, v := range rs.Result {
		bl = append(bl, &bp.BPBillerListResult{
			Description:     v.Description,
			Name:            v.Name,
			Code:            v.Code,
			Category:        v.Category,
			Type:            v.Type,
			Logo:            v.Logo,
			IsMultipleBills: int32(v.IsMultipleBills),
			IsCde:           int32(v.IsCde),
			IsAsync:         int32(v.IsAsync),
		})
	}
	res = &bp.BPBillerListResponse{
		Code:    int32(rs.Code),
		Message: rs.Message,
		Result:  bl,
	}

	return res, nil
}
