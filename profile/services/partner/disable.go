package partner

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/dsa/v2/partner"
)

func (s *Svc) DisablePartner(ctx context.Context, req *ppb.DisablePartnerRequest) (*ppb.DisablePartnerResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, validation.Required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.core.DisablePartner(ctx, req.OrgID, req.Partner.String()); err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to disable partner")
	}
	return &ppb.DisablePartnerResponse{}, nil
}
