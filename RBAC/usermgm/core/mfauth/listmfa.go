package mfauth

import (
	"context"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) ListMFA(ctx context.Context, user string) ([]core.MFA, error) {
	log := logging.FromContext(ctx).WithField("method", "core.mfauth.mfauth")

	m, err := s.st.GetMFAByUserID(ctx, user)
	if err != nil {
		logging.WithError(err, log).Error("list mfa from storage")
		if err == storage.NotFound {
			return nil, nil
		}
		return nil, err
	}

	recCd := false
	mfa := make([]core.MFA, 0, len(m))
	for _, mf := range m {
		if mf.MFAType == storage.Recovery {
			if recCd {
				continue
			}
			recCd = true
		}
		a := core.MFA{
			UserID:    user,
			Type:      mf.MFAType,
			MFAID:     mf.ID,
			Confirmed: mf.Confirmed.Time,
			Updated:   mf.Updated,
			Revoked:   mf.Revoked.Time,
		}
		switch mf.MFAType {
		case storage.Pass, storage.TOTP, storage.PINCode, storage.Recovery:
		case storage.SMS, storage.EMail:
			a.Source = mf.Token
		}
		mfa = append(mfa, a)
	}
	return mfa, nil
}
