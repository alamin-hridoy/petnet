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

func (s *Svc) ValidateReceiveMoney(ctx context.Context, req *bpa.ValidateReceiveMoneyRequest) (res *bpa.ValidateReceiveMoneyResponse, err error) {
	log := logging.FromContext(ctx)
	dtls, err := json.Marshal(req)
	if err != nil {
		dtls = []byte{}
	}

	pvrN := ""
	defer func() {
		_, err := util.RecordRemittanceHistory(ctx, s.st, storage.PerahubRemittanceHistory{
			DsaID:                         phmw.GetDSA(ctx),
			UserID:                        phmw.GetUserID(ctx),
			Phrn:                          req.GetPhrn(),
			PayoutValidateReferenceNumber: pvrN,
			TxnStatus:                     storage.VALIDATE_RECEIVE,
			Details:                       dtls,
			TxnCreatedTime:                time.Now(),
			Total:                         0,
			PayHisErr:                     err,
		})
		if err != nil {
			logging.WithError(err, log).Error("ValidateReceiveMoney: record remittance history db error")
		}
	}()
	rvsm, err := s.ph.RemitanceValidateReceiveMoney(ctx, perahub.RemitanceValidateReceiveMoneyReq{
		Phrn:                  req.GetPhrn(),
		PrincipalAmount:       req.GetPrincipalAmount(),
		IsoOriginatingCountry: req.GetIsoOriginatingCountry(),
		IsoDestinationCountry: req.GetIsoDestinationCountry(),
		SenderLastName:        req.GetSenderLastName(),
		SenderFirstName:       req.GetSenderFirstName(),
		SenderMiddleName:      req.GetSenderMiddleName(),
		ReceiverLastName:      req.GetReceiverLastName(),
		ReceiverFirstName:     req.GetReceiverFirstName(),
		ReceiverMiddleName:    req.GetReceiverMiddleName(),
		PayoutPartnerCode:     req.GetPayoutPartnerCode(),
	})
	if err != nil {
		logging.WithError(err, log).Error("ValidateReceiveMoney error")
		return nil, handlePerahubError(err)
	}

	pvrN = rvsm.Result.PayoutValidateReferenceNumber

	return &bpa.ValidateReceiveMoneyResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result:  &bpa.ValidateReceiveMoneyResult{PayoutValidateReferenceNumber: rvsm.Result.PayoutValidateReferenceNumber},
	}, nil
}
