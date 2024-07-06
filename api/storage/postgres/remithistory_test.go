package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bojanz/currency"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"brank.as/petnet/api/storage"
)

func TestRemitHistory(t *testing.T) {
	ts := newTestStorage(t)

	amt1, err := currency.NewMinor("100000", "PHP")
	if err != nil {
		t.Fatal(err)
	}
	amt2, err := currency.NewMinor("200000", "PHP")
	if err != nil {
		t.Fatal(err)
	}
	sum1, err := currency.NewMinor("300000", "PHP")
	if err != nil {
		t.Fatal(err)
	}

	amt3, err := currency.NewMinor("200000", "PHP")
	if err != nil {
		t.Fatal(err)
	}
	amt4, err := currency.NewMinor("300000", "PHP")
	if err != nil {
		t.Fatal(err)
	}
	sum2, err := currency.NewMinor("500000", "PHP")
	if err != nil {
		t.Fatal(err)
	}

	rm := storage.Remittance{
		RemcoAltControlNo: "refalt-ref",
		TxnType:           "SO",
		Remitter: storage.Contact{
			FirstName:     "first",
			MiddleName:    "middle",
			LastName:      "last",
			RemcoMemberID: "remco-id",
			Email:         "email@mail.com",
			Address1:      "addr1",
			Address2:      "addr2",
			City:          "city",
			State:         "state",
			PostalCode:    "12345",
			Country:       "PH",
			Province:      "province",
			Zone:          "zone",
			PhoneCty:      "PH",
			Phone:         "123456",
			MobileCty:     "PH",
			Mobile:        "123456",
			Message:       "msg",
		},
		Receiver: storage.Contact{
			FirstName:     "first2",
			MiddleName:    "middle2",
			LastName:      "last2",
			RemcoMemberID: "remco-id2",
			Email:         "email@mail.com2",
			Address1:      "addr12",
			Address2:      "addr22",
			City:          "city2",
			State:         "state2",
			PostalCode:    "123452",
			Country:       "PH2",
			Province:      "province2",
			Zone:          "zone2",
			PhoneCty:      "PH2",
			Phone:         "1234562",
			MobileCty:     "PH2",
			Mobile:        "1234562",
			Message:       "msg2",
		},
		Business: storage.Business{
			Name:      "name",
			Account:   "account",
			ControlNo: "ref-no",
			Country:   "PH",
		},
		Account: storage.Account{
			BIC:     "111",
			AcctNo:  "2222",
			AcctSfx: "suf",
		},
		GrossTotal: storage.GrossTotal{Minor: amt1},
		SourceAmt:  amt1,
		DestAmt:    amt2,
		Taxes: map[string]currency.Minor{
			"tax1": amt1,
			"tax2": amt2,
		},
		Tax: sum1,
		Charges: map[string]currency.Minor{
			"fee1": amt1,
			"fee2": amt2,
		},
		Charge: sum1,
	}

	dsaID := uuid.NewString()
	ptnrs := map[string]struct {
		in        *storage.RemitHistory
		recName   string
		totAmount string
	}{
		"P1": {
			recName:   "a",
			totAmount: "10000",
			in: &storage.RemitHistory{
				DsaOrderID:     uuid.NewString(),
				TxnID:          uuid.NewString(),
				DsaID:          dsaID,
				UserID:         "a",
				RemcoID:        "P1",
				RemcoControlNo: "1",
				Remittance:     rm,
				TxnStagedTime:  sql.NullTime{Time: time.Date(1999, 12, 31, 12, 0, 0, 0, time.Local), Valid: true},
				TxnStatus:      string(storage.SuccessStatus),
				TxnStep:        string(storage.StageStep),
				TransactionType: sql.NullString{
					String: "Digital",
					Valid:  true,
				},
			},
		},
		"P2": {
			recName:   "b",
			totAmount: "20000",
			in: &storage.RemitHistory{
				DsaOrderID:     uuid.NewString(),
				TxnID:          uuid.NewString(),
				DsaID:          dsaID,
				UserID:         "b",
				RemcoID:        "P2",
				RemcoControlNo: "2",
				Remittance:     rm,
				TxnStagedTime:  sql.NullTime{Time: time.Date(1999, 12, 31, 12, 0, 0, 0, time.Local), Valid: true},
				TxnStatus:      string(storage.SuccessStatus),
				TxnStep:        string(storage.ConfirmStep),
				TransactionType: sql.NullString{
					String: "OTC",
					Valid:  true,
				},
			},
		},
		"P3": {
			recName:   "c",
			totAmount: "30000",
			in: &storage.RemitHistory{
				DsaOrderID:     uuid.NewString(),
				TxnID:          uuid.NewString(),
				DsaID:          uuid.NewString(),
				UserID:         "c",
				RemcoID:        "P3",
				RemcoControlNo: "3",
				Remittance:     rm,
				TxnStagedTime:  sql.NullTime{Time: time.Date(1999, 12, 31, 12, 0, 0, 0, time.Local), Valid: true},
				TxnStatus:      string(storage.SuccessStatus),
				TxnStep:        string(storage.ConfirmStep),
				TransactionType: sql.NullString{
					String: "Digital",
					Valid:  true,
				},
			},
		},
		"P4": {
			recName:   "d",
			totAmount: "40000",
			in: &storage.RemitHistory{
				DsaOrderID:     uuid.NewString(),
				TxnID:          uuid.NewString(),
				DsaID:          uuid.NewString(),
				UserID:         "d",
				RemcoID:        "P4",
				RemcoControlNo: "4",
				Remittance:     rm,
				TxnStagedTime:  sql.NullTime{Time: time.Date(1999, 12, 31, 12, 0, 0, 0, time.Local), Valid: true},
				TxnStatus:      string(storage.SuccessStatus),
				TxnStep:        string(storage.ConfirmStep),
				TransactionType: sql.NullString{
					String: "Digital",
					Valid:  true,
				},
			},
		},
		"P5": {
			recName:   "e",
			totAmount: "50000",
			in: &storage.RemitHistory{
				DsaOrderID:     uuid.NewString(),
				TxnID:          uuid.NewString(),
				DsaID:          uuid.NewString(),
				UserID:         "e",
				RemcoID:        "P5",
				RemcoControlNo: "5",
				Remittance:     rm,
				TxnStagedTime:  sql.NullTime{Time: time.Date(1999, 12, 31, 12, 0, 0, 0, time.Local), Valid: true},
				TxnStatus:      string(storage.SuccessStatus),
				TxnStep:        string(storage.StageStep),
				TransactionType: sql.NullString{
					String: "Digital",
					Valid:  true,
				},
			},
		},
		"P6": {
			recName:   "f",
			totAmount: "60000",
			in: &storage.RemitHistory{
				DsaOrderID:     uuid.NewString(),
				TxnID:          uuid.NewString(),
				DsaID:          uuid.NewString(),
				UserID:         "f",
				RemcoID:        "P6",
				RemcoControlNo: "6",
				Remittance:     rm,
				TxnStagedTime:  sql.NullTime{Time: time.Date(1999, 12, 31, 12, 0, 0, 0, time.Local), Valid: true},
				TxnStatus:      string(storage.FailStatus),
				TxnStep:        string(storage.StageStep),
				ErrorCode:      "400",
				ErrorMsg:       "error",
				ErrorTime:      sql.NullTime{Time: time.Date(1999, 12, 31, 12, 0, 0, 0, time.Local), Valid: true},
				ErrorType:      "NONEX",
				TransactionType: sql.NullString{
					String: "Digital",
					Valid:  true,
				},
			},
		},
		"P7": {
			recName:   "g",
			totAmount: "70000",
			in: &storage.RemitHistory{
				DsaOrderID:     uuid.NewString(),
				TxnID:          uuid.NewString(),
				DsaID:          uuid.NewString(),
				UserID:         "g",
				RemcoID:        "P7",
				RemcoControlNo: "7",
				Remittance:     rm,
				TxnStagedTime:  sql.NullTime{Time: time.Date(1999, 12, 31, 12, 0, 0, 0, time.Local), Valid: true},
				TxnStatus:      string(storage.FailStatus),
				TxnStep:        string(storage.ConfirmStep),
				ErrorCode:      "400",
				ErrorMsg:       "error",
				ErrorTime:      sql.NullTime{Time: time.Date(1999, 12, 31, 12, 0, 0, 0, time.Local), Valid: true},
				ErrorType:      "DRP",
				TransactionType: sql.NullString{
					String: "Digital",
					Valid:  true,
				},
			},
		},
	}

	o := cmp.Options{
		cmpopts.IgnoreFields(storage.RemitHistory{}, "Updated"),
	}
	for _, ptnr := range ptnrs {
		srcAmt, err := currency.NewMinor(ptnr.totAmount, "PHP")
		if err != nil {
			t.Fatal(err)
		}

		ptnr.in.Remittance.SourceAmt = srcAmt
		ptnr.in.Remittance.Receiver.FirstName = ptnr.recName
		rh, err := ts.CreateRemitHistory(context.TODO(), *ptnr.in)
		if err != nil {
			t.Fatalf("CreateRemitHistory() = got error %v, want nil", err)
		}
		if rh.Updated.IsZero() {
			t.Fatal("CreateRemitHistory() = returned empty Updated date")
		}

		grh, err := ts.GetRemitHistory(context.TODO(), ptnr.in.TxnID)
		if err != nil {
			t.Fatalf("GetRemitHistory() = got error %v, want nil", err)
		}

		if !cmp.Equal(*ptnr.in, *grh, o) {
			t.Fatal(cmp.Diff(*ptnr.in, *grh, o))
		}
		if grh.Updated.IsZero() {
			t.Fatal("GetRemitHistory() = returned empty Updated date")
		}
		_, err = ts.GetTransactionReport(context.TODO(), &storage.LRHFilter{
			TxnStatus:       ptnr.in.TxnStatus,
			TxnStep:         ptnr.in.TxnStep,
			DsaID:           dsaID,
			From:            time.Now(),
			Until:           time.Now(),
			Transactiontype: "DIGITAL",
		})
		if err != nil {
			t.Fatalf("GetTransactionReport() = got error %v, want nil", err)
		}
	}
	ptnrs["P1"] = struct {
		in        *storage.RemitHistory
		recName   string
		totAmount string
	}{
		in: &storage.RemitHistory{
			TxnID:            ptnrs["P1"].in.TxnID,
			DsaID:            dsaID,
			DsaOrderID:       ptnrs["P1"].in.DsaOrderID,
			UserID:           "a",
			RemcoID:          "P1",
			TxnStagedTime:    sql.NullTime{Time: time.Date(1999, 12, 31, 12, 0, 0, 0, time.Local), Valid: true},
			TxnCompletedTime: sql.NullTime{Time: time.Date(2000, 0o1, 0o1, 12, 0o0, 0o0, 0o00000000, time.Local), Valid: true},
			TxnStatus:        string(storage.SuccessStatus),
			TxnStep:          string(storage.ConfirmStep),
			Remittance: storage.Remittance{
				RemcoAltControlNo: "remaltref-u",
				TxnType:           "SO",
				Remitter: storage.Contact{
					FirstName:     "first-u",
					MiddleName:    "middle-u",
					LastName:      "last-u",
					RemcoMemberID: "remco-id-u",
					Email:         "email@mail.com-u",
					Address1:      "addr1-u",
					Address2:      "addr2-u",
					City:          "city-u",
					State:         "state-u",
					PostalCode:    "12345-u",
					Country:       "PH-u",
					Province:      "province-u",
					Zone:          "zone-u",
					PhoneCty:      "PH-u",
					Phone:         "123456-u",
					MobileCty:     "PH-u",
					Mobile:        "123456-u",
					Message:       "msg-u",
				},
				Receiver: storage.Contact{
					MiddleName:    "middle2-u",
					LastName:      "last2-u",
					RemcoMemberID: "remco-id2-u",
					Email:         "email@mail.com2-u",
					Address1:      "addr12-u",
					Address2:      "addr22-u",
					City:          "city2-u",
					State:         "state2-u",
					PostalCode:    "123452-u",
					Country:       "PH2-u",
					Province:      "province2-u",
					Zone:          "zone2-u",
					PhoneCty:      "PH2-u",
					Phone:         "1234562-u",
					MobileCty:     "PH2-u",
					Mobile:        "1234562-u",
					Message:       "msg2-u",
				},
				Business: storage.Business{
					Name:      "name-u",
					Account:   "account-u",
					ControlNo: "ref-no-u",
					Country:   "PH-u",
				},
				Account: storage.Account{
					BIC:     "111-u",
					AcctNo:  "2222-u",
					AcctSfx: "suf-u",
				},
				GrossTotal: storage.GrossTotal{Minor: amt1},
				SourceAmt:  amt1,
				DestAmt:    amt2,
				Taxes: map[string]currency.Minor{
					"tax1": amt3,
					"tax2": amt4,
				},
				Tax: sum2,
				Charges: map[string]currency.Minor{
					"fee1": amt3,
					"fee2": amt4,
				},
				Charge: sum2,
			},
		},
	}

	ptnrs["P2"].in.TxnCompletedTime = sql.NullTime{Time: time.Date(2000, 1, 2, 12, 0, 0, 0, time.Local), Valid: true}
	ptnrs["P2"] = struct {
		in        *storage.RemitHistory
		recName   string
		totAmount string
	}{
		in: ptnrs["P2"].in,
	}

	ptnrs["P3"].in.TxnCompletedTime = sql.NullTime{Time: time.Date(2000, 1, 3, 12, 0, 0, 0, time.Local), Valid: true}
	ptnrs["P3"] = struct {
		in        *storage.RemitHistory
		recName   string
		totAmount string
	}{
		in: ptnrs["P3"].in,
	}

	ptnrs["P4"].in.TxnCompletedTime = sql.NullTime{Time: time.Date(2000, 1, 4, 12, 0, 0, 0, time.Local), Valid: true}
	ptnrs["P4"] = struct {
		in        *storage.RemitHistory
		recName   string
		totAmount string
	}{
		in: ptnrs["P4"].in,
	}

	ptnrs["P5"].in.TxnCompletedTime = sql.NullTime{Time: time.Date(2000, 1, 5, 12, 0, 0, 0, time.Local), Valid: true}
	ptnrs["P5"] = struct {
		in        *storage.RemitHistory
		recName   string
		totAmount string
	}{
		in: ptnrs["P5"].in,
	}

	for k, v := range ptnrs {
		rhu, err := ts.UpdateRemitHistory(context.TODO(), *v.in)
		if err != nil {
			t.Fatalf("UpdateRemitHistory() = got error %v, want nil", err)
		}
		if k == "P1" {
			if !cmp.Equal(*ptnrs["P1"].in, *rhu, o) {
				t.Error(cmp.Diff(*ptnrs["P1"].in, *rhu, o))
			}
		}
	}

	h, err := ts.GetRemitHistory(context.TODO(), ptnrs["P1"].in.TxnID)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(*ptnrs["P1"].in, *h, o) {
		t.Error(cmp.Diff(*ptnrs["P1"].in, *h, o))
	}

	tests := []struct {
		name string
		f    storage.LRHFilter
		want []string
	}{
		{
			name: "All",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1", "P2", "P3", "P4"},
		},
		{
			name: "Ref",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				ControlNo: []string{
					ptnrs["P1"].in.RemcoControlNo,
					ptnrs["P2"].in.RemcoControlNo,
				},
				TxnStep:   string(storage.ConfirmStep),
				TxnStatus: string(storage.SuccessStatus),
			},
			want: []string{"P1", "P2"},
		},
		{
			name: "DsaID",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				DsaOrgID:     dsaID,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1", "P2"},
		},
		{
			name: "DsaID/Ref",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				DsaOrgID:     dsaID,
				ControlNo:    []string{ptnrs["P1"].in.RemcoControlNo},
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1"},
		},
		{
			name: "From",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				From:         time.Date(2000, 1, 2, 12, 0, 0, 0, time.Local),
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P2", "P3", "P4"},
		},
		{
			name: "Until",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				Until:        time.Date(2000, 1, 2, 12, 0, 0, 0, time.Local),
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1", "P2"},
		},
		{
			name: "From Until",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				From:         time.Date(2000, 1, 2, 12, 0, 0, 0, time.Local),
				Until:        time.Date(2000, 1, 3, 12, 0, 0, 0, time.Local),
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P2", "P3"},
		},
		{
			name: "After/Before/DsaID Within Date Range",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				From:         time.Date(2000, 1, 2, 12, 0, 0, 0, time.Local),
				Until:        time.Date(2000, 1, 3, 12, 0, 0, 0, time.Local),
				DsaOrgID:     dsaID,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P2"},
		},
		{
			name: "After/Before/DsaID Outside Date Range",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				From:         time.Date(2000, 1, 3, 12, 0, 0, 0, time.Local),
				Until:        time.Date(2000, 1, 4, 12, 0, 0, 0, time.Local),
				DsaOrgID:     dsaID,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{},
		},
		{
			name: "After/Before/Ref Within Date Range",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				From:         time.Date(2000, 1, 3, 12, 0, 0, 0, time.Local),
				Until:        time.Date(2000, 1, 4, 12, 0, 0, 0, time.Local),
				ControlNo:    []string{ptnrs["P4"].in.RemcoControlNo},
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P4"},
		},
		{
			name: "After/Before/Ref Outside Date Range",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				From:         time.Date(2000, 1, 1, 12, 0, 0, 0, time.Local),
				Until:        time.Date(2000, 1, 2, 12, 0, 0, 0, time.Local),
				ControlNo:    []string{ptnrs["P4"].in.RemcoControlNo},
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{},
		},
		{
			name: "Partner",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				Partner:      "P1",
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1"},
		},
		{
			name: "Partner/After/Before/Ref Within Date Range",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				From:         time.Date(2000, 1, 2, 12, 0, 0, 0, time.Local),
				Until:        time.Date(2000, 1, 3, 12, 0, 0, 0, time.Local),
				ControlNo:    []string{ptnrs["P3"].in.RemcoControlNo},
				Partner:      "P3",
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P3"},
		},
		{
			name: "Partner/After/Before/Ref Outside Date Range",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				From:         time.Date(2000, 1, 1, 12, 0, 0, 0, time.Local),
				Until:        time.Date(2000, 1, 2, 12, 0, 0, 0, time.Local),
				ControlNo:    []string{ptnrs["P3"].in.RemcoControlNo},
				Partner:      "P3",
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{},
		},
		{
			name: "ExcludePartner",
			f: storage.LRHFilter{
				SortByColumn:   storage.UserIDCol,
				SortOrder:      storage.Asc,
				ExcludePartner: "P1",
				TxnStep:        string(storage.ConfirmStep),
				TxnStatus:      string(storage.SuccessStatus),
			},
			want: []string{"P2", "P3", "P4"},
		},
		{
			name: "ExcludePartner/After/Before",
			f: storage.LRHFilter{
				SortByColumn:   storage.UserIDCol,
				SortOrder:      storage.Asc,
				ExcludePartner: "P1",
				From:           time.Date(2000, 1, 1, 12, 0, 0, 0, time.Local),
				Until:          time.Date(2000, 1, 2, 12, 0, 0, 0, time.Local),
				TxnStep:        string(storage.ConfirmStep),
				TxnStatus:      string(storage.SuccessStatus),
			},
			want: []string{"P2"},
		},
		{
			name: "ExcludePartner/DSAOrgID",
			f: storage.LRHFilter{
				SortByColumn:   storage.UserIDCol,
				SortOrder:      storage.Asc,
				ExcludePartner: "P1",
				DsaOrgID:       dsaID,
				From:           time.Date(2000, 1, 1, 12, 0, 0, 0, time.Local),
				Until:          time.Date(2000, 1, 2, 12, 0, 0, 0, time.Local),
				TxnStep:        string(storage.ConfirmStep),
				TxnStatus:      string(storage.SuccessStatus),
			},
			want: []string{"P2"},
		},
		{
			name: "SortByColumn Partner Asc",
			f: storage.LRHFilter{
				SortByColumn: storage.PartnerCol,
				SortOrder:    storage.Asc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1", "P2", "P3", "P4"},
		},
		{
			name: "SortByColumn Partner Desc",
			f: storage.LRHFilter{
				SortByColumn: storage.PartnerCol,
				SortOrder:    storage.Desc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P4", "P3", "P2", "P1"},
		},
		{
			name: "SortByColumn RefNo Asc",
			f: storage.LRHFilter{
				SortByColumn: storage.ControlNumberCol,
				SortOrder:    storage.Asc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1", "P2", "P3", "P4"},
		},
		{
			name: "SortByColumn RefNo Desc",
			f: storage.LRHFilter{
				SortByColumn: storage.ControlNumberCol,
				SortOrder:    storage.Desc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P4", "P3", "P2", "P1"},
		},
		{
			name: "SortByColumn RemittedTo Asc",
			f: storage.LRHFilter{
				SortByColumn: storage.RemittedToCol,
				SortOrder:    storage.Asc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1", "P2", "P3", "P4"},
		},
		{
			name: "SortByColumn RemittedTo Desc",
			f: storage.LRHFilter{
				SortByColumn: storage.RemittedToCol,
				SortOrder:    storage.Desc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P4", "P3", "P2", "P1"},
		},
		{
			name: "SortByColumn TotalRemittedAmount Asc",
			f: storage.LRHFilter{
				SortByColumn: storage.TotalRemittedAmountCol,
				SortOrder:    storage.Asc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1", "P2", "P3", "P4"},
		},
		{
			name: "SortByColumn TotalRemittedAmount Desc",
			f: storage.LRHFilter{
				SortByColumn: storage.TotalRemittedAmountCol,
				SortOrder:    storage.Desc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P4", "P3", "P2", "P1"},
		},
		{
			name: "SortByColumn TransactionTime Asc",
			f: storage.LRHFilter{
				SortByColumn: storage.TransactionTimeCol,
				SortOrder:    storage.Asc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1", "P2", "P3", "P4"},
		},
		{
			name: "SortByColumn TransactionTime Desc",
			f: storage.LRHFilter{
				SortByColumn: storage.TransactionTimeCol,
				SortOrder:    storage.Desc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P4", "P3", "P2", "P1"},
		},
		{
			name: "Limit Asc",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				Limit:        1,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P1"},
		},
		{
			name: "Limit Desc",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Desc,
				Limit:        1,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P4"},
		},
		{
			name: "Offset Asc",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Asc,
				Offset:       2,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P3", "P4"},
		},
		{
			name: "Offset Desc",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Desc,
				Offset:       2,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P2", "P1"},
		},
		{
			name: "Limit/Offset",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Desc,
				Limit:        1,
				Offset:       2,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P2"},
		},
		{
			name: "Stage Success Txns",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Desc,
				TxnStep:      string(storage.StageStep),
				TxnStatus:    string(storage.SuccessStatus),
			},
			want: []string{"P5"},
		},
		{
			name: "Stage Fail Txns",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Desc,
				TxnStep:      string(storage.StageStep),
				TxnStatus:    string(storage.FailStatus),
			},
			want: []string{"P6"},
		},
		{
			name: "Confirm Fail Txns",
			f: storage.LRHFilter{
				SortByColumn: storage.UserIDCol,
				SortOrder:    storage.Desc,
				TxnStep:      string(storage.ConfirmStep),
				TxnStatus:    string(storage.FailStatus),
			},
			want: []string{"P7"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			grh, err := ts.ListRemitHistory(context.TODO(), test.f)
			if err != nil {
				t.Fatalf("ListRemitHistory() = got error %v, want nil", err)
			}

			ps := make([]string, len(grh))
			for i, h := range grh {
				ps[i] = h.RemcoID
			}
			if !cmp.Equal(test.want, ps) {
				t.Fatal(cmp.Diff(test.want, ps))
			}
		})
	}
}
