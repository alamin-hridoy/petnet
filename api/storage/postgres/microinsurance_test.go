package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"brank.as/petnet/api/storage"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

func TestMicroInsuranceHistory(t *testing.T) {
	ts := newTestStorage(t)
	person := storage.MicroInsurancePerson{
		FirstName:     "Lily",
		LastName:      "Wood",
		MiddleName:    "",
		NoMiddleName:  true,
		ContactNumber: "09638832662",
		BirthDate:     "1974-07-08",
		Relationship:  "PAR",
	}

	bPersons, err := json.Marshal([]storage.MicroInsurancePerson{person})
	if err != nil {
		bPersons = []byte{}
	}

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		in := &storage.MicroInsuranceHistory{
			ID:               uuid.NewString(),
			DsaID:            uuid.NewString(),
			Coy:              "drp",
			LocationID:       "RF24058023",
			UserCode:         "RF2405800123",
			TrxDate:          time.Now(),
			PromoAmount:      "0",
			PromoCode:        "",
			Amount:           "500",
			CoverageCount:    "12",
			ProductCode:      "UCPB04",
			ProcessingBranch: "000",
			ProcessedBy:      "designex",
			UserEmail:        "test@petnet.com",
			LastName:         "Fiaschi",
			FirstName:        "Julian",
			MiddleName:       "Meoni",
			Gender:           "M",
			Birthdate:        time.Now(),
			MobileNumber:     "09745488965",
			ProvinceCode:     "1",
			CityCode:         "1",
			Address:          "1484 Rikid Grove",
			MaritalStatus:    "M",
			Occupation:       "VICE PRESIDENT",
			CardNumber:       "",
			NumberUnits:      "1",
			Beneficiaries:    bPersons,
			Dependents:       bPersons,
			TrxStatus:        string(storage.TRANSACTION_FAIL),
			TraceNumber:      sql.NullString{Valid: true, String: "INSPNIHEDK6795317"},
			InsuranceDetails: []byte{},
			ErrorCode:        "07",
			ErrorMsg:         "Invalid Age",
			ErrorType:        "MICROINSURANCE",
			ErrorTime:        sql.NullTime{Valid: true, Time: time.Now()},
			OrgID:            uuid.NewString(),
		}

		r, err := ts.CreateMicroInsuranceHistory(context.TODO(), *in)
		if err != nil {
			t.Fatalf("CreateMicroInsuranceHistory() = got error %v, want nil", err)
		}

		if r.ID == "" {
			t.Fatal("CreateMicroInsuranceHistory() = returned empty ID")
		}

		getIDRes, err := ts.GetMicroInsuranceHistoryByID(context.TODO(), r.ID)
		if err != nil {
			t.Fatalf("GetMicroInsuranceHistoryByID() = got error %v, want nil", err)
		}

		if !cmp.Equal(*r, *getIDRes) {
			t.Fatal(cmp.Diff(*r, *getIDRes))
		}

		getTraceNoRes, err := ts.GetMicroInsuranceHistoryByTraceNumber(context.TODO(), r.TraceNumber.String)
		if err != nil {
			t.Fatalf("GetMicroInsuranceHistoryByTraceNumber() = got error %v, want nil", err)
		}

		if !cmp.Equal(*r, *getTraceNoRes) {
			t.Fatal(cmp.Diff(*r, *getTraceNoRes))
		}

		listRes, err := ts.ListMicroInsuranceHistory(context.TODO(), storage.MicroInsuranceFilter{
			TraceNumber: r.TraceNumber.String,
			DsaID:       r.DsaID,
			UserCode:    r.UserCode,
			TrxStatus:   r.TrxStatus,
			Limit:       10,
			Offset:      0,
			OrgID:       r.OrgID,
		})
		if err != nil {
			t.Fatalf("ListMicroInsuranceHistory() = got error %v, want nil", err)
		}

		opt := cmpopts.IgnoreFields(storage.MicroInsuranceHistory{}, "Total")
		if !cmp.Equal(*r, listRes[0], opt) {
			t.Fatal(cmp.Diff(*r, listRes[0]))
		}

		toUpdate := listRes[0]
		toUpdate.TrxStatus = "SUCCESS"
		toUpdate.InsuranceDetails, _ = json.Marshal(&migunk.Insurance{
			SessionID:      "34234234",
			StatusCode:     "24334",
			StatusDesc:     "Test",
			InsProductID:   "TESTP1",
			InsProductDesc: "TEST prod desc",
			TrnDate:        "06/07/2022",
			TrnAmount:      330,
			TraceNumber:    r.TraceNumber.String,
		})

		upRes, err := ts.UpdateMicroInsuranceHistoryStatusByTraceNumber(context.TODO(), toUpdate)
		require.Nil(t, err)

		require.NotNil(t, upRes)
		assert.Equal(t, r.TraceNumber, upRes.TraceNumber)
		assert.Equal(t, r.ID, upRes.ID)
		assert.Equal(t, "SUCCESS", upRes.TrxStatus)
	})
}
