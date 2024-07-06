package terminal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"brank.as/petnet/api/util"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/pariz/gountries"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/remit"
	aya "brank.as/petnet/api/core/remit/ayannah"
	bp "brank.as/petnet/api/core/remit/bpi"
	ceb "brank.as/petnet/api/core/remit/cebuana"
	cebi "brank.as/petnet/api/core/remit/cebuanaint"
	ic "brank.as/petnet/api/core/remit/instacash"
	ie "brank.as/petnet/api/core/remit/intelexpress"
	ir "brank.as/petnet/api/core/remit/iremit"
	jr "brank.as/petnet/api/core/remit/japanremit"
	mb "brank.as/petnet/api/core/remit/metrobank"
	prmt "brank.as/petnet/api/core/remit/perahubremit"
	rm "brank.as/petnet/api/core/remit/remitly"
	ria "brank.as/petnet/api/core/remit/ria"
	tf "brank.as/petnet/api/core/remit/transfast"
	unt "brank.as/petnet/api/core/remit/uniteller"
	usc "brank.as/petnet/api/core/remit/ussc"
	ws "brank.as/petnet/api/core/remit/wise"
	"brank.as/petnet/api/core/remit/wu"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	pfpb "brank.as/petnet/gunk/drp/v1/profile"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/svcutil/random"
)

var q *gountries.Query

func init() {
	q = gountries.New()
}

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
	vs := NewValidators(q)
	for _, v := range vs {
		for _, sp := range static.Partners["PH"] {
			if v.Kind() == sp.Code {
				ts, ok := v.(tester)
				if !ok {
					t.Fatal("add Txn testcase for partner: ", v.Kind())
				}
				ts.Test(t)
			}
		}
	}
}

