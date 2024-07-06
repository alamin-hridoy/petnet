package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sirupsen/logrus"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	"github.com/google/uuid"

	client "brank.as/rbac/svcutil/hydraclient"

	pfpm "brank.as/petnet/profile/permission"
	pfstorage "brank.as/petnet/profile/storage"
	idstorage "brank.as/rbac/usermgm/storage"

	pfpg "brank.as/petnet/profile/storage/postgres"
	rbipb "brank.as/rbac/gunk/v1/invite"
	rbmpb "brank.as/rbac/gunk/v1/mfa"
	rbopb "brank.as/rbac/gunk/v1/organization"
	rbppb "brank.as/rbac/gunk/v1/permissions"
	rbsapb "brank.as/rbac/gunk/v1/serviceaccount"
	rbupb "brank.as/rbac/gunk/v1/user"
	idpg "brank.as/rbac/usermgm/storage/postgres"
)

var svcName = "local-setup"

func main() {
	c := newConfig()
	log := logging.NewLogger(c).WithFields(logrus.Fields{"service": svcName})
	log.Info("starting local-setup-util service")

	cs := newConns(c, log)
	defer cs.close()
	cl := newSvcClients(cs)
	st := newStores(c, log)
	hy := newHydra(c, log)
	s := newSvc(c, log, cl, st, hy)

	ctx := context.Background()
	if err := s.createCMSAuthClient(ctx, c.GetString("cms.clientID"), c.GetString("cms.clientSecret"), "cms"); err != nil {
		if err == pfstorage.Conflict {
			log.Info("already setup skipping")
			return
		}
		logging.WithError(err, log).Fatal("unable to create cms hydra client")
		return
	}
	log.Info("created auth client")

	if err := s.setupAdmin(ctx, c, log); err != nil {
		logging.WithError(err, log).Fatal("unable to bootstrap local petnet admin")
		return
	}
	log.Info("setup admin")
	uid, oid, err := s.setupDSA(ctx, c, log)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to bootstrap local dsa")
		return
	}
	log.Info("setup dsa")

	if err := s.createAPIKey(ctx,
		c.GetString("dsa.apiKeyStg"), c.GetString("dsa.apiSecretStg"),
		uid, oid, "sandbox", "sandbox",
	); err != nil {
		logging.WithError(err, log).Fatal("unable to create dsa api key")
		return
	}
	if err := s.createAPIKey(ctx,
		c.GetString("dsa.apiKeyLive"), c.GetString("dsa.apiSecretLive"),
		uid, oid, "live", "live",
	); err != nil {
		logging.WithError(err, log).Fatal("unable to create dsa api key")
		return
	}

	if err := s.createDSAAuthClient(ctx, c, "sandbox", uid, oid); err != nil {
		logging.WithError(err, log).Fatal("unable to create dsa auth client")
		return
	}

	if err := pfpm.BootstrapAdminPermissions(ctx, log, rbppb.NewPermissionServiceClient(cs.idExt), rbppb.NewProductServiceClient(cs.idExt), st.pf); err != nil {
		logging.WithError(err, log).Fatal("unable to bootstrap local dsa")
		return
	}
	log.Info("done")
}

