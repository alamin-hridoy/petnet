package remittance

import (
	"context"
	"encoding/json"
	"time"

	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) CancelSendMoney(ctx context.Context, req *bpa.CancelSendMoneyRequest) (res *bpa.CancelSendMoneyResponse, err error) {
	log := logging.FromContext(ctx)
	dtls, err := json.Marshal(req)
	if err != nil {
		dtls = []byte{}
	}
	Phrn := ""
	Csrn := ""
	defer func() {
		_, err := util.RecordRemittanceHistory(ctx, s.st, storage.PerahubRemittanceHistory{
			DsaID:                     phmw.GetDSA(ctx),
			UserID:                    phmw.GetUserID(ctx),
			Phrn:                      Phrn,
			CancelSendReferenceNumber: Csrn,
			TxnStatus:                 storage.CANCEL_SEND,
			Details:                   dtls,
			Remarks:                   req.GetRemarks(),
			TxnUpdatedTime:            time.Now(),
			PayHisErr:                 err,
		})
		if err != nil {
			logging.WithError(err, log).Error("CancelSendMoney: record remittance history db error")
		}
	}()
	rvsm, err := s.ph.RemitanceCancelSendMoney(ctx, perahub.RemitanceCancelSendMoneyReq{
		Phrn:        req.GetPhrn(),
		PartnerCode: req.GetPartnerCode(),
		Remarks:     req.GetRemarks(),
	})
	if err != nil {
		logging.WithError(err, log).Error("CancelSendMoney error")
		return nil, handlePerahubError(err)
	}

	Phrn = rvsm.Result.Phrn
	Csrn = rvsm.Result.CancelSendReferenceNumber
	return &bpa.CancelSendMoneyResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.CancelSendMoneyResult{
			Phrn:                      rvsm.Result.Phrn,
			CancelSendDate:            rvsm.Result.CancelSendDate,
			CancelSendReferenceNumber: rvsm.Result.CancelSendReferenceNumber,
		},
	}, nil
}
