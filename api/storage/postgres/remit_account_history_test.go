package postgres

import (
	"context"
	"sort"
	"testing"
	"time"

	"brank.as/petnet/api/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func TestRemitToAccountHistory(t *testing.T) {
	ts := newTestStorage(t)
	oid := "10000000-0000-0000-0000-000000000000"
	oid2 := "20000000-0000-0000-0000-000000000000"
	want := []storage.RemitToAccountHistory{
		{
			ID:                          "",
			OrgID:                       oid,
			Partner:                     "Test Partner",
			ReferenceNumber:             "ABCDEFGHIJ",
			TrxDate:                     time.Time{},
			AccountNumber:               "0663718000107",
			Currency:                    "1",
			ServiceCharge:               "100",
			Remarks:                     "Meow Meow",
			Particulars:                 "Transfer Particulars",
			MerchantName:                "Perahub",
			BankID:                      30,
			LocationID:                  371,
			UserID:                      5188,
			CurrencyID:                  "1",
			CustomerID:                  "6925597",
			FormType:                    "OAR",
			FormNumber:                  "HOA0021942",
			TrxType:                     "1",
			RemoteLocationID:            371,
			RemoteUserID:                5188,
			BillerName:                  "BPI",
			TrxTime:                     "16:59:51",
			TotalAmount:                 "110",
			AccountName:                 "Naparate Jerica Reas",
			BeneficiaryAddress:          "1953 PH 3B BLOCK 6 LOT 9 CAMARIN, 175",
			BeneficiaryBirthDate:        "1996-08-10",
			BeneficiaryCity:             "UNIVERSITY OF THE PH",
			BeneficiaryCivil:            "S",
			BeneficiaryCountry:          "Philippines",
			BeneficiaryCustomerType:     "I",
			BeneficiaryFirstName:        "JOSE",
			BeneficiaryLastName:         "SISON",
			BeneficiaryMiddleName:       "MARIA",
			BeneficiaryTin:              "000000000000000",
			BeneficiarySex:              "F",
			BeneficiaryState:            "QUEZON CITY",
			CurrencyCodePrincipalAmount: "PHP",
			PrincipalAmount:             "10",
			RecordType:                  "01",
			RemitterAddress:             "1953 PH 3B BLOCK 6 LOT 9 CAMARIN, 175",
			RemitterBirthDate:           "1996-08-10",
			RemitterCity:                "UNIVERSITY OF THE PH",
			RemitterCivil:               "S",
			RemitterCountry:             "PH",
			RemitterCustomerType:        "I",
			RemitterFirstName:           "Jerica",
			RemitterGender:              "F",
			RemitterID:                  6925597,
			RemitterLastName:            "Naparate",
			RemitterMiddleName:          "Reas",
			RemitterState:               "QUEZON CITY",
			SettlementMode:              "03",
			Notification:                false,
			BeneZipCode:                 "1000",
			Info:                        []byte{},
			Details:                     []byte{},
			TxnStatus:                   "SUCCESS",
			ErrorCode:                   "",
			ErrorMessage:                "",
			ErrorTime:                   "",
			ErrorType:                   "",
			CreatedBy:                   uuid.NewString(),
			UpdatedBy:                   uuid.NewString(),
			Created:                     time.Now(),
			Updated:                     time.Now(),
			SortByColumn:                "",
			SortOrder:                   "",
			Limit:                       0,
			Offset:                      0,
			Total:                       0,
		},
		{
			ID:                          "",
			OrgID:                       oid2,
			Partner:                     "Test Partner",
			ReferenceNumber:             "ABCDEFGHIJ",
			TrxDate:                     time.Time{},
			AccountNumber:               "0663718000107",
			Currency:                    "1",
			ServiceCharge:               "100",
			Remarks:                     "Meow Meow",
			Particulars:                 "Transfer Particulars",
			MerchantName:                "Perahub",
			BankID:                      30,
			LocationID:                  371,
			UserID:                      5188,
			CurrencyID:                  "1",
			CustomerID:                  "6925597",
			FormType:                    "OAR",
			FormNumber:                  "HOA0021942",
			TrxType:                     "1",
			RemoteLocationID:            371,
			RemoteUserID:                5188,
			BillerName:                  "BPI",
			TrxTime:                     "16:59:51",
			TotalAmount:                 "110",
			AccountName:                 "Naparate Jerica Reas",
			BeneficiaryAddress:          "1953 PH 3B BLOCK 6 LOT 9 CAMARIN, 175",
			BeneficiaryBirthDate:        "1996-08-10",
			BeneficiaryCity:             "UNIVERSITY OF THE PH",
			BeneficiaryCivil:            "S",
			BeneficiaryCountry:          "Philippines",
			BeneficiaryCustomerType:     "I",
			BeneficiaryFirstName:        "JOSE",
			BeneficiaryLastName:         "SISON",
			BeneficiaryMiddleName:       "MARIA",
			BeneficiaryTin:              "000000000000000",
			BeneficiarySex:              "F",
			BeneficiaryState:            "QUEZON CITY",
			CurrencyCodePrincipalAmount: "PHP",
			PrincipalAmount:             "10",
			RecordType:                  "01",
			RemitterAddress:             "1953 PH 3B BLOCK 6 LOT 9 CAMARIN, 175",
			RemitterBirthDate:           "1996-08-10",
			RemitterCity:                "UNIVERSITY OF THE PH",
			RemitterCivil:               "S",
			RemitterCountry:             "PH",
			RemitterCustomerType:        "I",
			RemitterFirstName:           "Jerica",
			RemitterGender:              "F",
			RemitterID:                  6925597,
			RemitterLastName:            "Naparate",
			RemitterMiddleName:          "Reas",
			RemitterState:               "QUEZON CITY",
			SettlementMode:              "03",
			Notification:                false,
			BeneZipCode:                 "1000",
			Info:                        []byte{},
			Details:                     []byte{},
			TxnStatus:                   "SUCCESS",
			ErrorCode:                   "",
			ErrorMessage:                "",
			ErrorTime:                   "",
			ErrorType:                   "",
			CreatedBy:                   uuid.NewString(),
			UpdatedBy:                   uuid.NewString(),
			Created:                     time.Now(),
			Updated:                     time.Now(),
			SortByColumn:                "",
			SortOrder:                   "",
			Limit:                       0,
			Offset:                      0,
			Total:                       0,
		},
	}
	logr := logging.NewLogger(nil)
	logr.SetFormatter(&logrus.JSONFormatter{})
	ctx := logging.WithLogger(context.TODO(), logr)
	_, err := ts.CreateRTAHistory(ctx, want[0])
	if err != nil {
		t.Fatal(err)
	}
	_, err = ts.CreateRTAHistory(ctx, want[1])
	if err != nil {
		t.Fatal(err)
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(storage.RemitToAccountHistory{}, "ID", "Created", "Updated", "Total", "Info", "Details"),
	}
	_, err = ts.GetRTAHistory(ctx, want[0].OrgID)
	if err != nil {
		t.Fatalf("Get RTA History = got error %v, want nil", err)
	}
	gotlist, err := ts.ListRTAHistory(ctx, want[0])
	if err != nil {
		t.Fatal(err)
	}
	upData := storage.RemitToAccountHistory{
		OrgID: want[0].OrgID,
	}
	upid, err := ts.UpdateRTAHistory(ctx, upData)
	if err != nil {
		t.Fatal(err)
	}
	if upid.ID != gotlist[0].ID {
		t.Error("id mismatch")
	}
	sort.Slice(want, func(i, j int) bool {
		return want[i].OrgID < want[j].OrgID
	})
	sort.Slice(gotlist, func(i, j int) bool {
		return gotlist[i].OrgID < gotlist[j].OrgID
	})
	for i, pf := range gotlist {
		if !cmp.Equal(want[i], pf, tOps...) {
			t.Error("(-want +got): ", cmp.Diff(want[i], pf, tOps...))
		}
		if pf.OrgID == "" {
			t.Error("org id should not be empty")
		}

	}
}
