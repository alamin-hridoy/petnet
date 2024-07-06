package mainpkg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"brank.as/rbac/serviceutil/logging"
)

const defElectorPoll = 10 * time.Second

// Leader functions must respect context cancellation, to avoid running after leader lease is lost.
type LeaderFunc = func(context.Context) error

type leader struct {
	le   Leader
	ival time.Duration
	cron []cronFunc
}

func newLead(sock string, lead func(string) (Leader, error), mock string, f map[time.Duration][]LeaderFunc, cronList []cronFunc) (*leader, error) {
	for k, v := range f {
		if k < 0 { // enforce minimum ticker duration
			f[0] = append(f[0], v...)
			delete(f, k)
		}
	}
	for k, v := range f {
		delete(f, k)
		if k == 0 {
			for i, cf := range v {
				cronList = append(cronList, cronFunc{
					name: fmt.Sprintf("persistent job %d", i),
					cron: runner{},
					f:    cf,
				})
			}
			continue
		}
		for i, cf := range v {
			cronList = append(cronList, cronFunc{
				name: fmt.Sprintf("%s job %d", k.String(), i),
				cron: NewCronTicker(k),
				f:    cf,
			})
		}
	}
	for i, c := range cronList {
		if c.cron == nil {
			return nil, fmt.Errorf("leader - invalid cronjob %q: no Schedule", c.name)
		}
		if err := cronList[i].cron.Parse(); err != nil {
			return nil, fmt.Errorf("leader - invalid cronjob schedule %q: %w", c.name, err)
		}
	}
	leadr := &leader{ival: defElectorPoll, cron: cronList}

	if mock != "" {
		mk, err := newMock(mock)
		if err != nil {
			return nil, err
		}
		leadr.le = mk
		return leadr, nil
	}
	if lead == nil {
		return nil, fmt.Errorf("missing elector. use client.WithElector() option")
	}
	cl, err := lead(sock)
	if err != nil {
		return nil, err
	}
	leadr.le = cl
	return leadr, nil
}

func newMock(mock string) (*mockLead, error) {
	mk := &mockLead{}
	l, err := strconv.ParseBool(mock)
	if err != nil {
		d, err := time.ParseDuration(mock)
		if err != nil {
			return nil, fmt.Errorf("invalid leader mock config (boolean or duration) %w", err)
		}
		mk.tick = time.NewTicker(d)
	}
	mk.lead = l
	if mk.tick != nil {
		mk.mu = &sync.RWMutex{}
		mk.cl = make(chan struct{})
		go mk.start()
	}
	return mk, nil
}

func (l *leader) Start(ctx context.Context, log *logrus.Entry) error {
	if l == nil {
		return nil
	}
	if l.le == nil {
		return fmt.Errorf("leader election client not configured")
	}
	defer l.le.Close()
	if len(l.cron) == 0 {
		logging.FromContext(ctx).Debug("no leader workers configured")
		return nil
	}

	return l.startCron(ctx, log)
}

// startCron ...
func (l *leader) startCron(ctx context.Context, log *logrus.Entry) error {
	if len(l.cron) == 0 {
		return nil
	}
	tk := time.NewTicker(l.ival)
	var (
		eg   *errgroup.Group
		gCF  context.CancelFunc
		gctx context.Context
	)
	log = log.WithField("worker", "leader-crontab")
	c := cronJobs{crontab: l.cron}
	for {
		select {
		case <-ctx.Done():
			log.WithError(ctx.Err()).Trace("task context")
			if gCF != nil {
				gCF()
				return eg.Wait()
			}
			return nil
		case <-tk.C:
			// Cancel when leadership is lost.
			if !l.le.IsLead(ctx) {
				if gCF != nil {
					gCF()
					if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
						return err
					}
					gCF, eg = nil, nil
				}
				continue
			}
			if eg != nil { // already running
				continue
			}
			gctx, gCF = context.WithCancel(ctx)
			eg, gctx = errgroup.WithContext(logging.WithLogger(gctx, log))
			eg.Go(func() error {
				switch err := c.startCron(gctx, log); err {
				case nil, context.Canceled, context.DeadlineExceeded:
					return nil
				default:
					return err
				}
			})
		}
	}
}

type Leader interface {
	Close() error
	IsLead(context.Context) bool
}

type sigLead struct {
	lead    bool
	mu      *sync.RWMutex
	sig     chan bool
	sigCron chan bool
	Leader
}

func sigWrap(lc Leader) *sigLead {
	return &sigLead{mu: &sync.RWMutex{}, sig: make(chan bool, 1), Leader: lc}
}

func (s *sigLead) IsLead(ctx context.Context) bool {
	l := s.Leader.IsLead(ctx)
	defer func() {
		s.mu.Lock()
		s.lead = l
		s.mu.Unlock()
	}()
	s.mu.RLock()
	if l != s.lead {
		select {
		case s.sig <- l:
		default:
		}
		select {
		case s.sigCron <- l:
		default:
		}
	}
	s.mu.RUnlock()
	return l
}

type mockLead struct {
	lead bool
	mu   *sync.RWMutex
	tick *time.Ticker
	cl   chan struct{}
}

// start ticker for toggle configuration.
func (m *mockLead) start() {
	for {
		select {
		case <-m.tick.C:
			m.mu.Lock()
			m.lead = !m.lead
			m.mu.Unlock()
		case <-m.cl:
			return
		}
	}
}

func (m *mockLead) Close() error {
	if m.tick != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.tick.Stop()
		close(m.cl)
	}
	return nil
}

func (m *mockLead) IsLead(context.Context) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lead
}
