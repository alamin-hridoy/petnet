package microinsurance

import (
	"context"
	"encoding/json"
	"fmt"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/microinsurance"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

// GetOfferProduct ...
func (s *MICoreSvc) GetOfferProduct(ctx context.Context, req *migunk.GetOfferProductRequest) (*migunk.OfferProduct, error) {
	res, err := s.cl.GetOfferProduct(ctx, &microinsurance.GetOfferProductRequest{
		LastName:   req.LastName,
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		Birthdate:  req.Birthdate,
		Gender:     req.Gender,
		TrxType:    json.Number(fmt.Sprintf("%d", req.TrxType)),
		Amount:     json.Number(fmt.Sprintf("%f", req.Amount)),
	})
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	deps, _ := res.Dependents.Int64()
	benf, _ := res.Beneficiary.Int64()

	return &migunk.OfferProduct{
		ProductName:          res.ProductName,
		ProductCode:          res.ProductCode,
		ProductType:          res.ProductType,
		Dependents:           int32(deps),
		Beneficiary:          int32(benf),
		BeneficiaryPolicy:    toMinMax(res.BeneficiaryPolicy),
		AgePolicy:            toAgePolicy(res.AgePolicy),
		EndSpielsTitle:       res.EndSpielsTitle,
		EndSpielsDescription: res.EndSpielsDescription,
		SalesPitch:           res.SalesPitch,
		TermsAndCondition:    res.TermsAndCondition,
		DataPrivacy:          res.DataPrivacy,
	}, nil
}

func toAgePolicy(ap *microinsurance.AgePolicy) *migunk.AgePolicy {
	if ap == nil {
		return nil
	}

	return &migunk.AgePolicy{
		Insurer:    toMinMaxAge(ap.Insurer),
		Dependents: toDependentsPolicy(ap.Dependents),
	}
}

func toDependentsPolicy(d *microinsurance.DependentsPolicy) *migunk.DependentsPolicy {
	if d == nil {
		return nil
	}

	return &migunk.DependentsPolicy{
		Children: toMinMaxAge(d.Children),
		Parents:  toMinMaxAge(d.Parents),
		Siblings: toMinMaxAge(d.Siblings),
		Spouse:   toMinMaxAge(d.Spouse),
	}
}

func toMinMax(mm *microinsurance.MinMax) *migunk.MinMax {
	if mm == nil {
		return nil
	}

	max, _ := mm.Max.Int64()
	min, _ := mm.Min.Int64()

	return &migunk.MinMax{
		Max: int32(max),
		Min: int32(min),
	}
}

func toMinMaxAge(ma *microinsurance.MinMaxAge) *migunk.MinMaxAge {
	if ma == nil {
		return nil
	}

	max, _ := ma.MaxAge.Int64()
	min, _ := ma.MinAge.Int64()

	return &migunk.MinMaxAge{
		MaxAge: int32(max),
		MinAge: int32(min),
	}
}
