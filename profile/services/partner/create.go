package partner

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/dsa/v2/partner"
)

func (s *Svc) CreatePartners(ctx context.Context, req *ppb.CreatePartnersRequest) (*ppb.CreatePartnersResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Partners, validation.Required, validation.By(func(interface{}) error {
			svc := req.GetPartners()
			if err := validation.ValidateStruct(svc,
				validation.Field(&svc.OrgID, validation.Required, is.UUIDv4),
			); err != nil {
				return status.Error(codes.InvalidArgument, err.Error())
			}
			return nil
		})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	svc := req.GetPartners()
	if err := s.core.CreatePartners(ctx, svc); err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to store partner record")
	}
	return &ppb.CreatePartnersResponse{Partners: svc}, nil
}
