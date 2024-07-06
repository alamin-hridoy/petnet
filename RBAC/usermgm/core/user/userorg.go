package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/storage"
)

//  UserOrg returns the platform associated with the authenticated user.
func (h *Svc) UserOrg(ctx context.Context) (*storage.Organization, error) {
	log := logging.FromContext(ctx).WithField("method", "user.UserOrg")
	uID := hydra.ClientID(ctx)
	log.WithField("user id", uID).Info("org lookup")
	usr, err := h.usr.GetUserByID(ctx, uID)
	if err != nil {
		return nil, err
	}
	if usr.OrgID == "" {
		return &storage.Organization{
			ID:      "test-org",
			OrgName: "Test Org",
			Active:  true,
		}, nil
	}
	org, err := h.org.GetOrgByID(ctx, usr.OrgID)
	if err != nil {
		return nil, err
	}
	if !org.Active {
		return nil, status.Error(codes.FailedPrecondition, "platform is not active")
	}
	return org, nil
}
