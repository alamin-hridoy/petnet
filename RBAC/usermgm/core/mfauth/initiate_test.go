package mfauth

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"
)

func TestResend(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skipf("env not set %q", "DATABASE_CONNECTION")
	}
	t.Parallel()
	st, clean := postgres.NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)
	mail := &MockEmail{}
	conf := viper.New()
	conf.Set("project.mfaissuer", "retry test")
	conf.Set("project.mfatimeout", "1m")
	logr, _ := test.NewNullLogger()
	ctx := logging.WithLogger(context.Background(), logr.WithField("test", "retry"))

	o, err := st.CreateOrg(ctx, storage.Organization{
		OrgName:      "Test org",
		ContactEmail: "testing@email.com",
		ContactPhone: random.InvitationCode(10),
		Active:       true,
		MFALogin:     sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.CreateUser(ctx, storage.User{
		OrgID:         o,
		Username:      "testuser",
		FirstName:     "test",
		LastName:      "user",
		Email:         "user@example.com",
		EmailVerified: true,
		PreferredMFA:  storage.EMail,
		MFALogin:      true,
	}, storage.Credential{
		Username: random.InvitationCode(10),
		Password: random.InvitationCode(30),
	})
	if err != nil {
		t.Fatal(err)
	}

	svc, err := New(conf, st, mail)
	if err != nil {
		t.Fatal(err)
	}

	addr := "mfatest@example.com"
	m, err := svc.RegisterMFA(ctx, core.MFA{
		UserID: u.ID,
		Type:   storage.EMail,
		Source: addr,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(addr, mail.addr) {
		t.Error(cmp.Diff(addr, mail.addr))
	}

	// email MFA confirmed/activated
	c, err := svc.MFAuth(ctx, core.MFAChallenge{
		EventID:  m.ConfirmID,
		UserID:   m.UserID,
		SourceID: m.MFAID,
		Type:     m.Type,
		Token:    mail.code,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(m.MFAID, c.SourceID) {
		t.Error(cmp.Diff(m.MFAID, c.SourceID))
	}
	mail.EmailMFA("", "") // reset

	mi, err := svc.InitiateMFA(ctx, core.MFAChallenge{
		EventDesc: "Test Resend",
		UserID:    m.UserID,
		SourceID:  m.MFAID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(addr, mail.addr) {
		t.Error(cmp.Diff(addr, mail.addr))
	}
	if !cmp.Equal(mi.SourceID, c.SourceID) {
		t.Error(cmp.Diff(mi.SourceID, c.SourceID))
	}
	mail.EmailMFA("", "") // reset

	mr, err := svc.RestartMFA(ctx, core.MFAChallenge{
		EventID: mi.EventID,
		UserID:  mi.UserID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(addr, mail.addr) {
		t.Error(cmp.Diff(addr, mail.addr))
	}
	if !cmp.Equal(mr.SourceID, c.SourceID) {
		t.Error(cmp.Diff(mr.SourceID, c.SourceID))
	}
	mr.Token = mail.code

	scc, err := svc.MFAuth(ctx, *mr)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(addr, scc.Sources[0].Source) {
		t.Error(cmp.Diff(addr, scc.Sources[0].Source))
	}
}