func (IRVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := irDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemcoControlNo: "CTRL1",
				RemType:        string(storage.DisburseType),
				TxnStatus:      string(storage.SuccessStatus),
				TxnStep:        string(storage.StageStep),
				Remittance:     rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemcoControlNo: "CTRL1",
				RemType:        string(storage.DisburseType),
				TxnStatus:      string(storage.SuccessStatus),
				TxnStep:        string(storage.ConfirmStep),
				Remittance:     rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemcoControlNo: "CTRL1",
				RemType:        string(storage.DisburseType),
				TxnStatus:      string(storage.FailStatus),
				TxnStep:        string(storage.StageStep),
				ErrorCode:      "1",
				ErrorMsg:       "Transaction does not exists",
				ErrorType:      string(perahub.PartnerError),
				Remittance:     rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemcoControlNo: "CTRL1",
				RemType:        string(storage.DisburseType),
				TxnStatus:      string(storage.SuccessStatus),
				TxnStep:        string(storage.StageStep),
				Remittance:     rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemcoControlNo: "CTRL1",
				RemType:        string(storage.DisburseType),
				TxnStatus:      string(storage.FailStatus),
				TxnStep:        string(storage.ConfirmStep),
				ErrorCode:      "1",
				ErrorMsg:       "Transaction does not exists",
				ErrorType:      string(perahub.PartnerError),
				Remittance:     rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			cfDsbPtnrErr:  true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = uuid.NewString()
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.Remitter.FirstName = ""
					test.recordDsbHistory.Remittance.Remitter.MiddleName = ""
					test.recordDsbHistory.Remittance.Remitter.LastName = ""
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

func (TFVal) Test(t *testing.T) {
	st := newTestStorage(t)
	dr := tfDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			cfDsbPtnrErr:  true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.Remitter.FirstName = ""
					test.recordDsbHistory.Remittance.Remitter.MiddleName = ""
					test.recordDsbHistory.Remittance.Remitter.LastName = ""
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

func (WUVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := wuDisburseReq
	cr := wuCreateReq
	crRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "SO",
		Remitter: storage.Contact{
			FirstName:     cr.Remitter.ContactInfo.GetFirstName(),
			MiddleName:    cr.Remitter.ContactInfo.GetMiddleName(),
			LastName:      cr.Remitter.ContactInfo.GetLastName(),
			RemcoMemberID: cr.Remitter.GetPartnerMemberID(),
			Email:         cr.Remitter.GetEmail(),
			Address1:      cr.Remitter.ContactInfo.Address.GetAddress1(),
			Address2:      cr.Remitter.ContactInfo.Address.GetAddress2(),
			City:          cr.Remitter.ContactInfo.Address.GetCity(),
			State:         cr.Remitter.ContactInfo.Address.GetState(),
			PostalCode:    cr.Remitter.ContactInfo.Address.GetPostalCode(),
			Country:       cr.Remitter.ContactInfo.Address.GetCountry(),
			Province:      cr.Remitter.ContactInfo.Address.GetProvince(),
			Zone:          cr.Remitter.ContactInfo.Address.GetZone(),
			PhoneCty:      cr.Remitter.ContactInfo.Phone.GetCountryCode(),
			Phone:         cr.Remitter.ContactInfo.Phone.GetNumber(),
			MobileCty:     cr.Remitter.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        cr.Remitter.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:  cr.Receiver.ContactInfo.GetFirstName(),
			MiddleName: cr.Receiver.ContactInfo.GetMiddleName(),
			LastName:   cr.Receiver.ContactInfo.GetLastName(),
			Email:      cr.Receiver.ContactInfo.GetEmail(),
			Address1:   cr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:   cr.Receiver.ContactInfo.Address.GetAddress2(),
			City:       cr.Receiver.ContactInfo.Address.GetCity(),
			State:      cr.Receiver.ContactInfo.Address.GetState(),
			PostalCode: cr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:    cr.Receiver.ContactInfo.Address.GetCountry(),
			Province:   cr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:       cr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:   cr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:      cr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:  cr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:     cr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:    "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100100", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("100000", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("100", "PHP"),
	}

	dsbRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "addr1",
			Address2:      "",
			City:          "city",
			State:         "state",
			PostalCode:    "12345",
			Country:       "PH",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       cr.Receiver.ContactInfo.Address.GetCountry(),
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc                string
		createReq           *tpb.CreateRemitRequest
		confirmReq          *tpb.ConfirmRemitRequest
		disburseReq         *tpb.DisburseRemitRequest
		crPtnrErr           bool
		dsbPtnrErr          bool
		cfCrPtnrErr         bool
		cfDsbPtnrErr        bool
		cfCrNotFoundErr     bool
		cfDsbNotFoundErr    bool
		recordSendHistory   *storage.RemitHistory
		recordSendCfHistory *storage.RemitHistory
		recordDsbHistory    *storage.RemitHistory
		recordDsbCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			createReq:   cr,
			confirmReq:  confirmReq,
			disburseReq: dr,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:      "Create Partner Error",
			createReq: cr,
			crPtnrErr: true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "E0000",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Confirm Create Partner Error",
			createReq:   cr,
			confirmReq:  confirmReq,
			cfCrPtnrErr: true,
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "E0000",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:            "Confirm Create NotFound Error",
			createReq:       cr,
			confirmReq:      confirmReq,
			cfCrNotFoundErr: true,
			cfCrPtnrErr:     true,
		},
		{
			desc:        "Disburse Partner Error",
			createReq:   cr,
			confirmReq:  confirmReq,
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "E0000",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Disburse Partner Error",
			createReq:    cr,
			confirmReq:   confirmReq,
			disburseReq:  dr,
			cfDsbPtnrErr: true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "E0000",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:             "Confirm Disburse NotFound Error",
			createReq:        cr,
			confirmReq:       confirmReq,
			disburseReq:      dr,
			cfDsbNotFoundErr: true,
			cfCrPtnrErr:      true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
		cmpopts.IgnoreFields(
			storage.Remittance{}, "RemcoAltControlNo",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, test.crPtnrErr, test.dsbPtnrErr, test.cfCrPtnrErr, test.cfDsbPtnrErr)

			test.createReq.OrderID = uuid.New().String()
			if test.recordSendHistory != nil {
				test.recordSendHistory.DsaID = uid
				test.recordSendHistory.UserID = uid
				test.recordSendHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendHistory.RemcoID = cr.RemitPartner
			}
			crRes, crErr := h.CreateRemit(ctx, test.createReq)
			if err := checkError(t, crErr, test.crPtnrErr); err != nil {
				t.Fatal(err)
			}
			txnID := crRes.GetTransactionID()

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				Partner: test.createReq.RemitPartner,
				RemType: string(storage.SendType),
			})
			if err != nil {
				t.Fatal(err)
			}
			var rh *storage.RemitHistory
			for _, h := range r {
				if h.DsaOrderID == test.createReq.OrderID {
					rh = &h
				}
			}
			if rh == nil {
				t.Fatal("missing remit history")
			}

			if crErr != nil {
				test.recordSendHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				test.recordSendHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.DestAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.Charge = core.MustMinor("0", "PHP")
			}
			if test.recordSendHistory != nil {
				if !cmp.Equal(test.recordSendHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendHistory, rh, o))
				}
			}

			if test.crPtnrErr {
				return
			}

			test.confirmReq.TransactionID = txnID
			if test.recordSendCfHistory != nil {
				test.recordSendCfHistory.DsaID = uid
				test.recordSendCfHistory.UserID = uid
				test.recordSendCfHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendCfHistory.RemcoID = cr.RemitPartner
			}
			if test.cfCrNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			cfRes, err := h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfCrNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfCrPtnrErr); err != nil {
				t.Fatal(err)
			}

			if test.recordSendCfHistory != nil {
				if cfRes != nil {
					test.recordSendCfHistory.RemcoControlNo = cfRes.ControlNumber
				}
			}
			rh, err = st.GetRemitHistory(ctx, txnID)
			if err != nil {
				t.Fatal(err)
			}
			if test.recordSendCfHistory != nil {
				if !cmp.Equal(test.recordSendCfHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendCfHistory, rh, o))
				}
			}

			if test.cfCrPtnrErr {
				return
			}

			test.disburseReq.OrderID = uuid.NewString()
			test.disburseReq.ControlNumber = cfRes.ControlNumber
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
				RemType:   string(storage.DisburseType),
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 2 {
				t.Fatal("remit history should be 2")
			}
			rh2 := r[1]
			if test.recordDsbHistory != nil {
				rdh := test.recordDsbHistory
				if dsbErr != nil {
					rdh.Remittance.Remitter.FirstName = ""
					rdh.Remittance.Remitter.MiddleName = ""
					rdh.Remittance.Remitter.LastName = ""
					rdh.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					rdh.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				rdh.Remittance.Remitter.Address1 = ""
				rdh.Remittance.Remitter.Address2 = ""
				rdh.Remittance.Remitter.City = ""
				rdh.Remittance.Remitter.State = ""
				rdh.Remittance.Remitter.PostalCode = ""
				rdh.Remittance.Remitter.Country = ""
				if !cmp.Equal(*rdh, rh2, o) {
					t.Error(cmp.Diff(*rdh, rh2, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordDsbCfHistory != nil {
				test.recordDsbCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbCfHistory.DsaID = uid
				test.recordDsbCfHistory.UserID = uid
				test.recordDsbCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordDsbCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfDsbNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfDsbNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
				RemType:   string(storage.DisburseType),
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 2 {
				t.Fatal("remit history should be 2")
			}
			rh2 = r[1]
			if test.recordDsbCfHistory != nil {
				if !cmp.Equal(*test.recordDsbCfHistory, rh2, o) {
					t.Error(cmp.Diff(*test.recordDsbCfHistory, rh2, o))
				}
				if rh2.Remittance.RemcoAltControlNo == "" {
					t.Error("Remittance.RemcoAltControlNo should not be empty")
				}
			}
		})
	}
}

func (RMVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := rmDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			cfDsbPtnrErr:  true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.Remitter.FirstName = ""
					test.recordDsbHistory.Remittance.Remitter.MiddleName = ""
					test.recordDsbHistory.Remittance.Remitter.LastName = ""
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

func (RIAVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := riaDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			cfDsbPtnrErr:  true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.Remitter.FirstName = ""
					test.recordDsbHistory.Remittance.Remitter.MiddleName = ""
					test.recordDsbHistory.Remittance.Remitter.LastName = ""
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

func (MBVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := mbDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			cfDsbPtnrErr:  true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

func (BPIVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := bpDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.Remitter.FirstName = ""
					test.recordDsbHistory.Remittance.Remitter.MiddleName = ""
					test.recordDsbHistory.Remittance.Remitter.LastName = ""
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}
			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

func (USSCVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := uscDisburseReq
	cr := uscCreateReq
	crRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "SO",
		Remitter: storage.Contact{
			FirstName:     cr.Remitter.ContactInfo.GetFirstName(),
			MiddleName:    cr.Remitter.ContactInfo.GetMiddleName(),
			LastName:      cr.Remitter.ContactInfo.GetLastName(),
			RemcoMemberID: cr.Remitter.GetPartnerMemberID(),
			Email:         cr.Remitter.GetEmail(),
			Address1:      cr.Remitter.ContactInfo.Address.GetAddress1(),
			Address2:      cr.Remitter.ContactInfo.Address.GetAddress2(),
			City:          cr.Remitter.ContactInfo.Address.GetCity(),
			State:         cr.Remitter.ContactInfo.Address.GetState(),
			PostalCode:    cr.Remitter.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      cr.Remitter.ContactInfo.Address.GetProvince(),
			Zone:          cr.Remitter.ContactInfo.Address.GetZone(),
			PhoneCty:      cr.Remitter.ContactInfo.Phone.GetCountryCode(),
			Phone:         cr.Remitter.ContactInfo.Phone.GetNumber(),
			MobileCty:     cr.Remitter.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        cr.Remitter.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:  cr.Receiver.ContactInfo.GetFirstName(),
			MiddleName: cr.Receiver.ContactInfo.GetMiddleName(),
			LastName:   cr.Receiver.ContactInfo.GetLastName(),
			Phone:      cr.Receiver.ContactInfo.Phone.GetNumber(),
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100100", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("100", "PHP"),
	}

	dsbRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc                string
		createReq           *tpb.CreateRemitRequest
		disburseReq         *tpb.DisburseRemitRequest
		confirmReq          *tpb.ConfirmRemitRequest
		crPtnrErr           bool
		dsbPtnrErr          bool
		cfCrPtnrErr         bool
		cfDsbPtnrErr        bool
		cfCrNotFoundErr     bool
		cfDsbNotFoundErr    bool
		recordSendHistory   *storage.RemitHistory
		recordSendCfHistory *storage.RemitHistory
		recordDsbHistory    *storage.RemitHistory
		recordDsbCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			createReq:   cr,
			confirmReq:  confirmReq,
			disburseReq: dr,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:      "Create Partner Error",
			createReq: cr,
			crPtnrErr: true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Confirm Create Partner Error",
			createReq:   cr,
			confirmReq:  confirmReq,
			cfCrPtnrErr: true,
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:            "Confirm Create NotFound Error",
			createReq:       cr,
			confirmReq:      confirmReq,
			cfCrPtnrErr:     true,
			cfCrNotFoundErr: true,
		},
		{
			desc:        "Disburse Partner Error",
			createReq:   cr,
			confirmReq:  confirmReq,
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Disburse Partner Error",
			createReq:    cr,
			confirmReq:   confirmReq,
			disburseReq:  dr,
			cfDsbPtnrErr: true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:             "Confirm Disburse NotFound Error",
			createReq:        cr,
			confirmReq:       confirmReq,
			disburseReq:      dr,
			cfDsbNotFoundErr: true,
			cfDsbPtnrErr:     true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
		cmpopts.IgnoreFields(
			storage.Remittance{}, "RemcoAltControlNo",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, test.crPtnrErr, test.dsbPtnrErr, test.cfCrPtnrErr, test.cfDsbPtnrErr)

			test.createReq.OrderID = random.NumberString(18)
			if test.recordSendHistory != nil {
				test.recordSendHistory.DsaID = uid
				test.recordSendHistory.UserID = uid
				test.recordSendHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendHistory.RemcoID = cr.RemitPartner
			}
			crRes, crErr := h.CreateRemit(ctx, test.createReq)
			if err := checkError(t, crErr, test.crPtnrErr); err != nil {
				t.Fatal(err)
			}
			txnID := crRes.GetTransactionID()

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				Partner: test.createReq.RemitPartner,
				RemType: string(storage.SendType),
			})
			if err != nil {
				t.Fatal(err)
			}
			var rh *storage.RemitHistory
			for _, h := range r {
				if h.DsaOrderID == test.createReq.OrderID {
					rh = &h
				}
			}
			if rh == nil {
				t.Fatal("missing remit history")
			}

			if crErr != nil {
				test.recordSendHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				test.recordSendHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.DestAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.Charge = core.MustMinor("0", "PHP")
			}
			if test.recordSendHistory != nil {
				if !cmp.Equal(test.recordSendHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendHistory, rh, o))
				}
			}

			if test.crPtnrErr {
				return
			}

			test.confirmReq.TransactionID = txnID
			if test.recordSendCfHistory != nil {
				test.recordSendCfHistory.DsaID = uid
				test.recordSendCfHistory.UserID = uid
				test.recordSendCfHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendCfHistory.RemcoID = cr.RemitPartner
			}
			if test.cfCrNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			cfRes, err := h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfCrNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfCrPtnrErr); err != nil {
				t.Fatal(err)
			}

			if test.recordSendCfHistory != nil {
				if cfRes != nil {
					test.recordSendCfHistory.RemcoControlNo = cfRes.ControlNumber
				}
			}
			rh, err = st.GetRemitHistory(ctx, txnID)
			if err != nil {
				t.Fatal(err)
			}
			if test.recordSendCfHistory != nil {
				if !cmp.Equal(test.recordSendCfHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendCfHistory, rh, o))
				}
			}

			if test.cfCrPtnrErr {
				return
			}

			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = cfRes.ControlNumber
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
				RemType:   string(storage.DisburseType),
				Partner:   cr.RemitPartner,
				TxnStep:   string(storage.StageStep),
			})
			if err != nil {
				t.Fatal(err)
			}
			var rh2 storage.RemitHistory
			for _, h := range r {
				if h.DsaOrderID == test.disburseReq.OrderID {
					rh2 = h
				}
			}
			if test.recordDsbHistory != nil {
				rdh := test.recordDsbHistory
				if dsbErr != nil {
					rdh.Remittance.Remitter.FirstName = ""
					rdh.Remittance.Remitter.MiddleName = ""
					rdh.Remittance.Remitter.LastName = ""
					rdh.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					rdh.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				rdh.Remittance.Remitter.Address1 = ""
				rdh.Remittance.Remitter.Address2 = ""
				rdh.Remittance.Remitter.City = ""
				rdh.Remittance.Remitter.State = ""
				rdh.Remittance.Remitter.PostalCode = ""
				rdh.Remittance.Remitter.Country = ""
				if !cmp.Equal(*rdh, rh2, o) {
					t.Error(cmp.Diff(*rdh, rh2, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordDsbCfHistory != nil {
				test.recordDsbCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbCfHistory.DsaID = uid
				test.recordDsbCfHistory.UserID = uid
				test.recordDsbCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordDsbCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfDsbNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfDsbNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
				RemType:   string(storage.DisburseType),
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, h := range r {
				if h.DsaOrderID == test.disburseReq.OrderID {
					rh2 = h
				}
			}
			if test.recordDsbCfHistory != nil {
				if !cmp.Equal(*test.recordDsbCfHistory, rh2, o) {
					t.Error(cmp.Diff(*test.recordDsbCfHistory, rh2, o))
				}
			}
		})
	}
}

func (CEBVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := cebDisburseReq
	cr := cebCreateReq
	crRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "SO",
		Remitter: storage.Contact{
			FirstName:     cr.Remitter.ContactInfo.GetFirstName(),
			MiddleName:    cr.Remitter.ContactInfo.GetMiddleName(),
			LastName:      cr.Remitter.ContactInfo.GetLastName(),
			RemcoMemberID: cr.Remitter.GetPartnerMemberID(),
			Email:         cr.Remitter.GetEmail(),
			Address1:      cr.Remitter.ContactInfo.Address.GetAddress1(),
			Address2:      cr.Remitter.ContactInfo.Address.GetAddress2(),
			City:          cr.Remitter.ContactInfo.Address.GetCity(),
			State:         cr.Remitter.ContactInfo.Address.GetState(),
			PostalCode:    cr.Remitter.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      cr.Remitter.ContactInfo.Address.GetProvince(),
			Zone:          cr.Remitter.ContactInfo.Address.GetZone(),
			PhoneCty:      cr.Remitter.ContactInfo.Phone.GetCountryCode(),
			Phone:         cr.Remitter.ContactInfo.Phone.GetNumber(),
			MobileCty:     cr.Remitter.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        cr.Remitter.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:  cr.Receiver.ContactInfo.GetFirstName(),
			MiddleName: cr.Receiver.ContactInfo.GetMiddleName(),
			LastName:   cr.Receiver.ContactInfo.GetLastName(),
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100100", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("100", "PHP"),
	}

	dsbRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "PH",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc                string
		createReq           *tpb.CreateRemitRequest
		disburseReq         *tpb.DisburseRemitRequest
		confirmReq          *tpb.ConfirmRemitRequest
		crPtnrErr           bool
		dsbPtnrErr          bool
		cfCrPtnrErr         bool
		cfDsbPtnrErr        bool
		cfCrNotFoundErr     bool
		cfDsbNotFoundErr    bool
		recordSendHistory   *storage.RemitHistory
		recordSendCfHistory *storage.RemitHistory
		recordDsbHistory    *storage.RemitHistory
		recordDsbCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			createReq:   cr,
			confirmReq:  confirmReq,
			disburseReq: dr,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:      "Create Partner Error",
			createReq: cr,
			crPtnrErr: true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Confirm Create Partner Error",
			createReq:   cr,
			confirmReq:  confirmReq,
			cfCrPtnrErr: true,
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:            "Confirm Create NotFound Error",
			createReq:       cr,
			confirmReq:      confirmReq,
			cfCrNotFoundErr: true,
		},
		{
			desc:        "Disburse Partner Error",
			createReq:   cr,
			confirmReq:  confirmReq,
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Disburse Partner Error",
			createReq:    cr,
			confirmReq:   confirmReq,
			disburseReq:  dr,
			cfDsbPtnrErr: true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:             "Confirm Disburse NotFound Error",
			createReq:        cr,
			confirmReq:       confirmReq,
			disburseReq:      dr,
			cfDsbNotFoundErr: true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
		cmpopts.IgnoreFields(
			storage.Remittance{}, "RemcoAltControlNo",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, test.crPtnrErr, test.dsbPtnrErr, test.cfCrPtnrErr, test.cfDsbPtnrErr)

			test.createReq.OrderID = random.NumberString(18)
			if test.recordSendHistory != nil {
				test.recordSendHistory.DsaID = uid
				test.recordSendHistory.UserID = uid
				test.recordSendHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendHistory.RemcoID = cr.RemitPartner
			}
			crRes, crErr := h.CreateRemit(ctx, test.createReq)
			if err := checkError(t, crErr, test.crPtnrErr); err != nil {
				t.Fatal(err)
			}
			txnID := crRes.GetTransactionID()

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				Partner: test.createReq.RemitPartner,
				RemType: string(storage.SendType),
			})
			if err != nil {
				t.Fatal(err)
			}
			var rh *storage.RemitHistory
			for _, h := range r {
				if h.DsaOrderID == test.createReq.OrderID {
					rh = &h
				}
			}
			if rh == nil {
				t.Fatal("missing remit history")
			}

			if crErr != nil {
				test.recordSendHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				test.recordSendHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.DestAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.Charge = core.MustMinor("0", "PHP")
			}
			if test.recordSendHistory != nil {
				if !cmp.Equal(test.recordSendHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendHistory, rh, o))
				}
			}

			if test.crPtnrErr {
				return
			}

			test.confirmReq.TransactionID = txnID
			if test.recordSendCfHistory != nil {
				test.recordSendCfHistory.DsaID = uid
				test.recordSendCfHistory.UserID = uid
				test.recordSendCfHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendCfHistory.RemcoID = cr.RemitPartner
			}
			if test.cfCrNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			cfRes, err := h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfCrNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfCrPtnrErr); err != nil {
				t.Fatal(err)
			}

			if test.recordSendCfHistory != nil {
				if cfRes != nil {
					test.recordSendCfHistory.RemcoControlNo = cfRes.ControlNumber
				}
			}
			rh, err = st.GetRemitHistory(ctx, txnID)
			if err != nil {
				t.Fatal(err)
			}
			if test.recordSendCfHistory != nil {
				if !cmp.Equal(test.recordSendCfHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendCfHistory, rh, o))
				}
			}

			if test.cfCrPtnrErr {
				return
			}

			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = cfRes.ControlNumber
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
				RemType:   string(storage.DisburseType),
				Partner:   cr.RemitPartner,
				TxnStep:   string(storage.StageStep),
			})
			if err != nil {
				t.Fatal(err)
			}
			var rh2 storage.RemitHistory
			for _, h := range r {
				if h.DsaOrderID == test.disburseReq.OrderID {
					rh2 = h
				}
			}
			if test.recordDsbHistory != nil {
				rdh := test.recordDsbHistory
				if dsbErr != nil {
					rdh.Remittance.Remitter.FirstName = ""
					rdh.Remittance.Remitter.MiddleName = ""
					rdh.Remittance.Remitter.LastName = ""
					rdh.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					rdh.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				rdh.Remittance.Remitter.Address1 = ""
				rdh.Remittance.Remitter.Address2 = ""
				rdh.Remittance.Remitter.City = ""
				rdh.Remittance.Remitter.State = ""
				rdh.Remittance.Remitter.PostalCode = ""
				rdh.Remittance.Remitter.Country = ""
				if !cmp.Equal(*rdh, rh2, o) {
					t.Error(cmp.Diff(*rdh, rh2, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordDsbCfHistory != nil {
				test.recordDsbCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbCfHistory.DsaID = uid
				test.recordDsbCfHistory.UserID = uid
				test.recordDsbCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordDsbCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfDsbNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfDsbNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
				RemType:   string(storage.DisburseType),
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, h := range r {
				if h.DsaOrderID == test.disburseReq.OrderID {
					rh2 = h
				}
			}
			if test.recordDsbCfHistory != nil {
				if !cmp.Equal(*test.recordDsbCfHistory, rh2, o) {
					t.Error(cmp.Diff(*test.recordDsbCfHistory, rh2, o))
				}
			}
		})
	}
}

func (ICVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := icDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.Remitter.FirstName = ""
					test.recordDsbHistory.Remittance.Remitter.MiddleName = ""
					test.recordDsbHistory.Remittance.Remitter.LastName = ""
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

func (JPRVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := jprDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.Remitter.FirstName = ""
					test.recordDsbHistory.Remittance.Remitter.MiddleName = ""
					test.recordDsbHistory.Remittance.Remitter.LastName = ""
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

func (UNTVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := untDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       dr.Receiver.ContactInfo.Address.GetCountry(),
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.Remitter.FirstName = ""
					test.recordDsbHistory.Remittance.Remitter.MiddleName = ""
					test.recordDsbHistory.Remittance.Remitter.LastName = ""
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

func (CEBIVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := cebiDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.Remitter.FirstName = ""
					test.recordDsbHistory.Remittance.Remitter.MiddleName = ""
					test.recordDsbHistory.Remittance.Remitter.LastName = ""
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

func (AYAVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := ayaDisburseReq
	cr := ayaCreateReq
	crRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "SO",
		Remitter: storage.Contact{
			FirstName:     cr.Remitter.ContactInfo.GetFirstName(),
			MiddleName:    cr.Remitter.ContactInfo.GetMiddleName(),
			LastName:      cr.Remitter.ContactInfo.GetLastName(),
			RemcoMemberID: cr.Remitter.GetPartnerMemberID(),
			Email:         cr.Remitter.GetEmail(),
			Address1:      cr.Remitter.ContactInfo.Address.GetAddress1(),
			Address2:      cr.Remitter.ContactInfo.Address.GetAddress2(),
			City:          cr.Remitter.ContactInfo.Address.GetCity(),
			State:         cr.Remitter.ContactInfo.Address.GetState(),
			PostalCode:    cr.Remitter.ContactInfo.Address.GetPostalCode(),
			Country:       cr.Remitter.ContactInfo.Address.GetCountry(),
			Province:      cr.Remitter.ContactInfo.Address.GetProvince(),
			Zone:          cr.Remitter.ContactInfo.Address.GetZone(),
			PhoneCty:      cr.Remitter.ContactInfo.Phone.GetCountryCode(),
			Phone:         cr.Remitter.ContactInfo.Phone.GetNumber(),
			MobileCty:     cr.Remitter.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        cr.Remitter.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:  cr.Receiver.ContactInfo.GetFirstName(),
			MiddleName: cr.Receiver.ContactInfo.GetMiddleName(),
			LastName:   cr.Receiver.ContactInfo.GetLastName(),
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	dsbRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc                string
		createReq           *tpb.CreateRemitRequest
		disburseReq         *tpb.DisburseRemitRequest
		confirmReq          *tpb.ConfirmRemitRequest
		crPtnrErr           bool
		dsbPtnrErr          bool
		cfCrPtnrErr         bool
		cfDsbPtnrErr        bool
		cfCrNotFoundErr     bool
		cfDsbNotFoundErr    bool
		recordSendHistory   *storage.RemitHistory
		recordSendCfHistory *storage.RemitHistory
		recordDsbHistory    *storage.RemitHistory
		recordDsbCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			createReq:   cr,
			confirmReq:  confirmReq,
			disburseReq: dr,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Confirm Create Partner Error",
			createReq:   cr,
			confirmReq:  confirmReq,
			cfCrPtnrErr: true,
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:            "Confirm Create NotFound Error",
			createReq:       ayaCreateReq,
			confirmReq:      confirmReq,
			cfCrNotFoundErr: true,
		},
		{
			desc:        "Disburse Partner Error",
			createReq:   cr,
			confirmReq:  confirmReq,
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Disburse Partner Error",
			createReq:    cr,
			confirmReq:   confirmReq,
			disburseReq:  dr,
			cfDsbPtnrErr: true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:             "Confirm Disburse NotFound Error",
			createReq:        cr,
			confirmReq:       confirmReq,
			disburseReq:      dr,
			cfDsbNotFoundErr: true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
		cmpopts.IgnoreFields(
			storage.Remittance{}, "RemcoAltControlNo",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, test.crPtnrErr, test.dsbPtnrErr, test.cfCrPtnrErr, test.cfDsbPtnrErr)

			test.createReq.OrderID = random.NumberString(18)
			if test.recordSendHistory != nil {
				test.recordSendHistory.DsaID = uid
				test.recordSendHistory.UserID = uid
				test.recordSendHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendHistory.RemcoID = cr.RemitPartner
			}
			crRes, crErr := h.CreateRemit(ctx, test.createReq)
			if err := checkError(t, crErr, test.crPtnrErr); err != nil {
				t.Fatal(err)
			}
			txnID := crRes.GetTransactionID()

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				Partner: test.createReq.RemitPartner,
				RemType: string(storage.SendType),
			})
			if err != nil {
				t.Fatal(err)
			}
			var rh *storage.RemitHistory
			for _, h := range r {
				if h.DsaOrderID == test.createReq.OrderID {
					rh = &h
				}
			}
			if rh == nil {
				t.Fatal("missing remit history")
			}

			if crErr != nil {
				test.recordSendHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.DestAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.Charge = core.MustMinor("0", "PHP")
			}
			if test.recordSendHistory != nil {
				r := *rh
				r.RemcoControlNo = ""
				if !cmp.Equal(test.recordSendHistory, &r, o) {
					t.Error(cmp.Diff(test.recordSendHistory, &r, o))
				}
			}

			if test.crPtnrErr {
				return
			}

			test.confirmReq.TransactionID = txnID
			if test.recordSendCfHistory != nil {
				test.recordSendCfHistory.DsaID = uid
				test.recordSendCfHistory.UserID = uid
				test.recordSendCfHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendCfHistory.RemcoID = cr.RemitPartner
			}
			if test.cfCrNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			cfRes, err := h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfCrNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfCrPtnrErr); err != nil {
				t.Fatal(err)
			}

			if test.recordSendCfHistory != nil {
				if cfRes != nil {
					test.recordSendCfHistory.RemcoControlNo = cfRes.ControlNumber
				}
			}
			rh, err = st.GetRemitHistory(ctx, txnID)
			if err != nil {
				t.Fatal(err)
			}
			if test.recordSendCfHistory != nil {
				if test.cfCrPtnrErr {
					rh.RemcoControlNo = ""
				}
				if !cmp.Equal(test.recordSendCfHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendCfHistory, rh, o))
				}
			}

			if test.cfCrPtnrErr {
				return
			}

			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = cfRes.ControlNumber
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
				RemType:   string(storage.DisburseType),
				Partner:   test.disburseReq.RemitPartner,
			})
			if err != nil {
				t.Fatal(err)
			}
			var rh2 storage.RemitHistory
			for _, h := range r {
				if h.DsaOrderID == test.disburseReq.OrderID {
					rh2 = h
				}
			}
			if test.recordDsbHistory != nil {
				rdh := test.recordDsbHistory
				if dsbErr != nil {
					rdh.Remittance.Remitter.FirstName = ""
					rdh.Remittance.Remitter.MiddleName = ""
					rdh.Remittance.Remitter.LastName = ""
					rdh.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					rdh.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				rdh.Remittance.Remitter.Address1 = ""
				rdh.Remittance.Remitter.Address2 = ""
				rdh.Remittance.Remitter.City = ""
				rdh.Remittance.Remitter.State = ""
				rdh.Remittance.Remitter.PostalCode = ""
				rdh.Remittance.Remitter.Country = ""
				if !cmp.Equal(*rdh, rh2, o) {
					t.Error(cmp.Diff(*rdh, rh2, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordDsbCfHistory != nil {
				test.recordDsbCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbCfHistory.DsaID = uid
				test.recordDsbCfHistory.UserID = uid
				test.recordDsbCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordDsbCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfDsbNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfDsbNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
				RemType:   string(storage.DisburseType),
				Partner:   test.disburseReq.RemitPartner,
				TxnStep:   string(storage.ConfirmStep),
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, h := range r {
				if h.DsaOrderID == test.disburseReq.OrderID {
					rh2 = h
				}
			}
			if test.recordDsbCfHistory != nil {
				if !cmp.Equal(*test.recordDsbCfHistory, rh2, o) {
					t.Error(cmp.Diff(*test.recordDsbCfHistory, rh2, o))
				}
			}
		})
	}
}

func (IEVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := ieDisburseReq
	cr := ieCreateReq
	crRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "SO",
		Remitter: storage.Contact{
			FirstName:     cr.Remitter.ContactInfo.GetFirstName(),
			MiddleName:    cr.Remitter.ContactInfo.GetMiddleName(),
			LastName:      cr.Remitter.ContactInfo.GetLastName(),
			RemcoMemberID: cr.Remitter.GetPartnerMemberID(),
			Email:         cr.Remitter.GetEmail(),
			Address1:      cr.Remitter.ContactInfo.Address.GetAddress1(),
			Address2:      cr.Remitter.ContactInfo.Address.GetAddress2(),
			City:          cr.Remitter.ContactInfo.Address.GetCity(),
			State:         cr.Remitter.ContactInfo.Address.GetState(),
			PostalCode:    cr.Remitter.ContactInfo.Address.GetPostalCode(),
			Country:       "PH",
			Province:      cr.Remitter.ContactInfo.Address.GetProvince(),
			Zone:          cr.Remitter.ContactInfo.Address.GetZone(),
			PhoneCty:      cr.Remitter.ContactInfo.Phone.GetCountryCode(),
			Phone:         cr.Remitter.ContactInfo.Phone.GetNumber(),
			MobileCty:     cr.Remitter.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        cr.Remitter.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:  cr.Receiver.ContactInfo.GetFirstName(),
			MiddleName: cr.Receiver.ContactInfo.GetMiddleName(),
			LastName:   cr.Receiver.ContactInfo.GetLastName(),
			Email:      cr.Receiver.ContactInfo.GetEmail(),
			Address1:   cr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:   cr.Receiver.ContactInfo.Address.GetAddress2(),
			City:       cr.Receiver.ContactInfo.Address.GetCity(),
			State:      cr.Receiver.ContactInfo.Address.GetState(),
			PostalCode: cr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:    "Philippines",
			Province:   cr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:       cr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:   cr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:      cr.Receiver.ContactInfo.Phone.GetNumber(),
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	dsbRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "John",
			MiddleName:    "Michael",
			LastName:      "Doe",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "Philippines",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
		SourceAmt:  core.MustMinor("100000", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc                string
		createReq           *tpb.CreateRemitRequest
		disburseReq         *tpb.DisburseRemitRequest
		confirmReq          *tpb.ConfirmRemitRequest
		crPtnrErr           bool
		dsbPtnrErr          bool
		cfCrPtnrErr         bool
		cfDsbPtnrErr        bool
		cfCrNotFoundErr     bool
		cfDsbNotFoundErr    bool
		recordSendHistory   *storage.RemitHistory
		recordSendCfHistory *storage.RemitHistory
		recordDsbHistory    *storage.RemitHistory
		recordDsbCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			createReq:   cr,
			confirmReq:  confirmReq,
			disburseReq: dr,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:            "Confirm Create Partner Error",
			createReq:       cr,
			confirmReq:      confirmReq,
			cfCrPtnrErr:     true,
			cfCrNotFoundErr: true,
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:            "Confirm Create NotFound Error",
			createReq:       cr,
			confirmReq:      confirmReq,
			cfCrNotFoundErr: true,
		},
		{
			desc:        "Disburse Partner Error",
			createReq:   cr,
			confirmReq:  confirmReq,
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Disburse Partner Error",
			createReq:    cr,
			confirmReq:   confirmReq,
			disburseReq:  dr,
			cfDsbPtnrErr: true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordDsbCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: dsbRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:             "Confirm Disburse NotFound Error",
			createReq:        cr,
			confirmReq:       confirmReq,
			disburseReq:      dr,
			cfDsbNotFoundErr: true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
		cmpopts.IgnoreFields(
			storage.Remittance{}, "RemcoAltControlNo",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, test.crPtnrErr, test.dsbPtnrErr, test.cfCrPtnrErr, test.cfDsbPtnrErr)

			test.createReq.OrderID = random.NumberString(18)
			if test.recordSendHistory != nil {
				test.recordSendHistory.DsaID = uid
				test.recordSendHistory.UserID = uid
				test.recordSendHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendHistory.RemcoID = cr.RemitPartner
			}
			crRes, crErr := h.CreateRemit(ctx, test.createReq)
			if err := checkError(t, crErr, test.crPtnrErr); err != nil {
				t.Fatal(err)
			}
			txnID := crRes.GetTransactionID()

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				Partner: test.createReq.RemitPartner,
				RemType: string(storage.SendType),
			})
			if err != nil {
				t.Fatal(err)
			}
			var rh *storage.RemitHistory
			for _, h := range r {
				if h.DsaOrderID == test.createReq.OrderID {
					rh = &h
				}
			}
			if rh == nil {
				t.Fatal("missing remit history")
			}

			if crErr != nil {
				test.recordSendHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.DestAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.Charge = core.MustMinor("0", "PHP")
			}
			if test.recordSendHistory != nil {
				if !cmp.Equal(test.recordSendHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendHistory, rh, o))
				}
			}

			if test.crPtnrErr {
				return
			}

			test.confirmReq.TransactionID = txnID
			if test.recordSendCfHistory != nil {
				test.recordSendCfHistory.DsaID = uid
				test.recordSendCfHistory.UserID = uid
				test.recordSendCfHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendCfHistory.RemcoID = cr.RemitPartner
			}
			if test.cfCrNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			cfRes, err := h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfCrNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfCrPtnrErr); err != nil {
				t.Fatal(err)
			}

			if test.recordSendCfHistory != nil {
				if cfRes != nil {
					test.recordSendCfHistory.RemcoControlNo = cfRes.ControlNumber
				}
			}
			rh, err = st.GetRemitHistory(ctx, txnID)
			if err != nil {
				t.Fatal(err)
			}
			if test.recordSendCfHistory != nil {
				if !cmp.Equal(test.recordSendCfHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendCfHistory, rh, o))
				}
			}

			if test.cfCrPtnrErr {
				return
			}

			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = cfRes.ControlNumber
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
				RemType:   string(storage.DisburseType),
				Partner:   cr.RemitPartner,
				TxnStep:   string(storage.StageStep),
			})
			if err != nil {
				t.Fatal(err)
			}
			var rh2 storage.RemitHistory
			for _, h := range r {
				if h.DsaOrderID == test.disburseReq.OrderID {
					rh2 = h
				}
			}
			if test.recordDsbHistory != nil {
				rdh := test.recordDsbHistory
				if dsbErr != nil {
					rdh.Remittance.Remitter.FirstName = ""
					rdh.Remittance.Remitter.MiddleName = ""
					rdh.Remittance.Remitter.LastName = ""
					rdh.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					rdh.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				}
				rdh.Remittance.Remitter.Address1 = ""
				rdh.Remittance.Remitter.Address2 = ""
				rdh.Remittance.Remitter.City = ""
				rdh.Remittance.Remitter.State = ""
				rdh.Remittance.Remitter.PostalCode = ""
				rdh.Remittance.Remitter.Country = ""
				if !cmp.Equal(*rdh, rh2, o) {
					t.Error(cmp.Diff(*rdh, rh2, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordDsbCfHistory != nil {
				test.recordDsbCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbCfHistory.DsaID = uid
				test.recordDsbCfHistory.UserID = uid
				test.recordDsbCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordDsbCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfDsbNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfDsbNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
				RemType:   string(storage.DisburseType),
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, h := range r {
				if h.DsaOrderID == test.disburseReq.OrderID {
					rh2 = h
				}
			}
			if test.recordDsbCfHistory != nil {
				if !cmp.Equal(*test.recordDsbCfHistory, rh2, o) {
					t.Error(cmp.Diff(*test.recordDsbCfHistory, rh2, o))
				}
			}
		})
	}
}

func TestListRemit(t *testing.T) {
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)

	ptnrs := map[string]struct {
		createReq   *tpb.CreateRemitRequest
		disburseReq *tpb.DisburseRemitRequest
	}{
		static.WUCode: {
			createReq:   wuCreateReq,
			disburseReq: wuDisburseReq,
		},
		static.IRCode: {
			disburseReq: irDisburseReq,
		},
		static.TFCode: {
			disburseReq: tfDisburseReq,
		},
		static.RMCode: {
			disburseReq: rmDisburseReq,
		},
		static.RIACode: {
			disburseReq: riaDisburseReq,
		},
		static.MBCode: {
			disburseReq: mbDisburseReq,
		},
		static.BPICode: {
			disburseReq: bpDisburseReq,
		},
		static.USSCCode: {
			createReq:   uscCreateReq,
			disburseReq: uscDisburseReq,
		},
		static.ICCode: {
			disburseReq: icDisburseReq,
		},
		static.JPRCode: {
			disburseReq: jprDisburseReq,
		},
		static.WISECode: {
			createReq: wiseCreateReq,
		},
		static.UNTCode: {
			disburseReq: untDisburseReq,
		},
		static.CEBCode: {
			createReq:   cebCreateReq,
			disburseReq: cebDisburseReq,
		},
		static.CEBINTCode: {
			disburseReq: cebiDisburseReq,
		},
		static.AYACode: {
			createReq:   ayaCreateReq,
			disburseReq: ayaDisburseReq,
		},
		static.IECode: {
			createReq:   ieCreateReq,
			disburseReq: ieDisburseReq,
		},
		static.PerahubRemit: {
			disburseReq: phDisburseReq,
		},
	}

	vs := NewValidators(q)
	for _, v := range vs {
		testExists := false
		for ptnr := range ptnrs {
			if v.Kind() == ptnr {
				testExists = true
			}
		}
		if !testExists {
			t.Fatal("add List testcase for partner: ", v.Kind())
		}
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid, hydra.OrgIDKey, uid))
	ctx := nmd.ToIncoming(context.Background())

	h := newTestSvc(t, st, false, false, false, false)
	for _, test := range ptnrs {
		cf := &tpb.ConfirmRemitRequest{
			AuthSource: "User Review",
			AuthCode:   "Manual",
		}

		if test.createReq != nil {
			test.createReq.OrderID = uuid.New().String()
			crRes, err := h.CreateRemit(ctx, test.createReq)
			if err != nil {
				t.Fatal(err)
			}

			cf.TransactionID = crRes.TransactionID
			cfRes, err := h.ConfirmRemit(ctx, cf)
			if err != nil {
				t.Fatal(err)
			}

			if test.disburseReq != nil {
				test.disburseReq.ControlNumber = cfRes.GetControlNumber()
			}
		}
		if test.disburseReq != nil && test.disburseReq.ControlNumber == "" {
			test.disburseReq.ControlNumber = time.Now().Local().Format("20060102") + "PHB" + random.NumberString(9)
		}

		if test.disburseReq != nil {
			test.disburseReq.OrderID = random.NumberString(10)
			disRes, err := h.DisburseRemit(ctx, test.disburseReq)
			if err != nil {
				t.Fatal(err)
			}

			cf.TransactionID = disRes.TransactionID
			_, err = h.ConfirmRemit(ctx, cf)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	// update the transaction dates of some of the remits to test filtering
	rhs, err := st.ListRemitHistory(ctx, storage.LRHFilter{
		TxnStep:   string(storage.ConfirmStep),
		TxnStatus: string(storage.SuccessStatus),
	})
	if err != nil {
		t.Fatal(err)
	}
	var irRef string
	var tfRef string
	for _, rh := range rhs {
		switch rh.RemcoID {
		case static.IRCode, static.TFCode:
			ts := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
			if err := st.UpdateRemitHistoryDate(ctx, rh.TxnID, ts); err != nil {
				t.Fatal(err)
			}
			if rh.RemcoID == static.IRCode {
				irRef = rh.RemcoControlNo
			} else {
				tfRef = rh.RemcoControlNo
			}
		case static.RMCode, static.RIACode, static.USSCCode, static.ICCode, static.JPRCode, static.UNTCode:
			ts := time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC)
			if err := st.UpdateRemitHistoryDate(ctx, rh.TxnID, ts); err != nil {
				t.Fatal(err)
			}
		}
	}
	rhs, err = st.ListRemitHistory(ctx, storage.LRHFilter{
		TxnStep:   string(storage.ConfirmStep),
		TxnStatus: string(storage.SuccessStatus),
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := h.ListRemit(ctx, &tpb.ListRemitRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res.GetRemittances() {
		for _, rh := range rhs {
			if r.RemitPartner == rh.RemcoID && r.RemitType == rh.RemType && r.ControlNumber == rh.RemcoControlNo {
				want := &tpb.Remittance{
					RemitPartner:             rh.RemcoID,
					ControlNumber:            rh.RemcoControlNo,
					RemitType:                rh.RemType,
					TransactionCompletedTime: timestamppb.New(rh.TxnCompletedTime.Time),
					GrossAmount: &tpb.Amount{
						Amount:   rh.Remittance.GrossTotal.Number(),
						Currency: rh.Remittance.GrossTotal.CurrencyCode(),
					},
					RemitAmount: &tpb.Amount{
						Amount:   rh.Remittance.SourceAmt.Number(),
						Currency: rh.Remittance.SourceAmt.CurrencyCode(),
					},
					Remitter: &tpb.Contact{
						FirstName:  rh.Remittance.Remitter.FirstName,
						MiddleName: rh.Remittance.Remitter.MiddleName,
						LastName:   rh.Remittance.Remitter.LastName,
						Email:      rh.Remittance.Remitter.Email,
						Address: &tpb.Address{
							Address1:   rh.Remittance.Remitter.Address1,
							Address2:   rh.Remittance.Remitter.Address2,
							City:       rh.Remittance.Remitter.City,
							State:      rh.Remittance.Remitter.State,
							PostalCode: rh.Remittance.Remitter.PostalCode,
							Country:    rh.Remittance.Remitter.Country,
							Zone:       rh.Remittance.Remitter.Zone,
							Province:   rh.Remittance.Remitter.Province,
						},
						Mobile: &ppb.PhoneNumber{
							CountryCode: rh.Remittance.Remitter.MobileCty,
							Number:      rh.Remittance.Remitter.Mobile,
						},
						Phone: &ppb.PhoneNumber{
							CountryCode: rh.Remittance.Remitter.PhoneCty,
							Number:      rh.Remittance.Remitter.Phone,
						},
					},
					Receiver: &tpb.Contact{
						FirstName:  rh.Remittance.Receiver.FirstName,
						MiddleName: rh.Remittance.Receiver.MiddleName,
						LastName:   rh.Remittance.Receiver.LastName,
						Email:      rh.Remittance.Receiver.Email,
						Address: &tpb.Address{
							Address1:   rh.Remittance.Receiver.Address1,
							Address2:   rh.Remittance.Receiver.Address2,
							City:       rh.Remittance.Receiver.City,
							State:      rh.Remittance.Receiver.State,
							PostalCode: rh.Remittance.Receiver.PostalCode,
							Country:    rh.Remittance.Receiver.Country,
							Zone:       rh.Remittance.Receiver.Zone,
							Province:   rh.Remittance.Receiver.Province,
						},
						Mobile: &ppb.PhoneNumber{
							CountryCode: rh.Remittance.Receiver.MobileCty,
							Number:      rh.Remittance.Receiver.Mobile,
						},
						Phone: &ppb.PhoneNumber{
							CountryCode: rh.Remittance.Receiver.PhoneCty,
							Number:      rh.Remittance.Receiver.Phone,
						},
					},
				}

				o := cmp.Options{
					cmpopts.IgnoreFields(tpb.Remittance{}, "TransactionCompletedTime", "TransactionStagedTime"),
					cmpopts.IgnoreUnexported(
						tpb.Remittance{},
						tpb.Contact{},
						tpb.Amount{},
						tpb.Address{},
						ppb.PhoneNumber{},
						timestamppb.Timestamp{},
					),
				}
				if want.TransactionCompletedTime.AsTime() != r.TransactionCompletedTime.AsTime() {
					t.Errorf("transaction time mismatch for partner: %v, remType: %v", r.RemitPartner, r.RemitType)
				}
				if !cmp.Equal(want, r, o) {
					t.Error("(-want +got): ", cmp.Diff(want, r, o))
				}
			}
		}
	}

	tests := []struct {
		desc    string
		listReq *tpb.ListRemitRequest
		want    []string
	}{
		{
			// sorting by partner name to get consistency in the listing
			// when adding a test case for new partner add it to the list
			// of partners in alphabetic order
			desc: "All",
			listReq: &tpb.ListRemitRequest{
				SortOrder:    tpb.SortOrder_ASC,
				SortByColumn: tpb.SortByColumn_Partner,
			},
			want: []string{static.AYACode, static.AYACode, static.BPICode, static.CEBCode, static.CEBCode, static.CEBINTCode, static.ICCode, static.IECode, static.IECode, static.IRCode, static.JPRCode, static.MBCode, static.PerahubRemit, static.RIACode, static.RMCode, static.TFCode, static.UNTCode, static.USSCCode, static.USSCCode, static.WISECode, static.WUCode, static.WUCode},
		},
		{
			desc: "From",
			listReq: &tpb.ListRemitRequest{
				From:         "2010-01-01",
				SortOrder:    tpb.SortOrder_ASC,
				SortByColumn: tpb.SortByColumn_Partner,
			},
			want: []string{static.AYACode, static.AYACode, static.BPICode, static.CEBCode, static.CEBCode, static.CEBINTCode, static.ICCode, static.IECode, static.IECode, static.JPRCode, static.MBCode, static.PerahubRemit, static.RIACode, static.RMCode, static.UNTCode, static.USSCCode, static.USSCCode, static.WISECode, static.WUCode, static.WUCode},
		},
		{
			desc: "Until",
			listReq: &tpb.ListRemitRequest{
				Until:        "2000-01-01",
				SortOrder:    tpb.SortOrder_ASC,
				SortByColumn: tpb.SortByColumn_Partner,
			},
			want: []string{static.IRCode, static.TFCode},
		},
		{
			desc: "From/Until",
			listReq: &tpb.ListRemitRequest{
				From:         "2010-01-01",
				Until:        "2010-01-01",
				SortOrder:    tpb.SortOrder_ASC,
				SortByColumn: tpb.SortByColumn_Partner,
			},
			want: []string{static.ICCode, static.JPRCode, static.RIACode, static.RMCode, static.UNTCode, static.USSCCode, static.USSCCode},
		},
		{
			desc: "Limit",
			listReq: &tpb.ListRemitRequest{
				Limit:        2,
				SortOrder:    tpb.SortOrder_ASC,
				SortByColumn: tpb.SortByColumn_Partner,
			},
			want: []string{static.AYACode, static.AYACode},
		},
		{
			desc: "Offset",
			listReq: &tpb.ListRemitRequest{
				Offset:       2,
				SortOrder:    tpb.SortOrder_ASC,
				SortByColumn: tpb.SortByColumn_Partner,
			},
			want: []string{static.BPICode, static.CEBCode, static.CEBCode, static.CEBINTCode, static.ICCode, static.IECode, static.IECode, static.IRCode, static.JPRCode, static.MBCode, static.PerahubRemit, static.RIACode, static.RMCode, static.TFCode, static.UNTCode, static.USSCCode, static.USSCCode, static.WISECode, static.WUCode, static.WUCode},
		},
		{
			desc: "Limit/Offset",
			listReq: &tpb.ListRemitRequest{
				Limit:        2,
				Offset:       3,
				SortOrder:    tpb.SortOrder_ASC,
				SortByColumn: tpb.SortByColumn_Partner,
			},
			want: []string{static.CEBCode, static.CEBCode},
		},
		{
			desc: "SortByColumn Asc",
			listReq: &tpb.ListRemitRequest{
				SortOrder:    tpb.SortOrder_ASC,
				SortByColumn: tpb.SortByColumn_Partner,
			},
			want: []string{static.AYACode, static.AYACode, static.BPICode, static.CEBCode, static.CEBCode, static.CEBINTCode, static.ICCode, static.IECode, static.IECode, static.IRCode, static.JPRCode, static.MBCode, static.PerahubRemit, static.RIACode, static.RMCode, static.TFCode, static.UNTCode, static.USSCCode, static.USSCCode, static.WISECode, static.WUCode, static.WUCode},
		},
		{
			desc: "SortByColumn Desc",
			listReq: &tpb.ListRemitRequest{
				SortOrder:    tpb.SortOrder_DESC,
				SortByColumn: tpb.SortByColumn_Partner,
			},
			want: []string{static.WUCode, static.WUCode, static.WISECode, static.USSCCode, static.USSCCode, static.UNTCode, static.TFCode, static.RMCode, static.RIACode, static.PerahubRemit, static.MBCode, static.JPRCode, static.IRCode, static.IECode, static.IECode, static.ICCode, static.CEBINTCode, static.CEBCode, static.CEBCode, static.BPICode, static.AYACode, static.AYACode},
		},
		{
			desc: "By RefNumbers",
			listReq: &tpb.ListRemitRequest{
				ControlNumbers: []string{irRef, tfRef},
				SortOrder:      tpb.SortOrder_ASC,
				SortByColumn:   tpb.SortByColumn_Partner,
			},
			want: []string{static.IRCode, static.TFCode},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			res, err := h.ListRemit(ctx, test.listReq)
			if err != nil {
				t.Fatal(err)
			}
			ps := make([]string, len(res.GetRemittances()))
			for i, r := range res.GetRemittances() {
				ps[i] = r.GetRemitPartner()
			}
			if !cmp.Equal(test.want, ps) {
				t.Fatal(cmp.Diff(test.want, ps))
			}
		})
	}
}

func (WISEVal) Test(t *testing.T) {
	st := newTestStorage(t)

	cr := wiseCreateReq
	crRmt := storage.Remittance{
		CustomerTxnID:     "",
		RemcoAltControlNo: "",
		TxnType:           "SO",
		Remitter: storage.Contact{
			FirstName:     cr.Remitter.ContactInfo.GetFirstName(),
			MiddleName:    cr.Remitter.ContactInfo.GetMiddleName(),
			LastName:      cr.Remitter.ContactInfo.GetLastName(),
			RemcoMemberID: cr.Remitter.GetPartnerMemberID(),
			Email:         cr.Remitter.GetEmail(),
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "hello testing",
		},
		Receiver: storage.Contact{
			FirstName:  cr.Receiver.ContactInfo.GetFirstName(),
			MiddleName: cr.Receiver.ContactInfo.GetMiddleName(),
			LastName:   cr.Receiver.ContactInfo.GetLastName(),
			Email:      cr.Receiver.ContactInfo.GetEmail(),
			Address1:   "",
			Address2:   "",
			City:       "",
			State:      "",
			PostalCode: "",
			Country:    "",
			Province:   "",
			Zone:       "",
			PhoneCty:   "",
			Phone:      "",
			MobileCty:  "",
			Mobile:     "",
			Message:    "hello testing",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("150000", "PHP")},
		SourceAmt:  core.MustMinor("150000", "PHP"),
		DestAmt:    core.MustMinor("2110", "GBP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("7966", "PHP"),
	}

	tests := []struct {
		desc                string
		createReq           *tpb.CreateRemitRequest
		confirmReq          *tpb.ConfirmRemitRequest
		crPtnrErr           bool
		cfCrPtnrErr         bool
		cfCrNotFoundErr     bool
		recordSendHistory   *storage.RemitHistory
		recordSendCfHistory *storage.RemitHistory
	}{
		{
			desc:       "Success",
			createReq:  cr,
			confirmReq: confirmReq,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:      "Create Partner Error",
			createReq: cr,
			crPtnrErr: true,
			recordSendHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Confirm Create Partner Error",
			createReq:   cr,
			confirmReq:  confirmReq,
			cfCrPtnrErr: true,
			recordSendCfHistory: &storage.RemitHistory{
				RemType:    string(storage.SendType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: crRmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:            "Confirm Create NotFound Error",
			createReq:       wiseCreateReq,
			confirmReq:      confirmReq,
			cfCrNotFoundErr: true,
			cfCrPtnrErr:     true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
		cmpopts.IgnoreFields(
			storage.Remittance{}, "RemcoAltControlNo",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, test.crPtnrErr, false, test.cfCrPtnrErr, false)

			test.createReq.OrderID = uuid.New().String()
			if test.recordSendHistory != nil {
				test.recordSendHistory.DsaID = uid
				test.recordSendHistory.UserID = uid
				test.recordSendHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendHistory.RemcoID = cr.RemitPartner
			}
			crRes, crErr := h.CreateRemit(ctx, test.createReq)
			if err := checkError(t, crErr, test.crPtnrErr); err != nil {
				t.Fatal(err)
			}
			txnID := crRes.GetTransactionID()

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				Partner: test.createReq.RemitPartner,
				RemType: string(storage.SendType),
			})
			if err != nil {
				t.Fatal(err)
			}
			var rh *storage.RemitHistory
			for _, h := range r {
				if h.DsaOrderID == test.createReq.OrderID {
					rh = &h
				}
			}
			if rh == nil {
				t.Fatal("missing remit history")
			}

			if crErr != nil {
				test.recordSendHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
				test.recordSendHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.DestAmt = core.MustMinor("0", "PHP")
				test.recordSendHistory.Remittance.Charge = core.MustMinor("0", "PHP")
			}
			if test.recordSendHistory != nil {
				if !cmp.Equal(test.recordSendHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendHistory, rh, o))
				}
			}

			if test.crPtnrErr {
				return
			}

			test.confirmReq.TransactionID = txnID
			if test.recordSendCfHistory != nil {
				test.recordSendCfHistory.DsaID = uid
				test.recordSendCfHistory.UserID = uid
				test.recordSendCfHistory.DsaOrderID = test.createReq.OrderID
				test.recordSendCfHistory.RemcoID = cr.RemitPartner
			}
			if test.cfCrNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			cfRes, err := h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfCrNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfCrPtnrErr); err != nil {
				t.Fatal(err)
			}

			if test.recordSendCfHistory != nil {
				if cfRes != nil {
					test.recordSendCfHistory.RemcoControlNo = cfRes.ControlNumber
				}
			}
			rh, err = st.GetRemitHistory(ctx, txnID)
			if err != nil {
				t.Fatal(err)
			}
			if test.recordSendCfHistory != nil {
				if !test.cfCrPtnrErr {
					test.recordSendCfHistory.Remittance.CustomerTxnID = "aecd179d"
				}
				if !cmp.Equal(test.recordSendCfHistory, rh, o) {
					t.Error(cmp.Diff(test.recordSendCfHistory, rh, o))
				}
			}
		})
	}
}

