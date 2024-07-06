package mainpkg

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"brank.as/rbac/serviceutil/logging"
)

type Schedule interface {
	cron.Schedule
	Parse() error
}

type ticker struct {
	cron.Schedule
	d time.Duration
}

func NewCronTicker(d time.Duration) Schedule { return &ticker{d: d} }
func (t *ticker) Parse() error {
	if t.d < 0 || t.d < time.Second {
		return fmt.Errorf("invalid tick %q must be positive and at least one second", t.d.String())
	}
	t.Schedule = cron.Every(t.d)
	return nil
}

type crontab struct {
	cron.Schedule
	s string
}

// NewCrontab creates a schedule from a cron string.
func NewCrontab(cron string) Schedule { return &crontab{s: cron} }

func (c *crontab) Parse() error {
	s, err := cron.ParseStandard(c.s)
	if err != nil {
		return fmt.Errorf("invalid cron string %q: %w", c.s, err)
	}
	c.Schedule = s
	return nil
}

// runner is a run-once cron schedule.
type runner struct{}

func (runner) Parse() error             { return nil }
func (runner) Next(time.Time) time.Time { return time.Time{} }

// WithRecovery adds a recovery handler in case of panic.
func WithRecovery(f LeaderFunc) LeaderFunc {
	return func(ctx context.Context) (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, false)]
				err = fmt.Errorf("panic (%s): %s", rec, stack)
			}
		}()
		return f(ctx)
	}
}

// WithFuncTimeout wraps a function with a timeout context.
// Use for setting explicit timeouts in leader or cron functions.
func WithFuncTimeout(to time.Duration, f func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, to)
		defer cancel()
		return f(ctx)
	}
}

type cronJobs struct {
	crontab []cronFunc
}

// Error implements error
func (*cronJobs) Error() string {
	panic("unimplemented")
}

type tickKey struct{}

// GetCronTime returns the time when the cron function was triggered.
func GetCronTime(ctx context.Context) time.Time {
	if tk, ok := ctx.Value(&tickKey{}).(time.Time); ok {
		return tk
	}
	return time.Time{}
}

func setCronTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, &tickKey{}, t)
}

func (c *cronJobs) startCron(ctx context.Context, log *logrus.Entry) error {
	const timeFmt = "2006-01-02 15:04:05"
	if len(c.crontab) == 0 {
		return nil
	}
	crontab := c.crontab
	now := time.Now()
	for i, c := range c.crontab {
		crontab[i].next = c.cron.Next(now)
	}
	sort.SliceStable(crontab, func(i, j int) bool {
		return crontab[i].next.Before(crontab[j].next)
	})
	tickr := time.NewTimer(0)
	<-tickr.C
	wg, ctx := errgroup.WithContext(ctx)
	zcnt := 0
	for n := crontab[zcnt].next; ; n = crontab[zcnt].next {
		wait := time.Until(n)
		if wait < 0 {
			wait = 0
		}
		tickr.Reset(wait)
		select {
		case <-ctx.Done():
			if !tickr.Stop() {
				<-tickr.C
			}
			return wg.Wait()
		case tk := <-tickr.C:
			tkLog := log.WithField("cron trigger", tk.Format(timeFmt))
			tkCtx := logging.WithLogger(setCronTime(ctx, tk), tkLog)
			// skip long-running tasks after they've been started
			for i := zcnt; i < len(crontab); i++ {
				c := crontab[i]
				if c.next.After(tk) {
					break
				}
				wg.Go(func() error {
					if err := c.f(tkCtx); err != nil {
						logging.WithError(err, tkLog).
							WithField("cronjob", c.name).
							Error("cron job failed")
					}
					return nil
				})
				n := c.cron.Next(tk)
				if !n.After(tk) && !n.IsZero() {
					// Rate limit to one per second
					n = tk.Add(time.Second)
				}
				if n.IsZero() {
					zcnt++
					if zcnt == len(crontab) {
						return wg.Wait()
					}
				}
				crontab[i].next = n
			}
		}
		sort.SliceStable(crontab, func(i, j int) bool {
			return crontab[i].next.Before(crontab[j].next)
		})
	}
}

type serverKey struct{}

func (s *server) withServer(ctx context.Context) context.Context {
	return context.WithValue(ctx, &serverKey{}, s)
}

// ServiceReady is for use in Leader or Cron functions.
// Sets the service readiness on the parent server.
func ServiceReady(ctx context.Context, name string, ready readiness) {
	s, ok := ctx.Value(&serverKey{}).(*server)
	if !ok {
		return
	}
	s.Ready(name, ready)
}
