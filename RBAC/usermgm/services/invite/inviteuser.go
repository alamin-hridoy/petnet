package invite

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"

	ipb "brank.as/rbac/gunk/v1/invite"
)

// InviteUser creates a new user.
func (h *Svc) InviteUser(ctx context.Context, req *ipb.InviteUserRequest) (*ipb.InviteUserResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.invite.InviteUser")
	log.Trace("request received")

	uid, orgID := hydra.ClientID(ctx), hydra.OrgID(ctx)

	if err := validation.ValidateStruct(req,
		validation.Field(&req.FirstName, validation.Required, validation.Length(1, 70)),
		validation.Field(&req.LastName, validation.Required, validation.Length(1, 70)),
		validation.Field(&req.Email, is.EmailFormat, validation.Required, validation.Length(3, 0)),
		validation.Field(&req.Role, is.UUIDv4),
		validation.Field(&req.OrgID, is.UUIDv4),
		validation.Field(&req.OrgName,
			// Empty org id signals that a new org must be created.
			// Can only happen with an authenticated admin account.
			// TODO(Chad): Validate permission to invite new orgs when
			// autoOrg/autoApprove aren't enabled.
			validation.When(req.GetOrgID() == "" && hydra.OrgID(ctx) == "",
				validation.Required)),
	); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	inv, err := h.inv.InviteUser(ctx, core.Invite{
		InvOrgID:        orgID,
		InvUserID:       uid,
		OrgID:           req.OrgID,
		OrgName:         req.OrgName,
		FName:           req.FirstName,
		LName:           req.LastName,
		RoleID:          req.Role,
		Email:           req.Email,
		CustomEmailData: req.GetCustomEmailData(),
	})
	if err != nil {
		logging.WithError(err, log).Error("inviting")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to send invite")
	}

	return &ipb.InviteUserResponse{
		ID:             inv.ID,
		InvitationCode: inv.Code,
	}, nil
}
