package remittance

import (
	"context"
	"encoding/json"

	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Inquire(ctx context.Context, req *bpa.InquireRequest) (res *bpa.InquireResponse, err error) {
	log := logging.FromContext(ctx)
	dtls, err := json.Marshal(req)
	if err != nil {
		dtls = []byte{}
	}

	phrn := ""
	defer func() {
		_, err := util.RecordRemittanceHistory(ctx, s.st, storage.PerahubRemittanceHistory{
			DsaID:     phmw.GetDSA(ctx),
			UserID:    phmw.GetUserID(ctx),
			Phrn:      phrn,
			Details:   dtls,
			PayHisErr: err,
		})
		if err != nil {
			logging.WithError(err, log).Error("Inquire: record remittance history db error")
		}
	}()
	rvsm, err := s.ph.RemitanceInquire(ctx, perahub.RemitanceInquireReq{
		Phrn: req.GetPhrn(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Inquire error")
		return nil, handlePerahubError(err)
	}

	phrn = rvsm.Result.Phrn

	return &bpa.InquireResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.InquireResult{
			Phrn:                  rvsm.Result.Phrn,
			PrincipalAmount:       int32(rvsm.Result.PrincipalAmount),
			IsoCurrency:           rvsm.Result.IsoCurrency,
			ConversionRate:        int32(rvsm.Result.ConversionRate),
			IsoOriginatingCountry: rvsm.Result.IsoOriginatingCountry,
			IsoDestinationCountry: rvsm.Result.IsoDestinationCountry,
			SenderLastName:        rvsm.Result.SenderLastName,
			SenderFirstName:       rvsm.Result.SenderFirstName,
			SenderMiddleName:      rvsm.Result.SenderMiddleName,
			ReceiverLastName:      rvsm.Result.ReceiverLastName,
			ReceiverFirstName:     rvsm.Result.ReceiverFirstName,
			ReceiverMiddleName:    rvsm.Result.ReceiverMiddleName,
		},
	}, nil
}
