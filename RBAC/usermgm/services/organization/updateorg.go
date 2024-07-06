package organization

import (
	"context"
	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	opb "brank.as/rbac/gunk/v1/organization"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/storage"
)

func (h *Svc) UpdateOrganization(ctx context.Context, req *opb.UpdateOrganizationRequest) (*opb.UpdateOrganizationResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.org.updateorg")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrganizationID, validation.Required, is.UUIDv4),
		validation.Field(&req.Email, is.EmailFormat),
		validation.Field(&req.Name, validation.Length(3, 256)),
		validation.Field(&req.Phone, is.Int, validation.Length(10, 20)),
		validation.Field(&req.LoginMFA, validation.By(func(interface{}) error {
			return validation.Validate(
				opb.EnableOpt_name[opb.EnableOpt_value[req.LoginMFA.String()]],
				validation.Required)
		})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	org, err := h.org.UpdateOrg(ctx, storage.Organization{
		ID:           req.OrganizationID,
		OrgName:      req.Name,
		ContactEmail: req.Email,
		ContactPhone: req.Phone,
		MFALogin: sql.NullBool{
			Bool:  req.LoginMFA != opb.EnableOpt_Disable,
			Valid: req.LoginMFA != opb.EnableOpt_NoChange,
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("update org")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "falied to apply updates")
	}

	return &opb.UpdateOrganizationResponse{Updated: tspb.New(org.Updated)}, nil
}
