package user

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	pnpb "brank.as/petnet/gunk/drp/v1/partner"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) CreateProfile(ctx context.Context, req *ppb.CreateProfileRequest) (*ppb.CreateProfileResponse, error) {
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

	creq, err := vp.CreateProfileValidate(ctx, req)
	if err != nil {
		return nil, err
	}

	if _, err := s.user.CreateProfile(ctx, *creq); err != nil {
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
	return &ppb.CreateProfileResponse{}, nil
}

func (*WISEVal) CreateProfileValidate(ctx context.Context, req *ppb.CreateProfileRequest) (*core.CreateProfileReq, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Type, validation.Required),
		validation.Field(&req.FirstName, validation.Required, is.ASCII),
		validation.Field(&req.LastName, validation.Required, is.ASCII),
		validation.Field(&req.BirthDate, validation.Required, validation.By(func(interface{}) error {
			if _, err := time.Parse("2006-01-02", req.BirthDate); err != nil {
				return fmt.Errorf("wrong format for birth_date, want: yyyy-mm-dd")
			}
			return nil
		})),
		validation.Field(&req.Phone, validation.Required, validation.By(func(interface{}) error {
			r := req.Phone
			return validation.ValidateStruct(r,
				validation.Field(&r.CountryCode, validation.Required, is.Digit),
				validation.Field(&r.Number, validation.Required, is.Digit),
			)
		})),
		validation.Field(&req.Address, validation.Required, validation.By(func(interface{}) error {
			r := req.Address
			return validation.ValidateStruct(r,
				validation.Field(&r.Address1, validation.Required),
				validation.Field(&r.City, validation.Required),
				validation.Field(&r.Country, validation.Required),
				validation.Field(&r.PostalCode, validation.Required, is.Alphanumeric),
			)
		})),
		validation.Field(&req.Occupation, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &core.CreateProfileReq{
		Email:     req.Email,
		Type:      req.Type,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
		Phone: core.PhoneNumber{
			CtyCode: req.Phone.CountryCode,
			Number:  req.Phone.Number,
		},
		Address: core.Address{
			Address1:   req.Address.Address1,
			City:       req.Address.City,
			Country:    req.Address.Country,
			PostalCode: req.Address.PostalCode,
		},
		Occupation: req.Occupation,
	}, nil
}

func (*CEBVal) CreateProfileValidate(ctx context.Context, req *ppb.CreateProfileRequest) (*core.CreateProfileReq, error) {
	return nil, fmt.Errorf("service not available for Cebuana")
}
