package microinsurance

import (
	"context"

	coreerror "brank.as/petnet/api/core/error"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

// GetProductList ...
func (s *MICoreSvc) GetProductList(ctx context.Context) (*migunk.GetProductListResult, error) {
	res, err := s.cl.GetProductList(ctx)
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	prods := make([]*migunk.ActiveProduct, 0, len(res))
	for _, p := range res {
		prods = append(prods, toActiveProduct(&p))
	}

	return &migunk.GetProductListResult{
		Products: prods,
	}, nil
}
