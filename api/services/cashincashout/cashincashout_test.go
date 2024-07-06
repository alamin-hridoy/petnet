package cashincashout

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"

	bpc "brank.as/petnet/api/core/cashincashout"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/services"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	bpa "brank.as/petnet/gunk/drp/v1/cashincashout"
	"brank.as/petnet/serviceutil/auth/hydra"
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

func newTestSvc(t *testing.T, st *postgres.Storage) (*Svc, *services.Mock) {
	cl := perahub.NewTestHTTPMock(st, perahub.MockConfig{})
	log := logrus.New().WithField("stage", "testing")
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
		perahub.WithLogger(log),
		perahub.WithPerahubDefaultAPIKey("api-key"),
		perahub.WithCiCoURL("https://privatedrp.dev.perahub.com.ph/v1/cico/wrapper/"),
		perahub.WithPHRemittanceURL("https://privatedrp.dev.perahub.com.ph/v1/remit/dmt/"),
	)
	if err != nil {
		t.Fatal("setting up perahub integration: ", err)
	}
	m := &services.Mock{}
	return New(bpc.New(ph, st)), m
}

func TestCashInCashOut(t *testing.T) {
	CURRENTTIME := time.Date(2023, 8, 15, 14, 30, 45, 100, time.Local)
	st := newTestStorage(t)
	tests := []struct {
		apiName string
		desc    string
		in      interface{}
		want    interface{}
		tops    cmp.Options
		env     string
	}{
		{
			apiName: "inquire",
			desc:    "Inquire Transaction",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.CicoInquireResult{}, "Expiry"), cmpopts.IgnoreUnexported(bpa.CiCoInquireResponse{}, bpa.CicoInquireResult{})},
			in: &bpa.CiCoInquireRequest{
				PartnerCode:      "DSA",
				Provider:         "GCASH",
				TrxType:          "Cash In",
				ReferenceNumber:  "09654767706",
				PetnetTrackingno: "3115bc3f587d747cf8f5",
			},
			want: &bpa.CiCoInquireResponse{
				Code:    200,
				Message: "Successful",
				Result: &bpa.CicoInquireResult{
					StatusMessage:    "SUCCESSFUL CASHIN",
					PetnetTrackingno: "238a8006885b57765cd8",
					TrxType:          "Cash In",
					ReferenceNumber:  "09654767706",
				},
			},
		},
		{
			apiName: "execute",
			desc:    "Execute Transaction",
			tops:    cmp.Options{cmpopts.IgnoreUnexported(bpa.CiCoExecuteResponse{}, bpa.CicoExecuteResult{})},
			in: &bpa.CiCoExecuteRequest{
				PartnerCode:      "DSA",
				PetnetTrackingno: "3115bc3f587d747cf8f5",
				TrxDate:          "2022-05-17",
				Provider:         "GCASH",
			},
			want: &bpa.CiCoExecuteResponse{
				Code:    200,
				Message: "Successful",
				Result: &bpa.CicoExecuteResult{
					PartnerCode:        "DSA",
					Provider:           "GCASH",
					PetnetTrackingno:   "5a269417e107691f3d7c",
					TrxDate:            "2022-05-17",
					TrxType:            "Cash In",
					ProviderTrackingno: "7000001521345",
					ReferenceNumber:    "09654767706",
					PrincipalAmount:    10,
					Charges:            0,
					TotalAmount:        10,
				},
			},
		},
		{
			apiName: "retry",
			desc:    "Retry Transaction",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.CicoRetryResult{}, "OTPPayload"), cmpopts.IgnoreUnexported(bpa.CiCoRetryResponse{}, bpa.CicoRetryResult{})},
			in: &bpa.CiCoRetryRequest{
				PartnerCode:      "DSA",
				PetnetTrackingno: "5a269417e107691f3d7c",
				TrxDate:          "2022-05-17",
				Provider:         "GCASH",
			},
			want: &bpa.CiCoRetryResponse{
				Code:    200,
				Message: "SUCCESS TRANSACTION.",
				Result: &bpa.CicoRetryResult{
					PartnerCode:        "DSA",
					Provider:           "GCASH",
					PetnetTrackingno:   "5a269417e107691f3d7c",
					TrxDate:            "2022-05-17",
					TrxType:            "Cash In",
					ProviderTrackingno: "09654767706",
					ReferenceNumber:    "09654767706",
					PrincipalAmount:    10,
					Charges:            0,
					TotalAmount:        10,
					OTPPayload: &bpa.OTPPayload{
						CommandID: 27897,
						Payload:   "yZJdFQ+TDbjuSr3ZeLAgFj+S3NnYwkYwRyp1RxZKR20B36kH5Fr52y9I+tg+y5TGrUHRxQAjdVE6rWgJoV1VoLJK0SBHSVNuOts1fRSsXdIPJuxt6v/auPm0gZqyaUXWS+Dtl2OiVpbBPtNB2H6v+bbs7ldWzIDa+47EsWUnUEuVeq8nMM4TPKU0zILbf4lXv6dr2EZCmTX1eNvnyK44QaNBgN68Jb1i50PoD6Gqb61T9CS28btwgIjTlZ/U0s9by4Q8MBmfsEYSejlHpypj/nt0/v9+o8zG9r1kefGt3h4vgtkH3QNoaC7YdWMsAuPJTo1VJTv4ufufBdWD+E+NGwKliGpvUvn4OZ504kyFa5ltLytebihUND70r5S7aI4aYSpEQw==",
					},
				},
			},
		},
		{
			apiName: "otp",
			desc:    "OTP Confirm Transaction",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.CicoOTPConfirmResult{}), cmpopts.IgnoreUnexported(bpa.CiCoOTPConfirmResponse{}, bpa.CicoOTPConfirmResult{})},
			in: &bpa.CiCoOTPConfirmRequest{
				PartnerCode:      "DSA",
				PetnetTrackingno: "8340b7bdb171cbdf6350",
				TrxDate:          "2022-06-03",
				OTP:              "123456",
				Provider:         "GCASH",
				OTPPayload: &bpa.OTPPayload{
					CommandID: 27897,
					Payload:   "yZJdFQ+TDbjuSr3ZeLAgFj+S3NnYwkYwRyp1RxZKR20B36kH5Fr52y9I+tg+y5TGrUHRxQAjdVE6rWgJoV1VoLJK0SBHSVNuOts1fRSsXdIPJuxt6v/auPm0gZqyaUXWS+Dtl2OiVpbBPtNB2H6v+bbs7ldWzIDa+47EsWUnUEuVeq8nMM4TPKU0zILbf4lXv6dr2EZCmTX1eNvnyK44QaNBgN68Jb1i50PoD6Gqb61T9CS28btwgIjTlZ/U0s9by4Q8MBmfsEYSejlHpypj/nt0/v9+o8zG9r1kefGt3h4vgtkH3QNoaC7YdWMsAuPJTo1VJTv4ufufBdWD+E+NGwKliGpvUvn4OZ504kyFa5ltLytebihUND70r5S7aI4aYSpEQw==",
				},
			},
			want: &bpa.CiCoOTPConfirmResponse{
				Code:    200,
				Message: "SUCCESS TRANSACTION.",
				Result: &bpa.CicoOTPConfirmResult{
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
		{
			apiName: "validate",
			desc:    "Validate Transaction",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.CicoValidateResult{}, "Timestamp"), cmpopts.IgnoreUnexported(bpa.CiCoValidateResponse{}, bpa.CicoValidateResult{})},
			in: &bpa.CiCoValidateRequest{
				PartnerCode: "DSA",
				Trx: &bpa.CicoValidateTrx{
					Provider:        "GCASH",
					ReferenceNumber: "09654767706",
					TrxType:         "Cash In",
					PrincipalAmount: 10,
				},
				Customer: &bpa.CicoValidateCustomer{
					CustomerID:        "7085257",
					CustomerFirstname: "Jerica",
					CustomerLastname:  "Naparate",
					CurrAddress:       "1953 Ph 3b Blk 6 Lot 9",
					CurrBarangay:      "Barangay 175",
					CurrCity:          "CALOOCAN CITY",
					CurrProvince:      "METRO MANILA",
					CurrCountry:       "Philippines",
					BirthDate:         "1996-08-10",
					BirthPlace:        "MANILA , METRO MANILA",
					BirthCountry:      "Philippines",
					ContactNo:         "09089488896",
					IDType:            "Passport",
					IDNumber:          "IDTest1233453",
				},
			},
			want: &bpa.CiCoValidateResponse{
				Code:    200,
				Message: "Successful",
				Result: &bpa.CicoValidateResult{
					PetnetTrackingno:   "5a269417e107691f3d7c",
					TrxDate:            "2022-05-17",
					TrxType:            "Cash In",
					Provider:           "GCASH",
					ProviderTrackingno: "7000001521345",
					ReferenceNumber:    "09654767706",
					PrincipalAmount:    10,
					Charges:            0,
					TotalAmount:        10,
					Timestamp:          "",
				},
			},
		},
		{
			apiName: "cicotrxlist",
			desc:    "List CICO Transaction",
			tops: cmp.Options{
				cmpopts.IgnoreFields(bpa.CICOTransact{}, "TransactionCompletedTime"),
				cmpopts.IgnoreUnexported(
					bpa.CICOTransactListRequest{},
					bpa.CICOTransactListResponse{},
					bpa.CICOTransact{},
					bpa.Amount{},
					timestamppb.Timestamp{},
				),
			},
			in: &bpa.CICOTransactListRequest{
				From:   "2023-01-30",
				Until:  "2024-01-30",
				Limit:  5,
				Offset: 0,
			},
			want: &bpa.CICOTransactListResponse{
				Next: 1,
				CICOTransacts: []*bpa.CICOTransact{
					{
						ReferenceNumber: "FWYXWP65",
						Provider:        "DRAGONPAY",
						TotalAmount: &bpa.Amount{
							Amount:   "102",
							Currency: "PHP",
						},
						TransactFee: &bpa.Amount{
							Amount:   "2",
							Currency: "PHP",
						},
						TransactCommission: &bpa.Amount{
							Amount:   "0",
							Currency: "PHP",
						},
						TransactionCompletedTime: timestamppb.New(CURRENTTIME),
					},
				},
				Total: 1,
			},
			env: "TEST",
		},
	}
	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	for _, test := range tests {
		test := test
		oid := uuid.New().String()
		nmd2 := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, hydra.OrgIDKey, oid, "owner", uid, "environment", test.env))
		ctx2 := nmd2.ToIncoming(context.Background())
		t.Run(test.desc, func(t *testing.T) {
			h, _ := newTestSvc(t, st)
			switch test.apiName {
			case "inquire":
				req, ok := test.in.(*bpa.CiCoInquireRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.CiCoInquire(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "execute":
				req, ok := test.in.(*bpa.CiCoExecuteRequest)
				if !ok {
					t.Error("request type conversion error")
				}
				got, err := h.CiCoExecute(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "retry":
				req, ok := test.in.(*bpa.CiCoRetryRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.CiCoRetry(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "otp":
				req, ok := test.in.(*bpa.CiCoOTPConfirmRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.CiCoOTPConfirm(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "validate":
				req, ok := test.in.(*bpa.CiCoValidateRequest)
				if !ok {
					t.Error("request type conversion error")
				}
				got, err := h.CiCoValidate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "cicotrxlist":
				st.CreateCICOHistory(ctx2, storage.CashInCashOutHistory{
					OrgID:              oid,
					PartnerCode:        "DSA",
					Provider:           "DRAGONPAY",
					PetnetTrackingNo:   "ASDVDAW",
					TrxType:            "TRX",
					ProviderTrackingNo: "UVWXYZ",
					ReferenceNumber:    "FWYXWP65",
					PrincipalAmount:    100,
					Charges:            2,
					TotalAmount:        102,
					SvcProvider:        "DRAGONPAY",
					TrxDate:            CURRENTTIME,
					Details:            []byte{},
					TxnStatus:          "SUCCESS",
					ErrorCode:          "",
					ErrorMessage:       "",
					ErrorTime:          "",
					ErrorType:          "",
					CreatedBy:          uuid.NewString(),
					UpdatedBy:          uuid.NewString(),
					Created:            CURRENTTIME,
					Updated:            CURRENTTIME,
				})

				req, ok := test.in.(*bpa.CICOTransactListRequest)
				if !ok {
					t.Error("request type conversion error")
				}
				req.OrgID = oid
				got, err := h.CICOTransactList(ctx2, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			}
		})
	}
}
