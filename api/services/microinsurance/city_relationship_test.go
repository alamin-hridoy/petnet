package microinsurance

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMicroInsuranceCityRelationship(t *testing.T) {
	st := newTestStorage(t)
	s := newTestSvc(t, st)

	for _, tc := range []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "GetRelationships",
			testFunc: func(t *testing.T) {
				res, err := s.GetRelationships(context.TODO(), nil)

				require.Nil(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Relationships)
			},
		},
		{
			name: "GetAllCities",
			testFunc: func(t *testing.T) {
				res, err := s.GetAllCities(context.TODO(), nil)

				require.Nil(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Cities)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}
