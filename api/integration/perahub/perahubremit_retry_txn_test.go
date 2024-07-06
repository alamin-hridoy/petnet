package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPerahubRemitconfirm(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          PerahubRemitRetryRequest
		expectedReq PerahubRemitRetryRequest
		want        *PerahubRemitRetryResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: PerahubRemitRetryRequest{
				ID:          7461,
				PartnerCode: "DRP",
			},
			expectedReq: PerahubRemitRetryRequest{
				ID:          7461,
				PartnerCode: "DRP",
			},
			want: &PerahubRemitRetryResponseBody{
				Code:    "200",
				Message: "Good",
				Result: PerahubRemitRetryResult{
					ID:                 7461,
					LocationID:         0,
					UserID:             5188,
					TrxDate:            "2022-06-15",
					CurrencyID:         1,
					RemcoID:            25,
					TrxType:            2,
					IsDomestic:         1,
					CustomerID:         6925902,
					CustomerName:       "Soto, Blanche, G",
					ControlNumber:      "PH1655176065",
					SenderName:         "Sauer, Mittie, O",
					ReceiverName:       "Soto, Blanche, G",
					PrincipalAmount:    "179.00",
					ServiceCharge:      "50.00",
					DstAmount:          "1.00",
					TotalAmount:        "229.00",
					McRate:             "1.00",
					BuyBackAmount:      "1.00",
					RateCategory:       "OPTIONAL",
					McRateID:           1,
					OriginatingCountry: "PH",
					DestinationCountry: "PH",
					PurposeTransaction: "REQUIRED",
					SourceFund:         "REQUIRED",
					Occupation:         "EMPLOYED",
					RelationTo:         "REQUIRED",
					BirthPlace:         "MANILA",
					BirthCountry:       "PH",
					IDType:             "PASSPORT",
					IDNumber:           "0001292",
					Address:            "ADDRESS",
					Barangay:           "1",
					City:               "MANILA",
					Province:           "METRO",
					Country:            "PH",
					ContactNumber:      "09999999999",
					RiskScore:          0,
					RiskCriteria:       "",
					ClientReferenceNo:  "U306WYJXYCPF18",
					FormType:           "OPTIONAL",
					FormNumber:         "OPTIONAL",
					PayoutType:         1,
					RemoteLocationID:   1,
					RemoteUserID:       1,
					RemoteIPAddress:    "OPTIONAL",
					IPAddress:          "OPTIONAL",
					CreatedAt:          "2022-06-15T02:41:40.000000Z",
					UpdatedAt:          "2022-06-15T02:43:09.000000Z",
					ReferenceNumber:    "796e98e940b5982ffa5cee4f89816ad8",
					ZipCode:            "1013",
					Status:             1,
					CurrentAddress: PerahubRemitAddress{
						City:     "MANILA",
						Country:  "PH",
						Barangay: "Barangay 1",
						Province: "METRO MANILA",
						ZipCode:  "1013",
						Address1: "Marcos Highway",
					},
					PermanentAddress: PerahubRemitAddress{
						City:     "MANILA",
						Country:  "PH",
						Barangay: "Barangay 1",
						Province: "METRO MANILA",
						ZipCode:  "1013",
						Address1: "Marcos Highway",
					},
					APIRequest: PerahubRemitAPIRequest{
						City:               "MANILA",
						Address:            "ADDRESS",
						Country:            "PH",
						IDType:             "PASSPORT",
						McRate:             "1",
						UserID:             "5188",
						Barangay:           "Barangay 1",
						Province:           "METRO MANILA",
						RemcoID:            "25",
						TrxDate:            "2022-06-15",
						TrxType:            "2",
						ZipCode:            "1013",
						FormType:           "OPTIONAL",
						IDNumber:           "0001292",
						DstAmount:          "1",
						IPAddress:          "OPTIONAL",
						McRateID:           "1",
						Occupation:         "EMPLOYED",
						BirthPlace:         "MANILA",
						CurrencyID:         "1",
						CustomerID:         "6925902",
						FormNumber:         "OPTIONAL",
						IsDomestic:         "1",
						LocationID:         0,
						PayoutType:         "1",
						RelationTo:         "REQUIRED",
						SenderName:         "Sauer, Mittie, O",
						SourceFund:         "REQUIRED",
						PartnerCode:        "DRP",
						TotalAmount:        "229",
						BirthCountry:       "PH",
						CustomerName:       "Soto, Blanche, G",
						RateCategory:       "OPTIONAL",
						ReceiverName:       "Soto, Blanche, G",
						ContactNumber:      "09999999999",
						ControlNumber:      "PH1655176065",
						RemoteUserID:       "1",
						ServiceCharge:      "50",
						BuyBackAmount:      "1",
						PrincipalAmount:    179,
						ReferenceNumber:    "796e98e940b5982ffa5cee4f89816ad8",
						SenderLastName:     "Sauer",
						RemoteIPAddress:    "OPTIONAL",
						SenderFirstName:    "Mittie",
						ReceiverLastName:   "Soto",
						RemoteLocationID:   "1",
						SenderMiddleName:   "O",
						ClientReferenceNo:  "U306WYJXYCPF18",
						DestinationCountry: "PH",
						OriginatingCountry: "PH",
						PurposeTransaction: "REQUIRED",
						ReceiverFirstName:  "Blanche",
						ReceiverMiddleName: "G",
						APIRequestCurrentAddress: PerahubRemitAddress{
							City:     "MANILA",
							Country:  "PH",
							Barangay: "Barangay 1",
							Province: "METRO MANILA",
							ZipCode:  "1013",
							Address1: "Marcos Highway",
						},
						APIRequestPermanentAddress: PerahubRemitAddress{
							City:     "MANILA",
							Country:  "PH",
							Barangay: "Barangay 1",
							Province: "METRO MANILA",
							ZipCode:  "1013",
							Address1: "Marcos Highway",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.PerahubRemitRetry(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("PerahubRemitRetry() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq PerahubRemitRetryRequest
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
