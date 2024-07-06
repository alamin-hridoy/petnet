package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"

	"brank.as/petnet/serviceutil/logging"
	client "brank.as/rbac/svcutil/hydraclient"

	pfpg "brank.as/petnet/profile/storage/postgres"
	rbipb "brank.as/rbac/gunk/v1/invite"
	rbmpb "brank.as/rbac/gunk/v1/mfa"
	rbopb "brank.as/rbac/gunk/v1/organization"
	rbppb "brank.as/rbac/gunk/v1/permissions"
	rbsapb "brank.as/rbac/gunk/v1/serviceaccount"
	rbupb "brank.as/rbac/gunk/v1/user"
	idpg "brank.as/rbac/usermgm/storage/postgres"
)

func newProfileDBFromConfig(config *viper.Viper, log *logrus.Entry) *pfpg.Storage {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("profileUser")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("profileDBname")
	if password := cf("profilePassword"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslMode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransacionTimeout")
	db, err := pfpg.NewStorage(dbParams)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to configure storage")
	}
	return db
}

func newIdentityDBFromConfig(config *viper.Viper, log *logrus.Entry) *idpg.Storage {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("identityUser")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("identityDBname")
	if password := cf("identityPassword"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslMode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransacionTimeout")
	db, err := idpg.NewStorage(dbParams)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to configure storage")
	}
	return db
}

func newSvc(c *viper.Viper, log *logrus.Entry, cl cl, st st, hy *client.AdminClient) *svc {
	s := &svc{cl: cl, st: st, hy: hy}
	return s
}

func systemAuth(c *viper.Viper) grpc.UnaryClientInterceptor {
	cred := clientcredentials.Config{
		ClientID:     c.GetString("auth.clientID"),
		ClientSecret: c.GetString("auth.clientSecret"),
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

func newConns(c *viper.Viper, log *logrus.Entry) *conns {
	log.WithField("host", c.GetString("identity.internal")).Println("dialing identity internal with system auth")
	idInt, err := grpc.Dial(
		c.GetString("identity.internal"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(systemAuth(c)),
	)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to connect to identity internal")
	}

	log.WithField("host", c.GetString("identity.external")).Println("dialing identity external with system auth")
	idExt, err := grpc.Dial(
		c.GetString("identity.external"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(systemAuth(c)),
	)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to connect to identity external with system auth")
	}

	log.WithField("host", c.GetString("profile.internal")).Println("dialing profile internal with system auth")
	pfInt, err := grpc.Dial(
		c.GetString("profile.internal"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(systemAuth(c)),
	)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to connect to profile internal")
	}
	return &conns{
		idInt: idInt,
		idExt: idExt,
		pfInt: pfInt,
	}
}

func (cs *conns) close() {
	cs.idInt.Close()
	cs.idExt.Close()
	cs.pfInt.Close()
}

func newConfig() *viper.Viper {
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
	return c
}

func newSvcClients(cs *conns) cl {
	return cl{
		rbac: struct {
			rbupb.SignupClient
			rbopb.OrganizationServiceClient
			rbmpb.MFAServiceClient
			rbipb.InviteServiceClient
			rbppb.PermissionServiceClient
			rbppb.ProductServiceClient
			rbsapb.SvcAccountServiceClient
		}{
			SignupClient:              rbupb.NewSignupClient(cs.idExt),
			OrganizationServiceClient: rbopb.NewOrganizationServiceClient(cs.idExt),
			MFAServiceClient:          rbmpb.NewMFAServiceClient(cs.idExt),
			InviteServiceClient:       rbipb.NewInviteServiceClient(cs.idExt),
			PermissionServiceClient:   rbppb.NewPermissionServiceClient(cs.idExt),
			ProductServiceClient:      rbppb.NewProductServiceClient(cs.idExt),
			SvcAccountServiceClient:   rbsapb.NewSvcAccountServiceClient(cs.idInt),
		},
	}
}

func newStores(c *viper.Viper, log *logrus.Entry) st {
	return st{
		pf: newProfileDBFromConfig(c, log),
		id: newIdentityDBFromConfig(c, log),
	}
}

func newHydra(c *viper.Viper, log *logrus.Entry) *client.AdminClient {
	cl, err := client.NewAdminClient(c.GetString("hydra.adminurl"))
	if err != nil {
		log.Fatalf("setting up hydra: %v", err)
	}
	return cl
}

func fakeTLS(r http.RoundTripper) http.RoundTripper {
	return rt(func(req *http.Request) (*http.Response, error) {
		rq := req.Clone(req.Context())
		rq.Header.Set("X-Forwarded-Proto", "https")
		return r.RoundTrip(rq)
	})
}

type rt func(*http.Request) (*http.Response, error)

func (t rt) RoundTrip(req *http.Request) (*http.Response, error) { return t(req) }
