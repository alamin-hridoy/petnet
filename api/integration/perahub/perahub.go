package perahub

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/pkcs12"

	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/storage/postgres"
)

type Svc struct {
	baseUrl              *url.URL
	nonexUrl             *url.URL
	billerUrl            *url.URL
	billsUrl             *url.URL
	cicoUrl              *url.URL
	phRemittanceUrl      *url.URL
	phTransactUrl        *url.URL
	nonexAPIKey          string
	perahubDefaultAPIKey string
	phEnv                string
	cl                   HTTPClient
	token                string
	signKey              *rsa.PrivateKey

	ptnrAuthCreds map[string]OAuthCreds

	curr currencyCache
	exp  time.Duration

	log *logrus.Entry
}

type OAuthCreds struct {
	ClientID     string
	ClientSecret string
}

// SvcOption is options type for Svc
type SvcOption func(s *Svc)

// WithPerahubDefaultAPIKey ...
func WithPerahubDefaultAPIKey(apiKey string) SvcOption {
	return func(s *Svc) {
		s.perahubDefaultAPIKey = apiKey
	}
}

// WithCiCoURL ...
func WithCiCoURL(cicoUrl string) SvcOption {
	return func(s *Svc) {
		u, err := url.Parse(cicoUrl)
		if err != nil {
			s.log.Error(err)
			return
		}

		s.cicoUrl = u
	}
}

// WithPHRemittanceURL ...
func WithPHRemittanceURL(remittanceUrl string) SvcOption {
	return func(s *Svc) {
		u, err := url.Parse(remittanceUrl)
		if err != nil {
			s.log.Error(err)
			return
		}

		s.phRemittanceUrl = u
	}
}

// WithLogger ...
func WithLogger(log *logrus.Entry) SvcOption {
	return func(s *Svc) {
		s.log = log
	}
}

func New(cl HTTPClient, phEnv, baseUrl, nonexUrl, billerUrl, billsUrl, phTransactUrl, partnerID, clientKey, nonexAPIKey, serverIP, cert string, creds map[string]OAuthCreds, opts ...SvcOption) (*Svc, error) {
	if cl == nil {
		cl = &http.Client{}
	}

	bu, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	nu, err := url.Parse(nonexUrl)
	if err != nil {
		return nil, err
	}

	biu, err := url.Parse(billerUrl)
	if err != nil {
		return nil, err
	}

	bpu, err := url.Parse(billsUrl)
	if err != nil {
		return nil, err
	}

	ptu, err := url.Parse(phTransactUrl)
	if err != nil {
		return nil, err
	}

	var sgnKey *rsa.PrivateKey
	if cert != "" {
		f, err := os.ReadFile(cert)
		if err != nil {
			return nil, err
		}
		pkey, _, err := pkcs12.Decode(f, "")
		if err != nil {
			return nil, err
		}
		pk, ok := pkey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("invalid rsa private key")
		}
		sgnKey = pk
	}

	hasher := md5.Sum([]byte(partnerID + serverIP + clientKey))
	s := &Svc{
		baseUrl:       bu,
		nonexUrl:      nu,
		billerUrl:     biu,
		phTransactUrl: ptu,
		// cicoUrl and phRemittanceUrl should be overwritten by options
		cicoUrl:         nu,
		phRemittanceUrl: nu,
		billsUrl:        bpu,
		nonexAPIKey:     nonexAPIKey,
		// perahubDefaultAPIKey should be overwritten by options
		perahubDefaultAPIKey: nonexAPIKey,
		phEnv:                phEnv,
		cl:                   cl,
		token:                hex.EncodeToString(hasher[:]),
		signKey:              sgnKey,
		ptnrAuthCreds:        creds,

		curr: currencyCache{
			Mutex: &sync.Mutex{},
			lst:   map[string]country{},
		},
		exp: time.Hour,
		log: logrus.New().WithField("stage", "init"),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.log.WithField("baser url", s.baseUrl).
		WithField("nonex url", s.nonexUrl).
		WithField("cico url", s.cicoUrl).
		WithField("biller url", s.billerUrl).
		WithField("bills pay url", s.billsUrl).
		WithField("transact url", s.phTransactUrl).
		WithField("remittance url", s.phRemittanceUrl).
		Debug("init perahub service")

	return s, nil
}

func newTestSvc(t *testing.T, store *postgres.Storage) (*Svc, *HTTPMock) {
	cl := NewTestHTTPMock(store, MockConfig{
		SaveReq: true,
	})
	log := logrus.New().WithField("stage", "testing")
	ph, err := New(cl,
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
		map[string]OAuthCreds{
			static.WISECode: {
				ClientID:     "wise-id",
				ClientSecret: "wise-secret",
			},
		},
		WithLogger(log),
		WithPerahubDefaultAPIKey("api-key"),
		WithCiCoURL("https://privatedrp.dev.perahub.com.ph/v1/cico/wrapper/"),
		WithPHRemittanceURL("https://privatedrp.dev.perahub.com.ph/v1/remit/dmt/"),
	)
	if err != nil {
		t.Fatal("setting up perahub integration: ", err)
	}
	return ph, cl
}

func (s *Svc) SetMock(cl HTTPClient) *Svc {
	sc := *s
	sc.cl = cl
	return &sc
}

func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func FormatName(name string) (string, string, string) {
	sname := strings.FieldsFunc(name, func(r rune) bool {
		return r == ',' || r == ' '
	})
	switch len(sname) {
	case 1:
		return sname[0], "", ""
	case 2:
		return sname[0], sname[1], ""
	case 3:
		return sname[0], sname[1], sname[2]
	}
	return "", "", ""
}

func GenderChar(g string) string {
	gender := "M"
	if g == "Female" {
		gender = "F"
	}
	return gender
}

func CombinedName(f, m, l string) string {
	n := f + ", " + l
	if m != "" {
		n = f + ", " + m + " " + l
	}
	return n
}

func CurrencyNumber(c string) string {
	var ccur string
	switch c {
	case "PHP":
		ccur = "1"
	case "USD":
		ccur = "2"
	default:
		fmt.Println("currency not equal two one of the allowed currencies, PHP and USD, got: ", ccur)
	}
	return ccur
}

func FormatAddress(addr1, addr2 string) string {
	switch {
	case addr1 != "" && addr2 != "":
		return addr1 + ", " + addr2
	case addr1 != "":
		return addr1
	}
	return ""
}

func IsDomestic(srcCtry, destCtry string) int {
	isd := int(0)
	if srcCtry == destCtry {
		isd = 1
	}
	return isd
}

type NonexAddress struct {
	Address1 string `json:"address_1"`
	Address2 string `json:"address_2"`
	Barangay string `json:"barangay"`
	City     string `json:"city"`
	Province string `json:"province"`
	ZipCode  string `json:"zip_code"`
	Country  string `json:"country"`
}
