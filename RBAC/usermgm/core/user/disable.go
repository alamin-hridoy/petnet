package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/svcutil/mw"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/integrations/email"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) DisableUser(ctx context.Context, req core.UserActivation) error {
	log := logging.FromContext(ctx).WithField("method", "core.user.DisableUser")
	u, err := s.usr.GetUserByID(ctx, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("fetch user")
		if err == storage.NotFound {
			return status.Error(codes.NotFound, "user not found")
		}
		return status.Error(codes.Internal, "failed to read user record")
	}
	// todo: we should probably have a way for a superuser to be able to
	// disable users in other organizations, adding this for safety for now
	// until there is need for something else
	if mw.GetOrg(ctx) != u.OrgID {
		return status.Error(codes.PermissionDenied, "failed disabling user")
	}
	if err := s.usr.DisableUser(ctx, req.ID); err != nil {
		logging.WithError(err, log).Error("Failed to disable User.")
		if err == storage.NotFound {
			return status.Error(codes.NotFound, "user source not found")
		}
		return err
	}

	if !s.notifyDisable {
		return nil
	}

	o, err := s.org.GetOrgByID(ctx, u.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("fetch org")
		if err == storage.NotFound {
			return status.Error(codes.NotFound, "org not found")
		}
		return status.Error(codes.Internal, "failed to read org record")
	}

	dis := email.User{
		CustomEmailData: req.CustomEmailData,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		OrgName:         o.OrgName,
	}
	if err = s.mail.DisableUser(u.Email, dis); err != nil {
		const errMsg = "failed to send disable user email"
		logging.WithError(err, log).Error(errMsg)
		if !s.devEnv {
			return status.Error(codes.Internal, errMsg)
		}
	}

	return nil
}
