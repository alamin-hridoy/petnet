package rbsession

import (
	"context"
	"database/sql"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	apb "brank.as/rbac/gunk/v1/authenticate"
	upb "brank.as/rbac/gunk/v1/user"
	"brank.as/rbac/usermgm/errors/session"
)

const (
	OrgType     = "orgtype"
	OrgID       = "orgid"
	PetnetOwner = "petnetowner"
	Provider    = "provider"
)

func (s *Svc) Login(ctx context.Context, req *apb.LoginRequest) (*apb.Session, error) {
	log := logging.FromContext(ctx)
	log.Info("received")
	if req.GetMFAToken() == "" {
		upf, err := s.core.GetUserProfileByEmail(ctx, req.GetUsername())
		if err != nil {
			return s.ErrorHandler(ctx, err)
		}
		nt := sql.NullTime{}
		if upf.Deleted != nt {
			return s.ErrorHandler(ctx, status.Error(codes.NotFound, "User Profile is not available."))
		}
		if s.core.SessionExists(ctx, upf.UserID) {
			return nil, status.Error(codes.AlreadyExists, "session already exists")
		}
	}
	res, err := s.scl.Login(ctx, req)
	if err != nil {
		logging.WithError(err, log).
			WithField("username", req.GetUsername()).Error("failed to authenticate user")
		return s.ErrorHandler(ctx, err)
	}

	if res.GetOrgID() == "" {
		logging.WithError(err, log).Error("user session missing org id")
		return nil, status.Error(codes.NotFound, "user profile not found")
	}

	if res.MFAEventID != "" {
		return &apb.Session{
			UserID:     res.UserID,
			OrgID:      res.OrgID,
			MFAEventID: res.MFAEventID,
			MFAType:    res.MFAType,
		}, nil
	}

	opf, err := s.core.GetOrgProfile(ctx, res.OrgID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to get profile")
	}

	isPetnetOwner := "false"
	if opf.OrgType == int(profile.OrgType_PetNet) && res.GetUserID() == opf.UserID {
		isPetnetOwner = "true"
	}
	sess := &apb.Session{
		UserID: res.UserID,
		OrgID:  res.OrgID,
		Session: map[string]string{
			OrgType:     strconv.Itoa(opf.OrgType),
			PetnetOwner: isPetnetOwner,
			Provider: func() string {
				if opf.IsProvider {
					return "true"
				}
				return "false"
			}(),
		},
	}
	log.WithField("session", sess).Debug("authenticated")
	return sess, nil
}

func (s *Svc) GetSession(ctx context.Context, req *apb.GetSessionRequest) (*apb.Session, error) {
	log := logging.FromContext(ctx)

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.ucl.GetUser(ctx, &upb.GetUserRequest{
		ID: req.UserID,
	})
	if err != nil {
		logging.WithError(err, log).Error("getting user")
		return nil, status.Error(codes.NotFound, "user not found")
	}
	orgType := profile.OrgType_UnknownOrgType

	pf, _ := s.core.GetOrgProfile(ctx, res.User.OrgID)
	if err == nil {
		orgType = profile.OrgType(pf.OrgType)
	}

	isPetnetOwner := "false"
	if pf.OrgType == int(profile.OrgType_PetNet) && res.GetUser().GetID() == pf.UserID {
		isPetnetOwner = "true"
	}
	sess := &apb.Session{
		UserID: res.User.ID,
		OrgID:  res.User.OrgID,
		Session: map[string]string{
			OrgType:     strconv.Itoa(int(orgType)),
			PetnetOwner: isPetnetOwner,
		},
	}
	log.WithField("session", sess).Debug("refreshed")
	return sess, nil
}

func (s *Svc) ErrorHandler(ctx context.Context, err error) (*apb.Session, error) {
	switch status.Code(err) {
	case codes.NotFound:
		return nil, session.Error(codes.NotFound, "NotFound", &apb.SessionError{
			Message:           "User Account Not Found",
			RemainingAttempts: 1,
			ErrorDetails: map[string]string{
				"NotFoundError": err.Error(),
			},
		})
	case codes.PermissionDenied:
		return nil, session.Error(codes.PermissionDenied, "PermissionDenied", &apb.SessionError{
			Message:           "Permission Denied",
			RemainingAttempts: 1,
			ErrorDetails: map[string]string{
				"PermissionError": err.Error(),
			},
		})
	case codes.ResourceExhausted:
		return nil, session.Error(codes.ResourceExhausted, "ResourceExhausted", &apb.SessionError{
			Message:           "Resource Exhausted",
			RemainingAttempts: 1,
			ErrorDetails: map[string]string{
				"ResourceError": err.Error(),
			},
		})
	default:
		return nil, session.Error(codes.PermissionDenied, "default", &apb.SessionError{
			Message:           "You entered an invalid OTP. Please try again",
			RemainingAttempts: 1,
			ErrorDetails: map[string]string{
				"DefaultError": err.Error(),
			},
		})
	}
}
