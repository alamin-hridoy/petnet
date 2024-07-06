package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemitanceValidateSendMoney(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := RemitanceValidateSendMoneyReq{
		PartnerReferenceNumber: "12JEKDWW213",
		PrincipalAmount:        "10000",
		ServiceFee:             "50",
		IsoCurrency:            "PHP",
		ConversionRate:         "1",
		IsoOriginatingCountry:  "PHP",
		IsoDestinationCountry:  "PHP",
		SenderLastName:         "HERMO",
		SenderFirstName:        "IRENE",
		SenderMiddleName:       "M",
		ReceiverLastName:       "HERMO",
		ReceiverFirstName:      "SONNY",
		ReceiverMiddleName:     "D",
		SenderBirthDate:        "1981-06-12",
		SenderBirthPlace:       "TARLAC",
		SenderBirthCountry:     "PH",
		SenderGender:           "FEMALE",
		SenderRelationship:     "SPOUSE",
		SenderPurpose:          "GIFT",
		SenderOfFund:           "SALARY",
		SenderOccupation:       "DOCTOR",
		SenderEmploymentNature: "IT",
		SendPartnerCode:        "USP",
		SenderSourceOfFund:     "SALARY",
	}
	tests := []struct {
		name        string
		in          RemitanceValidateSendMoneyReq
		expectedReq RemitanceValidateSendMoneyReq
		want        *RemitanceValidateSendMoneyRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RemitanceValidateSendMoneyRes{
				Code:    200,
				Message: "Good",
				Result: RemitanceValidateSendMoneyResult{
					SendValidateReferenceNumber: "1653296685161",
				},
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.RemitanceValidateSendMoney(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("RemitanceValidateSendMoney() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RemitanceValidateSendMoneyReq
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
