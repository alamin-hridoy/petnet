package auth

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/core/mfauth"
	"brank.as/rbac/usermgm/integrations/email"
	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"
)

func TestLogin(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing env 'DATABASE_CONNECTION'")
	}
	config := viper.New()
	config.Set("project.mfaissuer", "testlogin")
	config.Set("project.mfatimeout", "10m")
	config.Set("user.lockoutCount", "2")
	st, cl := postgres.NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cl)

	ctx := context.Background()
	org, err := st.CreateOrg(ctx, storage.Organization{
		OrgName: "Test Login Lock",
		Active:  true,
	})
	if err != nil {
		t.Fatal(err)
	}

	usrname := "testuser"
	pass := "apasswordthatworks"
	usr, err := st.CreateUser(ctx, storage.User{
		OrgID:         org,
		FirstName:     "That",
		LastName:      "Test",
		Username:      usrname,
		Email:         "test@email.org",
		EmailVerified: true,
	}, storage.Credential{
		Password: pass,
	})
	if err != nil {
		t.Fatal(err)
	}

	ma, err := mfauth.New(config, st, &email.MockSender{})
	if err != nil {
		t.Fatal(err)
	}
	au := New(config, st, ma, st)
	t.Run("login", func(t *testing.T) {
		logr, _ := test.NewNullLogger()
		ctx := logging.WithLogger(ctx, logr)
		ur, err := st.GetUserByID(ctx, usr.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(usr.ID, ur.ID) {
			t.Error(cmp.Diff(usr.ID, ur.ID))
		}

		u, err := au.AuthUser(ctx, core.AuthCredential{
			Username: usrname,
			Password: pass,
		})
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(usr.ID, u.ID) {
			t.Error(cmp.Diff(usr.ID, u.ID))
		}

		badpass := "notthepassword"
		u, err = au.AuthUser(ctx, core.AuthCredential{
			Username: usrname,
			Password: badpass,
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if !cmp.Equal(u, &core.Identity{Retries: 1, TrackRetries: true}) {
			t.Error(cmp.Diff(u, &core.Identity{Retries: 1, TrackRetries: true}))
		}

		u, err = au.AuthUser(ctx, core.AuthCredential{
			Username: usrname,
			Password: badpass,
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if !cmp.Equal(u, &core.Identity{Locked: true}) {
			t.Error(cmp.Diff(u, &core.Identity{Locked: true}))
		}
	})

	t.Run("reset", func(t *testing.T) {
		logr, _ := test.NewNullLogger()
		ctx := logging.WithLogger(ctx, logr)
		if err := st.UnlockUser(ctx, usr.ID); err != nil {
			t.Fatal(err)
		}

		u, err := au.AuthUser(ctx, core.AuthCredential{
			Username: usrname,
			Password: pass,
		})
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(usr.ID, u.ID) {
			t.Error(cmp.Diff(usr.ID, u.ID))
		}

		badpass := "notthepassword"
		u, err = au.AuthUser(ctx, core.AuthCredential{
			Username: usrname,
			Password: badpass,
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if !cmp.Equal(u, &core.Identity{Retries: 1, TrackRetries: true}) {
			t.Error(cmp.Diff(u, &core.Identity{Retries: 1, TrackRetries: true}))
		}

		u, err = au.AuthUser(ctx, core.AuthCredential{
			Username: usrname,
			Password: pass,
		})
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(usr.ID, u.ID) {
			t.Error(cmp.Diff(usr.ID, u.ID))
		}

		u, err = au.AuthUser(ctx, core.AuthCredential{
			Username: usrname,
			Password: badpass,
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if !cmp.Equal(u, &core.Identity{Retries: 1, TrackRetries: true}) {
			t.Error(cmp.Diff(u, &core.Identity{Retries: 1, TrackRetries: true}))
		}

		u, err = au.AuthUser(ctx, core.AuthCredential{
			Username: usrname,
			Password: badpass,
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if !cmp.Equal(u, &core.Identity{Locked: true}) {
			t.Error(cmp.Diff(u, &core.Identity{Locked: true}))
		}
	})
}

func TestLoginMFA(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing env 'DATABASE_CONNECTION'")
	}
	config := viper.New()
	config.Set("project.mfaissuer", "testlogin")
	config.Set("project.mfatimeout", "10m")
	config.Set("user.lockoutCount", "2")
	st, cl := postgres.NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cl)

	ctx := context.Background()
	org, err := st.CreateOrg(ctx, storage.Organization{
		OrgName: "Test Login Lock",
		Active:  true,
	})
	if err != nil {
		t.Fatal(err)
	}

	usrname := "testuser"
	pass := "apasswordthatworks"
	usr, err := st.CreateUser(ctx, storage.User{
		OrgID:         org,
		FirstName:     "That",
		LastName:      "Test",
		Username:      usrname,
		Email:         "test@email.org",
		EmailVerified: true,
	}, storage.Credential{
		Password: pass,
	})
	if err != nil {
		t.Fatal(err)
	}

	ma, err := mfauth.New(config, st, &email.MockSender{})
	if err != nil {
		t.Fatal(err)
	}
	au := New(config, st, ma, st)

	m, err := ma.RegisterMFA(ctx, core.MFA{
		UserID: usr.ID,
		Type:   storage.SMS,
		Source: "+0987654321",
	})
	if err != nil {
		t.Fatal(err)
	}
	chg, err := ma.MFAuth(ctx, core.MFAChallenge{
		EventID: m.ConfirmID,
		UserID:  m.UserID,
		Type:    m.Type,
		Token:   m.Source,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(m.UserID, chg.UserID) {
		t.Error(cmp.Diff(m.UserID, chg.UserID))
	}

	{
		usr, err := st.GetUserByID(ctx, m.UserID)
		if err != nil {
			t.Fatal(err)
		}
		usr.MFALogin = true
		usr, err = st.UpdateUserByID(ctx, *usr)
		if err != nil {
			t.Fatal(err)
		}
	}

	u, err := au.AuthUser(ctx, core.AuthCredential{
		Username: usrname,
		Password: pass,
	})
	if err != nil {
		t.Fatal(err)
	}
	if u.EventID == "" {
		t.Fatal("missing MFA event id")
	}
	u, err = au.AuthUser(ctx, core.AuthCredential{
		MFA: &core.MFAChallenge{
			EventID: u.EventID,
			UserID:  u.ID,
			Token:   u.Token,
			Type:    u.MFA,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(usr.ID, u.ID) {
		t.Error(cmp.Diff(usr.ID, u.ID))
	}
}
