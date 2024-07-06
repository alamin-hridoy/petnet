package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	bpa "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecordBillPayment(ctx context.Context, st *postgres.Storage, req *bpa.BPTransactRequest, res *bpa.BPTransactResponse, er error) (*storage.BillPayment, error) {
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

	layout := "2006-01-02"
	td, err := time.Parse(layout, req.GetTrxDate())
	if err != nil {
		td = time.Now()
	}
	bls, err := json.Marshal(storage.Bills{
		Code:    res.GetCode(),
		Message: res.GetMessage(),
		Result: storage.BillDetails{
			Message:         res.GetResult().GetMessage(),
			Timestamp:       res.GetResult().GetTimestamp(),
			ReferenceNumber: res.GetResult().GetReferenceNumber(),
			Status:          res.GetResult().GetStatus(),
			ServiceCharge:   core.MustMinor(string(req.GetServiceCharge()), "PHP"),
			TransactionID:   res.GetResult().GetTransactionID(),
			ClientReference: res.GetResult().GetClientReference(),
			BillerReference: res.GetResult().GetBillerReference(),
			PaymentMethod:   res.GetResult().GetPaymentMethod(),
			Amount:          core.MustMinor(req.GetAmount(), "PHP"),
			OtherCharges:    core.MustMinor(res.GetResult().GetOtherCharges(), "PHP"),
			CreatedAt:       res.GetResult().GetCreatedAt(),
			URL:             res.GetResult().GetURL(),
		},
		RemcoID: int(res.GetRemcoID()),
	})
	if err != nil {
		bls = []byte{}
	}
	billID, err := strconv.Atoi(req.GetBillID())
	if err != nil {
		logging.WithError(err, log).Error("Cant convert bill id")
	}

	othInf := storage.BillsOtherInfo{
		LastName:        req.GetOtherInfo().GetLastName(),
		FirstName:       req.GetOtherInfo().GetFirstName(),
		MiddleName:      req.GetOtherInfo().GetMiddleName(),
		PaymentType:     req.GetOtherInfo().GetPaymentType(),
		Course:          req.GetOtherInfo().GetCourse(),
		TotalAssessment: req.GetOtherInfo().GetTotalAssessment(),
		SchoolYear:      req.GetOtherInfo().GetSchoolYear(),
		Term:            req.GetOtherInfo().GetTerm(),
	}
	othrInfo, err := json.Marshal(othInf)
	if err != nil {
		logging.WithError(err, log).Error("Can't marshal OtherInfo")
	}
	amnt, _ := core.MustMinor(req.GetAmount(), "PHP").MarshalJSON()
	bp := storage.BillPayment{
		BillID:                  int32(billID),
		UserID:                  req.GetUserID(),
		BillPaymentStatus:       string(bps),
		ErrorCode:               errCode,
		ErrorMsg:                errMsg,
		ErrorType:               string(errType),
		Bills:                   bls,
		BillerTag:               req.GetBillerTag(),
		LocationID:              req.GetLocationID(),
		CurrencyID:              req.GetCurrencyID(),
		AccountNumber:           req.GetAccountNumber(),
		Amount:                  amnt,
		Identifier:              req.GetIdentifier(),
		Coy:                     req.GetCoy(),
		ServiceCharge:           amnt,
		TotalAmount:             amnt,
		BillPaymentDate:         td,
		BillerName:              req.GetBillerName(),
		RemoteUserID:            req.GetRemoteUserID(),
		CustomerID:              strconv.Itoa(int(req.GetCustomerID())),
		RemoteLocationID:        req.GetRemoteLocationID(),
		LocationName:            req.GetLocationName(),
		FormType:                req.GetFormType(),
		FormNumber:              req.GetFormNumber(),
		PaymentMethod:           req.GetPaymentMethod(),
		OtherInfo:               othrInfo,
		TrxDate:                 td,
		Created:                 time.Now(),
		Updated:                 time.Now(),
		ClientRefNumber:         req.GetClientReferenceNumber(),
		PartnerCharge:           req.GetPartnerCharge(),
		ReferenceNumber:         req.GetReferenceNumber(),
		ValidationNumber:        req.GetValidationNumber(),
		ReceiptValidationNumber: req.GetReceiptValidationNumber(),
		TpaID:                   req.GetTpaID(),
		Type:                    req.GetType(),
		TxnID:                   req.GetTxnid(),
		PartnerID:               req.GetPartner(),
		OrgID:                   phmw.GetDSAOrgID(ctx),
	}

	rs, err := st.CreateBillPayment(ctx, bp)
	if err != nil {
		logging.WithError(err, log).Error("creating bill payment")
		if err == storage.Conflict {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("a bill payment transaction with bill id: %d, already exists", bp.BillID))
		}
		return nil, err
	}
	return rs, nil
}
