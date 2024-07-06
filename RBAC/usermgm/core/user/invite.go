package user

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/integrations/email"
	"brank.as/rbac/usermgm/storage"
)

var sgTz = time.FixedZone("Asia/Manila", 8*3600)

func (s *Svc) InviteUser(ctx context.Context, req core.Invite) (*core.Invite, error) {
	log := logging.FromContext(ctx).WithField("method", "core.user.inviteuser")

	if req.OrgID == "" {
		orgID, err := s.org.CreateOrg(ctx, storage.Organization{
			OrgName:      req.OrgName,
			ContactEmail: req.Email,
		})
		if err != nil {
			logging.WithError(err, log).Error("create invited org")
			return nil, status.Error(codes.Internal, "creating records failed")
		}
		req.OrgID = orgID
	} else if _, err := s.org.GetOrgByID(ctx, req.OrgID); err != nil {
		if err == storage.NotFound {
			return nil, status.Errorf(codes.InvalidArgument, "invalid OrgID")
		}
		logging.WithError(err, log).Error("failed to get org recerd")
		return nil, status.Error(codes.Internal, "record not found")
	}

	invExpiry := time.Now().AddDate(0, 0, 3) // invitation expires after 3 days
	usr := storage.User{
		OrgID:        req.OrgID,
		FirstName:    req.FName,
		LastName:     req.LName,
		Email:        req.Email,
		InviteStatus: storage.InviteSent,
		InviteSender: req.InvUserID,
		InviteExpiry: invExpiry,
	}

	usr.InviteCode = random.InvitationCode(16)
	cred := storage.Credential{}

	// TODO(robin): for now set the username to email, fix later, usernames should be able to be empty and not cause collision
	usr.Username = usr.Email
	u, err := s.usr.CreateUser(ctx, usr, cred)
	for err != nil {
		logging.WithError(err, log).Error("failed to create user storage entry")
		switch err {
		case storage.InvCodeExists:
			log.Trace("code already exists, recreating..")
		case storage.EmailExists:
			return nil, status.Errorf(codes.AlreadyExists, "email already exists")
		case storage.UsernameExists:
			return nil, status.Errorf(codes.AlreadyExists, "username already exists")
		default:
			return nil, status.Errorf(codes.Internal, "unable to create user")
		}
		usr.InviteCode = random.InvitationCode(16)
		u, err = s.usr.CreateUser(ctx, usr, cred)
	}

	exp := time.Now().In(sgTz).AddDate(0, 0, 3)
	inv := email.Invite{
		Username:        usr.FirstName + " " + usr.LastName,
		UserEmail:       usr.Email,
		Duration:        "3 days",
		ExpiryDate:      fmt.Sprintf("%s %d", exp.Month().String(), exp.Day()),
		CustomEmailData: req.CustomEmailData,
	}
	if err = s.mail.InviteUser(req.Email, usr.InviteCode, inv); err != nil {
		const errMsg = "failed to send invite email"
		logging.WithError(err, log).Error(errMsg)
		if !s.devEnv {
			return nil, status.Error(codes.Internal, errMsg)
		}
	}
	req.ID = u.ID
	req.Code = usr.InviteCode
	return &req, nil
}
