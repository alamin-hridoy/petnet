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

func (s *Svc) ValidateSendMoney(ctx context.Context, req *bpa.ValidateSendMoneyRequest) (res *bpa.ValidateSendMoneyResponse, err error) {
	log := logging.FromContext(ctx)
	dtls, err := json.Marshal(req)
	if err != nil {
		dtls = []byte{}
	}

	svrN := ""
	defer func() {
		_, err := util.RecordRemittanceHistory(ctx, s.st, storage.PerahubRemittanceHistory{
			RemittanceHistoryID:         req.GetPartnerReferenceNumber(),
			DsaID:                       phmw.GetDSA(ctx),
			UserID:                      phmw.GetUserID(ctx),
			SendValidateReferenceNumber: svrN,
			TxnStatus:                   storage.VALIDATE_SEND,
			Details:                     dtls,
			TxnCreatedTime:              time.Now(),
			Total:                       0,
			PayHisErr:                   err,
		})
		if err != nil {
			logging.WithError(err, log).Error("ValidateSendMoney: record remittance history db error")
		}
	}()
	rvsm, err := s.ph.RemitanceValidateSendMoney(ctx, perahub.RemitanceValidateSendMoneyReq{
		PartnerReferenceNumber: req.GetPartnerReferenceNumber(),
		PrincipalAmount:        req.GetPrincipalAmount(),
		ServiceFee:             req.GetServiceFee(),
		IsoCurrency:            req.GetIsoCurrency(),
		ConversionRate:         req.GetConversionRate(),
		IsoOriginatingCountry:  req.GetIsoOriginatingCountry(),
		IsoDestinationCountry:  req.GetIsoDestinationCountry(),
		SenderLastName:         req.GetSenderLastName(),
		SenderFirstName:        req.GetSenderFirstName(),
		SenderMiddleName:       req.GetSenderMiddleName(),
		ReceiverLastName:       req.GetReceiverLastName(),
		ReceiverFirstName:      req.GetReceiverFirstName(),
		ReceiverMiddleName:     req.GetReceiverMiddleName(),
		SenderBirthDate:        req.GetSenderBirthDate(),
		SenderBirthPlace:       req.GetSenderBirthPlace(),
		SenderBirthCountry:     req.GetSenderBirthCountry(),
		SenderGender:           req.GetSenderGender(),
		SenderRelationship:     req.GetSenderRelationship(),
		SenderPurpose:          req.GetSenderPurpose(),
		SenderOccupation:       req.GetSenderOccupation(),
		SenderEmploymentNature: req.GetSenderEmploymentNature(),
		SendPartnerCode:        req.GetSendPartnerCode(),
		SenderSourceOfFund:     req.GetSenderSourceOfFund(),
	})
	if err != nil {
		logging.WithError(err, log).Error("ValidateSendMoney error")
		return nil, handlePerahubError(err)
	}

	svrN = rvsm.Result.SendValidateReferenceNumber
	return &bpa.ValidateSendMoneyResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result:  &bpa.ValidateSendMoneyResult{SendValidateReferenceNumber: rvsm.Result.SendValidateReferenceNumber},
	}, nil
}
