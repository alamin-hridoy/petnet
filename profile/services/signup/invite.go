package signup

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ipb "brank.as/petnet/gunk/dsa/v1/user"
)

func (h *Svc) RetrieveInvite(ctx context.Context, req *ipb.RetrieveInviteRequest) (*ipb.RetrieveInviteResponse, error) {
	err := validation.ValidateStruct(req,
		validation.Field(&req.Code, validation.Required),
	)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	usr := &ipb.RetrieveInviteResponse{}
	if req.Code != "" {
		usr, err = h.core.RetrieveInvite(ctx, req.Code)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return usr, nil
}
