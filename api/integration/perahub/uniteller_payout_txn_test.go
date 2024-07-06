package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUNTPayout(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := UNTPayoutRequest{
		ClientReferenceNo:  "7884447474",
		ReferenceNumber:    "848jjfu23u333",
		LocationCode:       "PERAH-00001",
		LocationID:         "371",
		LocationName:       "INFORMATION TECHONOLOGY",
		Gender:             "M",
		ControlNumber:      "5142096205",
		Currency:           "PHP",
		PrincipalAmount:    "3673.5",
		IDNumber:           "B83180608851",
		IDType:             "123456",
		IDIssuedBY:         "PH",
		IDDOIssue:          "2016-05-24",
		IDExpDate:          "2023-05-24",
		ContactNumber:      "09516738640",
		Address:            "MAIN ST",
		City:               "AKLAN CITY",
		Province:           "ILOCOS SUR",
		Country:            "PH",
		ZipCode:            "36989",
		State:              "PH-00",
		Nationality:        "PH",
		BirthDate:          "1995-04-14",
		BirthCountry:       "PH",
		Occupation:         "OTH",
		UserID:             5500,
		TrxDate:            "2021-06-03",
		CustomerID:         "6925594",
		CurrencyID:         "1",
		RemcoID:            "20",
		TrxType:            "2",
		IsDomestic:         "1",
		CustomerName:       "Levy Robert Sogocio",
		ServiceCharge:      "0",
		RemoteLocationID:   "371",
		DstAmount:          "0",
		TotalAmount:        "4898",
		BuyBackAmount:      "0",
		McRateId:           "0",
		McRate:             "0",
		RemoteIPAddress:    "130.211.2.187",
		RemoteUserID:       5684,
		OriginatingCountry: "Philippines",
		DestinationCountry: "PH",
		PurposeTransaction: "Family Support/Living Expenses",
		SourceFund:         "Savings",
		RelationTo:         "Family",
		BirthPlace:         "TAGUDIN,ILOCOS SUR",
		Barangay:           "TALLAOEN",
		RiskScore:          "0",
		RiskCriteria:       "0",
		FormType:           "0",
		FormNumber:         "0",
		PayoutType:         "1",
		SenderName:         "Reyes, Ana,",
		ReceiverName:       "Bayuga, Mary Monica",
		SendFName:          "Ana",
		SendMName:          "null",
		SendLName:          "Reyes",
		RecFName:           "Mary Monica",
		RecMName:           "null",
		RecLName:           "Bayuga",
		DeviceID:           "dc29d6674a776db145af78f5ac20a293409a6c1f807885bbb5",
		AgentID:            "84424911",
		AgentCode:          "TIS",
		OrderNumber:        "CA642308560",
		IPAddress:          "130.211.2.187",
		RateCategory:       "0",
		DsaCode:            "TEST_DSA",
		DsaTrxType:         "digital",
	}

	tests := []struct {
		name        string
		in          UNTPayoutRequest
		expectedReq UNTPayoutRequest
		want        *UNTPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &UNTPayoutResponseBody{
				Code:    "00000000",
				Message: "Success",
				Result: UNTPayoutResult{
					ResponseCode:       "00000000",
					ControlNumber:      "CTRL1",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					CreationDate:       "2021-06-15T19:10:16.000-0400",
					ReceiverName:       "RAMSES COMODO",
					Address:            "TEST",
					City:               "MORONG",
					Country:            "PH",
					SenderName:         "FERDINAND CORTES",
					ZipCode:            "36978",
					OriginatingCountry: "US",
					DestinationCountry: "PH",
					ContactNumber:      "1540254852",
					FmtSenderName:      "CORTES, FERDINAND, ",
					FmtReceiverName:    "COMODO, RAMSES, ",
				},
				RemcoID: "20",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.UNTPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("UNTPayout() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq UNTPayoutRequest
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
