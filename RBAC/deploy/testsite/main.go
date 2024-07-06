package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kenshaw/sentinel"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"brank.as/rbac/serviceutil/logging"
	client "brank.as/rbac/svcutil/hydraclient"
	"brank.as/rbac/svcutil/mw"
	"brank.as/rbac/svcutil/random"

	cpb "brank.as/rbac/gunk/v1/consent"
)

//go:embed assets
var assets embed.FS

func main() {
	log := logging.NewLogger().WithField("service", "testsite")
	config := viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer(".", "_")))
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}
	if err := setDefaults(config); err != nil {
		log.Fatal(err)
	}

	hy, err := mw.NewHydra(config, []string{"/login", "/oauth2/callback"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	hy.Scopes = append(config.GetStringSlice("auth.scopes"), "offline_access", "openid")

	s, err := NewServer(config, log, hy, assets)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.Listen("tcp", ":"+config.GetString("server.port"))
	if err != nil {
		log.Fatal(err)
	}
	ss, _ := sentinel.WithContext(context.Background(), os.Interrupt)
	if err := ss.ManageHTTP(l, s); err != nil {
		log.Fatal(err)
	}
	log.Infof("starting server on port :%s", config.GetString("server.port"))
	if err := ss.Run(log, 10*time.Second); err != nil {
		log.Fatal(err)
	}
}

func setDefaults(config *viper.Viper) error {
	cookieSecret := config.GetString("auth.cookieSecret")
	if cookieSecret == "" {
		s, err := random.String(32)
		if err != nil {
			return err
		}
		config.Set("auth.cookieSecret", s)
	}
	cookieName := config.GetString("auth.cookieName")
	if cookieName == "" {
		config.Set("auth.cookieName", "testsite")
	}
	fmt.Printf("cookie name %q\n", cookieName)
	clientID := config.GetString("auth.clientid")
	if clientID == "" {
		s, err := random.String(12)
		if err != nil {
			return err
		}
		config.Set("auth.clientid", s)
		clientID = s
	}
	clientSecret := config.GetString("auth.clientSecret")
	if clientSecret == "" {
		s, err := random.String(20)
		if err != nil {
			return err
		}
		config.Set("auth.clientSecret", s)
		clientSecret = s
	}
	cl, err := client.NewAdminClient(config.GetString("hydra.adminurl"))
	if err != nil {
		return err
	}
	ctx, canc := context.WithTimeout(context.Background(), 5*time.Second)
	defer canc()
	if cl, err := cl.GetClient(ctx, clientID); err == nil {
		fmt.Println("client:", cl)
		return nil
	}

	fmt.Println(clientID, clientSecret)
	fmt.Println("redirect url", config.GetString("auth.redirecturl"))

	u := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort("127.0.0.1", config.GetString("authclient.port")),
	}
	if c, err := cl.CreateClient(ctx, client.AuthClient{
		OwnerID:                uuid.NewString(),
		ClientID:               clientID,
		ClientName:             "Testing Auth Client",
		RedirectURIs:           []string{u.String() + "/oauth2/callback"},
		PostLogoutRedirectURIs: []string{u.String() + "/"},
		GrantTypes:             []string{"authorization_code", "offline_access", "openid"},
		ResponseTypes:          []string{"code", "refresh_token"},
		Scopes: append([]string{"offline_access", "openid"},
			config.GetStringSlice("authclient.scopes")...),
		Secret: clientSecret,
		AuthConfig: client.AuthConfig{
			LoginTmpl:       config.GetString("authclient.logintmpl"),
			OTPTmpl:         config.GetString("authclient.otptmpl"),
			ConsentTmpl:     config.GetString("authclient.consenttmpl"),
			RememberConsent: config.GetBool("authclient.rememeber"),
			SessionDuration: config.GetDuration("authclient.sessiondur"),
			Authenticator:   config.GetString("authclient.authenticator"),
		},
	}); err != nil {
		return err
	} else {
		fmt.Printf("client created: %+v\n", c)
	}
	return scopes(ctx, config)
}

func scopes(ctx context.Context, conf *viper.Viper) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, conf.GetString("scopes.host"),
		grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	cl := cpb.NewScopeServiceClient(conn)
	for _, v := range []struct {
		scope string
		name  string
		desc  string
		group string
	}{
		{
			scope: "https://product.bnk.to/service.read",
			name:  "Service Read",
			desc:  "Read Service Objects",
			group: "ServiceName",
		},
		{
			scope: "https://product.bnk.to/service.write",
			name:  "Service Write",
			desc:  "Write Service Objects",
			group: "ServiceName",
		},
	} {
		s, err := cl.UpsertScope(ctx, &cpb.UpsertScopeRequest{
			Scope:       v.scope,
			Name:        v.name,
			GroupName:   v.group,
			Description: v.desc,
		})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("%w (check that rbac usermgm is running)", err)
			}
			return err
		}
		fmt.Println("scope registered", v.scope, s.Updated.AsTime())
	}
	return nil
}
