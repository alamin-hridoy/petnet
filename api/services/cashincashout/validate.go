package cashincashout

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	cico "brank.as/petnet/gunk/drp/v1/cashincashout"
	"brank.as/petnet/serviceutil/logging"
)

const (
	DRAGONPAY string = "DRAGONPAY"
	PERAHUB   string = "PERAHUB"
)

func (s *Svc) CiCoValidate(ctx context.Context, req *cico.CiCoValidateRequest) (*cico.CiCoValidateResponse, error) {
	if err := s.CiCoValidateValidation(ctx, req); err != nil {
		return nil, util.HandleServiceErr(err)
	}
	res, err := s.core.CiCoValidate(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}
	return res, err
}

func (s *Svc) CiCoValidateValidation(ctx context.Context, req *cico.CiCoValidateRequest) error {
	log := logging.FromContext(ctx)
	required := validation.Required
	if req == nil {
		return status.Error(codes.InvalidArgument, "please provide required fields")
	}
	if req.Trx == nil {
		return status.Error(codes.InvalidArgument, "trx is required")
	}
	if req.Customer == nil {
		return status.Error(codes.InvalidArgument, "customer is required")
	}
	if err := validation.ValidateStruct(req,
		validation.Field(&req.PartnerCode, required),
		validation.Field(&req.Trx, validation.By(func(interface{}) error {
			t := req.Trx
			return validation.ValidateStruct(t,
				validation.Field(&t.Provider, required),
				validation.Field(&t.ReferenceNumber, required),
				validation.Field(&t.TrxType, required),
				validation.Field(&t.PrincipalAmount, validation.When(req.Trx.Provider != DRAGONPAY && req.Trx.Provider != PERAHUB, required)),
			)
		})),
		validation.Field(&req.Customer, required, validation.By(func(interface{}) error {
			c := req.Customer
			return validation.ValidateStruct(c,
				validation.Field(&c.CustomerID, required),
				validation.Field(&c.CustomerFirstname, required),
				validation.Field(&c.CustomerLastname, required),
				validation.Field(&c.CurrAddress, required),
				validation.Field(&c.CurrBarangay, required),
				validation.Field(&c.CurrCity, required),
				validation.Field(&c.CurrProvince, required),
				validation.Field(&c.CurrCountry, required),
				validation.Field(&c.BirthDate, required),
				validation.Field(&c.BirthPlace, required),
				validation.Field(&c.BirthCountry, required),
				validation.Field(&c.ContactNo, required),
				validation.Field(&c.IDType, required),
				validation.Field(&c.IDNumber, required),
			)
		})),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}
