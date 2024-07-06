package fee

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"brank.as/petnet/api/core/fee"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/serviceutil/auth/hydra"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	fpb "brank.as/petnet/gunk/drp/v1/fee"
	pnpb "brank.as/petnet/gunk/drp/v1/partner"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
)

var _testStorage *postgres.Storage

func TestMain(m *testing.M) {
	const dbConnEnv = "DATABASE_CONNECTION"
	ddlConnStr := os.Getenv(dbConnEnv)
	if ddlConnStr == "" {
		log.Printf("%s is not set, skipping", dbConnEnv)
		return
	}

	var teardown func()
	_testStorage, teardown = postgres.NewTestStorage(ddlConnStr, filepath.Join("..", "..", "migrations", "sql"))

	exitCode := m.Run()

	if teardown != nil {
		teardown()
	}
	os.Exit(exitCode)
}

func newTestStorage(tb testing.TB) *postgres.Storage {
	if testing.Short() {
		tb.Skip("skipping tests that use postgres on -short")
	}

	return _testStorage
}

type tester interface {
	Test(t *testing.T)
}

func TestEnforce(t *testing.T) {
	vs := NewValidators()
	for _, v := range vs {
		for _, sp := range static.Partners["PH"] {
			if v.Kind() == sp.Code {
				ts, ok := v.(tester)
				if !ok {
					t.Fatal("add fee testcase for partner: ", v.Kind())
				}
				ts.Test(t)
			}
		}
	}
}

func (WUVal) Test(t *testing.T) {
	st := newTestStorage(t)

	tests := []struct {
		desc          string
		feeInquiryReq *fpb.FeeInquiryRequest
		want          *fpb.FeeInquiryResponse
		ptnrErr       bool
	}{
		{
			desc:          "Success",
			feeInquiryReq: wuFeeInquiryReq,
			want: &fpb.FeeInquiryResponse{
				Fees: map[string]string{
					"base_message_charge":            "5000",
					"base_message_limit":             "10",
					"canadian_dollar_exchange_fee":   "0",
					"charges":                        "300",
					"county_tax":                     "0",
					"destination_principal_amount":   "100000",
					"exchange_rate":                  "1.0000000",
					"gross_total_amount":             "100000",
					"incremental_message_charge":     "500",
					"incremental_message_limit":      "10",
					"message_charge":                 "0",
					"municipal_tax":                  "0",
					"originating_currency_principal": "",
					"originators_principal_amount":   "100000",
					"pay_amount":                     "100000",
					"plus_charges_amount":            "0",
					"promo_code_description":         "",
					"promo_discount_amount":          "0",
					"promo_name":                     "",
					"promo_sequence_no":              "",
					"state_tax":                      "0",
					"tax_rate":                       "0",
					"tax_worksheet":                  "",
					"tolls":                          "0",
				},
			},
		},
		{
			desc:          "Partner Error",
			feeInquiryReq: wuFeeInquiryReq,
			ptnrErr:       true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, test.ptnrErr)
			got, err := h.FeeInquiry(ctx, test.feeInquiryReq)
			if err := checkError(t, err, test.ptnrErr); err != nil {
				t.Fatal(err)
			}

			o := cmp.Options{
				cmpopts.IgnoreUnexported(
					fpb.FeeInquiryResponse{},
				),
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
			}
		})
	}
}

func (USSCVal) Test(t *testing.T) {
	st := newTestStorage(t)

	tests := []struct {
		desc          string
		feeInquiryReq *fpb.FeeInquiryRequest
		want          *fpb.FeeInquiryResponse
		ptnrErr       bool
	}{
		{
			desc:          "Success",
			feeInquiryReq: usscFeeInquiryReq,
			want: &fpb.FeeInquiryResponse{
				Fees: map[string]string{
					"code":             "0",
					"journal_no":       "000000202",
					"message":          "",
					"new_screen":       "0",
					"principal_amount": "100000",
					"process_date":     "null",
					"reference_number": "1",
					"send_otp":         "Y",
					"service_charge":   "100",
					"total_amount":     "100100",
				},
			},
		},
		{
			desc:          "Partner Error",
			feeInquiryReq: usscFeeInquiryReq,
			ptnrErr:       true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, test.ptnrErr)
			got, err := h.FeeInquiry(ctx, test.feeInquiryReq)
			if err := checkError(t, err, test.ptnrErr); err != nil {
				t.Fatal(err)
			}

			o := cmp.Options{
				cmpopts.IgnoreUnexported(
					fpb.FeeInquiryResponse{},
				),
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
			}
		})
	}
}

var wuFeeInquiryReq = &fpb.FeeInquiryRequest{
	RemitPartner: static.WUCode,
	RemitType:    "Send",
	Amount: &tpb.SendAmount{
		SourceCurrency:      "PHP",
		DestinationCountry:  "PH",
		DestinationCurrency: "PHP",
		DestinationAmount:   true,
		Amount:              "100000",
	},
	Promo:   "promo",
	Message: "msg",
}

var usscFeeInquiryReq = &fpb.FeeInquiryRequest{
	RemitPartner: static.USSCCode,
	Amount: &tpb.SendAmount{
		Amount: "1000",
	},
}

func checkError(t *testing.T, err error, wantErr bool) error {
	if !wantErr && err == nil {
		return nil
	}
	if !wantErr && err != nil {
		return fmt.Errorf("don't want error got: %v", err)
	}
	if wantErr && err == nil {
		return fmt.Errorf("want error got nil")
	}
	st := status.Convert(err)
	if len(st.Details()) != 1 {
		return fmt.Errorf("partner error not returned, got: %v", err)
	}
	det := st.Details()[0]
	switch d := det.(type) {
	case *pnpb.Error:
		if (d.Code != "1" && d.Code != "E0000" && d.Code != "400") ||
			(d.Message != "Transaction does not exists" &&
				d.Message != "Transaction Already Claimed" &&
				d.Message != "Something went wrong") {
			t.Fatalf(`wrong error code or message
got code: %v, got msg: %v
want code: %v, want msg: %v
				`, d.Code, d.Message, "1", "some error")
		}
	default:
		t.Fatalf("should be pnpb.Error")
	}
	return nil
}

func newTestSvc(t *testing.T, store *postgres.Storage, ptnrErr bool) *Svc {
	cl := perahub.NewTestHTTPMock(store, perahub.MockConfig{
		CrPtnrErr:    ptnrErr,
		DsbPtnrErr:   false,
		CfCrPtnrErr:  false,
		CfDsbPtnrErr: false,
	})
	ph, err := perahub.New(cl,
		"dev",
		"https://newkycgateway.dev.perahub.com.ph/gateway/",
		"https://privatedrp.dev.perahub.com.ph/v1/remit/nonex/",
		"https://privatedrp.dev.perahub.com.ph/v1/billspay/wrapper/api/",
		"https://privatedrp.dev.perahub.com.ph/v1/transactions/api/",
		"https://privatedrp.dev.perahub.com.ph/v1/billspay/",
		"partner-id",
		"client-key",
		"api-key",
		"",
		"",
		nil,
	)
	if err != nil {
		t.Fatal("setting up perahub integration: ", err)
	}

	st := static.New(ph, store)
	fc := fee.New(store, ph)

	tvs := NewValidators()
	h := New(st, fc, tvs)
	return h
}