func (PHUBVal) Test(t *testing.T) {
	st := newTestStorage(t)

	dr := phDisburseReq
	rmt := storage.Remittance{
		CustomerTxnID:     "7461",
		RemcoAltControlNo: "",
		TxnType:           "PO",
		Remitter: storage.Contact{
			FirstName:     "Mittie",
			MiddleName:    "O",
			LastName:      "Sauer",
			RemcoMemberID: "",
			Email:         "",
			Address1:      "",
			Address2:      "",
			City:          "",
			State:         "",
			PostalCode:    "",
			Country:       "",
			Province:      "",
			Zone:          "",
			PhoneCty:      "",
			Phone:         "",
			MobileCty:     "",
			Mobile:        "",
			Message:       "",
		},
		Receiver: storage.Contact{
			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
			LastName:      dr.Receiver.ContactInfo.GetLastName(),
			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
			Email:         dr.Receiver.ContactInfo.GetEmail(),
			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
			City:          dr.Receiver.ContactInfo.Address.GetCity(),
			State:         dr.Receiver.ContactInfo.Address.GetState(),
			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
			Country:       "PH",
			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
			Message:       "",
		},
		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("0", "PHP")},
		SourceAmt:  core.MustMinor("17900", "PHP"),
		DestAmt:    core.MustMinor("0", "PHP"),
		Taxes:      nil,
		Tax:        core.MustMinor("0", "PHP"),
		Charges:    nil,
		Charge:     core.MustMinor("0", "PHP"),
	}

	tests := []struct {
		desc             string
		disburseReq      *tpb.DisburseRemitRequest
		confirmReq       *tpb.ConfirmRemitRequest
		dsbPtnrErr       bool
		cfDsbPtnrErr     bool
		cfNotFoundErr    bool
		recordDsbHistory *storage.RemitHistory
		recordCfHistory  *storage.RemitHistory
	}{
		{
			desc:        "Success",
			disburseReq: dr,
			confirmReq:  confirmReq,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.ConfirmStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:        "Disburse Partner Error",
			disburseReq: dr,
			dsbPtnrErr:  true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.StageStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:         "Confirm Partner Error",
			disburseReq:  dr,
			confirmReq:   confirmReq,
			cfDsbPtnrErr: true,
			recordDsbHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.SuccessStatus),
				TxnStep:    string(storage.StageStep),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			recordCfHistory: &storage.RemitHistory{
				RemType:    string(storage.DisburseType),
				TxnStatus:  string(storage.FailStatus),
				TxnStep:    string(storage.ConfirmStep),
				ErrorCode:  "1",
				ErrorMsg:   "Transaction does not exists",
				ErrorType:  string(perahub.PartnerError),
				Remittance: rmt,
				TransactionType: sql.NullString{
					String: util.PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
		},
		{
			desc:          "Confirm NotFound Error",
			disburseReq:   dr,
			confirmReq:    confirmReq,
			cfNotFoundErr: true,
			cfDsbPtnrErr:  true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h := newTestSvc(t, st, false, test.dsbPtnrErr, false, test.cfDsbPtnrErr)
			test.disburseReq.OrderID = random.NumberString(18)
			test.disburseReq.ControlNumber = random.InvitationCode(20)
			if test.recordDsbHistory != nil {
				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordDsbHistory.DsaID = uid
				test.recordDsbHistory.UserID = uid
				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
				test.recordDsbHistory.RemcoID = dr.RemitPartner
			}
			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh := r[0]
			if test.recordDsbHistory != nil {
				if dsbErr != nil {
					test.recordDsbHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
					test.recordDsbHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
					test.recordDsbHistory.Remittance.CustomerTxnID = ""
				}
				if !cmp.Equal(*test.recordDsbHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordDsbHistory, rh, o))
				}
			}
			if test.dsbPtnrErr {
				return
			}
			if res.TransactionID == "" {
				t.Fatal("h.DisburseRemit: ", err)
			}

			test.confirmReq.TransactionID = res.TransactionID
			if test.recordCfHistory != nil {
				test.recordCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
				test.recordCfHistory.DsaID = uid
				test.recordCfHistory.UserID = uid
				test.recordCfHistory.DsaOrderID = test.disburseReq.OrderID
				test.recordCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
				test.recordCfHistory.RemcoID = dr.RemitPartner
			}
			if test.cfNotFoundErr {
				test.confirmReq.TransactionID = uuid.New().String()
			}
			_, err = h.ConfirmRemit(ctx, test.confirmReq)
			if err != nil {
				if test.cfNotFoundErr && status.Code(err) == codes.NotFound {
					return
				}
			}
			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
				t.Fatal(err)
			}

			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
				ControlNo: []string{test.disburseReq.ControlNumber},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(r) != 1 {
				t.Fatal("remit history should be 1")
			}
			rh = r[0]
			if test.recordCfHistory != nil {
				if !cmp.Equal(*test.recordCfHistory, rh, o) {
					t.Error(cmp.Diff(*test.recordCfHistory, rh, o))
				}
			}
		})
	}
}

