package perahub

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCicoValidate(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      CicoValidateRequest
		want    *CicoValidateResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: CicoValidateRequest{
				PartnerCode: "",
				Trx: CicoValidateTrx{
					Provider:        "GCASH",
					ReferenceNumber: "09654767706",
					TrxType:         "Cash In",
					PrincipalAmount: 10,
				},
				Customer: CicoValidateCustomer{
					CustomerID:        "7085257",
					CustomerFirstname: "Jerica",
					CustomerLastname:  "Naparate",
					CurrAddress:       "1953 Ph 3b Blk 6 Lot 9",
					CurrBarangay:      "Barangay 175",
					CurrCity:          "CALOOCAN CITY",
					CurrProvince:      "METRO MANILA",
					CurrCountry:       "Philippines",
					BirthDate:         "1996-08-10",
					BirthPlace:        "MANILA , METRO MANILA",
					BirthCountry:      "Philippines",
					ContactNo:         "09089488896",
					IDType:            "Passport",
					IDNumber:          "IDTest1233453",
				},
			},
			want: &CicoValidateResponse{
				Code:    200,
				Message: "Successful",
				Result: &CicoValidateResult{
					PetnetTrackingno:   "5a269417e107691f3d7c",
					TrxDate:            "2022-05-17",
					TrxType:            "Cash In",
					Provider:           "GCASH",
					ProviderTrackingno: "7000001521345",
					ReferenceNumber:    "09654767706",
					PrincipalAmount:    10,
					Charges:            0,
					TotalAmount:        10,
					Timestamp:          "",
				},
			},
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(CicoValidateResult{}, "Timestamp"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.CicoValidate(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("CicoValidate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}
