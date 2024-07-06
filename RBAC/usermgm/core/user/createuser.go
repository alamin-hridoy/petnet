package user

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/svcutil/mw/hydra"
	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/storage"
)

// CreateUser will add a user account.
// Org will be created if AutoOrg is configured.
func (s *Svc) CreateUser(ctx context.Context, user storage.User, cred storage.Credential) (*storage.User, string, error) {
	log := logging.FromContext(ctx).WithField("method", "user.createuser")

	var usr *storage.User
	var code string
	var err error
	if user.InviteCode != "" {
		usr, err = s.invitedUser(ctx, user, cred)
		if err != nil {
			return nil, "", err
		}
		user = *usr
	} else {
		usr, err = s.existingUser(ctx, user, cred)
		if err != nil {
			return nil, "", err
		}
		if usr != nil { // org-created user
			return usr, "", err
		}
	}

	if s.autoOrg && user.OrgID == "" {
		id, err := s.orgInit.CreateOrg(ctx, storage.Organization{
			OrgName:      user.Username,
			ContactEmail: user.Email,
			Active:       s.autoApp,
		})
		if err != nil {
			logging.WithError(err, log).Error("creating org")
			return nil, "", status.Error(codes.Internal, "processing failed")
		}
		user.OrgID = id
	}
	if s.autoApp && !user.EmailVerified {
		user.EmailVerified = true
		user.InviteStatus = storage.Approved
	}

	for usr == nil {
		// set an invite code for direct signups
		user.InviteCode = random.InvitationCode(16)

		usr, err = s.usr.CreateUser(ctx, user, cred)
		if err != nil {
			logging.WithError(err, log).Error("failed to create signup storage entry")
			switch err {
			case storage.UsernameExists:
				return nil, "", status.Errorf(codes.InvalidArgument, "username invalid")
			case storage.EmailExists:
				return nil, "", status.Errorf(codes.InvalidArgument, "email already registered")
			case storage.InvCodeExists:
				log.Info("code already exists, recreating..")
				continue
			}
			return nil, "", err
		}
		if code, err = s.usr.CreateConfirmationCode(ctx, usr.ID); err != nil {
			logging.WithError(err, log).Error("failed to create confirmation code")
			return nil, "", err
		}

		user.ID = usr.ID
		break
	}

	if !s.autoApp {
		return &user, code, nil
	}

	// org approves itself
	ctx = metautils.ExtractIncoming(ctx).Set(hydra.OrgIDKey, user.OrgID).ToIncoming(ctx)
	if _, err := s.orgInit.ActivateOrg(ctx, user.ID, user.OrgID); err != nil {
		logging.WithError(err, log).Error("activating org")
		return nil, "", err
	}
	return &user, code, nil
}

func (s *Svc) invitedUser(ctx context.Context, user storage.User, cred storage.Credential) (*storage.User, error) {
	log := logging.FromContext(ctx).WithField("method", "core.user.inviteduser")

	dbUser, err := s.usr.GetUserByInvite(ctx, user.InviteCode)
	if err != nil {
		logging.WithError(err, log).Error("invite code doesn't exist")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	if dbUser.InviteStatus != storage.InviteSent {
		log.Error("Invite code already used")
		return nil, status.Error(codes.FailedPrecondition, "invite code already used")
	}
	user.ID = dbUser.ID
	user.EmailVerified = true
	if s.autoApp {
		user.InviteStatus = storage.Approved
	}
	if cred.Username != dbUser.Username {
		if err := s.usr.SetUsername(ctx, user.InviteCode, cred.Username); err != nil {
			switch err {
			case storage.NotFound:
				return nil, status.Error(codes.NotFound, "user not found")
			case storage.UsernameExists:
				return nil, status.Error(codes.AlreadyExists, "username already exists")
			default:
				return nil, status.Error(codes.Internal, "processing failed")
			}
		}
	}
	usr, err := s.usr.UpdateUserByID(ctx, user)
	if err != nil {
		logging.WithError(err, log).Error("updating user")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	if err := s.usr.SetPasswordByID(ctx, dbUser.ID, cred.Password); err != nil {
		logging.WithError(err, log).Error("set password record")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return usr, nil
}

func (s *Svc) existingUser(ctx context.Context, user storage.User, cred storage.Credential) (*storage.User, error) {
	log := logging.FromContext(ctx).WithField("method", "core.user.existinguser")

	// check for email record
	dbUser, err := s.usr.GetUserByEmail(ctx, user.Email)
	if err != nil && err != storage.NotFound {
		logging.WithError(err, log).Error("user record")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	if dbUser == nil {
		return nil, nil
	}
	// duplicate creation
	if dbUser.InviteStatus == storage.Approved || s.usrExistErr {
		log.WithField("userid", dbUser.ID).Error("already exists")
		return nil, status.Error(codes.AlreadyExists, "User already exists. Please log in.")
	}
	// org-created user.  set password only
	if err := s.usr.SetPasswordByID(ctx, dbUser.ID, cred.Password); err != nil {
		logging.WithError(err, log).Error("set password record")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return dbUser, nil
}

func (s *Svc) getUser(ctx context.Context, u storage.User) (*storage.User, error) {
	log := logging.FromContext(ctx).WithField("method", "user.getUser")
	user, err := s.usr.GetUserByEmail(ctx, u.Email)
	if err != nil && err != storage.NotFound {
		logging.WithError(err, log).Error("user record")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	if user != nil {
		return user, nil
	}
	if s.invReq && u.InviteCode == "" {
		return nil, status.Error(codes.InvalidArgument, "invite code required")
	}
	user, err = s.usr.GetUserByInvite(ctx, u.InviteCode)
	if err != nil {
		logging.WithError(err, log).Error("user invite")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return user, nil
}
