//go:build bootstrap
// +build bootstrap

package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"brank.as/rbac/serviceutil/logging"
	client "brank.as/rbac/svcutil/hydraclient"
	perm "brank.as/rbac/usermgm/core/permissions"
	"brank.as/rbac/usermgm/core/svcacct"
	"brank.as/rbac/usermgm/integrations/keto"
	"brank.as/rbac/usermgm/storage/postgres"
)

func TestBootstrap(t *testing.T) {
	const dbConnEnv = "DATABASE_CONNECTION"
	ddlConnStr := os.Getenv(dbConnEnv)
	if ddlConnStr == "" {
		t.Skipf("%s is not set, skipping", dbConnEnv)
	}
	st, cleanup := postgres.NewTestStorage(ddlConnStr, filepath.Join("migrations", "sql"))
	t.Cleanup(cleanup)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	v := viper.New()
	v.Set("bootstrap.email", "test@example.com")
	v.Set("bootstrap.name", "test")
	v.Set("bootstrap.envlist", []string{"sandbox", "live"})

	log := logging.NewLogger().WithField("test", "bootstrap")
	log.Logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp:  true,
		DisableHTMLEscape: false,
		PrettyPrint:       false,
	})
	ctx = logging.WithLogger(ctx, log)

	cl, err := client.NewAdminClient("http://localhost:4445")
	if err != nil {
		t.Fatal(err)
	}

	k := keto.New("localhost:4466")
	sa := svcacct.New(v, cl, st)
	p := perm.New(st, k)

	if err := Bootstrap(ctx, v, log, st, p, sa); err != nil {
		t.Error(err)
	}
}
