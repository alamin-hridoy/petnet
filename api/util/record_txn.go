package util

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/serviceutil/logging"
)

const (
	PerahubTrxTypeOTC        = "otc"
	PerahubTrxTypeDigital    = "digital"
	PerahubRemitTypeInbound  = "inbound"
	PerahubRemitTypeOutbound = "outbound"
)

type StageTxnOpts struct {
	TxnID       string
	TxnType     storage.TxnType
	PtnrRemType string
	TxnErr      error
}

type ConfirmTxnOpts struct {
	TxnID       string
	TxnType     storage.TxnType
	PtnrRemType string
	TxnErr      error
}

type errorType string

func GetPerahubTrxType(ctx context.Context) string {
	if strings.ToUpper(phmw.GetTransactionTypes(ctx)) == "DIGITAL" {
		return PerahubTrxTypeDigital
	}

	return PerahubTrxTypeOTC
}

func RecordStageTxn(ctx context.Context, st *postgres.Storage, rmt core.Remittance, o StageTxnOpts) (*storage.RemitHistory, error) {
	log := logging.FromContext(ctx)

	var errCode string
	var errMsg string
	var errType perahub.ErrType
	if o.TxnErr != nil {
		// setting zero amount to PHP by default on failed transactions
		zamt := core.MustMinor("0", "PHP")
		rmt.GrossTotal = zamt
		rmt.SourceAmount = zamt
		rmt.DestAmount = zamt
		rmt.Tax = zamt
		rmt.Charge = zamt

		switch t := o.TxnErr.(type) {
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
			errMsg = o.TxnErr.Error()
			errType = perahub.DRPError
		}
	}

	switch {
	case rmt.GrossTotal.CurrencyCode() == "":
		log.Error("GrossTotal cannot be empty")
		return nil, errors.New("GrossTotal cannot be empty")
	case rmt.SourceAmount.CurrencyCode() == "":
		log.Error("SourceAmount cannot be empty")
		return nil, errors.New("SourceAmount cannot be empty")
	case rmt.DestAmount.CurrencyCode() == "":
		log.Error("DestAmount cannot be empty")
		return nil, errors.New("DestAmount cannot be empty")
	case rmt.Tax.CurrencyCode() == "":
		log.Error("Tax cannot be empty")
		return nil, errors.New("Tax cannot be empty")
	case rmt.Charge.CurrencyCode() == "":
		log.Error("Charge cannot be empty")
		return nil, errors.New("Charge cannot be empty")
	case o.TxnType == "":
		log.Error("TxnType cannot be empty")
		return nil, errors.New("TxnType cannot be empty")
	case o.PtnrRemType == "":
		log.Error("PtnrRemType cannot be empty")
		return nil, errors.New("PtnrRemType cannot be empty")
	// if there is an error the cache hasn't been created and there is no txn id yet so it needs to be generated
	case o.TxnID == "":
		o.TxnID = uuid.NewString()
	}

	txnSts := storage.SuccessStatus
	if o.TxnErr != nil {
		txnSts = storage.FailStatus
	}

	trxType := GetPerahubTrxType(ctx)
	req := storage.RemitHistory{
		TxnID:          o.TxnID,
		DsaID:          rmt.DsaID,
		DsaOrderID:     rmt.DsaOrderID,
		UserID:         rmt.UserID,
		RemcoID:        rmt.RemitPartner,
		RemType:        string(o.TxnType),
		TxnStatus:      string(txnSts),
		TxnStep:        string(storage.StageStep),
		ReceiverID:     rmt.Receiver.PartnerMemberID,
		RemcoControlNo: rmt.ControlNo,
		TxnStagedTime:  sql.NullTime{Time: time.Now(), Valid: true},
		ErrorCode:      errCode,
		ErrorMsg:       errMsg,
		ErrorType:      string(errType),
		TransactionType: sql.NullString{
			String: trxType,
			Valid:  len(trxType) != 0,
		},
		Remittance: storage.Remittance{
			RemcoAltControlNo: rmt.RemcoAltControlNo,
			TxnType:           o.PtnrRemType,
			CustomerTxnID:     rmt.CustomerTxnID,
			Remitter: storage.Contact{
				FirstName:     rmt.Remitter.FName,
				MiddleName:    rmt.Remitter.MdName,
				LastName:      rmt.Remitter.LName,
				RemcoMemberID: rmt.Remitter.PartnerMemberID,
				Message:       rmt.Message,
				Email:         rmt.Remitter.Email,
				Address1:      rmt.Remitter.Address.Address1,
				Address2:      rmt.Remitter.Address.Address2,
				City:          rmt.Remitter.Address.City,
				State:         rmt.Remitter.Address.State,
				PostalCode:    rmt.Remitter.Address.PostalCode,
				Country:       rmt.Remitter.Address.Country,
				Province:      rmt.Remitter.Address.Province,
				Zone:          rmt.Remitter.Address.Zone,
				PhoneCty:      rmt.Remitter.Phone.CtyCode,
				Phone:         rmt.Remitter.Phone.Number,
				MobileCty:     rmt.Remitter.Mobile.CtyCode,
				Mobile:        rmt.Remitter.Mobile.Number,
			},
			Receiver: storage.Contact{
				FirstName:     rmt.Receiver.FName,
				MiddleName:    rmt.Receiver.MdName,
				LastName:      rmt.Receiver.LName,
				RemcoMemberID: rmt.Receiver.PartnerMemberID,
				Message:       rmt.Message,
				Email:         rmt.Receiver.Email,
				Address1:      rmt.Receiver.Address.Address1,
				Address2:      rmt.Receiver.Address.Address2,
				City:          rmt.Receiver.Address.City,
				State:         rmt.Receiver.Address.State,
				PostalCode:    rmt.Receiver.Address.PostalCode,
				Country:       rmt.Receiver.Address.Country,
				Province:      rmt.Receiver.Address.Province,
				Zone:          rmt.Receiver.Address.Zone,
				PhoneCty:      rmt.Receiver.Phone.CtyCode,
				Phone:         rmt.Receiver.Phone.Number,
				MobileCty:     rmt.Receiver.Mobile.CtyCode,
				Mobile:        rmt.Receiver.Mobile.Number,
			},
			GrossTotal: storage.GrossTotal{Minor: rmt.GrossTotal},
			SourceAmt:  rmt.SourceAmount,
			DestAmt:    rmt.DestAmount,
			Tax:        rmt.Tax,
			Charge:     rmt.Charge,
		},
	}
	if o.TxnErr != nil {
		req.ErrorTime = req.TxnStagedTime
	}

	res, err := st.CreateRemitHistory(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("creating remit history")
		if err == storage.Conflict {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("a transaction with order id: %s, already exists", rmt.DsaOrderID))
		}
		return nil, err
	}
	return res, nil
}

