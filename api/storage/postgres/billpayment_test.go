package postgres

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestListBillPayment(t *testing.T) {
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
	p1 := storage.BillPayment{
		BillPaymentID:           uuid.NewString(),
		BillID:                  1,
		UserID:                  "a",
		SenderMemberID:          "1",
		BillPaymentStatus:       string(storage.SuccessStatus),
		ErrorCode:               "",
		ErrorMsg:                "",
		ErrorType:               "",
		Bills:                   []byte{},
		BillerTag:               "testBiller",
		LocationID:              "test location",
		CurrencyID:              "test currency",
		AccountNumber:           "test ammount no",
		Amount:                  []byte{},
		Identifier:              "test identifier",
		Coy:                     "test coy",
		ServiceCharge:           []byte{},
		TotalAmount:             []byte{},
		BillPaymentDate:         time.Now(),
		PartnerID:               "ECPAY",
		BillerName:              "Biller Name",
		RemoteUserID:            "Remote User ID",
		CustomerID:              "Customer ID",
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
		BillPartnerID:           1,
		PartnerCharge:           "test partner charge",
		ReferenceNumber:         "test ref no",
		ValidationNumber:        "test validation no",
		ReceiptValidationNumber: "test receipt validation no",
		TpaID:                   "test tpaid",
		Type:                    "test type",
		TxnID:                   "test txn id",
		OrgID:                   "test org id",
	}
	p2 := storage.BillPayment{
		BillPaymentID:           uuid.NewString(),
		BillID:                  2,
		UserID:                  "a1",
		SenderMemberID:          "2",
		BillPaymentStatus:       string(storage.SuccessStatus),
		ErrorCode:               "",
		ErrorMsg:                "",
		ErrorType:               "",
		Bills:                   []byte{},
		BillerTag:               "testBiller1",
		LocationID:              "test location1",
		CurrencyID:              "test currency1",
		AccountNumber:           "test ammount no1",
		Amount:                  []byte{},
		Identifier:              "test identifier1",
		Coy:                     "test coy1",
		ServiceCharge:           []byte{},
		TotalAmount:             []byte{},
		BillPaymentDate:         time.Now(),
		PartnerID:               "ECPAY1",
		BillerName:              "Biller Name1",
		RemoteUserID:            "Remote User ID1",
		CustomerID:              "Customer ID1",
		RemoteLocationID:        "Remote Location ID1",
		LocationName:            "Location Name1",
		FormType:                "Form Type1",
		FormNumber:              "Form Number1",
		PaymentMethod:           "Payment Method1",
		OtherInfo:               othrInfo,
		TrxDate:                 time.Time{},
		Created:                 time.Time{},
		Updated:                 time.Time{},
		ClientRefNumber:         "test client ref no1",
		BillPartnerID:           2,
		PartnerCharge:           "test partner charge1",
		ReferenceNumber:         "test ref no1",
		ValidationNumber:        "test validation no1",
		ReceiptValidationNumber: "test receipt validation no1",
		TpaID:                   "test tpaid1",
		Type:                    "test type1",
		TxnID:                   "test txn id1",
		OrgID:                   "test org id1",
	}
	p3 := storage.BillPayment{
		BillPaymentID:           uuid.NewString(),
		BillID:                  3,
		UserID:                  "a2",
		SenderMemberID:          "3",
		BillPaymentStatus:       string(storage.SuccessStatus),
		ErrorCode:               "",
		ErrorMsg:                "",
		ErrorType:               "",
		Bills:                   []byte{},
		BillerTag:               "testBiller2",
		LocationID:              "test location2",
		CurrencyID:              "test currency2",
		AccountNumber:           "test ammount no2",
		Amount:                  []byte{},
		Identifier:              "test identifier2",
		Coy:                     "test coy2",
		ServiceCharge:           []byte{},
		TotalAmount:             []byte{},
		BillPaymentDate:         time.Now(),
		PartnerID:               "ECPAY2",
		BillerName:              "Biller Name2",
		RemoteUserID:            "Remote User ID2",
		CustomerID:              "Customer ID2",
		RemoteLocationID:        "Remote Location ID2",
		LocationName:            "Location Name2",
		FormType:                "Form Type2",
		FormNumber:              "Form Number2",
		PaymentMethod:           "Payment Method2",
		OtherInfo:               othrInfo,
		TrxDate:                 time.Time{},
		Created:                 time.Time{},
		Updated:                 time.Time{},
		ClientRefNumber:         "test client ref no2",
		BillPartnerID:           3,
		PartnerCharge:           "test partner charge2",
		ReferenceNumber:         "test ref no2",
		ValidationNumber:        "test validation no2",
		ReceiptValidationNumber: "test receipt validation no2",
		TpaID:                   "test tpaid2",
		Type:                    "test type2",
		TxnID:                   "test txn id2",
		OrgID:                   "test org id2",
	}
	p4 := storage.BillPayment{
		BillPaymentID:           uuid.NewString(),
		BillID:                  4,
		UserID:                  "a3",
		SenderMemberID:          "4",
		BillPaymentStatus:       string(storage.SuccessStatus),
		ErrorCode:               "",
		ErrorMsg:                "",
		ErrorType:               "",
		Bills:                   []byte{},
		BillerTag:               "testBiller3",
		LocationID:              "test location3",
		CurrencyID:              "test currency3",
		AccountNumber:           "test ammount no3",
		Amount:                  []byte{},
		Identifier:              "test identifier3",
		Coy:                     "test coy3",
		ServiceCharge:           []byte{},
		TotalAmount:             []byte{},
		BillPaymentDate:         time.Now(),
		PartnerID:               "ECPAY3",
		BillerName:              "Biller Name3",
		RemoteUserID:            "Remote User ID3",
		CustomerID:              "Customer ID3",
		RemoteLocationID:        "Remote Location ID3",
		LocationName:            "Location Name3",
		FormType:                "Form Type3",
		FormNumber:              "Form Number3",
		PaymentMethod:           "Payment Method3",
		OtherInfo:               othrInfo,
		TrxDate:                 time.Time{},
		Created:                 time.Time{},
		Updated:                 time.Time{},
		ClientRefNumber:         "test client ref no3",
		BillPartnerID:           4,
		PartnerCharge:           "test partner charge3",
		ReferenceNumber:         "test ref no3",
		ValidationNumber:        "test validation no3",
		ReceiptValidationNumber: "test receipt validation no3",
		TpaID:                   "test tpaid3",
		Type:                    "test type3",
		TxnID:                   "test txn id3",
		OrgID:                   "test org id3",
	}
	p5 := storage.BillPayment{
		BillPaymentID:           uuid.NewString(),
		BillID:                  5,
		UserID:                  "a4",
		SenderMemberID:          "5",
		BillPaymentStatus:       string(storage.SuccessStatus),
		ErrorCode:               "",
		ErrorMsg:                "",
		ErrorType:               "",
		Bills:                   []byte{},
		BillerTag:               "testBiller4",
		LocationID:              "test location4",
		CurrencyID:              "test currency4",
		AccountNumber:           "test ammount no4",
		Amount:                  []byte{},
		Identifier:              "test identifier4",
		Coy:                     "test coy4",
		ServiceCharge:           []byte{},
		TotalAmount:             []byte{},
		BillPaymentDate:         time.Now(),
		PartnerID:               "ECPAY4",
		BillerName:              "Biller Name4",
		RemoteUserID:            "Remote User ID4",
		CustomerID:              "Customer ID4",
		RemoteLocationID:        "Remote Location ID4",
		LocationName:            "Location Name4",
		FormType:                "Form Type4",
		FormNumber:              "Form Number4",
		PaymentMethod:           "Payment Method4",
		OtherInfo:               othrInfo,
		TrxDate:                 time.Time{},
		Created:                 time.Time{},
		Updated:                 time.Time{},
		ClientRefNumber:         "test client ref no4",
		BillPartnerID:           5,
		PartnerCharge:           "test partner charge4",
		ReferenceNumber:         "test ref no4",
		ValidationNumber:        "test validation no4",
		ReceiptValidationNumber: "test receipt validation no4",
		TpaID:                   "test tpaid4",
		Type:                    "test type4",
		TxnID:                   "test txn id4",
		OrgID:                   "test org id4",
	}
	reqs := []storage.BillPayment{p2, p1, p4, p3, p5}

	for _, r := range reqs {
		_, err := st.CreateBillPayment(context.TODO(), r)
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
	}

	tests := []struct {
		desc   string
		want   []storage.BillPayment
		filter storage.BillPaymentFilter
	}{
		{
			desc:   "All",
			filter: storage.BillPaymentFilter{},
			want:   []storage.BillPayment{p2, p1, p4, p3, p5},
		},
		{
			desc: "Sort UserCol ASC",
			filter: storage.BillPaymentFilter{
				SortByColumn: storage.UserCol,
				SortOrder:    "ASC",
			},
			want: []storage.BillPayment{p1, p2, p3, p4, p5},
		},
		{
			desc: "Sort UserCol DESC",
			filter: storage.BillPaymentFilter{
				SortByColumn: storage.UserCol,
				SortOrder:    "DESC",
			},
			want: []storage.BillPayment{p5, p4, p3, p2, p1},
		},
		{
			desc: "Sort BillIDCol ASC",
			filter: storage.BillPaymentFilter{
				SortByColumn: storage.BillIDCol,
				SortOrder:    "ASC",
			},
			want: []storage.BillPayment{p1, p2, p3, p4, p5},
		},
		{
			desc: "Sort BillIDCol DESC",
			filter: storage.BillPaymentFilter{
				SortByColumn: storage.BillIDCol,
				SortOrder:    "DESC",
			},
			want: []storage.BillPayment{p5, p4, p3, p2, p1},
		},
		{
			desc: "Sort SenderMemberIDCol ASC",
			filter: storage.BillPaymentFilter{
				SortByColumn: storage.SenderMemberIDCol,
				SortOrder:    "ASC",
			},
			want: []storage.BillPayment{p1, p2, p3, p4, p5},
		},
		{
			desc: "Sort SenderMemberIDCol DESC",
			filter: storage.BillPaymentFilter{
				SortByColumn: storage.SenderMemberIDCol,
				SortOrder:    "DESC",
			},
			want: []storage.BillPayment{p5, p4, p3, p2, p1},
		},
		{
			desc: "Limit OrgID",
			filter: storage.BillPaymentFilter{
				OrgID:        "test org id4",
				SortOrder:    "ASC",
				SortByColumn: storage.SenderMemberIDCol,
			},
			want: []storage.BillPayment{p5},
		},
		{
			desc: "Limit BillID",
			filter: storage.BillPaymentFilter{
				OrgID:        "test org id3",
				BillID:       "4",
				SortOrder:    "ASC",
				SortByColumn: storage.BillIDCol,
			},
			want: []storage.BillPayment{p4},
		},
		{
			desc: "Limit",
			filter: storage.BillPaymentFilter{
				OrgID:        "test org id2",
				Limit:        3,
				SortOrder:    "ASC",
				SortByColumn: storage.SenderMemberIDCol,
			},
			want: []storage.BillPayment{p3},
		},
		{
			desc: "Offset",
			filter: storage.BillPaymentFilter{
				OrgID:        "test org id1",
				Offset:       1,
				SortOrder:    "ASC",
				SortByColumn: storage.SenderMemberIDCol,
			},
			want: []storage.BillPayment{},
		},
		{
			desc: "Limit/Offset",
			filter: storage.BillPaymentFilter{
				OrgID:        "test org id1",
				Limit:        3,
				Offset:       1,
				SortOrder:    "ASC",
				SortByColumn: storage.SenderMemberIDCol,
			},
			want: []storage.BillPayment{},
		},
	}

	o := cmp.Options{
		cmpopts.IgnoreFields(storage.BillPayment{}, "Updated", "Created", "OrgID", "BillPaymentID", "Bills", "OtherInfo", "BillPaymentDate", "TrxDate", "Total", "SenderMemberID", "Amount", "ServiceCharge", "TotalAmount"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			got, err := st.ListBillPayment(context.TODO(), test.filter)
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}

			if !cmp.Equal(test.want, got, o) {
				t.Fatal(cmp.Diff(test.want, got, o))
			}
		})
	}
}

