package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBCBillerInfo(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BCBillerInfoRequest
		want    *BCBillerInfoResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: BCBillerInfoRequest{
				Code:       "MWCOMs",
				UserID:     "5500",
				LocationID: "371",
			},
			want: &BCBillerInfoResponse{
				Code:    200,
				Message: "Success",
				Result: BCBillerInfoResult{
					Code:            "MWCOM",
					IsCde:           0,
					IsAsync:         0,
					Name:            "MANILA WATER",
					Description:     "Manila Water Company",
					Logo:            "https://stg-bc-api-images.s3-ap-southeast-1.amazonaws.com/biller-logos/250/MWCOM.png",
					Category:        "Water",
					Type:            "Batch",
					IsMultipleBills: 0,
					Parameters: Parameters{
						Verify: []Verify{
							{
								ReferenceNumber{
									Description: "Account Number",
									Rules: BillerInfoRules{
										Digits8: BCCM{
											Message: "The account number must be 8 digits.",
											Code:    5,
										},
										Required: BCCM{
											Message: "Please provide the account number.",
											Code:    4,
										},
									},
									Label: "Account Number",
								},
							},
						},
						Transact: []Transact{
							{
								ClientReference{
									Description: "Client unique transaction reference number",
									Rules: CRules{
										AlphaDash: BCCM{
											Message: "Please make sure that the client reference number is in alpha dash format.",
											Code:    9,
										},
										Required: BCCM{
											Message: "Please provide the client reference number.",
											Code:    4,
										},
										UniqueCrn: BCCM{
											Message: "This client reference number already exists.",
											Code:    11,
										},
									},
									Label: "Client Reference Number",
								},
							},
						},
					},
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
			got, err := s.BCBillerInfo(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BCBillerInfo() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}
