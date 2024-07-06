package microinsurance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/microinsurance"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

// Transact ...
func (s *MICoreSvc) Transact(ctx context.Context, req *migunk.TransactRequest) (*migunk.Insurance, error) {
	log := logging.FromContext(ctx)

	var (
		err error
		ins *migunk.Insurance
	)

	defer func() {
		_, recErr := util.RecordMicroInsurance(ctx, s.storage, req, ins, err)
		if recErr != nil {
			log.Error(recErr)
		}
	}()
	res, err := s.cl.Transact(ctx, toTransactRequest(req))
	if err != nil {
		logging.WithError(err, log).Error("microinsurance transact error")
		return nil, coreerror.ToCoreError(err)
	}

	ins = toInsurance(res)
	if ins == nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return ins, nil
}

func toTransactRequest(req *migunk.TransactRequest) *microinsurance.TransactRequest {
	trReq := &microinsurance.TransactRequest{
		Coy:              req.Coy,
		LocationID:       req.LocationID,
		UserCode:         req.UserCode,
		TrxDate:          req.TrxDate,
		PromoAmount:      json.Number(fmt.Sprintf("%f", req.PromoAmount)),
		PromoCode:        req.PromoCode,
		Amount:           req.Amount,
		CoverageCount:    req.CoverageCount,
		ProductCode:      req.ProductCode,
		ProcessingBranch: req.ProcessingBranch,
		ProcessedBy:      req.ProcessedBy,
		UserEmail:        req.UserEmail,
		LastName:         req.LastName,
		FirstName:        req.FirstName,
		MiddleName:       req.MiddleName,
		Gender:           req.Gender,
		Birthdate:        req.Birthdate,
		MobileNumber:     req.MobileNumber,
		ProvinceCode:     req.ProvinceCode,
		CityCode:         req.CityCode,
		Address:          req.Address,
		MaritalStatus:    req.MaritalStatus,
		Occupation:       req.Occupation,
		CardNumber:       req.CardNumber,
		NumberUnits:      req.NumberUnits,
		Dependents:       toDependents(req.Dependents),
	}

	addBeneficiaries(trReq, req.Beneficiaries)

	return trReq
}

func addBeneficiaries(r *microinsurance.TransactRequest, benArr []*migunk.Person) {
	if len(benArr) == 0 {
		return
	}

	cnt := 1
	for _, ben := range benArr {
		if ben == nil {
			continue
		}

		switch cnt {
		case 1:
			r.Ben1FirstName = ben.FirstName
			r.Ben1LastName = ben.LastName
			r.Ben1MiddleName = ben.MiddleName
			r.Ben1NoMiddleName = ben.MiddleName == ""
			r.Ben1ContactNumber = ben.ContactNumber
			r.Ben1Relationship = ben.Relationship
		case 2:
			r.Ben2FirstName = ben.FirstName
			r.Ben2LastName = ben.LastName
			r.Ben2MiddleName = ben.MiddleName
			r.Ben2NoMiddleName = ben.MiddleName == ""
			r.Ben2ContactNumber = ben.ContactNumber
			r.Ben2Relationship = ben.Relationship
		case 3:
			r.Ben3FirstName = ben.FirstName
			r.Ben3LastName = ben.LastName
			r.Ben3MiddleName = ben.MiddleName
			r.Ben3NoMiddleName = ben.MiddleName == ""
			r.Ben3ContactNumber = ben.ContactNumber
			r.Ben3Relationship = ben.Relationship
		case 4:
			r.Ben4FirstName = ben.FirstName
			r.Ben4LastName = ben.LastName
			r.Ben4MiddleName = ben.MiddleName
			r.Ben4NoMiddleName = ben.MiddleName == ""
			r.Ben4ContactNumber = ben.ContactNumber
			r.Ben4Relationship = ben.Relationship
		default:
			break
		}

		cnt++
	}
}

func toDependents(depArr []*migunk.Person) []microinsurance.Dependent {
	deps := make([]microinsurance.Dependent, 0, len(depArr))
	for _, d := range depArr {
		if d == nil {
			continue
		}

		deps = append(deps, microinsurance.Dependent{
			LastName:      d.FirstName,
			FirstName:     d.LastName,
			MiddleName:    d.MiddleName,
			NoMiddleName:  d.MiddleName == "",
			ContactNumber: d.ContactNumber,
			BirthDate:     d.BirthDate,
			Relationship:  d.Relationship,
		})
	}

	return deps
}

func toInsurance(i *microinsurance.Insurance) *migunk.Insurance {
	if i == nil {
		return nil
	}

	trnAmt, _ := i.TrnAmount.Float64()
	noUints, _ := i.NumUnits.Int64()
	partnerComm, _ := i.PartnerCommission.Float64()
	tellerComm, _ := i.TellerCommission.Float64()
	timeStamp, err := time.Parse("2006-01-02 15:04:05", i.Timestamp)
	if err == nil {
		timeStamp = time.Now()
	}

	return &migunk.Insurance{
		SessionID:         i.SessionID,
		StatusCode:        i.StatusCode,
		StatusDesc:        i.StatusDesc,
		InsProductID:      i.InsProductID,
		InsProductDesc:    i.InsProductDesc,
		TrnDate:           i.TrnDate,
		TrnAmount:         trnAmt,
		TraceNumber:       i.TraceNumber,
		ClientNo:          i.ClientNo,
		NumUnits:          int32(noUints),
		BegPolicyNo:       i.BegPolicyNo,
		EndPolicyNo:       i.EndPolicyNo,
		EffectiveDate:     i.EffectiveDate,
		ExpiryDate:        i.ExpiryDate,
		PocPDFLink:        i.ExpiryDate,
		CocPDFLink:        i.CocPDFLink,
		PartnerCommission: partnerComm,
		TellerCommission:  tellerComm,
		Timestamp:         timestamppb.New(timeStamp),
	}
}
