package userremit

import (
	"context"

	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) GetPartnerByTxnID(ctx context.Context, txnID string) (string, error) {
	log := logging.FromContext(ctx)
	rm, err := s.st.GetRemitCache(ctx, txnID)
	if err != nil || rm.RemcoID == "" {
		logging.WithError(err, log).Error("partner not found")
		return "", err
	}
	return rm.RemcoID, nil
}
