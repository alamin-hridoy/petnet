package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUSSCPayout(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := USSCPayoutRequest{
		LocationID:         "191",
		LocationName:       "MALOLOS",
		UserID:             5188,
		TrxDate:            "2021-10-06",
		CustomerID:         "6276847",
		CurrencyID:         "1",
		RemcoID:            "10",
		TrxType:            "2",
		IsDomestic:         "1",
		CustomerName:       "ALVIN CAMITAN",
		ServiceCharge:      "1",
		RemoteLocationID:   "191",
		DstAmount:          "0",
		TotalAmount:        "101",
		RemoteUserID:       "5188",
		RemoteIPAddress:    "130.211.2.203",
		OriginatingCountry: "Philippines",
		DestinationCountry: "PH",
		PurposeTransaction: "Family Support/Living Expenses",
		SourceFund:         "Savings",
		Occupation:         "Unemployed",
		RelationTo:         "Family",
		BirthDate:          "1995-09-06",
		BirthPlace:         "ERMITA,BOHOL",
		BirthCountry:       "Philippines",
		IDType:             "445",
		IDNumber:           "ASD292929291",
		Address:            "18 SITIO PULO",
		Barangay:           "BULIHAN",
		City:               "MALOLOS",
		Province:           "BULACAN",
		ZipCode:            "3000A",
		Country:            "Philippines",
		ContactNumber:      "09065595959",
		CurrentAddress: NonexAddress{
			Address1: "Block 1 Lot 12",
			Address2: "Placid Homes",
			Barangay: "BULIHAN",
			City:     "ANDA",
			Province: "BOHOL",
			ZipCode:  "1000A",
			Country:  "Philippines",
		},
		PermanentAddress: NonexAddress{
			Address1: "Block 1 Lot 12",
			Address2: "Placid Homes",
			Barangay: "BULIHAN",
			City:     "ANDA",
			Province: "BOHOL",
			ZipCode:  "1000A",
			Country:  "Philippines",
		},
		RiskScore:       "1",
		PayoutType:      "1",
		SenderName:      "Mercado, Marites Cueto",
		ReceiverName:    "Cortez, Fernando",
		PrincipalAmount: "3300",
		ControlNumber:   "2107130913218147",
		ReferenceNumber: "21181472909563000013",
		BranchCode:      "branch1",
		DsaCode:         "TEST_DSA",
		DsaTrxType:      "digital",
	}

	tests := []struct {
		name        string
		in          USSCPayoutRequest
		expectedReq USSCPayoutRequest
		want        *USSCPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &USSCPayoutResponseBody{
				Code:    "000000",
				Message: "Successfully Payout!",
				Result: USSCPayResult{
					Spcn:            "SP697342605141",
					SendPk:          "",
					SendPw:          "",
					SendDate:        "0",
					SendJournalNo:   "0",
					SendLastName:    "Doe",
					SendFirstName:   "John",
					SendMiddleName:  "Michael",
					PayAmount:       "1000.00",
					SendFee:         "0.00",
					SendVat:         "0.00",
					SendFeeAfterVat: "0.00",
					SendTotalAmount: "0.00",
					PayPk:           "",
					PayPw:           "",
					PayLastName:     "CAMITAN",
					PayFirstName:    "ALVIN",
					PayMiddleName:   "",
					Relationship:    "",
					Purpose:         "",
					PromoCode:       "",
					PayBranchCode:   "",
					Remarks:         "",
					OrNo:            "",
					OboBranchCode:   "",
					OboUserID:       "",
					Message:         "0000: ACCEPTED - PAID OUT SUCCESSFULLY",
					Code:            "0",
					NewScreen:       "0",
					JournalNo:       "011330407",
					ProcessDate:     "20200703",
					ReferenceNo:     "1",
				},
				RemcoID: 10,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.USSCPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("USSCPayout() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq USSCPayoutRequest
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
