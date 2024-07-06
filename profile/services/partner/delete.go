package partner

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/dsa/v2/partner"
)

func (s *Svc) DeletePartner(ctx context.Context, req *ppb.DeletePartnerRequest) (*ppb.DeletePartnerResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.core.DeletePartner(ctx, req.GetID()); err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to get service record")
	}
	return &ppb.DeletePartnerResponse{}, nil
}
