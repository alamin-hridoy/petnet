package microinsurance

import (
	"context"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	miCore "brank.as/petnet/api/core/microinsurance"
	micins_int "brank.as/petnet/api/integration/microinsurance"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/gunk/drp/v1/microinsurance"
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

func newTestSvc(t *testing.T, st *postgres.Storage) *Svc {
	cl := perahub.NewTestHTTPMock(st, perahub.MockConfig{})
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

	micInsBaseUrl, _ := url.Parse("https://privatedrp.dev.perahub.com.ph/v1/insurance/ruralnet/")

	miClient := micins_int.NewMicroInsuranceClient(ph, micInsBaseUrl)

	return NewMicroInsuranceSvc(miCore.NewMicroInsuranceCoreSvc(st, miClient))
}

func TestMicroInsuranceTransactions(t *testing.T) {
	st := newTestStorage(t)
	s := newTestSvc(t, st)
	uid := uuid.New().String()
	oid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())
	traceNumber := ""

	for _, tc := range []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Transact",
			testFunc: func(t *testing.T) {
				req := getMockTransactReq()
				trans, err := s.Transact(ctx, req)

				require.NotNil(t, trans)
				require.Nil(t, err)

				traceNumber = trans.TraceNumber
			},
		},
		{
			name: "GetReprint",
			testFunc: func(t *testing.T) {
				require.NotEmpty(t, traceNumber)

				rpRes, err := s.GetReprint(ctx, &microinsurance.GetReprintRequest{
					TraceNumber: traceNumber,
				})

				require.NotNil(t, rpRes)
				require.Nil(t, err)
			},
		},
		{
			name: "GetTransactionList",
			testFunc: func(t *testing.T) {
				now := time.Now()
				listRes, err := s.GetTransactionList(ctx, &microinsurance.GetTransactionListRequest{
					DateFrom: now.Add(-time.Hour * 24).Format("2006-01-02"),
					DateTo:   now.Format("2006-01-02"),
					OrgID:    oid,
				})

				require.NotNil(t, listRes)
				require.Nil(t, err)
				require.False(t, len(listRes.Transactions) > 0)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

func getMockTransactReq() *microinsurance.TransactRequest {
	return &microinsurance.TransactRequest{
		Coy:              "drp",
		LocationID:       "RF24058023",
		UserCode:         "RF2405800123",
		TrxDate:          "2022-06-02",
		PromoAmount:      20,
		PromoCode:        "TESTPROMO",
		Amount:           "460.80",
		CoverageCount:    "12",
		ProductCode:      "UCPB04",
		ProcessingBranch: "000",
		ProcessedBy:      "designex",
		UserEmail:        "test@test.com",
		LastName:         "Test",
		FirstName:        "Test",
		MiddleName:       "test",
		Gender:           "M",
		Birthdate:        "1997-03-26",
		MobileNumber:     "234234",
		ProvinceCode:     "1",
		CityCode:         "1",
		Address:          "Test",
		MaritalStatus:    "M",
		Occupation:       "Test",
		CardNumber:       "test",
		NumberUnits:      "3",
		Beneficiaries: []*microinsurance.Person{
			{
				LastName:      "testb",
				FirstName:     "testb",
				MiddleName:    "",
				NoMiddleName:  true,
				ContactNumber: "23434334",
				BirthDate:     "1974-07-08",
				Relationship:  "PARENT",
			},
		},
		Dependents: []*microinsurance.Person{
			{
				LastName:      "testd",
				FirstName:     "testd",
				MiddleName:    "testd",
				NoMiddleName:  false,
				ContactNumber: "324524",
				BirthDate:     "1974-07-08",
				Relationship:  "SPOUSE",
			},
		},
	}
}
