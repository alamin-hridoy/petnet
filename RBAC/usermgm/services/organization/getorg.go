package organization

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	opb "brank.as/rbac/gunk/v1/organization"
)

func (s *Svc) GetOrganization(ctx context.Context, req *opb.GetOrganizationRequest) (*opb.GetOrganizationResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.org.getorg")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	o, err := s.org.GetOrg(ctx, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("fetch org record")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to read record")
	}

	return &opb.GetOrganizationResponse{
		Organization: []*opb.Organization{
			{
				ID:       o.ID,
				Name:     o.OrgName,
				Email:    o.ContactEmail,
				Phone:    o.ContactPhone,
				Active:   o.Active,
				LoginMFA: o.MFALogin.Valid && o.MFALogin.Bool,
				Created:  tspb.New(o.Created),
				Updated:  tspb.New(o.Updated),
			},
		},
	}, nil
}
