package user

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	pnpb "brank.as/petnet/gunk/drp/v1/partner"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) GetRecipients(ctx context.Context, req *ppb.GetRecipientsRequest) (*ppb.GetRecipientsResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetRemitPartner()

	if !static.PartnerExists(pn, "PH") {
		log.Error("partner doesn't exist")
		return nil, status.Error(codes.NotFound, "partner doesn't exist")
	}

	vp, ok := s.validators[pn]
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprint("missing service for ", pn))
	}

	if err := vp.GetRecipientsValidate(ctx, req); err != nil {
		return nil, err
	}

	res, err := s.user.GetRecipients(ctx, req)
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

func (*WISEVal) GetRecipientsValidate(ctx context.Context, req *ppb.GetRecipientsRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.SenderUserEmail, validation.Required, is.Email),
		validation.Field(&req.Currency, validation.Required, is.CurrencyCode),
	); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func (*CEBVal) GetRecipientsValidate(ctx context.Context, req *ppb.GetRecipientsRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.SenderUserID, validation.Required),
	); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}
