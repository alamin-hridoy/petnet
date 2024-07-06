package microinsurance

import (
	"context"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/microinsurance"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

// GetProduct ...
func (s *MICoreSvc) GetProduct(ctx context.Context, req *migunk.GetProductRequest) (*migunk.ProductResult, error) {
	if req == nil || req.ProductCode == "" {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "product code is required")
	}

	res, err := s.cl.GetProduct(ctx, &microinsurance.GetProductRequest{
		ProductCode: req.ProductCode,
	})
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	if res == nil || res.Product == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	coverages := make([]*migunk.Coverage, 0, len(res.Coverages))
	for _, c := range res.Coverages {
		coverages = append(coverages, toCoverage(c))
	}

	minAge, _ := res.Product.MinAge.Int64()
	maxAge, _ := res.Product.MaxAge.Int64()
	return &migunk.ProductResult{
		SessionID:  res.SessionID,
		StatusCode: res.StatusCode,
		StatusDesc: res.StatusDesc,
		Product: &migunk.InsuranceProduct{
			InsGroupID:          res.Product.InsGroupID,
			InsProductID:        res.Product.InsProductID,
			InsProductDesc:      res.Product.InsProductDesc,
			InsurerCode:         res.Product.InsurerCode,
			InsuranceType:       res.Product.InsuranceType,
			InsuranceCategory:   res.Product.InsuranceCategory,
			PolicyNo:            res.Product.PolicyNo,
			MinAge:              int32(minAge),
			MaxAge:              int32(maxAge),
			CoverageInMos:       res.Product.CoverageInMos.String(),
			ContestAbilityInMos: res.Product.ContestAbilityInMos.String(),
			ActivationDelay:     res.Product.ActivationDelay.String(),
			MaxUnits:            res.Product.MaxUnits.String(),
			PerUnitFee:          res.Product.PerUnitFee.String(),
			ProductName:         res.Product.ProductName,
		},
		Coverages: coverages,
	}, nil
}

func toCoverage(c microinsurance.Coverage) *migunk.Coverage {
	return &migunk.Coverage{
		InsCoverageID:     c.InsCoverageID,
		InsCoverageDesc:   c.InsCoverageDesc,
		InsCoverageIconID: c.InsCoverageIconID,
		InsCoverageType1:  c.InsCoverageType1,
		InsCoverageAmt1:   c.InsCoverageAmt1,
		InsCoverageType2:  c.InsCoverageType2,
		InsCoverageAmt2:   c.InsCoverageAmt2,
		InsCoverageType3:  c.InsCoverageType3,
		InsCoverageAmt3:   c.InsCoverageAmt3,
		InsCoverageType4:  c.InsCoverageType4,
		InsCoverageAmt4:   c.InsCoverageAmt4,
		InsCoverageType5:  c.InsCoverageType5,
		InsCoverageAmt5:   c.InsCoverageAmt5,
	}
}
