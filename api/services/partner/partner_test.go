package partner

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	pc "brank.as/petnet/api/core/partner"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/services"
	"brank.as/petnet/api/storage/postgres"
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
	ph, err := perahub.New(cl,
		"dev",
		"https://newkycgateway.dev.perahub.com.ph/gateway/",
		"https://privatedrp.dev.perahub.com.ph/v1/remit/nonex/",
		"https://privatedrp.dev.perahub.com.ph/v1/billspay/wrapper/api/",
		"https://privatedrp.dev.perahub.com.ph/v1/billspay/",
		"https://privatedrp.dev.perahub.com.ph/v1/transactions/api/",
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
	m := &services.Mock{}
	return New(static.New(ph, st), pc.New(st, ph), m, m, m, NewValidators()), m
}
