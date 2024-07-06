package bills_payment

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) BPValidate(ctx context.Context, req *bp.BPValidateRequest) (*bp.BPValidateResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetPartner()
	if pn == "" {
		return nil, util.HandleServiceErr(status.Error(codes.InvalidArgument, "partner is required."))
	}

	_, err := s.validators[pn].BPValidateValidate(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, util.HandleServiceErr(err)
	}

	res, err := s.bpStore.BPValidate(ctx, req, pn)
	if err != nil {
		logging.WithError(err, log).Error("failed to get inquire")
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *ECPBPVal) BPValidateValidate(ctx context.Context, req *bp.BPValidateRequest) (*bp.BPValidateRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.AccountNo, required),
		validation.Field(&req.Identifier, required),
		validation.Field(&req.BillerTag, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}

func (s *BYCBPVal) BPValidateValidate(ctx context.Context, req *bp.BPValidateRequest) (*bp.BPValidateRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.BillPartnerID, required),
		validation.Field(&req.BillerTag, required),
		validation.Field(&req.Code, required),
		validation.Field(&req.AccountNumber, required),
		validation.Field(&req.AccountNo, required),
		validation.Field(&req.Identifier, required),
		validation.Field(&req.PaymentMethod, required),
		validation.Field(&req.OtherCharges, required),
		validation.Field(&req.Amount, required),
		validation.Field(&req.OtherInfo, required, validation.By(func(interface{}) error {
			OtherInfo := req.OtherInfo
			return validation.ValidateStruct(OtherInfo,
				validation.Field(&OtherInfo.LastName, required),
				validation.Field(&OtherInfo.FirstName, required),
				validation.Field(&OtherInfo.MiddleName, required),
				validation.Field(&OtherInfo.Name, required),
				validation.Field(&OtherInfo.PaymentType, required),
				validation.Field(&OtherInfo.Course, required),
				validation.Field(&OtherInfo.TotalAssessment, required),
				validation.Field(&OtherInfo.SchoolYear, required),
				validation.Field(&OtherInfo.Term, required),
			)
		})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}

func (s *MLPBPVal) BPValidateValidate(ctx context.Context, req *bp.BPValidateRequest) (*bp.BPValidateRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.AccountNumber, required),
		validation.Field(&req.Amount, required),
		validation.Field(&req.ContactNumber, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}
