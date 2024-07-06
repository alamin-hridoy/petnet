package quote

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	pnpb "brank.as/petnet/gunk/drp/v1/partner"
	qpb "brank.as/petnet/gunk/drp/v1/quote"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) QuoteRequirements(ctx context.Context, req *qpb.QuoteRequirementsRequest) (*qpb.QuoteRequirementsResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetRemitPartner()

	if !static.PartnerExists(pn, "PH") {
		log.Error("partner doesn't exist")
		return nil, status.Error(codes.NotFound, "partner doesn't exist")
	}

	if err := s.validators[pn].QuoteRequirementsValidate(ctx, req); err != nil {
		return nil, err
	}

	res, err := s.quote.QuoteRequirements(ctx, req)
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
	return res, nil
}

func (*WISEVal) QuoteRequirementsValidate(ctx context.Context, req *qpb.QuoteRequirementsRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Amount, validation.Required, validation.By(func(interface{}) error {
			r := req.Amount
			return validation.ValidateStruct(r,
				validation.Field(&r.SourceAmount, validation.Required, is.Digit),
				validation.Field(&r.SourceCurrency, validation.Required, is.CurrencyCode),
				validation.Field(&r.DestinationCurrency, validation.Required, is.CurrencyCode),
			)
		})),
	); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}
