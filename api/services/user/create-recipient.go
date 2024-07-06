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

func (s *Svc) CreateRecipient(ctx context.Context, req *ppb.CreateRecipientRequest) (*ppb.CreateRecipientResponse, error) {
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

	if err := vp.CreateRecipientValidate(ctx, req); err != nil {
		return nil, err
	}

	res, err := s.user.CreateRecipient(ctx, req)
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

func (*WISEVal) CreateRecipientValidate(ctx context.Context, req *ppb.CreateRecipientRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Currency, validation.Required, is.CurrencyCode),
		validation.Field(&req.Type, validation.Required),
		validation.Field(&req.AccountHolderName, validation.Required, is.ASCII),
		validation.Field(&req.Requirements, validation.Each(
			validation.By(func(r interface{}) error {
				i, _ := r.(ppb.Requirement)
				return validation.ValidateStruct(&i,
					validation.Field(&i.Name, validation.Required, is.ASCII),
					validation.Field(&i.Value, validation.Required.When(i.Values == nil), is.ASCII),
					validation.Field(&i.Values, validation.Required.When(i.Value == "")),
				)
			}),
		)),
	); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func (*CEBVal) CreateRecipientValidate(ctx context.Context, req *ppb.CreateRecipientRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.FirstName, validation.Required),
		validation.Field(&req.MiddleName, validation.Required),
		validation.Field(&req.LastName, validation.Required),
		validation.Field(&req.SenderUserID, validation.Required),
		validation.Field(&req.BirthDate, validation.Required),
		validation.Field(&req.MobileCountryID, validation.Required),
		validation.Field(&req.ContactNumber, validation.Required),
		validation.Field(&req.PhoneCountryID, validation.Required),
		validation.Field(&req.PhoneAreaCode, validation.Required),
		validation.Field(&req.PhoneNumber, validation.Required),
		validation.Field(&req.CountryAddressID, validation.Required),
		validation.Field(&req.BirthCountryID, validation.Required),
		validation.Field(&req.ProvinceAddress, validation.Required),
		validation.Field(&req.Address, validation.Required),
		validation.Field(&req.UserID, validation.Required),
		validation.Field(&req.Occupation, validation.Required),
		validation.Field(&req.PostalCode, validation.Required),
		validation.Field(&req.StateIDAddress, validation.Required),
		validation.Field(&req.Tin, validation.Required),
	); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}
