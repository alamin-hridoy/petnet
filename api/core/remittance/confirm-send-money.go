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

func (s *Svc) ConfirmSendMoney(ctx context.Context, req *bpa.ConfirmSendMoneyRequest) (res *bpa.ConfirmSendMoneyResponse, err error) {
	log := logging.FromContext(ctx)
	dtls, err := json.Marshal(req)
	if err != nil {
		dtls = []byte{}
	}

	svrN := ""
	defer func() {
		_, err := util.RecordRemittanceHistory(ctx, s.st, storage.PerahubRemittanceHistory{
			SendValidateReferenceNumber: req.GetSendValidateReferenceNumber(),
			DsaID:                       phmw.GetDSA(ctx),
			UserID:                      phmw.GetUserID(ctx),
			Phrn:                        svrN,
			TxnStatus:                   storage.CONFIRM_SEND,
			Details:                     dtls,
			TxnConfirmTime:              time.Now(),
			PayHisErr:                   err,
		})
		if err != nil {
			logging.WithError(err, log).Error("confirm send money: record remittance history db error")
		}
	}()
	rvsm, err := s.ph.RemitanceConfirmSendMoney(ctx, perahub.RemitanceConfirmSendMoneyReq{
		SendValidateReferenceNumber: req.GetSendValidateReferenceNumber(),
	})
	if err != nil {
		logging.WithError(err, log).Error("ConfirmSendMoney error")
		return nil, handlePerahubError(err)
	}

	svrN = rvsm.Result.Phrn
	return &bpa.ConfirmSendMoneyResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result:  &bpa.ConfirmSendMoneyResult{Phrn: rvsm.Result.Phrn},
	}, nil
}
