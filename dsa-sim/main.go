package main

import (
	"context"
	"embed"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/kenshaw/sentinel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"

	pnpb "brank.as/petnet/gunk/drp/v1/partner"
	pf "brank.as/petnet/gunk/drp/v1/profile"
	ct "brank.as/petnet/gunk/drp/v1/quote"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
)

//go:embed assets
var assets embed.FS

func main() {
	c := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	c.SetConfigFile("env/config")
	c.SetConfigType("ini")
	c.AutomaticEnv()
	if err := c.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	log := logging.NewLogger(c).WithField("service", "dsa-sim")
	rand.Seed(time.Now().UnixNano())

	hy, err := mw.NewHydra(c, mw.Config{
		IgnorePaths: []string{"/login-basic", "/login-oauth", "/oauth2/callback", "/auth-option"},
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
	hy.Scopes = append(c.GetStringSlice("auth.scopes"), "offline_access", "openid")

	cs := newConns(log, c)
	defer cs.close()

	cl := newSvcClients(cs)
	s, err := NewServer(c, log, hy, assets, cl)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.Listen("tcp", ":"+c.GetString("server.port"))
	if err != nil {
		log.Fatal(err)
	}
	ss, _ := sentinel.WithContext(context.Background(), os.Interrupt)
	if err := ss.ManageHTTP(l, s); err != nil {
		log.Fatal(err)
	}
	log.Infof("starting server on port :%s", c.GetString("server.port"))
	if err := ss.Run(log, 10*time.Second); err != nil {
		log.Fatal(err)
	}
}

type conns struct {
	apiExtFwd    *grpc.ClientConn
	apiExtClient *grpc.ClientConn
}

func newConns(log *logrus.Entry, c *viper.Viper) *conns {
	log.WithField("host", c.GetString("api.external")).Println("dialing remit api external")
	apiExtFwd, err := grpc.Dial(
		c.GetString("api.external"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(mw.AuthForwarder()),
	)
	if err != nil {
		log.Fatal(err)
	}

	apiExtClient, err := grpc.Dial(
		c.GetString("api.external"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(WithClient(c)),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &conns{
		apiExtFwd:    apiExtFwd,
		apiExtClient: apiExtClient,
	}
}

type api interface {
	tpb.TerminalServiceClient
	pnpb.RemitPartnerServiceClient
	pf.ProfileServiceClient
	ct.QuoteServiceClient
}

type cl struct {
	apiFwd    api
	apiClient api
}

func newSvcClients(cs *conns) cl {
	return cl{
		apiFwd: struct {
			tpb.TerminalServiceClient
			pnpb.RemitPartnerServiceClient
			pf.ProfileServiceClient
			ct.QuoteServiceClient
		}{
			TerminalServiceClient:     tpb.NewTerminalServiceClient(cs.apiExtFwd),
			RemitPartnerServiceClient: pnpb.NewRemitPartnerServiceClient(cs.apiExtFwd),
			ProfileServiceClient:      pf.NewProfileServiceClient(cs.apiExtFwd),
			QuoteServiceClient:        ct.NewQuoteServiceClient(cs.apiExtFwd),
		},
		apiClient: struct {
			tpb.TerminalServiceClient
			pnpb.RemitPartnerServiceClient
			pf.ProfileServiceClient
			ct.QuoteServiceClient
		}{
			TerminalServiceClient:     tpb.NewTerminalServiceClient(cs.apiExtClient),
			RemitPartnerServiceClient: pnpb.NewRemitPartnerServiceClient(cs.apiExtClient),
			ProfileServiceClient:      pf.NewProfileServiceClient(cs.apiExtClient),
			QuoteServiceClient:        ct.NewQuoteServiceClient(cs.apiExtClient),
		},
	}
}

func (cs *conns) close() {
	cs.apiExtFwd.Close()
	cs.apiExtClient.Close()
}

func WithClient(c *viper.Viper) grpc.UnaryClientInterceptor {
	cred := clientcredentials.Config{
		ClientID:     c.GetString("dsa.ClientID"),
		ClientSecret: c.GetString("dsa.ClientSecret"),
		TokenURL:     c.GetString("auth.url") + "/oauth2/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	ts := cred.TokenSource(context.Background())
	return func(c context.Context, m string, rq, rp interface{}, cc *grpc.ClientConn, inv grpc.UnaryInvoker, o ...grpc.CallOption) error {
		tok, err := ts.Token()
		if err != nil {
			return err
		}
		c = metautils.ExtractOutgoing(c).
			Set("authorization", "Bearer "+tok.AccessToken).ToOutgoing(c)
		return inv(c, m, rq, rp, cc, o...)
	}
}
