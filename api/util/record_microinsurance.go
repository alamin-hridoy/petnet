package util

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
	"brank.as/petnet/serviceutil/logging"
)

// RecordMicroInsurance ...
func RecordMicroInsurance(ctx context.Context,
	st *postgres.Storage,
	r *migunk.TransactRequest,
	i *migunk.Insurance,
	resErr error,
) (*storage.MicroInsuranceHistory, error) {
	log := logging.FromContext(ctx)
	trxDate, err := time.Parse("2006-01-02", r.TrxDate)
	if err != nil {
		logging.WithError(err, log).Error("parsing TrxDate " + r.TrxDate)
		trxDate = time.Now()
	}

	bDay, _ := time.Parse("2006-01-02", r.Birthdate)
	beneficiaries, _ := json.Marshal(r.Beneficiaries)
	dependents, _ := json.Marshal(r.Dependents)

	insDetails := []byte("{}")
	traceNo := ""
	trxStatus := "Failed"
	if i != nil {
		traceNo = i.TraceNumber
		trxStatus = i.StatusDesc
		if insData, err := json.Marshal(i); err == nil {
			insDetails = insData
		}

		insTrnDate, err := time.Parse("01/02/2006", i.TrnDate)
		if err == nil && !insTrnDate.IsZero() {
			trxDate = insTrnDate
		}
	}

	errCode, errMsg, errType, errTime := parseMicroInsuranceError(resErr)

	h, err := st.CreateMicroInsuranceHistory(ctx, storage.MicroInsuranceHistory{
		DsaID:            phmw.GetDSAOrgID(ctx),
		Coy:              r.Coy,
		LocationID:       r.LocationID,
		UserCode:         r.UserCode,
		TrxDate:          trxDate,
		PromoAmount:      fmt.Sprintf("%f", r.PromoAmount),
		PromoCode:        r.PromoCode,
		Amount:           r.Amount,
		CoverageCount:    r.CoverageCount,
		ProductCode:      r.ProductCode,
		ProcessingBranch: r.ProcessingBranch,
		ProcessedBy:      r.ProcessedBy,
		UserEmail:        r.UserEmail,
		LastName:         r.LastName,
		FirstName:        r.FirstName,
		MiddleName:       r.MiddleName,
		Gender:           r.Gender,
		Birthdate:        bDay,
		MobileNumber:     r.MobileNumber,
		ProvinceCode:     r.ProvinceCode,
		CityCode:         r.CityCode,
		Address:          r.Address,
		MaritalStatus:    r.MaritalStatus,
		Occupation:       r.Occupation,
		CardNumber:       r.CardNumber,
		NumberUnits:      r.NumberUnits,
		Beneficiaries:    beneficiaries,
		Dependents:       dependents,
		TrxStatus:        trxStatus,
		TraceNumber: sql.NullString{
			String: traceNo,
			Valid:  traceNo != "",
		},
		InsuranceDetails: insDetails,
		ErrorCode:        errCode,
		ErrorMsg:         errMsg,
		ErrorType:        errType,
		ErrorTime:        errTime,
		OrgID:            phmw.GetDSAOrgID(ctx),
	})
	if err != nil {
		logging.WithError(err, log).Error("creating microinsurance history")
		return nil, status.Error(codes.Internal, "db error")
	}

	return h, nil
}

// UpdateMicroInsurance ...
func UpdateMicroInsurance(ctx context.Context, st *postgres.Storage, i *migunk.Insurance) (*storage.MicroInsuranceHistory, error) {
	log := logging.FromContext(ctx)
	if i == nil {
		log.Error("empty insurance")
		return nil, status.Error(codes.Internal, "empty insurance response")
	}

	if i.TraceNumber == "" {
		log.Error("empty trace number")
		return nil, status.Error(codes.Internal, "empty trace number")
	}

	micIns, err := st.GetMicroInsuranceHistoryByTraceNumber(ctx, i.TraceNumber)
	if err != nil {
		if err == storage.ErrNotFound {
			trxDate, e := time.Parse("02/01/2006", i.TrnDate)
			if e != nil {
				trxDate = time.Now()
			}

			r := &migunk.TransactRequest{
				TrxDate:     trxDate.Format("2006-01-02"),
				ProductCode: i.InsProductID,
				NumberUnits: fmt.Sprintf("%d", i.NumUnits),
			}

			return RecordMicroInsurance(ctx, st, r, i, nil)
		}

		logging.WithError(err, log).Error("getting microinsurance db for update")
		return nil, status.Error(codes.Internal, "db error")
	}

	micIns.TrxStatus = i.StatusDesc
	micIns.InsuranceDetails = []byte("{}")
	insData, err := json.Marshal(i)
	if err != nil {
		logging.WithError(err, log).Error("unmarshal insurance object for update")
	}

	if insData != nil {
		micIns.InsuranceDetails = insData
		micIns.ErrorCode = ""
		micIns.ErrorMsg = ""
		micIns.ErrorType = ""
		micIns.ErrorTime = sql.NullTime{Valid: false}
	}

	h, err := st.UpdateMicroInsuranceHistoryStatusByTraceNumber(ctx, *micIns)
	if err != nil {
		logging.WithError(err, log).Error("updating microinsurance history")
		return nil, status.Error(codes.Internal, "db error")
	}

	return h, nil
}

func parseMicroInsuranceError(err error) (string, string, string, sql.NullTime) {
	if err == nil {
		return "", "", "", sql.NullTime{Valid: false}
	}

	errTime := time.Now()
	switch err.(type) {
	case *perahub.Error:
		if pErr, ok := err.(*perahub.Error); ok && pErr != nil {
			msg := pErr.Msg
			if pErr.UnknownErr != "" {
				msg = pErr.UnknownErr
			}

			return pErr.Code, msg, string(pErr.Type), sql.NullTime{Time: errTime, Valid: true}
		}

	case *coreerror.Error:
		if cErr, ok := err.(*coreerror.Error); ok && cErr != nil {
			return cErr.Code.String(), cErr.Message, string(perahub.MicroInsuranceError), sql.NullTime{Time: errTime, Valid: true}
		}

	case interface{ GRPCStatus() *status.Status }:
		if grpcErr, ok := err.(interface{ GRPCStatus() *status.Status }); ok && grpcErr != nil {
			st := grpcErr.GRPCStatus()
			if st != nil {
				return st.Code().String(), st.Message(), string(perahub.MicroInsuranceError), sql.NullTime{Time: errTime, Valid: true}
			}
		}
	}

	return codes.Unknown.String(), err.Error(), string(perahub.MicroInsuranceError), sql.NullTime{Time: errTime, Valid: true}
}
