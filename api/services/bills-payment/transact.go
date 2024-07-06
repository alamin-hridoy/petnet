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

func (s *Svc) BPTransact(ctx context.Context, req *bp.BPTransactRequest) (*bp.BPTransactResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetPartner()
	if pn == "" {
		return nil, util.HandleServiceErr(status.Error(codes.InvalidArgument, "partner is required."))
	}

	_, err := s.validators[pn].BPTransactValidate(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, util.HandleServiceErr(err)
	}

	res, err := s.bpStore.BPTransact(ctx, req, pn)
	if err != nil {
		logging.WithError(err, log).Error("failed to get inquire")
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *ECPBPVal) BPTransactValidate(ctx context.Context, req *bp.BPTransactRequest) (*bp.BPTransactRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.BillID, required),
		validation.Field(&req.BillerTag, required),
		validation.Field(&req.TrxDate, required),
		validation.Field(&req.UserID, required),
		validation.Field(&req.RemoteUserID, required),
		validation.Field(&req.CustomerID, required),
		validation.Field(&req.LocationID, required),
		validation.Field(&req.RemoteLocationID, required),
		validation.Field(&req.LocationName, required),
		validation.Field(&req.Coy, required),
		validation.Field(&req.CurrencyID, required),
		validation.Field(&req.FormType, required),
		validation.Field(&req.FormNumber, required),
		validation.Field(&req.AccountNumber, required),
		validation.Field(&req.Identifier, required),
		validation.Field(&req.Amount, required),
		validation.Field(&req.ServiceCharge, required),
		validation.Field(&req.TotalAmount, required),
		validation.Field(&req.ClientReferenceNumber, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}

func (s *BYCBPVal) BPTransactValidate(ctx context.Context, req *bp.BPTransactRequest) (*bp.BPTransactRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, required),
		validation.Field(&req.CustomerID, required),
		validation.Field(&req.LocationID, required),
		validation.Field(&req.LocationName, required),
		validation.Field(&req.Coy, required),
		validation.Field(&req.CallbackURL, required),
		validation.Field(&req.BillID, required),
		validation.Field(&req.BillerTag, required),
		validation.Field(&req.BillerName, required),
		validation.Field(&req.TrxDate, required),
		validation.Field(&req.Amount, required),
		validation.Field(&req.ServiceCharge, required),
		validation.Field(&req.PartnerCharge, required),
		validation.Field(&req.TotalAmount, required),
		validation.Field(&req.Identifier, required),
		validation.Field(&req.AccountNumber, required),
		validation.Field(&req.PaymentMethod, required),
		validation.Field(&req.ClientReferenceNumber, required),
		validation.Field(&req.ReferenceNumber, required),
		validation.Field(&req.ValidationNumber, required),
		validation.Field(&req.ReceiptValidationNumber, required),
		validation.Field(&req.TpaID, required),
		validation.Field(&req.CurrencyID, required),
		validation.Field(&req.FormType, required),
		validation.Field(&req.FormNumber, required),
		validation.Field(&req.OtherInfo, required, validation.By(func(interface{}) error {
			OtherInfo := req.OtherInfo
			return validation.ValidateStruct(OtherInfo,
				validation.Field(&OtherInfo.LastName, required),
				validation.Field(&OtherInfo.FirstName, required),
				validation.Field(&OtherInfo.MiddleName, required),
				validation.Field(&OtherInfo.PaymentType, required),
				validation.Field(&OtherInfo.Course, required),
				validation.Field(&OtherInfo.TotalAssessment, required),
				validation.Field(&OtherInfo.SchoolYear, required),
				validation.Field(&OtherInfo.Term, required),
			)
		})),
		validation.Field(&req.Type, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}

func (s *MLPBPVal) BPTransactValidate(ctx context.Context, req *bp.BPTransactRequest) (*bp.BPTransactRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Amount, required),
		validation.Field(&req.Txnid, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}
