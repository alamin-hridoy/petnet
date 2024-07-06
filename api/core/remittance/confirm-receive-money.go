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

func (s *Svc) ConfirmReceiveMoney(ctx context.Context, req *bpa.ConfirmReceiveMoneyRequest) (res *bpa.ConfirmReceiveMoneyResponse, err error) {
	log := logging.FromContext(ctx)
	dtls, err := json.Marshal(req)
	if err != nil {
		dtls = []byte{}
	}
	phrN := ""
	defer func() {
		_, err := util.RecordRemittanceHistory(ctx, s.st, storage.PerahubRemittanceHistory{
			RemittanceHistoryID:           req.GetPayoutValidateReferenceNumber(),
			DsaID:                         phmw.GetDSA(ctx),
			UserID:                        phmw.GetUserID(ctx),
			Phrn:                          phrN,
			PayoutValidateReferenceNumber: req.GetPayoutValidateReferenceNumber(),
			TxnStatus:                     storage.CONFIRM_RECEIVE,
			Details:                       dtls,
			TxnCreatedTime:                time.Now(),
			PayHisErr:                     err,
		})
		if err != nil {
			logging.WithError(err, log).Error("CancelSendMoney: record remittance history db error")
		}
	}()
	rvsm, err := s.ph.RemitanceConfirmReceiveMoney(ctx, perahub.RemitanceConfirmReceiveMoneyReq{
		PayoutValidateReferenceNumber: req.GetPayoutValidateReferenceNumber(),
	})
	if err != nil {
		logging.WithError(err, log).Error("ConfirmReceiveMoney error")
		return nil, handlePerahubError(err)
	}

	phrN = rvsm.Result.Phrn
	return &bpa.ConfirmReceiveMoneyResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.ConfirmReceiveMoneyResult{
			Phrn:                  rvsm.Result.Phrn,
			PrincipalAmount:       int32(rvsm.Result.PrincipalAmount),
			IsoOriginatingCountry: rvsm.Result.IsoOriginatingCountry,
			IsoDestinationCountry: rvsm.Result.IsoDestinationCountry,
			SenderLastName:        rvsm.Result.SenderLastName,
			SenderFirstName:       rvsm.Result.SenderFirstName,
			SenderMiddleName:      rvsm.Result.SenderMiddleName,
			ReceiverLastName:      rvsm.Result.ReceiverLastName,
			ReceiverFirstName:     rvsm.Result.ReceiverFirstName,
			ReceiverMiddleName:    rvsm.Result.ReceiverFirstName,
		},
	}, nil
}