func TestBillPayment(t *testing.T) {
	ts := newTestStorage(t)
	var err error
	var bj []byte
	bj, err = json.Marshal(storage.Bills{
		Code:    "404",
		Message: "test-msg",
		Result: storage.BillDetails{
			Message:         "test-msg",
			Timestamp:       "2006-01-02",
			ReferenceNumber: "test-reference-number",
		},
	})
	if err != nil {
		bj = []byte{}
	}
	amnt, _ := core.MustMinor("1.5", "PHP").MarshalJSON()
	svcChrg, _ := core.MustMinor("1", "PHP").MarshalJSON()
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
	tests := []struct {
		desc string
		in   *storage.BillPayment
		want interface{}
		tops cmp.Options
	}{
		{
			desc: "billpayment success",
			tops: cmp.Options{cmpopts.IgnoreFields(storage.BillPayment{}, "BillPaymentDate", "SenderMemberID", "BillPaymentID", "Bills", "Amount", "Updated", "OtherInfo", "Created")},
			in: &storage.BillPayment{
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
				AccountNumber:           "test ammount no",
				Amount:                  amnt,
				Identifier:              "test identifier",
				Coy:                     "test coy",
				ServiceCharge:           svcChrg,
				TotalAmount:             svcChrg,
				BillPaymentDate:         time.Now(),
				PartnerID:               "ECPAY",
				BillerName:              "Biller Name",
				RemoteUserID:            "Remote User ID",
				CustomerID:              "Customer ID",
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
				BillPartnerID:           1,
				PartnerCharge:           "test partner charge",
				ReferenceNumber:         "test ref no",
				ValidationNumber:        "test validation no",
				ReceiptValidationNumber: "test receipt validation no",
				TpaID:                   "test tpaid",
				Type:                    "test type",
				TxnID:                   "test txn id",
				OrgID:                   "test org id",
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
				Bills:                   bj,
				BillerTag:               "testBiller",
				LocationID:              "test location",
				CurrencyID:              "test currency",
				AccountNumber:           "test ammount no",
				Amount:                  amnt,
				Identifier:              "test identifier",
				Coy:                     "test coy",
				ServiceCharge:           svcChrg,
				TotalAmount:             svcChrg,
				BillPaymentDate:         time.Now(),
				PartnerID:               "ECPAY",
				BillerName:              "Biller Name",
				RemoteUserID:            "Remote User ID",
				CustomerID:              "Customer ID",
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
				BillPartnerID:           1,
				PartnerCharge:           "test partner charge",
				ReferenceNumber:         "test ref no",
				ValidationNumber:        "test validation no",
				ReceiptValidationNumber: "test receipt validation no",
				TpaID:                   "test tpaid",
				Type:                    "test type",
				TxnID:                   "test txn id",
				OrgID:                   "test org id",
			},
		},
		{
			desc: "billpayment invalid argument",
			in:   &storage.BillPayment{},
			want: storage.ErrInvalid,
			tops: []cmp.Option{},
		},
		{
			desc: "billpayment error",
			tops: cmp.Options{cmpopts.IgnoreFields(storage.BillPayment{}, "BillPaymentDate", "SenderMemberID", "BillPaymentID", "Bills", "Amount", "Updated", "OtherInfo", "Created")},
			in: &storage.BillPayment{
				BillPaymentID:     uuid.NewString(),
				BillID:            1,
				PartnerID:         "1",
				UserID:            "a",
				SenderMemberID:    uuid.NewString(),
				BillPaymentStatus: string(storage.FailStatus),
				ErrorCode:         "404",
				ErrorMsg:          "Error",
				ErrorType:         "Not Found",
				Bills:             []byte{},
				ServiceCharge:     svcChrg,
				TotalAmount:       amnt,
			},
			want: &storage.BillPayment{
				BillPaymentID:     uuid.NewString(),
				BillID:            1,
				PartnerID:         "1",
				UserID:            "a",
				SenderMemberID:    uuid.NewString(),
				BillPaymentStatus: string(storage.FailStatus),
				ErrorCode:         "404",
				ErrorMsg:          "Error",
				ErrorType:         "Not Found",
				Bills:             []byte{},
				ServiceCharge:     svcChrg,
				TotalAmount:       amnt,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			switch test.desc {
			case "billpayment success":
				r, err := ts.CreateBillPayment(context.TODO(), *test.in)
				if err != nil {
					t.Fatalf("CreateBillPayment() = got error %v, want nil", err)
				}
				if r.BillPaymentID == "" {
					t.Fatal("CreateBillPayment() = returned empty ID")
				}

				r.Amount = amnt
				r.Bills = bj
				ru, err := ts.UpdateBillPayment(context.TODO(), *r)
				if err != nil {
					t.Fatalf("UpdateBillPayment() = got error %v, want nil", err)
				}

				if !cmp.Equal(r.Amount, ru.Amount) {
					t.Fatal(cmp.Diff(r.Amount, ru.Amount))
				}

				gr, err := ts.GetBillPayment(context.TODO(), r.BillPaymentID)
				if err != nil {
					t.Fatalf("GetBillPayment() = got error %v, want nil", err)
				}

				wnt, ok := test.want.(*storage.BillPayment)
				if !ok {
					t.Error("want type conversion error")
				}

				if !cmp.Equal(*wnt, *gr, test.tops) {
					t.Fatal(cmp.Diff(*wnt, *gr, test.tops))
				}
			case "billpayment invalid argument":
				_, err := ts.CreateBillPayment(context.TODO(), *test.in)
				wnt, ok := test.want.(error)
				if !ok {
					t.Error("want type conversion error")
				}
				if wnt != err {
					t.Fatal(cmp.Diff(wnt, err))
				}
			case "billpayment error":
				r, err := ts.CreateBillPayment(context.TODO(), *test.in)
				if err != nil {
					t.Fatalf("CreateBillPayment() = got error %v, want nil", err)
				}
				if r.BillPaymentID == "" {
					t.Fatal("CreateBillPayment() = returned empty ID")
				}

				gr, err := ts.GetBillPayment(context.TODO(), r.BillPaymentID)
				if err != nil {
					t.Fatalf("GetBillPayment() = got error %v, want nil", err)
				}
				wnt, ok := test.want.(*storage.BillPayment)
				if !ok {
					t.Error("want type conversion error")
				}

				if !cmp.Equal(*wnt, *gr, test.tops) {
					t.Fatal(cmp.Diff(*wnt, *gr, test.tops))
				}
			}
		})
	}
}
