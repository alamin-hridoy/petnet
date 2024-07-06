package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCebFindBF(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          CebFindBFReq
		expectedReq CebFindBFReq
		want        *CebFindBFRes
		wantErr     bool
	}{
		{
			name: "Success",
			in: CebFindBFReq{
				SenderClientId: "5632",
			},
			expectedReq: CebFindBFReq{
				SenderClientId: "5632",
			},
			want: &CebFindBFRes{
				Code:    "0",
				Message: "Successful",
				Result: FBFResult{
					Beneficiary: []Beneficiary{
						{
							BeneficiaryID:  "3663",
							FirstName:      "GREMAR",
							MiddleName:     "REAS",
							LastName:       "NAPARATE",
							BirthDate:      "1994-11-04T00:00:00",
							StateIDAddress: "0",
							CPCountry: CrtyID{
								CountryID: 0,
							},
							TPCountry: CrtyID{
								CountryID: 0,
							},
							CtryAddress: CrtyID{
								CountryID: 0,
							},
							BirthCountry: CrtyID{
								CountryID: 0,
							},
						},
					},
				},
				RemcoID: "9",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.CebFindBF(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("CebFindBF() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq CebFindBFReq
			if err := json.Unmarshal(m.GetMockRequest(), &newReq); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.expectedReq, newReq) {
				t.Error(cmp.Diff(test.expectedReq, newReq))
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}