func RecordConfirmTxn(ctx context.Context, st *postgres.Storage, rmt core.Remittance, o ConfirmTxnOpts) (*storage.RemitHistory, error) {
	log := logging.FromContext(ctx)

	var errCode string
	var errMsg string
	var errType perahub.ErrType
	if o.TxnErr != nil {
		switch t := o.TxnErr.(type) {
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
			errMsg = o.TxnErr.Error()
			errType = perahub.DRPError
		}
	}

	switch {
	case rmt.GrossTotal.CurrencyCode() == "":
		log.Error("GrossTotal cannot be empty")
		return nil, errors.New("GrossTotal cannot be empty")
	case rmt.SourceAmount.CurrencyCode() == "":
		log.Error("SourceAmount cannot be empty")
		return nil, errors.New("SourceAmount cannot be empty")
	case rmt.DestAmount.CurrencyCode() == "":
		log.Error("DestAmount cannot be empty")
		return nil, errors.New("DestAmount cannot be empty")
	case rmt.Tax.CurrencyCode() == "":
		log.Error("Tax cannot be empty")
		return nil, errors.New("Tax cannot be empty")
	case rmt.Charge.CurrencyCode() == "":
		log.Error("Charge cannot be empty")
		return nil, errors.New("Charge cannot be empty")
	case o.TxnType == "":
		log.Error("TxnType cannot be empty")
		return nil, errors.New("TxnType cannot be empty")
	case o.PtnrRemType == "":
		log.Error("PtnrRemType cannot be empty")
		return nil, errors.New("PtnrRemType cannot be empty")
	case o.TxnID == "":
		log.Error("TxnID cannot be empty")
		return nil, errors.New("TxnID cannot be empty")
	}

	txnSts := storage.SuccessStatus
	if o.TxnErr != nil {
		txnSts = storage.FailStatus
	}

	trxType := GetPerahubTrxType(ctx)
	req := storage.RemitHistory{
		TxnID:            o.TxnID,
		RemType:          string(o.TxnType),
		TxnStatus:        string(txnSts),
		TxnStep:          string(storage.ConfirmStep),
		RemcoControlNo:   rmt.ControlNo,
		TxnCompletedTime: sql.NullTime{Time: time.Now(), Valid: true},
		ErrorCode:        errCode,
		ErrorMsg:         errMsg,
		ErrorType:        string(errType),
		TransactionType: sql.NullString{
			String: trxType,
			Valid:  len(trxType) != 0,
		},
		Remittance: storage.Remittance{
			RemcoAltControlNo: rmt.RemcoAltControlNo,
			TxnType:           o.PtnrRemType,
			CustomerTxnID:     rmt.CustomerTxnID,
			Remitter: storage.Contact{
				FirstName:     rmt.Remitter.FName,
				MiddleName:    rmt.Remitter.MdName,
				LastName:      rmt.Remitter.LName,
				RemcoMemberID: rmt.Remitter.PartnerMemberID,
				Message:       rmt.Message,
				Email:         rmt.Remitter.Email,
				Address1:      rmt.Remitter.Address.Address1,
				Address2:      rmt.Remitter.Address.Address2,
				City:          rmt.Remitter.Address.City,
				State:         rmt.Remitter.Address.State,
				PostalCode:    rmt.Remitter.Address.PostalCode,
				Country:       rmt.Remitter.Address.Country,
				Province:      rmt.Remitter.Address.Province,
				Zone:          rmt.Remitter.Address.Zone,
				PhoneCty:      rmt.Remitter.Phone.CtyCode,
				Phone:         rmt.Remitter.Phone.Number,
				MobileCty:     rmt.Remitter.Mobile.CtyCode,
				Mobile:        rmt.Remitter.Mobile.Number,
			},
			Receiver: storage.Contact{
				FirstName:     rmt.Receiver.FName,
				MiddleName:    rmt.Receiver.MdName,
				LastName:      rmt.Receiver.LName,
				RemcoMemberID: rmt.Receiver.PartnerMemberID,
				Message:       rmt.Message,
				Email:         rmt.Receiver.Email,
				Address1:      rmt.Receiver.Address.Address1,
				Address2:      rmt.Receiver.Address.Address2,
				City:          rmt.Receiver.Address.City,
				State:         rmt.Receiver.Address.State,
				PostalCode:    rmt.Receiver.Address.PostalCode,
				Country:       rmt.Receiver.Address.Country,
				Province:      rmt.Receiver.Address.Province,
				Zone:          rmt.Receiver.Address.Zone,
				PhoneCty:      rmt.Receiver.Phone.CtyCode,
				Phone:         rmt.Receiver.Phone.Number,
				MobileCty:     rmt.Receiver.Mobile.CtyCode,
				Mobile:        rmt.Receiver.Mobile.Number,
			},
			GrossTotal: storage.GrossTotal{Minor: rmt.GrossTotal},
			SourceAmt:  rmt.SourceAmount,
			DestAmt:    rmt.DestAmount,
			Tax:        rmt.Tax,
			Charge:     rmt.Charge,
		},
	}
	if o.TxnErr != nil {
		req.ErrorTime = req.TxnCompletedTime
	}

	res, err := st.UpdateRemitHistory(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("creating remit history")
		if err == storage.Conflict {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("a transaction with order id: %s, already exists", rmt.DsaOrderID))
		}
		return nil, err
	}
	return res, nil
}