// func (PHUBVal) Test(t *testing.T) {
// 	st := newTestStorage(t)

// 	dr := phDisburseReq
// 	cr := phCreateReq
// 	crRmt := storage.Remittance{
// 		CustomerTxnID:     "",
// 		RemcoAltControlNo: "",
// 		TxnType:           "SO",
// 		Remitter: storage.Contact{
// 			FirstName:     cr.Remitter.ContactInfo.GetFirstName(),
// 			MiddleName:    cr.Remitter.ContactInfo.GetMiddleName(),
// 			LastName:      cr.Remitter.ContactInfo.GetLastName(),
// 			RemcoMemberID: cr.Remitter.GetPartnerMemberID(),
// 			Email:         cr.Remitter.GetEmail(),
// 			Address1:      cr.Remitter.ContactInfo.Address.GetAddress1(),
// 			Address2:      cr.Remitter.ContactInfo.Address.GetAddress2(),
// 			City:          cr.Remitter.ContactInfo.Address.GetCity(),
// 			State:         cr.Remitter.ContactInfo.Address.GetState(),
// 			PostalCode:    cr.Remitter.ContactInfo.Address.GetPostalCode(),
// 			Country:       "PH",
// 			Province:      cr.Remitter.ContactInfo.Address.GetProvince(),
// 			Zone:          cr.Remitter.ContactInfo.Address.GetZone(),
// 			PhoneCty:      cr.Remitter.ContactInfo.Phone.GetCountryCode(),
// 			Phone:         cr.Remitter.ContactInfo.Phone.GetNumber(),
// 			MobileCty:     cr.Remitter.ContactInfo.Mobile.GetCountryCode(),
// 			Mobile:        cr.Remitter.ContactInfo.Mobile.GetNumber(),
// 			Message:       "",
// 		},
// 		Receiver: storage.Contact{
// 			FirstName:  cr.Receiver.ContactInfo.GetFirstName(),
// 			MiddleName: cr.Receiver.ContactInfo.GetMiddleName(),
// 			LastName:   cr.Receiver.ContactInfo.GetLastName(),
// 			Address1:   cr.Receiver.ContactInfo.Address.GetAddress1(),
// 			Address2:   cr.Receiver.ContactInfo.Address.GetAddress2(),
// 			City:       cr.Receiver.ContactInfo.Address.GetCity(),
// 			State:      cr.Receiver.ContactInfo.Address.GetState(),
// 			PostalCode: cr.Receiver.ContactInfo.Address.GetPostalCode(),
// 			Country:    "PH",
// 			Province:   cr.Receiver.ContactInfo.Address.GetProvince(),
// 			Zone:       cr.Receiver.ContactInfo.Address.GetZone(),
// 			PhoneCty:   cr.Receiver.ContactInfo.Phone.GetCountryCode(),
// 			Phone:      cr.Receiver.ContactInfo.Phone.GetNumber(),
// 			MobileCty:  cr.Receiver.ContactInfo.Mobile.GetCountryCode(),
// 			Mobile:     cr.Receiver.ContactInfo.Mobile.GetNumber(),
// 			Message:    "",
// 		},
// 		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("100000", "PHP")},
// 		SourceAmt:  core.MustMinor("100000", "PHP"),
// 		DestAmt:    core.MustMinor("100000", "PHP"),
// 		Taxes:      nil,
// 		Tax:        core.MustMinor("0", "PHP"),
// 		Charges:    nil,
// 		Charge:     core.MustMinor("0", "PHP"),
// 	}

