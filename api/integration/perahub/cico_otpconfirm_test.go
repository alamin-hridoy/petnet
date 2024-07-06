package perahub

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCicoOTPConfirm(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      CicoOTPConfirmRequest
		want    *CicoOTPConfirmResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: CicoOTPConfirmRequest{
				PartnerCode:      "DSA",
				PetnetTrackingno: "8340b7bdb171cbdf6350",
				TrxDate:          "2022-06-03",
				OTP:              "123456",
				OTPPayload: OTPpayload{
					CommandID: 27897,
					Payload:   "yZJdFQ+TDbjuSr3ZeLAgFj+S3NnYwkYwRyp1RxZKR20B36kH5Fr52y9I+tg+y5TGrUHRxQAjdVE6rWgJoV1VoLJK0SBHSVNuOts1fRSsXdIPJuxt6v/auPm0gZqyaUXWS+Dtl2OiVpbBPtNB2H6v+bbs7ldWzIDa+47EsWUnUEuVeq8nMM4TPKU0zILbf4lXv6dr2EZCmTX1eNvnyK44QaNBgN68Jb1i50PoD6Gqb61T9CS28btwgIjTlZ/U0s9by4Q8MBmfsEYSejlHpypj/nt0/v9+o8zG9r1kefGt3h4vgtkH3QNoaC7YdWMsAuPJTo1VJTv4ufufBdWD+E+NGwKliGpvUvn4OZ504kyFa5ltLytebihUND70r5S7aI4aYSpEQw==",
				},
			},
			want: &CicoOTPConfirmResponse{
				Code:    200,
				Message: "SUCCESS TRANSACTION.",
				Result: &CicoOTPConfirmResult{
					PartnerCode:        "DSA",
					Provider:           "DiskarTech",
					PetnetTrackingno:   "8340b7bdb171cbdf6350",
					TrxDate:            "2022-06-03",
					TrxType:            "Cash Out",
					ProviderTrackingno: "",
					ReferenceNumber:    "220603-000003-1",
					PrincipalAmount:    200,
					Charges:            0,
					TotalAmount:        200,
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.CicoOTPConfirm(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("CicoOTPConfirm() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}
