package bills_payment

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"

	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
	"github.com/sirupsen/logrus"
)

type iPerahubReqService interface {
	BillsPost(ctx context.Context, url string, body interface{}) (json.RawMessage, error)
	BillsGet(ctx context.Context, url string) (json.RawMessage, error)
}

// perahub doc:https://bills-payment-microservice-dev-apidocs-bpsk2lcqiq-as.a.run.app/
type Client struct {
	baseUrl   *url.URL
	phService iPerahubReqService
}

func NewBillsClient(phService iPerahubReqService, baseUrl *url.URL) *Client {
	return &Client{
		baseUrl:   baseUrl,
		phService: phService,
	}
}

func (c *Client) getUrl(subUrl string) string {
	u := *c.baseUrl
	u.Path = path.Join(c.baseUrl.Path, subUrl)
	return u.String()
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

func newTestSvc(t *testing.T, st *postgres.Storage) (*Client, *perahub.HTTPMock) {
	cl := perahub.NewTestHTTPMock(st, perahub.MockConfig{
		SaveReq: true,
	})
	bilsUrl := "https://privatedrp.dev.perahub.com.ph/v1/billspay/"
	baseUrl, _ := url.Parse(bilsUrl)
	log := logrus.New().WithField("stage", "testing")
	ph, err := perahub.New(cl,
		"dev",
		"https://newkycgateway.dev.perahub.com.ph/gateway/",
		"https://privatedrp.dev.perahub.com.ph/v1/remit/nonex/",
		bilsUrl,
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
		perahub.WithLogger(log),
		perahub.WithPerahubDefaultAPIKey("api-key"),
		perahub.WithCiCoURL("https://privatedrp.dev.perahub.com.ph/v1/cico/wrapper/"),
		perahub.WithPHRemittanceURL("https://privatedrp.dev.perahub.com.ph/v1/remit/dmt/"),
	)
	if err != nil {
		t.Fatal("setting up perahub integration: ", err)
	}
	mck := perahub.NewHTTPMock(st)
	ph = ph.SetMock(mck)
	clnt := NewBillsClient(ph, baseUrl)
	return clnt, mck
}