// 	dsbRmt := storage.Remittance{
// 		CustomerTxnID:     "7461",
// 		RemcoAltControlNo: "",
// 		TxnType:           "PO",
// 		Remitter: storage.Contact{
// 			FirstName:     "Mittie",
// 			MiddleName:    "O",
// 			LastName:      "Sauer",
// 			RemcoMemberID: "",
// 			Email:         "",
// 			Address1:      "",
// 			Address2:      "",
// 			City:          "",
// 			State:         "",
// 			PostalCode:    "",
// 			Country:       "",
// 			Province:      "",
// 			Zone:          "",
// 			PhoneCty:      "",
// 			Phone:         "",
// 			MobileCty:     "",
// 			Mobile:        "",
// 			Message:       "",
// 		},
// 		Receiver: storage.Contact{
// 			FirstName:     dr.Receiver.ContactInfo.GetFirstName(),
// 			MiddleName:    dr.Receiver.ContactInfo.GetMiddleName(),
// 			LastName:      dr.Receiver.ContactInfo.GetLastName(),
// 			RemcoMemberID: dr.Receiver.GetPartnerMemberID(),
// 			Email:         dr.Receiver.ContactInfo.GetEmail(),
// 			Address1:      dr.Receiver.ContactInfo.Address.GetAddress1(),
// 			Address2:      dr.Receiver.ContactInfo.Address.GetAddress2(),
// 			City:          dr.Receiver.ContactInfo.Address.GetCity(),
// 			State:         dr.Receiver.ContactInfo.Address.GetState(),
// 			PostalCode:    dr.Receiver.ContactInfo.Address.GetPostalCode(),
// 			Country:       "PH",
// 			Province:      dr.Receiver.ContactInfo.Address.GetProvince(),
// 			Zone:          dr.Receiver.ContactInfo.Address.GetZone(),
// 			PhoneCty:      dr.Receiver.ContactInfo.Phone.GetCountryCode(),
// 			Phone:         dr.Receiver.ContactInfo.Phone.GetNumber(),
// 			MobileCty:     dr.Receiver.ContactInfo.Mobile.GetCountryCode(),
// 			Mobile:        dr.Receiver.ContactInfo.Mobile.GetNumber(),
// 			Message:       "",
// 		},
// 		GrossTotal: storage.GrossTotal{Minor: core.MustMinor("0", "PHP")},
// 		SourceAmt:  core.MustMinor("17900", "PHP"),
// 		DestAmt:    core.MustMinor("0", "PHP"),
// 		Taxes:      nil,
// 		Tax:        core.MustMinor("0", "PHP"),
// 		Charges:    nil,
// 		Charge:     core.MustMinor("0", "PHP"),
// 	}

