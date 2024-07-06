package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCebAddBf(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          CebAddBftReq
		expectedReq CebAddBftReq
		want        *CebAddBfResp
		wantErr     bool
	}{
		{
			name: "Success",
			in: CebAddBftReq{
				FirstName:          "Newwiee",
				MiddleName:         "Tay",
				LastName:           "Tawan",
				SenderClientID:     5632,
				BirthDate:          "1994-01-30",
				CellphoneCountryID: "63",
				ContactNumber:      "",
				TelephoneCountryID: "0",
				TelephoneAreaCode:  "",
				TelephoneNumber:    "",
				CountryAddressID:   "166",
				BirthCountryID:     "166",
				ProvinceAddress:    "",
				Address:            "1953 Ph 3b Block 6 Lot 9 Camarin Caloocan City",
				UserID:             5632,
				Occupation:         "",
				ZipCode:            "1850",
				StateIDAddress:     "0",
				Tin:                "",
			},
			expectedReq: CebAddBftReq{
				FirstName:          "Newwiee",
				MiddleName:         "Tay",
				LastName:           "Tawan",
				SenderClientID:     5632,
				BirthDate:          "1994-01-30",
				CellphoneCountryID: "63",
				ContactNumber:      "",
				TelephoneCountryID: "0",
				TelephoneAreaCode:  "",
				TelephoneNumber:    "",
				CountryAddressID:   "166",
				BirthCountryID:     "166",
				ProvinceAddress:    "",
				Address:            "1953 Ph 3b Block 6 Lot 9 Camarin Caloocan City",
				UserID:             5632,
				Occupation:         "",
				ZipCode:            "1850",
				StateIDAddress:     "0",
				Tin:                "",
			},
			want: &CebAddBfResp{
				Code:    0,
				Message: "Successful",
				Result: ABResult{
					ResultStatus:  "Successful",
					MessageID:     0,
					LogID:         0,
					BeneficiaryID: 8595,
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
			got, err := s.CebAddBf(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("CebAddBf() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq CebAddBftReq
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