func (s *svc) createDSAAuthClient(ctx context.Context, config *viper.Viper, env, uid, oid string) error {
	// cookieSecret := config.GetString("dsa.simCookieSecret")
	cookieName := config.GetString("dsa.simCookieName")
	fmt.Printf("cookie name %q\n", cookieName)
	clientID := config.GetString("dsa.simAuthID")
	clientSecret := config.GetString("dsa.simAuthSecret")
	u := config.GetString("server.url")

	fmt.Println(clientID, clientSecret)
	fmt.Println("redirect url", config.GetString("auth.redirecturl"))

	clName := "DSA Simulator"
	c, err := s.hy.CreateClient(ctx, client.AuthClient{
		OwnerID:                uuid.New().String(),
		ClientID:               clientID,
		ClientName:             clName,
		RedirectURIs:           []string{config.GetString("auth.redirecturl")},
		CORS:                   []string{u},
		PostLogoutRedirectURIs: []string{u},
		LogoURL:                u + "/logo",
		GrantTypes:             []string{"authorization_code", "offline_access", "openid"},
		ResponseTypes:          []string{"code", "refresh_token"},
		Scopes: []string{
			"offline_access", "openid",
			"https://product.bnk.to/service.read", "https://product.bnk.to/service.write",
		},
		Audience:    []string{"sandbox"},
		Secret:      clientSecret,
		SubjectType: "public",
		AuthMethod:  "client_secret_basic",
		AuthConfig: client.AuthConfig{
			Authenticator:   "perahub",
			LoginTmpl:       "",
			OTPTmpl:         "",
			ConsentTmpl:     "",
			RememberConsent: true,
			SessionDuration: 0,
		},
	})
	if err != nil {
		return err
	}
	fmt.Println("client created:", c)

	_, err = s.st.id.CreateOauthClient(ctx, idstorage.OAuthClient{
		OrgID:        oid,
		ClientID:     clientID,
		ClientName:   clName,
		CreateUserID: uid,
		UpdateUserID: uid,
		Environment:  env,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *svc) createCMSAuthClient(ctx context.Context, clID, clSecret, name string) error {
	_, err := s.hy.CreateClient(ctx, client.AuthClient{
		OwnerID:    uuid.NewString(),
		ClientID:   clID,
		ClientName: name,
		RedirectURIs: []string{
			"http://cms.localhost/oauth2/callback",
		},
		CORS: []string{},
		PostLogoutRedirectURIs: []string{
			"http://cms.localhost",
		},
		GrantTypes: []string{
			"authorization_code", "refresh_token", "client_credentials",
		},
		ResponseTypes: []string{"code", "id_token"},
		Scopes:        []string{"https://rbac.brank.as/read", "https://rbac.brank.as/write", "openid", "offline_access"},
		Secret:        clSecret,
	})
	if err != nil {
		return pfstorage.Conflict
	}
	return nil
}

func (s *svc) createAPIKey(ctx context.Context, clientID, clientSecret, uid, oid, env, name string) error {
	_, err := s.hy.CreateClient(ctx, client.AuthClient{
		ClientID:   clientID,
		Secret:     clientSecret,
		GrantTypes: []string{"client_credentials"},
		Scopes:     []string{"offline_access", "offline", "openid"},
	})
	if err != nil {
		return err
	}

	_, err = s.st.id.CreateSvcAccount(ctx, idstorage.SvcAccount{
		AuthType:     "oauth",
		OrgID:        oid,
		Environment:  env,
		ClientName:   name,
		ClientID:     clientID,
		CreateUserID: uid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *svc) setupAdmin(ctx context.Context, c *viper.Viper, log *logrus.Entry) error {
	e := c.GetString("admin.Email")
	res, err := s.cl.rbac.Signup(ctx, &rbupb.SignupRequest{
		Username:  e,
		FirstName: "placeholder",
		LastName:  "placeholder",
		Email:     e,
		Password:  c.GetString("admin.Password"),
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil
		}
		logging.WithError(err, log).Error("signing up")
		return err
	}

	o, err := s.cl.rbac.GetOrganization(ctx, &rbopb.GetOrganizationRequest{ID: res.GetOrgID()})
	if err != nil {
		logging.WithError(err, log).Error("get org")
	}

	if !o.Organization[0].Active {
		if _, err := s.cl.rbac.Approve(ctx, &rbipb.ApproveRequest{
			ID: res.UserID,
		}); err != nil {
			logging.WithError(err, log).Error("approve invitation")
		}
	}
	mfa := rbopb.EnableOpt_Enable
	if c.GetBool("common.disableLoginMFA") {
		mfa = rbopb.EnableOpt_Disable
	}
	if org, err := s.cl.rbac.UpdateOrganization(ctx, &rbopb.UpdateOrganizationRequest{
		OrganizationID: res.GetOrgID(),
		LoginMFA:       mfa,
	}); err != nil {
		logging.WithError(err, log).Error("update org")
	} else if org.MFAEventID != "" {
		log.WithField("event_id", org.MFAEventID).Error("auto-enable failed")
	}

	if _, err = s.st.pf.CreateOrgProfile(ctx, &pfstorage.OrgProfile{
		UserID:  res.UserID,
		OrgID:   res.OrgID,
		OrgType: int(ppb.OrgType(ppb.OrgType_PetNet)),
	}); err != nil {
		log.WithError(err).Error("creating org profile")
		return err
	}

	if _, err = s.st.pf.CreateUserProfile(ctx, &pfstorage.UserProfile{
		UserID: res.UserID,
		OrgID:  res.OrgID,
		Email:  e,
	}); err != nil {
		log.WithError(err).Error("creating user profile")
		return err
	}
	code, err := s.st.id.GetConfirmationCode(ctx, res.UserID)
	if err != nil {
		logging.WithError(err, log).Error("getting confirmation code")
		return err
	}

	if _, err := s.cl.rbac.EmailConfirmation(ctx, &rbupb.EmailConfirmationRequest{
		Code: code,
	}); err != nil {
		logging.WithError(err, log).Error("confirming email")
	}
	return nil
}

func (s *svc) setupDSA(ctx context.Context, c *viper.Viper, log *logrus.Entry) (string, string, error) {
	e := c.GetString("dsa.Email")
	res, err := s.cl.rbac.Signup(ctx, &rbupb.SignupRequest{
		Username:  e,
		FirstName: "placeholder",
		LastName:  "placeholder",
		Email:     e,
		Password:  c.GetString("dsa.Password"),
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return res.GetOrgID(), res.GetUserID(), nil
		}
		logging.WithError(err, log).Error("signing up")
		return "", "", err
	}
	oid := res.GetOrgID()
	uid := res.GetUserID()

	o, err := s.cl.rbac.GetOrganization(ctx, &rbopb.GetOrganizationRequest{ID: res.GetOrgID()})
	if err != nil {
		logging.WithError(err, log).Error("get org")
	}

	if !o.Organization[0].Active {
		if _, err := s.cl.rbac.Approve(ctx, &rbipb.ApproveRequest{
			ID: res.UserID,
		}); err != nil {
			logging.WithError(err, log).Error("approving invitation")
		}
	}

	mfa := rbopb.EnableOpt_Enable
	if c.GetBool("common.disableLoginMFA") {
		mfa = rbopb.EnableOpt_Disable
	}
	if org, err := s.cl.rbac.UpdateOrganization(ctx, &rbopb.UpdateOrganizationRequest{
		OrganizationID: res.GetOrgID(),
		LoginMFA:       mfa,
	}); err != nil {
		return uid, oid, err
	} else if org.MFAEventID != "" {
		log.WithField("event_id", org.MFAEventID).Error("auto-enable failed")
	}

	if _, err = s.st.pf.CreateOrgProfile(ctx, &pfstorage.OrgProfile{
		UserID:  res.UserID,
		OrgID:   res.OrgID,
		OrgType: int(ppb.OrgType_DSA),
		Status:  int(ppb.Status_Pending),
		BusinessInfo: pfstorage.BusinessInfo{
			CompanyName:   "dsa-company",
			StoreName:     "dsa-store-name",
			PhoneNumber:   "11111",
			FaxNumber:     "22222",
			Website:       "dsa.com",
			CompanyEmail:  e,
			ContactPerson: "Guybrush Threepwood",
			Position:      "CEO",
			Address: pfstorage.Address{
				Address1:   "addr1",
				City:       "city",
				State:      "state",
				PostalCode: "12345",
			},
		},
		AccountInfo: pfstorage.AccountInfo{
			Bank:                    "bank",
			BankAccountNumber:       "54321",
			BankAccountHolder:       "Guybrush Threepwood",
			AgreeTermsConditions:    int(ppb.Boolean_True),
			AgreeOnlineSupplierForm: int(ppb.Boolean_True),
			Currency:                int(ppb.Currency_PHP),
		},
		DateApplied: sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		return uid, oid, err
	}

	if _, err = s.st.pf.CreateUserProfile(ctx, &pfstorage.UserProfile{
		UserID: res.UserID,
		OrgID:  res.OrgID,
		Email:  e,
	}); err != nil {
		return uid, oid, err
	}

	code, err := s.st.id.GetConfirmationCode(ctx, res.UserID)
	if err != nil {
		return uid, oid, err
	}

	if _, err := s.cl.rbac.EmailConfirmation(ctx, &rbupb.EmailConfirmationRequest{
		Code: code,
	}); err != nil {
		return uid, oid, err
	}

	ts := []string{"WU", "IR", "TF"}
	for _, t := range ts {
		if _, err := s.st.pf.CreatePartner(ctx, &pfstorage.Partner{
			OrgID:     oid,
			Type:      t,
			Partner:   "{}",
			UpdatedBy: uid,
		}); err != nil {
			return uid, oid, err
		}
		if err := s.st.pf.EnablePartner(ctx, oid, t); err != nil {
			return uid, oid, err
		}
	}
	return uid, oid, nil
}

type svc struct {
	cl cl
	st st
	hy *client.AdminClient
}

type conns struct {
	idInt *grpc.ClientConn
	idExt *grpc.ClientConn
	pfInt *grpc.ClientConn
}

type cl struct {
	rbac identity
}

type st struct {
	pf *pfpg.Storage
	id *idpg.Storage
}

type identity interface {
	rbipb.InviteServiceClient
	rbupb.SignupClient
	rbopb.OrganizationServiceClient
	rbmpb.MFAServiceClient
	rbppb.PermissionServiceClient
	rbppb.ProductServiceClient
	rbsapb.SvcAccountServiceClient
}
