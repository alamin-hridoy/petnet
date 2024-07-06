package util

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	bpa "brank.as/petnet/gunk/drp/v1/bills-payment"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestRecordBillPayment(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	othrInfo, _ := json.Marshal(storage.BillsOtherInfo{
		LastName:        "Serato",
		FirstName:       "Mike Edward",
		MiddleName:      "Secret",
		PaymentType:     "B2",
		Course:          "BSCpE",
		TotalAssessment: "0.00",
		SchoolYear:      "2021-2022",
		Term:            "1",
	})
	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.BillPayment{}, "TrxDate", "OtherInfo", "BillPaymentID", "Bills", "Created", "Updated", "SenderMemberID", "BillPaymentDate"),
	}
	amnt, _ := core.MustMinor("200", "PHP").MarshalJSON()
	tests := []struct {
		name string
		in   *bpa.BPTransactRequest
		res  *bpa.BPTransactResponse
		want *storage.BillPayment
	}{
		{
			name: "Success Stage",
			in: &bpa.BPTransactRequest{
				BillerTag:               "testBiller",
				RemoteUserID:            "Remote User ID",
				LocationID:              "test location",
				RemoteLocationID:        "Remote Location ID",
				CurrencyID:              "test currency",
				FormType:                "Form Type",
				FormNumber:              "Form Number",
				Identifier:              "test identifier",
				TotalAmount:             1,
				ClientReferenceNumber:   "test client ref no",
				UserID:                  "a",
				CustomerID:              0,
				LocationName:            "Location Name",
				Coy:                     "test coy",
				CallbackURL:             "",
				BillID:                  "1",
				BillerName:              "Biller Name",
				TrxDate:                 "",
				Amount:                  "200",
				ServiceCharge:           "200",
				PartnerCharge:           "test partner charge",
				AccountNumber:           "test account no",
				PaymentMethod:           "Payment Method",
				ReferenceNumber:         "test ref no",
				ValidationNumber:        "test validation no",
				ReceiptValidationNumber: "test receipt validation no",
				TpaID:                   "test tpaid",
				OtherInfo: &bpa.BPTransactOtherInfo{
					LastName:        "Serato",
					FirstName:       "Mike Edward",
					MiddleName:      "Secret",
					PaymentType:     "B2",
					Course:          "BSCpE",
					TotalAssessment: "0.00",
					SchoolYear:      "2021-2022",
					Term:            "1",
				},
				Type:    "test type",
				Txnid:   "test txn id",
				Partner: "ECPAY",
			},
			res: &bpa.BPTransactResponse{
				Code:    "200",
				Message: "Success",
				Result: &bpa.BPTransactResult{
					Status:          "test status",
					Message:         "test-msg",
					ServiceCharge:   0,
					Timestamp:       "",
					ReferenceNumber: "test-reference-number",
					TransactionID:   "test id",
					ClientReference: "test client ref",
					BillerReference: "test biller ref",
					PaymentMethod:   "test payment method",
					Amount:          "test amount",
					OtherCharges:    "200",
					Details:         []*bpa.Details{},
					CreatedAt:       "",
					URL:             "test URL",
				},
				RemcoID: 2,
			},
			want: &storage.BillPayment{
				BillPaymentID:           uuid.NewString(),
				BillID:                  1,
				UserID:                  "a",
				SenderMemberID:          uuid.NewString(),
				BillPaymentStatus:       string(storage.SuccessStatus),
				ErrorCode:               "",
				ErrorMsg:                "",
				ErrorType:               "",
				Bills:                   []byte{},
				BillerTag:               "testBiller",
				LocationID:              "test location",
				CurrencyID:              "test currency",
				AccountNumber:           "test account no",
				Amount:                  amnt,
				Identifier:              "test identifier",
				Coy:                     "test coy",
				ServiceCharge:           amnt,
				TotalAmount:             amnt,
				BillPaymentDate:         time.Now(),
				PartnerID:               "ECPAY",
				BillerName:              "Biller Name",
				RemoteUserID:            "Remote User ID",
				CustomerID:              "0",
				RemoteLocationID:        "Remote Location ID",
				LocationName:            "Location Name",
				FormType:                "Form Type",
				FormNumber:              "Form Number",
				PaymentMethod:           "Payment Method",
				OtherInfo:               othrInfo,
				TrxDate:                 time.Time{},
				Created:                 time.Time{},
				Updated:                 time.Time{},
				ClientRefNumber:         "test client ref no",
				PartnerCharge:           "test partner charge",
				ReferenceNumber:         "test ref no",
				ValidationNumber:        "test validation no",
				ReceiptValidationNumber: "test receipt validation no",
				TpaID:                   "test tpaid",
				Type:                    "test type",
				TxnID:                   "test txn id",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			var err error

			aRes, aErr := RecordBillPayment(ctx, st, test.in, test.res, err)
			if aErr != nil {
				t.Fatal(aErr)
			}
			got, err := st.GetBillPayment(ctx, aRes.BillPaymentID)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error(cmp.Diff(test.want, got, o))
			}
			return
		})
	}
}
