package endtoend

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/svcutil/mw"
)

func envDefault() string {
	switch os.Getenv("DRONE_BRANCH") {
	case "production", "endtoend-production-test":
		return "production"
	}
	return "staging"
}

type option func(*md)

type md struct {
	clid       string
	token      string
	disableTLS bool
}

func withToken(tok string) option {
	return func(a *md) {
		a.token = tok
	}
}

func withClientID(clID string) option {
	return func(a *md) {
		a.clid = clID
	}
}

type util struct {
	t   *testing.T
	cnf *viper.Viper
}

func newConfig(t *testing.T) *viper.Viper {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile(filepath.Join("env", "sample.config"))
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		t.Fatalf("error loading configuration: %v", err)
	}
	return config
}

func (ut util) newUsermgmConn(u string, opts ...option) (*grpc.ClientConn, context.Context) {
	ao := &md{}
	for _, o := range opts {
		o(ao)
	}

	creds := credentials.NewClientTLSFromCert(nil, "")
	gopts := []grpc.DialOption{grpc.WithTransportCredentials(creds), grpc.WithBlock()}
	if ut.cnf.GetBool("usermgm.disableTLS") {
		gopts = []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	}
	conn, err := grpc.Dial(u, gopts...)
	if err != nil {
		ut.t.Fatal("dialing usermgm grpc: ", err)
	}
	ctx := context.Background()
	if ao.token != "" {
		md := metadata.Pairs("authorization", "Bearer "+ao.token)
		if ao.clid != "" {
			md.Append(hydra.ClientIDKey, ao.clid)
		}
		ctx = metadata.NewOutgoingContext(context.Background(), md)
	}
	return conn, ctx
}

const (
	sessionCookieName  = "proxtera-session"
	sessionCookieState = "state"
	sessionCookieToken = "token"
	authCodeURL        = "somerandomstring"
	cdpTimeout         = 4 * time.Minute
)

func (ut util) newCDP() (context.Context, context.CancelFunc) {
	ctx := testBrowser(ut.t)

	// create a timeout
	ctx, cancel := context.WithTimeout(ctx, cdpTimeout)

	return ctx, cancel
}

func (ut util) getToken(clid, clsec string) string {
	cred := clientcredentials.Config{
		ClientID:     clid,
		ClientSecret: clsec,
		TokenURL:     ut.cnf.GetString("auth.url") + "/oauth2/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	ts := oauth.TokenSource{TokenSource: cred.TokenSource(context.Background())}
	tok, err := ts.Token()
	if err != nil {
		ut.t.Fatal("token: ", err)
	}
	return tok.AccessToken
}

func (ut util) exchangeCode(ctx context.Context, code string) string {
	hmw, err := mw.NewHydra(ut.cnf, []string{"/oauth2/callback"}, nil)
	if err != nil {
		ut.t.Fatal("creating hydra: ", err)
	}
	token, err := hmw.Config.Exchange(ctx, code)
	if err != nil {
		ut.t.Fatal("exchanging code: ", err)
	}
	return token.AccessToken
}

func (ut util) getInvURL(t *testing.T, wHToken, authUrl string) (string, error) {
	retries := 15
retryGet:
	whURL := "https://webhook.site/token/" + wHToken + "/request/latest"

	resp, err := http.Get(whURL)
	if err != nil {
		t.Fatal("calling webhook.site: ", err)
	}
	defer resp.Body.Close()

	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("reading body: ", err)
	}
	b := string(bb)

	u, err := url.Parse(authUrl)
	if err != nil {
		t.Fatal("parsing auth url: ", err)
	}

	ss := "invite_code="
	i := strings.Index(b, ss)
	if i == -1 {
		if retries == 0 {
			return "", errors.New("receiving email, resend invite")
		}
		time.Sleep(1 * time.Second)
		retries--
		t.Log("Retrying getting invite code")
		goto retryGet
	}
	code := b[i+len(ss)+2 : i+len(ss)+18]
	invURL := "http://" + u.Host + "/signup?invite_code=" + code

	retries = 5
retryDel:
	whDeleteReqURL := "https://webhook.site/token/" + wHToken + "/request"
	req, err := http.NewRequest(http.MethodDelete, whDeleteReqURL, nil)
	if err != nil {
		t.Fatal("creating delete webhook request: ", err)
	}

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("deleting webook requests: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		if retries == 0 {
			t.Fatal("deleting webook requests")
		}
		time.Sleep(1 * time.Second)
		retries--
		t.Log("Retrying deleting webhook requests")
		goto retryDel
	}

	return invURL, nil
}

func (ut util) createWHToken(t *testing.T) string {
	d := []byte(`{"expiry":true}`)
	resp, err := http.Post("https://webhook.site/token", `Content-Type", "application/json`, bytes.NewBuffer(d))
	if err != nil {
		t.Fatal("connecting to webhook.site: ", err)
	}
	defer resp.Body.Close()
	token := &struct {
		UUID string `json:"uuid"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		t.Fatal("decoding body: ", err)
	}

	return token.UUID
}
