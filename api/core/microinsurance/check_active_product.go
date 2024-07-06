package microinsurance

import (
	"context"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/microinsurance"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

// CheckActiveProduct ...
func (s *MICoreSvc) CheckActiveProduct(ctx context.Context, req *migunk.CheckActiveProductRequest) (*migunk.ActiveProduct, error) {
	res, err := s.cl.CheckActiveProduct(ctx, &microinsurance.CheckActiveProductRequest{
		LastName:    req.LastName,
		FirstName:   req.FirstName,
		MiddleName:  req.MiddleName,
		Birthdate:   req.Birthdate,
		Gender:      req.Gender,
		ProductCode: req.ProductCode,
	})
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	return toActiveProduct(res), nil
}

func toActiveProduct(p *microinsurance.ActiveProduct) *migunk.ActiveProduct {
	if p == nil {
		return nil
	}

	deps, _ := p.Dependents.Int64()
	benf, _ := p.Beneficiary.Int64()

	return &migunk.ActiveProduct{
		ProductName:          p.ProductName,
		ProductCode:          p.ProductCode,
		ProductType:          p.ProductType,
		Dependents:           int32(deps),
		Beneficiary:          int32(benf),
		BeneficiaryPolicy:    toMinMax(p.BeneficiaryPolicy),
		DependentsPolicy:     toMinMax(p.DependentsPolicy),
		AgePolicy:            toAgePolicy(p.AgePolicy),
		EndSpielsTitle:       p.EndSpielsTitle,
		EndSpielsDescription: p.EndSpielsDescription,
		SalesPitch:           p.SalesPitch,
		TermsAndCondition:    p.TermsAndCondition,
		DataPrivacy:          p.DataPrivacy,
	}
}
