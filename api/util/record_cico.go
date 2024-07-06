package util

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecordCiCo(ctx context.Context, st *postgres.Storage, req storage.CashInCashOutHistory, res storage.CashInCashOutHistoryRes, er error) (*storage.CashInCashOutHistory, error) {
	log := logging.FromContext(ctx)
	var errCode string
	var errMsg string
	var errType perahub.ErrType
	bps := storage.SuccessStatus
	if er != nil {
		bps = storage.FailStatus
		switch t := er.(type) {
		case *perahub.Error:
			errMsg = t.Msg
			if t.UnknownErr != "" {
				errMsg = t.UnknownErr
			}
			if t.Errors != nil {
				b, err := json.Marshal(t.Errors)
				if err != nil {
					log.Error(err)
					return nil, err
				}
				errMsg = errMsg + ", " + string(b)
			}
			errCode = t.Code
			errType = t.Type
		default:
			errMsg = er.Error()
			errType = perahub.BillerError
		}
	}
	resDetails, err := json.Marshal(res)
	if err != nil {
		resDetails = []byte{}
	}
	req.Details = resDetails
	req.TotalAmount = res.Result.TotalAmount
	req.ReferenceNumber = res.Result.ReferenceNumber
	req.Provider = res.Result.Provider
	req.Charges = res.Result.Charges
	req.PartnerCode = res.Result.PartnerCode
	req.PetnetTrackingNo = res.Result.PetnetTrackingno
	req.TrxDate = res.Result.TrxDate
	req.TrxType = res.Result.TrxType
	req.ProviderTrackingNo = res.Result.ProviderTrackingNo
	req.PrincipalAmount = res.Result.PrincipalAmount
	req.ErrorCode = errCode
	req.ErrorMessage = errMsg
	req.ErrorType = string(errType)
	if errCode != "" {
		req.ErrorTime = time.Now().String()
	}
	req.TxnStatus = string(bps)
	rs, err := st.CreateCICOHistory(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("creating cash in cash out transaction")
		if err == storage.Conflict {
			rs, err = st.UpdateCICOHistory(ctx, req)
			if err != nil {
				logging.WithError(err, log).Error("updating cash in cash out transaction")
				return nil, status.Error(codes.Internal, fmt.Sprintf("transaction failed: %s", res.Result.PetnetTrackingno))
			}
			return rs, nil
		}
		return nil, err
	}
	return rs, nil
}
