//go:build crontest
// +build crontest

package mainpkg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"

	"brank.as/petnet/serviceutil/logging"
)

func TestCron(t *testing.T) {
	logr, hook := test.NewNullLogger()
	logr.SetFormatter(&logrus.JSONFormatter{})
	config := viper.New()
	config.Set("elector.mock_response", "10s")
	config.Set("runtime.environment", "development")

	var gotTicks, gotCron, gotCronCancel, gotLeaderCron, gotLeader int
	srv, err := Setup(config, logr.WithField("test", "cron"),
		WithPath("/tmp/test.socket"),
		WithCron("ticker", NewCronTicker(10*time.Second), func(ctx context.Context) error {
			logr := logging.FromContext(ctx)
			logr.Info("cron tick")
			gotTicks++
			return nil
		}),
		WithCron("cronsched", NewCrontab("TZ=Asia/Manila * * * * *"), func(ctx context.Context) error {
			logr := logging.FromContext(ctx)
			logr.Info("cron schedule")
			select {
			case <-ctx.Done():
				gotCronCancel++
				return ctx.Err()
			case <-time.After(5 * time.Second):
				gotCron++
			}
			return nil
		}),
		WithLeaderFunc(5*time.Second, func(ctx context.Context) error {
			logging.FromContext(ctx).Info("leader func")
			return nil
		}),
		WithLeaderFunc(0, func(ctx context.Context) error {
			logging.FromContext(ctx).Info("leader started")
			gotLeader++
			<-ctx.Done()
			logging.FromContext(ctx).Info("leader stopped")
			return nil
		}),
		WithLeaderCron("leadcron", NewCronTicker(2*time.Second), func(ctx context.Context) error {
			logging.FromContext(ctx).Info("leader task")
			gotLeaderCron++
			return nil
		}),
		WithCron("worker", runner{}, func(ctx context.Context) error {
			logging.FromContext(ctx).Info("worker started")
			<-ctx.Done()
			logging.FromContext(ctx).Info("worker stopped")
			return nil
		}),
		WithLeaderCron("lead-worker", runner{}, func(ctx context.Context) error {
			logging.FromContext(ctx).Info("leader worker started")
			<-ctx.Done()
			logging.FromContext(ctx).Info("leader worker stopped")
			return nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(srv.cron.crontab) != 3 {
		t.Errorf("unexpeted option fail %d", len(srv.cron.crontab))
	}

	cron := srv.cron
	n := time.Now().Add(2 * time.Minute).Truncate(time.Minute).Add(time.Second)
	wantTick := int(time.Until(n).Seconds()) / 10
	ctx, cancel := context.WithDeadline(context.Background(), n)
	t.Cleanup(cancel)
	rootCtx = ctx
	srv.Run()

	const wantCron, wantCronCancel = 1, 1
	switch gotTicks {
	case wantTick, wantTick + 1:
	default:
		t.Error("ticker", cmp.Diff(wantTick, gotTicks))
	}
	if !cmp.Equal(wantCron, gotCron) {
		t.Error("cron", cmp.Diff(wantCron, gotCron))
	}
	if !cmp.Equal(wantCronCancel, gotCronCancel) {
		t.Error("cron cancel", cmp.Diff(wantCronCancel, gotCronCancel))
	}
	if gotLeaderCron == 0 {
		t.Error("leader cron not triggered")
	}
	if gotLeader == 0 {
		t.Error("leader not triggered")
	}
	for _, c := range cron.crontab {
		if c.next.Before(n) && !c.next.IsZero() {
			t.Errorf("scheduling error %q: %v < %v", c.name, c.next, n)
		}
	}
	if t.Failed() || testing.Verbose() {
		for _, e := range hook.Entries {
			s, _ := e.String()
			fmt.Print(s)
		}
	}
}
