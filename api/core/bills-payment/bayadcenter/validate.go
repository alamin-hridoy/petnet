package bayadcenter

import (
	"context"
	"strconv"

	bpi "brank.as/petnet/api/integration/bills-payment"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Validate(ctx context.Context, req *bp.BPValidateRequest) (res *bp.BPValidateResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.billAcc.BCValidate(ctx, bpi.BCValidateRequest{
		BillPartnerID: int(req.GetBillPartnerID()),
		BillerTag:     req.GetBillerTag(),
		Code:          req.GetCode(),
		AccountNumber: req.GetAccountNumber(),
		AccountNo:     req.GetAccountNo(),
		Identifier:    req.GetIdentifier(),
		PaymentMethod: req.GetPaymentMethod(),
		OtherCharges:  req.GetOtherCharges(),
		Amount:        strconv.Itoa(int(req.GetAmount())),
		OtherInfo: bpi.BCOtherInfo{
			LastName:        req.OtherInfo.GetLastName(),
			FirstName:       req.OtherInfo.GetFirstName(),
			MiddleName:      req.OtherInfo.GetMiddleName(),
			PaymentType:     req.OtherInfo.GetPaymentType(),
			Course:          req.OtherInfo.GetCourse(),
			TotalAssessment: req.OtherInfo.GetTotalAssessment(),
			SchoolYear:      req.OtherInfo.GetSchoolYear(),
			Term:            req.OtherInfo.GetTerm(),
			Name:            req.OtherInfo.GetName(),
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("Bills payment bayadcenter validate failed.")
		return nil, handleBayadcenterError(err)
	}

	res = &bp.BPValidateResponse{
		Code:    int32(rs.Code),
		Message: rs.Message,
		Result: &bp.BPValidateResult{
			Valid:            rs.Result.Valid,
			Code:             int32(rs.Result.Code),
			Account:          rs.Result.Account,
			Details:          []*bp.Details{},
			ValidationNumber: rs.Result.ValidationNumber,
		},
		RemcoID: int32(rs.RemcoID),
	}

	return res, nil
}