// 	tests := []struct {
// 		desc                string
// 		createReq           *tpb.CreateRemitRequest
// 		disburseReq         *tpb.DisburseRemitRequest
// 		confirmReq          *tpb.ConfirmRemitRequest
// 		crPtnrErr           bool
// 		dsbPtnrErr          bool
// 		cfCrPtnrErr         bool
// 		cfDsbPtnrErr        bool
// 		cfCrNotFoundErr     bool
// 		cfDsbNotFoundErr    bool
// 		recordSendHistory   *storage.RemitHistory
// 		recordSendCfHistory *storage.RemitHistory
// 		recordDsbHistory    *storage.RemitHistory
// 		recordDsbCfHistory  *storage.RemitHistory
// 	}{
// 		{
// 			desc:        "Success",
// 			createReq:   cr,
// 			confirmReq:  confirmReq,
// 			disburseReq: dr,
// 			recordSendHistory: &storage.RemitHistory{
// 				RemType:    string(storage.SendType),
// 				TxnStatus:  string(storage.SuccessStatus),
// 				TxnStep:    string(storage.StageStep),
// 				Remittance: crRmt,
// 				TransactionType: sql.NullString{
// 					String: util.PerahubTrxTypeOTC,
// 					Valid:  true,
// 				},
// 			},
// 			recordSendCfHistory: &storage.RemitHistory{
// 				RemType:    string(storage.SendType),
// 				TxnStatus:  string(storage.SuccessStatus),
// 				TxnStep:    string(storage.ConfirmStep),
// 				Remittance: crRmt,
// 				TransactionType: sql.NullString{
// 					String: util.PerahubTrxTypeOTC,
// 					Valid:  true,
// 				},
// 			},
// 			recordDsbHistory: &storage.RemitHistory{
// 				RemType:    string(storage.DisburseType),
// 				TxnStatus:  string(storage.SuccessStatus),
// 				TxnStep:    string(storage.StageStep),
// 				Remittance: dsbRmt,
// 				TransactionType: sql.NullString{
// 					String: util.PerahubTrxTypeOTC,
// 					Valid:  true,
// 				},
// 			},
// 			recordDsbCfHistory: &storage.RemitHistory{
// 				RemType:    string(storage.DisburseType),
// 				TxnStatus:  string(storage.SuccessStatus),
// 				TxnStep:    string(storage.ConfirmStep),
// 				Remittance: dsbRmt,
// 				TransactionType: sql.NullString{
// 					String: util.PerahubTrxTypeOTC,
// 					Valid:  true,
// 				},
// 			},
// 		},
// 		{
// 			desc:      "Create Partner Error",
// 			createReq: cr,
// 			crPtnrErr: true,
// 			recordSendHistory: &storage.RemitHistory{
// 				RemType:    string(storage.SendType),
// 				TxnStatus:  string(storage.FailStatus),
// 				TxnStep:    string(storage.StageStep),
// 				ErrorCode:  "1",
// 				ErrorMsg:   "Transaction does not exists",
// 				ErrorType:  string(perahub.PartnerError),
// 				Remittance: crRmt,
// 				TransactionType: sql.NullString{
// 					String: util.PerahubTrxTypeOTC,
// 					Valid:  true,
// 				},
// 			},
// 		},
// 		{
// 			desc:        "Confirm Create Partner Error",
// 			createReq:   cr,
// 			confirmReq:  confirmReq,
// 			cfCrPtnrErr: true,
// 			recordSendCfHistory: &storage.RemitHistory{
// 				RemType:    string(storage.SendType),
// 				TxnStatus:  string(storage.FailStatus),
// 				TxnStep:    string(storage.ConfirmStep),
// 				ErrorCode:  "1",
// 				ErrorMsg:   "Transaction does not exists",
// 				ErrorType:  string(perahub.PartnerError),
// 				Remittance: crRmt,
// 				TransactionType: sql.NullString{
// 					String: util.PerahubTrxTypeOTC,
// 					Valid:  true,
// 				},
// 			},
// 		},
// 		{
// 			desc:            "Confirm Create NotFound Error",
// 			createReq:       cr,
// 			confirmReq:      confirmReq,
// 			cfCrNotFoundErr: true,
// 		},
// 		{
// 			desc:        "Disburse Partner Error",
// 			createReq:   cr,
// 			confirmReq:  confirmReq,
// 			disburseReq: dr,
// 			dsbPtnrErr:  true,
// 			recordSendHistory: &storage.RemitHistory{
// 				RemType:    string(storage.SendType),
// 				TxnStatus:  string(storage.SuccessStatus),
// 				TxnStep:    string(storage.StageStep),
// 				Remittance: crRmt,
// 				TransactionType: sql.NullString{
// 					String: util.PerahubTrxTypeOTC,
// 					Valid:  true,
// 				},
// 			},
// 			recordSendCfHistory: &storage.RemitHistory{
// 				RemType:    string(storage.SendType),
// 				TxnStatus:  string(storage.SuccessStatus),
// 				TxnStep:    string(storage.ConfirmStep),
// 				Remittance: crRmt,
// 				TransactionType: sql.NullString{
// 					String: util.PerahubTrxTypeOTC,
// 					Valid:  true,
// 				},
// 			},
// 			recordDsbHistory: &storage.RemitHistory{
// 				RemType:    string(storage.DisburseType),
// 				TxnStatus:  string(storage.FailStatus),
// 				TxnStep:    string(storage.StageStep),
// 				ErrorCode:  "1",
// 				ErrorMsg:   "Transaction does not exists",
// 				ErrorType:  string(perahub.PartnerError),
// 				Remittance: dsbRmt,
// 				TransactionType: sql.NullString{
// 					String: util.PerahubTrxTypeOTC,
// 					Valid:  true,
// 				},
// 			},
// 		},
// 		{
// 			desc:             "Confirm Disburse NotFound Error",
// 			createReq:        cr,
// 			confirmReq:       confirmReq,
// 			disburseReq:      dr,
// 			cfDsbNotFoundErr: true,
// 		},
// 	}

