package microinsurance

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

func TestMicroInsuranceProducts(t *testing.T) {
	st := newTestStorage(t)
	s := newTestSvc(t, st)

	for _, tc := range []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "GetProduct",
			testFunc: func(t *testing.T) {
				prodCode := "TEST123"
				prodRes, err := s.GetProduct(context.TODO(), &migunk.GetProductRequest{
					ProductCode: prodCode,
				})

				require.Nil(t, err)
				require.NotNil(t, prodRes)
				require.NotNil(t, prodRes.Product)
				require.Equal(t, prodRes.Product.InsProductID, prodCode)
			},
		},
		{
			name: "GetOfferProduct",
			testFunc: func(t *testing.T) {
				offProdRes, err := s.GetOfferProduct(context.TODO(), &migunk.GetOfferProductRequest{
					LastName:   "Test",
					FirstName:  "Test",
					MiddleName: "test",
					Birthdate:  "2006-01-02",
					Gender:     "M",
					TrxType:    1,
					Amount:     200,
				})

				require.Nil(t, err)
				require.NotNil(t, offProdRes)
				require.NotNil(t, offProdRes.AgePolicy)
			},
		},
		{
			name: "GetOfferProduct",
			testFunc: func(t *testing.T) {
				prodCode := "TEST123"
				prodRes, err := s.CheckActiveProduct(context.TODO(), &migunk.CheckActiveProductRequest{
					LastName:    "Test",
					FirstName:   "Test",
					MiddleName:  "test",
					Birthdate:   "2006-01-02",
					Gender:      "M",
					ProductCode: prodCode,
				})

				require.Nil(t, err)
				require.NotNil(t, prodRes)
				require.Equal(t, prodCode, prodRes.ProductCode)
			},
		},
		{
			name: "GetProductList",
			testFunc: func(t *testing.T) {
				prodRes, err := s.GetProductList(context.TODO(), nil)

				require.Nil(t, err)
				require.NotNil(t, prodRes)
				require.NotEmpty(t, prodRes.Products)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}
