package mw

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/status"

	spb "brank.as/rbac/gunk/v1/user"
	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	pm "brank.as/rbac/usermgm/core/permissions"
	"brank.as/rbac/usermgm/integrations/keto"
	"brank.as/rbac/usermgm/storage/postgres"
)

func testConfig(t *testing.T) *viper.Viper {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile(filepath.Join("..", "..", "usermgm", "env", "config"))
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		t.Fatalf("error loading configuration: %v", err)
	}
	return config
}

func TestPermValidation(t *testing.T) {
	baseURL := os.Getenv("KETO_URL")
	dbconn := os.Getenv("DATABASE_CONNECTION")
	switch "" {
	case baseURL:
		t.Skip("missing env 'KETO_URL'")
	case dbconn:
		t.Error("missing env 'DATABASE_CONNECTION'")
	}
	cnf := testConfig(t)
	db, err := postgres.NewStorage(dbconn)
	if err != nil {
		t.Fatal("creating new storage: ", err)
	}

	logr, _ := test.NewNullLogger()
	log := logr.WithField("test", "creating permission")
	ctx := logging.WithLogger(context.Background(), log)

	os, err := db.GetOrgs(ctx)
	if err != nil {
		t.Fatal("getting orgs: ", err)
	}
	var bsOrgID string
	for _, o := range os {
		if o.ContactEmail == "admin@example.com" {
			bsOrgID = o.ID
		}
	}
	sas, err := db.GetSvcAccountByOrgID(ctx, bsOrgID)
	if err != nil {
		t.Fatal("getting service account: ", err)
	}
	clID := sas[0].ClientID

	k := keto.New(baseURL)
	s := pm.New(db, k)
	pms, err := s.ListPermission(ctx, core.ListPermissionFilter{OrgID: bsOrgID})
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.CreateRole(ctx, core.Role{
		OrgID:     bsOrgID,
		Name:      "Test Signup",
		Desc:      "Signup testing",
		CreateUID: clID,
		Members:   []string{clID},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range pms {
		_, err := s.RoleGrant(ctx, core.Grant{
			RoleID:      r,
			GrantID:     p.ID,
			Environment: p.Environment,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	ctx, cc := setupGRPC(t, gRPCConfig{
		clientID:     clID,
		clientSecret: cnf.GetString("bootstrap.hydraSecret"),
		usermgmURL:   net.JoinHostPort(cnf.GetString("server.host"), cnf.GetString("server.port")),
		authURL:      cnf.GetString("hydra.publicURL") + "/oauth2/token",
	})
	defer cc.Close()
	cl := spb.NewSignupClient(cc)
	user, err := cl.Signup(ctx, &spb.SignupRequest{
		Username:   "user",
		FirstName:  "first",
		LastName:   "last",
		Email:      "email@mail.com",
		Password:   uuid.New().String(),
		InviteCode: "code",
		OrgID:      uuid.New().String(),
	})
	if err != nil {
		t.Fatal("signing up: ", err)
	}

	_, err = cl.ForgotPassword(ctx, &spb.ForgotPasswordRequest{})
	if err == nil {
		t.Fatal("status should be permission denied")
	}
	if e, ok := status.FromError(err); ok {
		if e.Code() != codes.PermissionDenied {
			t.Fatal("status should be permission denied")
		}
	}

	usvc := spb.NewUserServiceClient(cc)
	_, err = usvc.GetUser(ctx, &spb.GetUserRequest{ID: user.GetUserID()})
	if err != nil {
		t.Fatal("user should be granted to access it's own record")
	}
	if e, ok := status.FromError(err); ok {
		if e.Code() == codes.PermissionDenied {
			t.Fatal("user should be granted to access it's own record")
		}
	}
}

type gRPCConfig struct {
	usermgmURL   string
	clientID     string
	clientSecret string
	authURL      string
}

func setupGRPC(t *testing.T, cnf gRPCConfig) (context.Context, *grpc.ClientConn) {
	cred := clientcredentials.Config{
		ClientID:     cnf.clientID,
		ClientSecret: cnf.clientSecret,
		TokenURL:     cnf.authURL,
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	ts := oauth.TokenSource{TokenSource: cred.TokenSource(context.Background())}
	tok, err := ts.Token()
	if err != nil {
		t.Fatal("token: ", err)
	}
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	cc, err := grpc.Dial(cnf.usermgmURL, opts...)
	if err != nil {
		t.Fatal("dialing usermgm grpc: ", err)
	}

	md := metautils.NiceMD{
		"authorization":   {"Bearer " + tok.AccessToken},
		hydra.ClientIDKey: {cnf.clientID},
	}
	return md.ToOutgoing(context.Background()), cc
}