// 	uid := uuid.New().String()
// 	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
// 	ctx := nmd.ToIncoming(context.Background())

// 	o := cmp.Options{
// 		cmpopts.IgnoreFields(
// 			storage.RemitHistory{}, "TxnID", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
// 		),
// 		cmpopts.IgnoreFields(
// 			storage.Remittance{}, "RemcoAltControlNo", "TxnType",
// 		),
// 	}

// 	for _, test := range tests {
// 		test := test
// 		t.Run(test.desc, func(t *testing.T) {
// 			h := newTestSvc(t, st, test.crPtnrErr, test.dsbPtnrErr, test.cfCrPtnrErr, test.cfDsbPtnrErr)

// 			test.createReq.OrderID = random.NumberString(18)
// 			if test.recordSendHistory != nil {
// 				test.recordSendHistory.DsaID = uid
// 				test.recordSendHistory.UserID = uid
// 				test.recordSendHistory.DsaOrderID = test.createReq.OrderID
// 				test.recordSendHistory.RemcoID = cr.RemitPartner
// 			}
// 			crRes, crErr := h.CreateRemit(ctx, test.createReq)
// 			if err := checkError(t, crErr, test.crPtnrErr); err != nil {
// 				t.Fatal(err)
// 			}
// 			txnID := crRes.GetTransactionID()

// 			r, err := st.ListRemitHistory(ctx, storage.LRHFilter{
// 				Partner: test.createReq.RemitPartner,
// 				RemType: string(storage.SendType),
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			var rh *storage.RemitHistory
// 			for _, h := range r {
// 				if h.DsaOrderID == test.createReq.OrderID {
// 					rh = &h
// 				}
// 			}
// 			if rh == nil {
// 				t.Fatal("missing remit history")
// 			}

// 			if crErr != nil {
// 				test.recordSendHistory.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
// 				test.recordSendHistory.Remittance.SourceAmt = core.MustMinor("0", "PHP")
// 				test.recordSendHistory.Remittance.DestAmt = core.MustMinor("0", "PHP")
// 				test.recordSendHistory.Remittance.Charge = core.MustMinor("0", "PHP")
// 			}
// 			if test.recordSendHistory != nil {
// 				if !cmp.Equal(test.recordSendHistory, rh, o) {
// 					t.Error(cmp.Diff(test.recordSendHistory, rh, o))
// 				}
// 			}

// 			if test.crPtnrErr {
// 				return
// 			}

// 			test.confirmReq.TransactionID = txnID
// 			if test.recordSendCfHistory != nil {
// 				test.recordSendCfHistory.DsaID = uid
// 				test.recordSendCfHistory.UserID = uid
// 				test.recordSendCfHistory.DsaOrderID = test.createReq.OrderID
// 				test.recordSendCfHistory.RemcoID = cr.RemitPartner
// 			}
// 			if test.cfCrNotFoundErr {
// 				test.confirmReq.TransactionID = uuid.New().String()
// 			}
// 			cfRes, err := h.ConfirmRemit(ctx, test.confirmReq)
// 			if err != nil {
// 				if test.cfCrNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
// 					return
// 				}
// 			}
// 			if err := checkError(t, err, test.cfCrPtnrErr); err != nil {
// 				t.Fatal(err)
// 			}

// 			if test.recordSendCfHistory != nil {
// 				if cfRes != nil {
// 					test.recordSendCfHistory.RemcoControlNo = cfRes.ControlNumber
// 				}
// 			}
// 			rh, err = st.GetRemitHistory(ctx, txnID)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			if test.recordSendCfHistory != nil {
// 				if !cmp.Equal(test.recordSendCfHistory, rh, o) {
// 					t.Error(cmp.Diff(test.recordSendCfHistory, rh, o))
// 				}
// 			}

// 			if test.cfCrPtnrErr {
// 				return
// 			}

// 			test.disburseReq.OrderID = random.NumberString(18)
// 			test.disburseReq.ControlNumber = cfRes.ControlNumber
// 			if test.disburseReq.ControlNumber == "" {
// 				test.disburseReq.ControlNumber = random.NumberString(18)
// 			}
// 			if test.recordDsbHistory != nil {
// 				test.recordDsbHistory.RemcoControlNo = test.disburseReq.ControlNumber
// 				test.recordDsbHistory.DsaID = uid
// 				test.recordDsbHistory.UserID = uid
// 				test.recordDsbHistory.DsaOrderID = test.disburseReq.OrderID
// 				test.recordDsbHistory.ReceiverID = dr.Receiver.GetPartnerMemberID()
// 				test.recordDsbHistory.RemcoID = dr.RemitPartner
// 			}
// 			res, dsbErr := h.DisburseRemit(ctx, test.disburseReq)
// 			if err := checkError(t, dsbErr, test.dsbPtnrErr); err != nil {
// 				t.Fatal(err)
// 			}

// 			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
// 				ControlNo: []string{test.disburseReq.ControlNumber},
// 				RemType:   string(storage.DisburseType),
// 				Partner:   cr.RemitPartner,
// 				TxnStep:   string(storage.StageStep),
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			var rh2 storage.RemitHistory
// 			for _, h := range r {
// 				if h.DsaOrderID == test.disburseReq.OrderID {
// 					rh2 = h
// 				}
// 			}
// 			if test.recordDsbHistory != nil {
// 				rdh := test.recordDsbHistory
// 				if dsbErr != nil {
// 					rdh.Remittance.CustomerTxnID = ""
// 					rdh.Remittance.SourceAmt = core.MustMinor("0", "PHP")
// 					rdh.Remittance.GrossTotal.Amount = core.MustAmount("0", "PHP")
// 				}

// 				if !cmp.Equal(*rdh, rh2, o) {
// 					t.Error(cmp.Diff(*rdh, rh2, o))
// 				}
// 			}
// 			if test.dsbPtnrErr {
// 				return
// 			}
// 			if res.TransactionID == "" {
// 				t.Fatal("h.DisburseRemit: ", err)
// 			}

// 			test.confirmReq.TransactionID = res.TransactionID
// 			if test.recordDsbCfHistory != nil {
// 				test.recordDsbCfHistory.RemcoControlNo = test.disburseReq.ControlNumber
// 				test.recordDsbCfHistory.DsaID = uid
// 				test.recordDsbCfHistory.UserID = uid
// 				test.recordDsbCfHistory.DsaOrderID = test.disburseReq.OrderID
// 				test.recordDsbCfHistory.ReceiverID = dr.Receiver.PartnerMemberID
// 				test.recordDsbCfHistory.RemcoID = dr.RemitPartner
// 			}
// 			if test.cfDsbNotFoundErr {
// 				test.confirmReq.TransactionID = uuid.New().String()
// 			}
// 			_, err = h.ConfirmRemit(ctx, test.confirmReq)
// 			if err != nil {
// 				if test.cfDsbNotFoundErr && int(status.Convert(err).Code()) == http.StatusNotFound {
// 					return
// 				}
// 			}
// 			if err := checkError(t, err, test.cfDsbPtnrErr); err != nil {
// 				t.Fatal(err)
// 			}

// 			r, err = st.ListRemitHistory(ctx, storage.LRHFilter{
// 				ControlNo: []string{test.disburseReq.ControlNumber},
// 				RemType:   string(storage.DisburseType),
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			for _, h := range r {
// 				if h.DsaOrderID == test.disburseReq.OrderID {
// 					rh2 = h
// 				}
// 			}
// 			fmt.Println("rh2", rh2)
// 			if test.recordDsbCfHistory != nil {
// 				if !cmp.Equal(*test.recordDsbCfHistory, rh2, o) {
// 					t.Error(cmp.Diff(*test.recordDsbCfHistory, rh2, o))
// 				}
// 			}
// 		})
// 	}
// }

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

	// TODO: To find out purpose of this coplicated code
	//d := status.Convert(err)
	//if (int(d.Code()) != 500 &&
	//	int(d.Code()) != 404 &&
	//	int(d.Code()) != 409 &&
	//	int(d.Code()) != 503 &&
	//	int(d.Code()) != 422) || (d.Message() != coreerror.MsgDatabaseError &&
	//	d.Message() != coreerror.MsgInvalidInput &&
	//	d.Message() != coreerror.MsgControlNumberNotFound &&
	//	d.Message() != coreerror.MsgPartnerDoesntExist &&
	//	d.Message() != coreerror.MsgIdentifierAlreadyExists &&
	//	d.Message() != coreerror.MsgConnectionError) {
	//	t.Fatalf(`wrong error code or message
	//		got code: %v, got msg: %v
	//		want code: %v, want msg: %v
	//						`, d.Code(), d.Message(), "1", "transaction does not exists")
	//}
	return nil
}

var irDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  static.IRCode,
	ControlNumber: "3802032400050411",
	Receiver: &tpb.UserKYC{
		PartnerMemberID:    "1",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		TransactionPurpose: "Family Support/Living Expenses",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &pfpb.Identification{
			Type:   "Postal ID",
			Number: "PRND32200265569P",
		},
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
			Phone: &ppb.PhoneNumber{
				Number: "12345678",
			},
		},
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var tfDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  static.TFCode,
	ControlNumber: "33TF105956774",
	Receiver: &tpb.UserKYC{
		SendingReasonID:    "1133",
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "1",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		ProofOfAddress:     tpb.Bool_True,
		KYCVerified:        tpb.Bool_True,
		Gender:             tpb.Gender_Male,
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			OccupationID: "1",
			Occupation:   "Unemployed",
		},
		Identification: &pfpb.Identification{
			Type:   "1",
			Number: "24023497AB0877AAB20000",
			Expiration: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Issued: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Country: "PH",
		},
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
			Phone: &ppb.PhoneNumber{
				Number: "12345678",
			},
		},
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
}

var confirmReq = &tpb.ConfirmRemitRequest{
	AuthSource: "User Review",
	AuthCode:   "Manual",
}

var wuCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.WUCode,
	RemitType:    "Send",
	Agent: &tpb.Agent{
		UserID: 123,
	},
	Amount: &tpb.SendAmount{
		Amount:              "100000",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationCountry:  "PH",
		DestinationAmount:   true,
	},
	Remitter: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				PostalCode: "12345",
				Country:    "PH",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "87654321",
			},
		},
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		Gender:             tpb.Gender_Male,
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &pfpb.Identification{
			Type:    "M",
			Country: "PH",
			Number:  "24023497AB0877AAB20000",
			Expiration: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Issued: &pfpb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
		},
		Email:       "test@mail.com",
		Nationality: "PH",
	},
	Receiver: &tpb.Receiver{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				PostalCode: "12345",
				Country:    "PH",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "87654321",
			},
		},
	},
}

var wiseCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.WISECode,
	RemitType:    "Send",
	Remitter: &tpb.UserKYC{
		Email: "sender2@brankas.com",
	},
	Receiver: &tpb.Receiver{
		RecipientID:         "12345",
		AccountHolderName:   "Brankas Receiver",
		SourceAccountNumber: "54321",
	},
	Message: "hello testing",
}

var uscCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.USSCCode,
	RemitType:    "Send",
	OrderID:      "5142096205",
	Amount: &tpb.SendAmount{
		Amount:              "100000",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationCountry:  "PH",
		DestinationAmount:   true,
	},
	Receiver: &tpb.Receiver{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "12345678",
			},
		},
	},
	Remitter: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "bit",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		KYCVerified:        tpb.Bool_True,
		Gender:             tpb.Gender_Male,
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &pfpb.Identification{
			Type:    "M",
			Country: "PH",
			Number:  "24023497AB0877AAB20000",
		},
	},
	Agent: &tpb.Agent{
		UserID:    1123,
		IPAddress: "::1",
	},
}

var cebCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.CEBCode,
	RemitType:    "Send",
	OrderID:      "5142096205",
	Amount: &tpb.SendAmount{
		Amount:              "100000",
		SourceCountry:       "PH",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationCountry:  "PH",
		DestinationAmount:   true,
	},
	Receiver: &tpb.Receiver{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
		RecipientID: "12345",
	},
	Remitter: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "bit",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &pfpb.Identification{
			Type:    "M",
			Country: "PH",
			Number:  "24023497AB0877AAB20000",
		},
	},
	Agent: &tpb.Agent{
		UserID:    54321,
		IPAddress: "::1",
	},
}

var cebDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.CEBCode,
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var wuDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.WUCode,
	DisburseCurrency: "PHP",
	Agent: &tpb.Agent{
		UserID: 123,
	},
	Receiver: &tpb.UserKYC{
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		Gender:             tpb.Gender_Male,
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &pfpb.Identification{
			Type:    "M",
			Country: "PH",
			Number:  "24023497AB0877AAB20000",
			Expiration: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Issued: &pfpb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
		},
		ContactInfo: &tpb.Contact{
			FirstName:  "first",
			MiddleName: "middle",
			LastName:   "last",
			Email:      "test@mail.com",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				PostalCode: "4123A",
				Country:    "PH",
			},
			Phone: &pfpb.PhoneNumber{
				CountryCode: "43",
				Number:      "12345",
			},
			Mobile: &pfpb.PhoneNumber{
				CountryCode: "43",
				Number:      "54321",
			},
		},
		Email:       "test@mail.com",
		Nationality: "PH",
	},
}

var jprDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.JPRCode,
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
}

var rmDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.RMCode,
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:   "GOVERNMENT_ISSUED_ID",
			Number: "B83180608851",
			Issued: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Expiration: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation:   "OTH",
			OccupationID: "1",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		SendingReasonID:    "1133",
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		Gender:             tpb.Gender_Male,
		Nationality:        "PH",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
		DeviceID:  "5500",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var riaDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  static.RIACode,
	ControlNumber: "5142096205",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:   "GOVERNMENT_ISSUED_ID",
			Number: "B83180608851",
			Expiration: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Issued: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		Gender:             tpb.Gender_Male,
		Nationality:        "PH",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
		DeviceID:  "5500",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var mbDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  static.MBCode,
	ControlNumber: "5142096205",
	OrderID:       "5142096205",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:   "GOVERNMENT_ISSUED_ID",
			Number: "B83180608851",
			Expiration: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Issued: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		Gender:             tpb.Gender_Male,
	},
	Remitter: &tpb.Contact{
		FirstName:  "John",
		MiddleName: "Michael",
		LastName:   "Doe",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var bpDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  static.BPICode,
	ControlNumber: "5142096205",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var uscDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  static.USSCCode,
	ControlNumber: "5142096205",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
}

var icDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  static.ICCode,
	ControlNumber: "5142096205",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
}

var untDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.UNTCode,
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				State:      "state",
				Zone:       "zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:    "LICENSE",
			Number:  "B83180608851",
			Country: "PH",
			Expiration: &pfpb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Issued: &pfpb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		Gender:             tpb.Gender_Male,
		Nationality:        "PH",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
		DeviceID:  "5500",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var cebiDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.CEBINTCode,
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				State:      "state",
				Zone:       "zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:       "LICENSE",
			Number:     "B83180608851",
			Country:    "PH",
			Issued:     &pfpb.Date{Year: "2020", Month: "12", Day: "12"},
			Expiration: &ppb.Date{},
			City:       "city",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		Gender:             tpb.Gender_Male,
		Nationality:        "PH",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
		DeviceID:  "5500",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var ayaCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.AYACode,
	RemitType:    "Send",
	OrderID:      "5142096205",
	Amount: &tpb.SendAmount{
		Amount:              "100000",
		SourceCountry:       "PH",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationCountry:  "PH",
		DestinationAmount:   true,
	},
	Receiver: &tpb.Receiver{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
	},
	Remitter: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "bit",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &pfpb.Identification{
			Type:    "M",
			Country: "PH",
			Number:  "24023497AB0877AAB20000",
		},
	},
	Agent: &tpb.Agent{
		UserID:    12345,
		IPAddress: "::1",
	},
}

var ayaDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.AYACode,
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var ieCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.IECode,
	RemitType:    "Send",
	OrderID:      "5142096205",
	Amount: &tpb.SendAmount{
		Amount:              "100000",
		SourceCountry:       "PH",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationCountry:  "PH",
		DestinationAmount:   true,
	},
	Receiver: &tpb.Receiver{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
		Identification: &pfpb.Identification{
			Type:    "M",
			Country: "PH",
			Number:  "24023497AB0877AAB20000",
		},
	},
	Remitter: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "bit",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &pfpb.Identification{
			Type:    "M",
			Country: "PH",
			Number:  "24023497AB0877AAB20000",
		},
	},
	Agent: &tpb.Agent{
		UserID:    12345,
		IPAddress: "::1",
	},
}

var ieDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.IECode,
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

// var phCreateReq = &tpb.CreateRemitRequest{
// 	RemitPartner: static.PerahubRemit,
// 	RemitType:    "Send",
// 	OrderID:      "51420962056",
// 	Amount: &tpb.SendAmount{
// 		Amount:              "100000",
// 		SourceCountry:       "PH",
// 		SourceCurrency:      "PHP",
// 		DestinationCurrency: "PHP",
// 		DestinationCountry:  "PH",
// 		DestinationAmount:   true,
// 	},
// 	Receiver: &tpb.Receiver{
// 		ContactInfo: &tpb.Contact{
// 			FirstName:  "Blanche",
// 			MiddleName: "G",
// 			LastName:   "Soto",
// 			Address: &tpb.Address{
// 				Address1:   "addr1",
// 				Address2:   "addr2",
// 				City:       "city",
// 				Province:   "province",
// 				PostalCode: "12345",
// 				Country:    "PH",
// 				Zone:       "zone",
// 			},
// 			Phone: &ppb.PhoneNumber{
// 				CountryCode: "62",
// 				Number:      "12345678",
// 			},
// 			Mobile: &ppb.PhoneNumber{
// 				CountryCode: "62",
// 				Number:      "12345678",
// 			},
// 		},
// 		Identification: &pfpb.Identification{
// 			Type:    "M",
// 			Country: "PH",
// 			Number:  "24023497AB0877AAB20000",
// 		},
// 	},
// 	Remitter: &tpb.UserKYC{
// 		ContactInfo: &tpb.Contact{
// 			FirstName:  "Mittie",
// 			MiddleName: "O",
// 			LastName:   "Sauer",
// 			Address: &tpb.Address{
// 				Address1:   "addr1",
// 				Address2:   "addr2",
// 				City:       "city",
// 				Province:   "province",
// 				PostalCode: "12345",
// 				Country:    "PH",
// 				Zone:       "zone",
// 			},
// 			Phone: &ppb.PhoneNumber{
// 				CountryCode: "62",
// 				Number:      "12345678",
// 			},
// 		},
// 		PartnerMemberID:    "7712780",
// 		BirthCountry:       "PH",
// 		BirthPlace:         "bit",
// 		SourceFunds:        "Salary/Income",
// 		ReceiverRelation:   "Family",
// 		TransactionPurpose: "Gift",
// 		Birthdate: &pfpb.Date{
// 			Year:  "1950",
// 			Month: "12",
// 			Day:   "12",
// 		},
// 		Employment: &tpb.Employment{
// 			Occupation: "Unemployed",
// 		},
// 		Identification: &pfpb.Identification{
// 			Type:    "M",
// 			Country: "PH",
// 			Number:  "24023497AB0877AAB20000",
// 		},
// 	},
// 	Agent: &tpb.Agent{
// 		UserID:    12345,
// 		IPAddress: "::1",
// 	},
// }

var phDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.PerahubRemit,
	ControlNumber:    "51420962056",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "Blanche",
			MiddleName: "G",
			LastName:   "Soto",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &pfpb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &pfpb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

func newTestSvc(t *testing.T, store *postgres.Storage, crPtnrErr, dsbPtnrErr, cfCrPtnrErr, cfDsbPtnrErr bool) *Svc {
	cl := perahub.NewTestHTTPMock(store, perahub.MockConfig{
		CrPtnrErr:    crPtnrErr,
		DsbPtnrErr:   dsbPtnrErr,
		CfCrPtnrErr:  cfCrPtnrErr,
		CfDsbPtnrErr: cfDsbPtnrErr,
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
		map[string]perahub.OAuthCreds{
			static.WISECode: {
				ClientID:     "wise-id",
				ClientSecret: "wise-secret",
			},
		},
	)
	if err != nil {
		t.Fatal("setting up perahub integration: ", err)
	}

	st := static.New(ph, store)
	rms := []remit.Remitter{wu.New(store, ph, st, true), ir.New(store, ph, st), tf.New(store, ph, st), rm.New(store, ph, st), ria.New(store, ph, st), mb.New(store, ph, st), bp.New(store, ph, st), usc.New(store, ph, st), ic.New(store, ph, st), jr.New(store, ph, st), ws.New(store, ph, st), unt.New(store, ph, st), ceb.New(store, ph, st), cebi.New(store, ph, st), aya.New(store, ph, st), ie.New(store, ph, st), prmt.New(store, ph, st)}
	rt, err := remit.New(store, ph, st, rms)
	if err != nil {
		t.Fatal("setting up remit core: ", err)
	}

	tvs := NewValidators(q)
	h, err := New(rt, st, tvs)
	if err != nil {
		t.Fatal("setting up service: ", err)
	}
	return h
}
