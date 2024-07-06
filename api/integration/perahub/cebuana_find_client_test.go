package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCebFindClient(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          CebFindClientRequest
		expectedReq CebFindClientRequest
		want        *CebFindClientRespBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: CebFindClientRequest{
				FirstName:    "Gremar",
				LastName:     "Gremar",
				BirthDate:    "",
				ClientNumber: "",
			},
			expectedReq: CebFindClientRequest{
				FirstName:    "Gremar",
				LastName:     "Gremar",
				BirthDate:    "",
				ClientNumber: "",
			},
			want: &CebFindClientRespBody{
				Code:    0,
				Message: "Successful",
				Result: GUResult{
					Client: Client{
						ClientID:     3663,
						ClientNumber: "EWFHL0000070155824",
						FirstName:    "GREMAR",
						MiddleName:   "REAS",
						LastName:     "NAPARATE",
						BirthDate:    "1994-11-04T00:00:00",
						CPCountry: CrtyID{
							CountryID: 0,
						},
						TPCountry: CrtyID{
							CountryID: 0,
						},
						CtryAddress: CrtyID{
							CountryID: 0,
						},
						CSOfFund: CSOfFund{
							SourceOfFundID: 0,
						},
					},
				},
				RemcoID: 9,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.CebFindClient(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("CebFindClient() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq CebFindClientRequest
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
