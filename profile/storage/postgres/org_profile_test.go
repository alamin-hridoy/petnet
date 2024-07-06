package postgres_test

import (
	"context"
	"database/sql"
	"sort"
	"testing"
	"time"

	"brank.as/petnet/profile/core/profile"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus"
)

// todo(robin): split up to smaller tests
func TestCRUProfile(t *testing.T) {
	ts := newTestStorage(t)
	st := profile.New(ts)

	oid := "10000000-0000-0000-0000-000000000000"
	oid2 := "20000000-0000-0000-0000-000000000000"
	uid := "11000000-0000-0000-0000-000000000000"
	uid2 := "22000000-0000-0000-0000-000000000000"
	dsaCode := "dsa"
	want := []storage.OrgProfile{
		{
			OrgID:            oid,
			UserID:           uid,
			OrgType:          1,
			Status:           1,
			RiskScore:        1,
			TransactionTypes: "Digital,OverTheCounter",
			Partner:          "EasternUnion",
			BusinessInfo: storage.BusinessInfo{
				CompanyName:   "company-name",
				StoreName:     "store-name",
				PhoneNumber:   "123456789",
				FaxNumber:     "123456789",
				Website:       "website",
				CompanyEmail:  "company-email",
				ContactPerson: "contact-person",
				Position:      "position",
				Address: storage.Address{
					Address1:   "address1",
					City:       "city",
					State:      "state",
					PostalCode: "12345",
				},
			},
			AccountInfo: storage.AccountInfo{
				Bank:                    "bank",
				BankAccountNumber:       "bank-acct-no",
				BankAccountHolder:       "bank-acct-holder",
				AgreeTermsConditions:    1,
				AgreeOnlineSupplierForm: 1,
				Currency:                1,
			},
			DateApplied:       sql.NullTime{Time: time.Unix(1515151515, 0), Valid: true},
			ReminderSent:      1,
			DsaCode:           "123",
			TerminalIdOtc:     "456",
			TerminalIdDigital: "789",
			IsProvider:        false,
		},
		{
			OrgID:            oid2,
			UserID:           uid2,
			OrgType:          2,
			Status:           1,
			RiskScore:        1,
			TransactionTypes: "Digital,OverTheCounter",
			Partner:          "EasternUnion",
			BusinessInfo: storage.BusinessInfo{
				CompanyName:   "company-name",
				StoreName:     "store-name",
				PhoneNumber:   "123456789",
				FaxNumber:     "123456789",
				Website:       "website",
				CompanyEmail:  "company-email",
				ContactPerson: "contact-person",
				Position:      "position",
				Address: storage.Address{
					Address1:   "address1",
					City:       "city",
					State:      "state",
					PostalCode: "12345",
				},
			},
			AccountInfo: storage.AccountInfo{
				Bank:                    "bank",
				BankAccountNumber:       "bank-acct-no",
				BankAccountHolder:       "bank-acct-holder",
				AgreeTermsConditions:    1,
				AgreeOnlineSupplierForm: 1,
				Currency:                1,
			},
			DateApplied:       sql.NullTime{Time: time.Unix(1515151515, 0), Valid: true},
			ReminderSent:      1,
			DsaCode:           "1234",
			TerminalIdOtc:     "4567",
			TerminalIdDigital: "7890",
			IsProvider:        true,
		},
	}

	logr := logging.NewLogger(nil)
	logr.SetFormatter(&logrus.JSONFormatter{})
	ctx := logging.WithLogger(context.TODO(), logr)

	pid, err := st.CreateOrgProfile(ctx, want[0])
	if err != nil {
		t.Fatal("CreateOrgProfile: ", err)
	}

	if pid == "" {
		t.Fatal("CreateOrgProfile, profileID should not be empty")
	}
	_, err = ts.CreateOrgProfile(ctx, &want[0])
	if err != storage.Conflict {
		t.Fatal("CreateOrgProfile error should be conflict")
	}
	pid2, err := ts.CreateOrgProfile(ctx, &want[1])
	if err != nil {
		t.Fatal("CreateOrgProfile: ", err)
	}
	if pid2 == "" {
		t.Fatal("CreateOrgProfile, profileID should not be empty")
	}

	got, err := ts.GetOrgProfile(ctx, want[0].OrgID)
	if err != nil {
		t.Fatal("GetOrgProfile: ", err)
	}
	_, errDC := ts.GetProfileByDsaCode(ctx, want[0].DsaCode)
	if errDC != nil {
		t.Fatal("GetOrgProfile: ", errDC)
	}

	tOps := []cmp.Option{
		cmpopts.IgnoreFields(storage.OrgProfile{}, "ID", "Created", "Updated", "Deleted"),
	}
	if !cmp.Equal(&want[0], got, tOps...) {
		t.Fatal("GetOrgProfile (-want +got): ", cmp.Diff(&want[0], got, tOps...))
	}
	if got.ID == "" {
		t.Fatal("GetOrgProfile, profileID should not be empty")
	}
	if got.Created.IsZero() || got.Updated.IsZero() {
		t.Fatal("GetOrgProfile, created and updated shouldn't be empty")
	}
	if got.Deleted.Valid {
		t.Fatal("GetOrgProfile, deleted should be null")
	}

	gotlist, err := ts.GetOrgProfiles(context.TODO(), storage.FilterList{})
	if err != nil {
		t.Fatal("GetOrgProfiles: ", err)
	}
	t.Skip("Clean up tests and stop using 'Fatal' everywhere")

	sort.Slice(want, func(i, j int) bool {
		return want[i].OrgID < want[j].OrgID
	})
	sort.Slice(gotlist, func(i, j int) bool {
		return gotlist[i].OrgID < gotlist[j].OrgID
	})
	for i, pf := range gotlist {
		if !cmp.Equal(want[i], pf, tOps...) {
			t.Fatal("GetOrgProfiles (-want +got): ", cmp.Diff(want[i], pf, tOps...))
		}
		if pf.ID == "" {
			t.Fatal("GetOrgProfiles, profileID should not be empty")
		}
		if pf.Created.IsZero() || pf.Updated.IsZero() {
			t.Fatal("GetOrgProfiles, created and updated shouldn't be empty")
		}
		if pf.Deleted.Valid {
			t.Fatal("GetOrgProfiles, deleted should be null")
		}
	}

	wantup := &storage.OrgProfile{
		OrgID:            oid,
		OrgType:          2,
		Status:           2,
		RiskScore:        2,
		TransactionTypes: "Digital",
		Partner:          "EasternUnion",
		BusinessInfo: storage.BusinessInfo{
			CompanyName:   "company-name-u",
			StoreName:     "store-name-u",
			PhoneNumber:   "123456789-u",
			FaxNumber:     "123456789-u",
			Website:       "website-u",
			CompanyEmail:  "company-email-u",
			ContactPerson: "contact-person-u",
			Position:      "position-u",
			Address: storage.Address{
				Address1:   "address1-u",
				City:       "city-u",
				State:      "state-u",
				PostalCode: "12345-u",
			},
		},
		AccountInfo: storage.AccountInfo{
			Bank:                    "bank-u",
			BankAccountNumber:       "bank-acct-no-u",
			BankAccountHolder:       "bank-acct-holder-u",
			AgreeTermsConditions:    2,
			AgreeOnlineSupplierForm: 2,
			Currency:                2,
		},
		DateApplied:       sql.NullTime{Time: time.Unix(1515151515, 0), Valid: true},
		ReminderSent:      2,
		DsaCode:           "123",
		TerminalIdOtc:     "456",
		TerminalIdDigital: "789",
		Deleted:           sql.NullTime{Time: time.Unix(1515151515, 0), Valid: true},
		IsProvider:        true,
	}
	wantupp := &storage.UpdateOrgProfileOrgIDUserID{
		OldOrgID: oid,
		NewOrgID: "30000000-0000-0000-0000-000000000000",
		UserID:   "32000000-0000-0000-0000-000000000000",
	}
	upid, err := ts.UpdateOrgProfile(context.TODO(), wantup)
	if err != nil {
		t.Fatal("UpdateOrgProfileByOrgID: ", err)
	}
	if upid != pid {
		t.Fatal("profile id mismatch")
	}

	upidd, err := ts.UpdateOrgProfileUserID(context.TODO(), wantupp)
	if err != nil {
		t.Fatal("UpdateOrgProfileUserID: ", err)
	}
	if upidd != pid {
		t.Fatal("profile id mismatch")
	}

	got, err = ts.GetOrgProfile(context.TODO(), oid)
	if err != nil {
		t.Fatal("GetOrgProfileByOrgID: ", err)
	}
	_, err = ts.GetProfileByDsaCode(context.TODO(), dsaCode)
	if err != nil {
		t.Fatal("GetOrgProfileByOrgID: ", err)
	}
	tOps = []cmp.Option{
		cmpopts.IgnoreFields(storage.OrgProfile{}, "ID", "Created", "Updated"),
	}
	if !cmp.Equal(wantup, got, tOps...) {
		t.Fatal("GetOrgProfileByOrgID (-want +got): ", cmp.Diff(wantup, got, tOps...))
	}

	// make sure when set to zero values are not changed in db
	wantup2 := &storage.OrgProfile{
		OrgID:            oid,
		Status:           0,
		RiskScore:        0,
		TransactionTypes: "OverTheCounter",
		Partner:          "EasternUnion",
		BusinessInfo:     storage.BusinessInfo{},
		AccountInfo: storage.AccountInfo{
			AgreeTermsConditions:    0,
			AgreeOnlineSupplierForm: 0,
			Currency:                0,
		},
		DateApplied: sql.NullTime{Time: time.Unix(1515151515, 0), Valid: true},
		Deleted:     sql.NullTime{Time: time.Unix(1515151515, 0), Valid: true},
	}
	_, err = ts.UpdateOrgProfile(context.TODO(), wantup2)
	if err != nil {
		t.Fatal("UpdateOrgProfileByOrgID: ", err)
	}

	got, err = ts.GetOrgProfile(context.TODO(), oid)
	if err != nil {
		t.Fatal("GetOrgProfileByOrgID: ", err)
	}
	_, err = ts.GetProfileByDsaCode(context.TODO(), dsaCode)
	if err != nil {
		t.Fatal("GetOrgProfileByOrgID: ", err)
	}
	tOps = []cmp.Option{
		cmpopts.IgnoreFields(storage.OrgProfile{}, "ID", "Created", "Updated"),
	}
	if !cmp.Equal(wantup, got, tOps...) {
		t.Fatal("GetOrgProfileByOrgID (-want +got): ", cmp.Diff(wantup, got, tOps...))
	}
}
