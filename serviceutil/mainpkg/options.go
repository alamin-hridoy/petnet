package mainpkg

import (
	"context"
	"strings"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type cronFunc struct {
	name string
	cron Schedule
	next time.Time
	f    LeaderFunc
}

type opFunc struct {
	svcName string
	f       func(context.Context) error
}

func ServiceConfig(filename string) (*viper.Viper, error) {
	cfg := viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer(".", "_")))
	cfg.SetConfigFile(filename)
	cfg.SetConfigType("ini")
	cfg.AutomaticEnv()
	return cfg, cfg.ReadInConfig()
}

// WithCron schedules a background task.
// This will run on every replica, if only one pod needs to run, use WithLeaderCron.
// Especially useful for cache updates or other in-memory tasks.
func WithCron(name string, cron Schedule, f func(context.Context) error) Option {
	return func(conf *Config) {
		conf.cronFuncs = append(conf.cronFuncs, cronFunc{name: name, cron: cron, f: f})
	}
}

// WithPath option override the path in config.
func WithPath(path string) Option { return func(conf *Config) { conf.socketPath = path } }

// WithLeaderCron schedules a background task on the leader pod only.
// This will only run on one replica, if each pod needs the task, use WithCron.
// Especially useful for database jobs, scheduled reports, or other persistent state tasks.
//
// Leader functions must respect context cancellation, to avoid running after leader lease is lost.
func WithLeaderCron(name string, cron Schedule, f LeaderFunc) Option {
	return func(conf *Config) {
		conf.leadCron = append(conf.leadCron, cronFunc{name: name, cron: cron, f: f})
	}
}

// OptionList allows configuring a set of options for a given environment,
// avoiding the need to append all standard service options
// to an environment-specific set of options.
func OptionList(o []Option) Option {
	return func(c *Config) {
		for _, opt := range o {
			opt(c)
		}
	}
}

func WithServerOpts(opt ...grpc.ServerOption) Option {
	return func(conf *Config) { conf.grpcOpts = append(conf.grpcOpts, opt...) }
}

func WithKeepAlive(pingTime, timeout time.Duration) Option {
	return func(c *Config) {
		c.grpcOpts = append(c.grpcOpts,
			grpc.KeepaliveParams(keepalive.ServerParameters{
				Time:    pingTime,
				Timeout: timeout,
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				// MinTime is the minimum amount of time a client should wait before sending
				// a keepalive ping.
				MinTime:             pingTime,
				PermitWithoutStream: true,
			}),
		)
	}
}
