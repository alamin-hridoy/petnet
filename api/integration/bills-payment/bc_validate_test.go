package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBCValidate(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BCValidateRequest
		want    *BCValidateResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: BCValidateRequest{
				BillPartnerID: 2,
				BillerTag:     "ADMSN",
				Code:          "ADMSN",
				AccountNumber: "200820248",
				AccountNo:     "200820248",
				Identifier:    "Sample",
				PaymentMethod: "CASH",
				OtherCharges:  "0.00",
				Amount:        "1000.00",
				OtherInfo: BCOtherInfo{
					LastName:        "Serato",
					FirstName:       "Mike Edward",
					MiddleName:      "Secret",
					PaymentType:     "B2",
					Course:          "BSCpE",
					TotalAssessment: "0.00",
					SchoolYear:      "2021-2022",
					Term:            "1",
				},
			},
			want: &BCValidateResponse{
				Code:    200,
				Message: "Success",
				Result: BCValidateResult{
					Valid:            true,
					Code:             0,
					Account:          "200820248",
					Details:          []interface{}{},
					ValidationNumber: "48f2b647-0ff8-4ace-9928-346258f08df5",
				},
				RemcoID: 2,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BCValidate(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BCValidate() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}
