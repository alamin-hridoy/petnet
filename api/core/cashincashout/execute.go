package cashincashout

import (
	"context"
	"time"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/util"
	cio "brank.as/petnet/gunk/drp/v1/cashincashout"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
)

func (s *Svc) CiCoExecute(ctx context.Context, req *cio.CiCoExecuteRequest) (res *cio.CiCoExecuteResponse, err error) {
	log := logging.FromContext(ctx)
	trxdt, err := time.Parse("2006-01-02", req.TrxDate)
	if err != nil {
		trxdt = time.Time{}
	}
	defer func() {
		cd := 400
		msg := "execute failed"
		tp := ""
		ptn := ""
		rn := ""
		pamnt := 0
		chrg := 0
		tamnt := 0
		if res != nil && res.Result != nil {
			cd = int(res.GetCode())
			msg = res.GetMessage()
			tp = res.Result.GetTrxType()
			ptn = res.Result.GetProviderTrackingno()
			rn = res.Result.GetReferenceNumber()
			pamnt = int(res.Result.GetPrincipalAmount())
			chrg = int(res.Result.GetCharges())
			tamnt = int(res.Result.GetTotalAmount())
		}
		_, err := util.RecordCiCo(ctx, s.st, storage.CashInCashOutHistory{
			OrgID:            phmw.GetDSA(ctx),
			PartnerCode:      req.PartnerCode,
			PetnetTrackingNo: req.PetnetTrackingno,
			TrxDate:          trxdt,
		}, storage.CashInCashOutHistoryRes{
			Code:    cd,
			Message: msg,
			Result: storage.CashInCashOutHistoryDetails{
				PartnerCode:        req.GetPartnerCode(),
				Provider:           req.GetProvider(),
				PetnetTrackingno:   req.GetPetnetTrackingno(),
				TrxDate:            trxdt,
				TrxType:            tp,
				ProviderTrackingNo: ptn,
				ReferenceNumber:    rn,
				PrincipalAmount:    pamnt,
				Charges:            chrg,
				TotalAmount:        tamnt,
			},
		}, err)
		if err != nil {
			log.Error(err)
		}
	}()

	ci, err := s.ph.CicoExecute(ctx, perahub.CicoExecuteRequest{
		PartnerCode:      req.GetPartnerCode(),
		PetnetTrackingno: req.GetPetnetTrackingno(),
		TrxDate:          req.GetTrxDate(),
	})
	if err != nil {
		return nil, handleCiCoError(err)
	}

	if ci == nil || ci.Result == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, "not found")
	}

	res = &cio.CiCoExecuteResponse{
		Code:    int32(ci.Code),
		Message: ci.Message,
		Result: &cio.CicoExecuteResult{
			PartnerCode:        ci.Result.PartnerCode,
			Provider:           ci.Result.Provider,
			PetnetTrackingno:   ci.Result.PetnetTrackingno,
			TrxDate:            ci.Result.TrxDate,
			TrxType:            ci.Result.TrxType,
			ProviderTrackingno: ci.Result.ProviderTrackingno,
			ReferenceNumber:    ci.Result.ReferenceNumber,
			PrincipalAmount:    int32(ci.Result.PrincipalAmount),
			Charges:            int32(ci.Result.Charges),
			TotalAmount:        int32(ci.Result.TotalAmount),
		},
	}
	return res, nil
}
