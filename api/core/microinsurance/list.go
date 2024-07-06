package microinsurance

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
	"brank.as/petnet/serviceutil/logging"
)

// GetTransactionList ...
func (s *MICoreSvc) GetTransactionList(ctx context.Context, req *migunk.GetTransactionListRequest) (*migunk.TransactionListResult, error) {
	log := logging.FromContext(ctx)
	dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "invalid date from")
	}

	dateTo, err := time.Parse("2006-01-02", req.DateTo)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "invalid date to")
	}

	res, err := s.storage.ListMicroInsuranceHistory(ctx, storage.MicroInsuranceFilter{
		From:         dateFrom,
		Until:        dateTo,
		DsaID:        phmw.GetDSAOrgID(ctx),
		OrgID:        req.OrgID,
		SortOrder:    storage.SortOrder(req.GetSortOrder().String()),
		SortByColumn: storage.MicroInsuranceSort(req.GetSortByColumn().String()),
	})
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, coreerror.NewCoreError(codes.NotFound, "no transactions found")
		}

		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return toTransactionListResult(res, log), nil
}

func toTransactionListResult(res []storage.MicroInsuranceHistory, log *logrus.Entry) *migunk.TransactionListResult {
	list := make([]*migunk.InsuranceTransaction, 0, len(res))
	if len(res) == 0 {
		return &migunk.TransactionListResult{
			Transactions: list,
		}
	}

	var err error
	for _, h := range res {
		if h.InsuranceDetails == nil {
			log.WithField("id", h.ID).Error("insurance details empty")
			continue
		}

		var ins *migunk.Insurance
		err = json.Unmarshal(h.InsuranceDetails, &ins)
		if err != nil {
			logging.WithError(err, log).WithField("insurance db string", string(h.InsuranceDetails)).
				Error("unmarshalling insurance details")
			continue
		}

		var beneficiaries []*migunk.Person
		err = json.Unmarshal(h.Beneficiaries, &beneficiaries)
		if err != nil {
			logging.WithError(err, log).WithField("insurance db string", string(h.InsuranceDetails)).
				Error("unmarshalling beneficiaries details")
		}

		var dependents []*migunk.Person
		err = json.Unmarshal(h.Beneficiaries, &beneficiaries)
		if err != nil {
			logging.WithError(err, log).WithField("insurance db string", string(h.InsuranceDetails)).
				Error("unmarshalling dependents details")
		}

		ttamt := &migunk.Amount{
			Amount:   strconv.FormatFloat(ins.TrnAmount, 'f', -1, 64),
			Currency: "PHP",
		}
		trndate, err := time.Parse("01/02/2006", ins.TrnDate)
		if err != nil {
			logging.WithError(err, log).Error("cant parse TrnDate")
		}

		trans := &migunk.InsuranceTransaction{
			TrnDate:        timestamppb.New(trndate),
			TraceNumber:    ins.TraceNumber,
			ClientNo:       ins.ClientNo,
			LastName:       h.LastName,
			FirstName:      h.FirstName,
			MiddleName:     h.MiddleName,
			Gender:         h.Gender,
			BirthDate:      h.Birthdate.Format("2006-01-02"),
			MobileNumber:   h.MobileNumber,
			MaritalStatus:  h.MaritalStatus,
			Occupation:     h.Occupation,
			InsGroupID:     "",
			InsProductID:   ins.InsProductID,
			InsProductDesc: ins.InsProductDesc,
			InsurerCode:    "",
			InsuranceType:  "",
			BegPolicyNo:    ins.BegPolicyNo,
			EndPolicyNo:    ins.EndPolicyNo,
			CoverageInMos:  0,
			EffectiveDate:  ins.EffectiveDate,
			ExpiryDate:     ins.ExpiryDate,
			InsCardNo:      "",
			Beneficiaries:  beneficiaries,
			Dependents:     dependents,
			NumUnits:       ins.NumUnits,
			TotAmt:         ttamt,
			TrnStatus:      ins.StatusDesc,
			ProvinceCode:   h.ProvinceCode,
			CityCode:       h.CityCode,
			Address:        h.Address,
		}

		list = append(list, trans)
	}
	return &migunk.TransactionListResult{
		Transactions: list,
	}
}
