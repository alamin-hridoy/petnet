package fee

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/serviceutil/logging"

	fpb "brank.as/petnet/gunk/drp/v1/fee"
	pnpb "brank.as/petnet/gunk/drp/v1/partner"
)

func (s *Svc) FeeInquiry(ctx context.Context, req *fpb.FeeInquiryRequest) (*fpb.FeeInquiryResponse, error) {
	log := logging.FromContext(ctx)

	pn := req.GetRemitPartner()
	if !static.PartnerExists(pn, "PH") {
		log.Error("partner doesn't exist")
		return nil, status.Error(codes.NotFound, "partner doesn't exist")
	}

	r, err := s.validators[pn].FeeInquiryValidate(ctx, s.remit, req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	fs, err := s.fee.FeeInquiry(ctx, *r)
	if err != nil {
		switch t := err.(type) {
		case *perahub.Error:
			if t.Type == perahub.PartnerError {
				return nil, perahub.GRPCError(t.GRPCCode, "partner error", &pnpb.Error{
					Code:    t.Code,
					Message: t.Msg,
				})
			}
			return nil, status.Errorf(codes.Internal, "internal error occurred")
		}
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "internal error occurred")
	}
	return &fpb.FeeInquiryResponse{
		Fees: fs,
	}, nil
}

func (s *WUVal) FeeInquiryValidate(ctx context.Context, st RemcoStore, req *fpb.FeeInquiryRequest) (*core.FeeInquiryReq, error) {
	pn := req.GetRemitPartner()
	remType, err := st.SendRemitType(ctx, pn, req.GetRemitType(), false)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid remittance type %q for partner %q",
			req.GetRemitType(), req.GetRemitPartner(),
		)
	}

	if err := validation.ValidateStruct(req.Amount,
		validation.Field(&req.Amount.Amount, validation.Required, is.Int),
		validation.Field(&req.Amount.DestinationCurrency, validation.Required, is.CurrencyCode),
		validation.Field(&req.Amount.DestinationCountry, validation.Required, is.CountryCode2),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res := &core.FeeInquiryReq{
		RemitPartner:      req.RemitPartner,
		RemitType:         *remType,
		PrincipalAmount:   core.MustMinor(req.Amount.Amount, req.Amount.DestinationCurrency),
		DestinationAmount: req.Amount.DestinationAmount,
		DestCountry:       req.Amount.DestinationCountry,
		DestCurrency:      req.Amount.DestinationCurrency,
		Promo:             req.Promo,
		Message:           req.Message,
	}
	return res, nil
}

func (s *USSCVal) FeeInquiryValidate(ctx context.Context, st RemcoStore, req *fpb.FeeInquiryRequest) (*core.FeeInquiryReq, error) {
	if err := validation.ValidateStruct(req.Amount,
		validation.Field(&req.Amount.Amount, validation.Required, is.Int),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res := &core.FeeInquiryReq{
		RemitPartner:    req.RemitPartner,
		PrincipalAmount: core.MustMinor(req.Amount.Amount, "PHP"),
	}
	return res, nil
}
