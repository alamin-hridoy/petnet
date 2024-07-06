package partner

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/util"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) InputGuide(ctx context.Context, req *ppb.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetRemitPartner()
	if !static.PartnerExists(pn, "PH") {
		log.Error("partner doesn't exist")
		return nil, util.HandleServiceErr(coreerror.NewCoreError(codes.NotFound, "partner doesn't exist"))
	}

	v, ok := s.validators[pn]
	if !ok {
		return nil, util.HandleServiceErr(coreerror.NewCoreError(codes.NotFound, "no input guide for partner"))
	}

	r, err := v.InputGuideValidate(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}
	r.Ptnr = req.GetRemitPartner()

	ig, err := s.partner.InputGuide(ctx, *r)
	if err != nil {
		switch t := err.(type) {
		case *perahub.Error:
			if t.Type == perahub.PartnerError {
				return nil, perahub.GRPCError(t.GRPCCode, "partner error", &ppb.Error{
					Code:    t.Code,
					Message: t.Msg,
				})
			}
			return nil, util.HandleServiceErr(coreerror.NewCoreError(codes.Internal, "internal error occurred"))
		}
		if status.Code(err) != codes.Unknown {
			return nil, util.HandleServiceErr(err)
		}
		return nil, util.HandleServiceErr(coreerror.NewCoreError(codes.Internal, "internal error occurred"))
	}
	return ig, nil
}

func (*WUVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.SourceCountry, validation.Required, is.CountryCode2),
		validation.Field(&req.SourceCurrency, validation.Required, is.CurrencyCode),
	); err != nil {
		return nil, util.HandleServiceErr(coreerror.NewCoreError(codes.InvalidArgument, err.Error()))
	}
	return &core.InputGuideRequest{
		SrcCtry: req.SourceCountry,
		SrcCncy: req.SourceCurrency,
	}, nil
}

func (*TFVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*RMVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*WISEVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.CountryCode, validation.Required, is.CountryCode2),
	); err != nil {
		return nil, util.HandleServiceErr(coreerror.NewCoreError(codes.InvalidArgument, err.Error()))
	}
	return &core.InputGuideRequest{
		CtryCode: req.CountryCode,
	}, nil
}

func (*UNTVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*CEBVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.AgentCode, validation.Required),
	); err != nil {
		return nil, util.HandleServiceErr(coreerror.NewCoreError(codes.InvalidArgument, err.Error()))
	}
	return &core.InputGuideRequest{
		AgentCode: req.AgentCode,
	}, nil
}

func (*USSCVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*IRVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*RIAVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*MBVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*BPIVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*ICVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*JPRVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*AYAVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*CEBINTVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*IEVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	return &core.InputGuideRequest{}, nil
}

func (*PerahubRemitVal) InputGuideValidate(ctx context.Context, req *ppb.InputGuideRequest) (*core.InputGuideRequest, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, validation.Required),
		validation.Field(&req.ID, validation.Required),
		validation.Field(&req.City, validation.Required),
	); err != nil {
		return nil, util.HandleServiceErr(coreerror.NewCoreError(codes.InvalidArgument, err.Error()))
	}
	return &core.InputGuideRequest{
		Ptnr: req.GetRemitPartner(),
		ID:   int(req.GetID()),
		City: req.GetCity(),
	}, nil
}
