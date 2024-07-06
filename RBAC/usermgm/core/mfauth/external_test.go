package mfauth

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"
)

type MockEmail struct{ addr, code string }

func (m *MockEmail) EmailMFA(email, code string) error { m.addr, m.code = email, code; return nil }

func TestExternal(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skipf("missing $DATABASE_CONNECTION")
	}
	t.Parallel()
	conf := viper.New()
	conf.Set("project.mfaissuer", "retry test")
	conf.Set("project.mfatimeout", "1m")
	mail := &MockEmail{}
	logr, _ := test.NewNullLogger()
	ctx := logging.WithLogger(context.Background(), logr.WithField("test", "retry"))
	st, clean := postgres.NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)

	svc, err := New(conf, st, mail)
	if err != nil {
		t.Fatal(err)
	}

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
		PreferredMFA:  storage.SMS,
		MFALogin:      true,
	}, storage.Credential{
		Username: random.InvitationCode(10),
		Password: random.InvitationCode(30),
	})
	if err != nil {
		t.Fatal(err)
	}

	ph := random.NumString(13)
	m, err := svc.RegisterMFA(ctx, core.MFA{
		UserID: u.ID,
		Type:   storage.SMS,
		Source: ph,
	})
	if err != nil {
		t.Fatal(err)
	}

	// SMS MFA confirmed/activated
	c, err := svc.MFAuth(ctx, core.MFAChallenge{
		EventID:  m.ConfirmID,
		UserID:   m.UserID,
		SourceID: m.MFAID,
		Type:     m.Type,
		Token:    m.Source,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(m.MFAID, c.SourceID) {
		t.Error(cmp.Diff(m.MFAID, c.SourceID))
	}

	mi, err := svc.InitiateMFA(ctx, core.MFAChallenge{
		EventDesc: "Test External",
		UserID:    m.UserID,
		SourceID:  m.MFAID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(mi.SourceID, c.SourceID) {
		t.Error(cmp.Diff(mi.SourceID, c.SourceID))
	}

	newtok := random.NumString(8)
	ex, err := svc.ExternalMFA(ctx, core.MFAChallenge{
		EventID:    mi.EventID,
		ExternalID: uuid.NewString(),
		Token:      newtok,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(mi.EventID, ex.EventID) {
		t.Error(cmp.Diff(mi.EventID, ex.EventID))
	}

	// SMS MFA confirmed/activated
	cex, err := svc.MFAuth(ctx, core.MFAChallenge{
		EventID:  mi.EventID,
		UserID:   m.UserID,
		SourceID: m.MFAID,
		Type:     m.Type,
		Token:    newtok,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(m.MFAID, cex.SourceID) {
		t.Error(cmp.Diff(m.MFAID, cex.SourceID))
	}
}
