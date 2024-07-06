package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecordRemittanceHistory(ctx context.Context, st *postgres.Storage, ph storage.PerahubRemittanceHistory) (res *storage.PerahubRemittanceHistory, err error) {
	log := logging.FromContext(ctx)
	var errCode string
	var errMsg string
	var errType perahub.ErrType
	if ph.PayHisErr != nil {
		switch t := ph.PayHisErr.(type) {
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
			errMsg = ph.PayHisErr.Error()
			errType = perahub.DRPError
		}
	}

	switch {
	case ph.DsaID == "":
		log.Error("DsaID cannot be empty")
		return nil, errors.New("DsaID cannot be empty")
	case ph.UserID == "":
		log.Error("UserID cannot be empty")
		return nil, errors.New("UserID cannot be empty")
	}
	if ph.PayHisErr != nil && ph.TxnStatus == storage.VALIDATE_SEND {
		res, err = st.CreateRemittanceHistory(ctx, ph)
		if err != nil {
			logging.WithError(err, log).Error("creating validate send remittance history")
			if err == storage.Conflict {
				return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("a transaction with remittance history id: %s, already exists", ph.RemittanceHistoryID))
			}
			return
		}
		ph.RemittanceHistoryID = res.RemittanceHistoryID
	}

	if ph.PayHisErr != nil {
		ph.ErrorCode = errCode
		ph.ErrorMessage = errMsg
		ph.ErrorType = string(errType)
		ph.ErrorTime = time.Now().String()
		ph.TxnStatus = storage.TRANSACTION_FAIL
	}
	switch ph.TxnStatus {
	case storage.VALIDATE_SEND:
		res, err = st.CreateRemittanceHistory(ctx, ph)

		if err != nil {
			logging.WithError(err, log).Error("creating validate send remittance history")
			if err == storage.Conflict {
				return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("a transaction with remittance history id: %s, already exists", ph.RemittanceHistoryID))
			}
			return
		}
	case storage.CONFIRM_SEND:
		res, err = st.ConfirmRemittanceHistory(ctx, ph)
		if err != nil {
			logging.WithError(err, log).Error("confirm send remittance history")
			return
		}
	case storage.CANCEL_SEND:
		res, err = st.CancelRemittanceHistory(ctx, ph)
		if err != nil {
			logging.WithError(err, log).Error("cancel remittance history")
			return
		}
	case storage.VALIDATE_RECEIVE:
		res, err = st.ValidateReceiveRemittanceHistory(ctx, ph)
		if err != nil {
			logging.WithError(err, log).Error("confirm validate remittance history")
			return
		}
	case storage.CONFIRM_RECEIVE:
		res, err = st.ConfirmReceiveRemittanceHistory(ctx, ph)
		if err != nil {
			logging.WithError(err, log).Error("confirm receive remittance history")
			return
		}
	case storage.TRANSACTION_FAIL:
		res, err = st.UpdateRemittanceHistory(ctx, ph)
		if err != nil {
			logging.WithError(err, log).Error("update to failed remittance history")
			return
		}
	}
	return
}
