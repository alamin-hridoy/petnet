package mfauth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) MFAuth(ctx context.Context, m core.MFAChallenge) (*core.MFAChallenge, error) {
	log := logging.FromContext(ctx).WithField("method", "core.mfauth.mfauth")

	switch m.Type {
	case storage.Pass:
		_, _, err := s.st.ValidateUserPass(ctx, m.UserID, m.Token)
		if err != nil {
			logging.WithError(err, log).Error("password validation")
			return nil, status.Error(codes.InvalidArgument, "invalid password")
		}
	case storage.TOTP, storage.SMS, storage.EMail, storage.PINCode, storage.Recovery:
		mfa, err := s.st.ConfirmMFAEvent(ctx, storage.MFAEvent{
			EventID: m.EventID,
			UserID:  m.UserID,
			MFAType: m.Type,
			Token:   m.Token,
		})
		if err != nil {
			logging.WithError(err, log).Error("confirmation failed")
			return nil, status.Error(codes.InvalidArgument, "invalid token")
		}
		m.EventID = mfa.EventID
		if mfa.Validation {
			// validation skipped for PIN and Recovery codes
			ma, err := s.st.GetMFAByType(ctx, m.UserID, mfa.MFAType)
			if err != nil {
				logging.WithError(err, log).Error("listing type")
				break
			}
			if _, err := s.st.EnableMFA(ctx, m.UserID, mfa.MFAID); err != nil {
				logging.WithError(err, log).Error("mfa activation")
				break
			}
			for _, m := range ma {
				if !m.Active {
					continue
				}
				if _, err := s.st.DisableMFA(ctx, m.ID); err != nil {
					logging.WithError(err, log).Error("mfa activation")
				}
			}
			if u, err := s.st.GetUserByID(ctx, m.UserID); err == nil && u.PreferredMFA == "" {
				u.PreferredMFA = m.Type
				if _, err := s.st.UpdateUserByID(ctx, *u); err != nil {
					logging.WithError(err, log).Error("update preference")
				}
			}
		}
	default:
		return nil, status.Error(codes.Unimplemented, "MFA type not supported")
	}
	return &m, nil
}
