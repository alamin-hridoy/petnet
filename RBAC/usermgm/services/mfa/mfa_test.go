package mfa

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"

	mpb "brank.as/rbac/gunk/v1/mfa"
	core "brank.as/rbac/usermgm/core/mfauth"
	"brank.as/rbac/usermgm/integrations/email"
	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"
)

func TestEmailMFA(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cnf := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	cnf.SetConfigFile(filepath.Join("..", "..", "env", "sample.config"))
	cnf.SetConfigType("ini")
	cnf.AutomaticEnv()
	if err := cnf.ReadInConfig(); err != nil {
		t.Fatalf("error loading configuration: %v", err)
	}

	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing env 'DATABASE_CONNECTION'")
	}

	test.NewNullLogger()
	st, cu := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cu)

	em := "test@mail.com"
	u, err := st.CreateUser(ctx,
		storage.User{
			OrgID:     uuid.New().String(),
			Email:     em,
			FirstName: "first",
			LastName:  "last",
		},
		storage.Credential{
			OrgID:    uuid.New().String(),
			Username: em,
			Password: "password",
		})
	if err != nil {
		t.Fatal(err)
	}

	m := &email.MockSender{}
	c, err := core.New(cnf, st, m)
	if err != nil {
		t.Fatal(err)
	}

	s := New(c)
	res, err := s.EnableMFA(ctx, &mpb.EnableMFARequest{
		UserID: u.ID,
		Type:   mpb.MFA_EMAIL,
		Source: em,
	})
	if err != nil {
		t.Fatal("create service: ", err)
	}
	if m.Email != em {
		t.Errorf("email mismatch want: %s, got: %s", em, m.Email)
	}
	if m.Code == "" {
		t.Error("email code is empty")
	}

	_, err = s.ValidateMFA(ctx, &mpb.ValidateMFARequest{
		UserID:  u.ID,
		Type:    mpb.MFA_EMAIL,
		Token:   m.Code,
		EventID: res.GetEventID(),
	})
	if err != nil {
		t.Fatal("create service: ", err)
	}

	res2, err := s.InitiateMFA(ctx, &mpb.InitiateMFARequest{
		UserID:      u.ID,
		Type:        mpb.MFA_EMAIL,
		Description: "test test test",
	})
	if err != nil {
		t.Fatal("create service: ", err)
	}
	if m.Email != em {
		t.Errorf("email mismatch want: %s, got: %s", em, m.Email)
	}
	if m.Code == "" {
		t.Error("email code is empty")
	}

	_, err = s.ValidateMFA(ctx, &mpb.ValidateMFARequest{
		UserID:  u.ID,
		Type:    mpb.MFA_EMAIL,
		Token:   m.Code,
		EventID: res2.GetEventID(),
	})
	if err != nil {
		t.Fatal("create service: ", err)
	}
}
