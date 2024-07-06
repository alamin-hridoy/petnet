package auth

import (
	"context"

	"github.com/spf13/viper"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"
)

type Svc struct {
	usr UserStore
	mfa MFAStore
	db  *postgres.Storage

	// lock is how many times we can try login, without user being locked.
	// If set to 0, login attempts are still tracked in DB, but user is not locked.
	lock int

	requireEmail bool
}

type UserStore interface {
	SetPasswordByID(ctx context.Context, id, pw string) error
	GetUser(ctx context.Context, username, password string) (*storage.User, error)
	GetUserByID(context.Context, string) (*storage.User, error)
	GetMFAEventByID(ctx context.Context, id string) (*storage.MFAEvent, error)
}

type MFAStore interface {
	InitiateMFA(ctx context.Context, c core.MFAChallenge) (*core.MFAChallenge, error)
	MFAuth(ctx context.Context, m core.MFAChallenge) (*core.MFAChallenge, error)
	RestartMFA(ctx context.Context, c core.MFAChallenge) (*core.MFAChallenge, error)
}

func New(config *viper.Viper, usr UserStore, mfa MFAStore, st *postgres.Storage) *Svc {
	return &Svc{
		lock:         config.GetInt("user.lockoutCount"),
		usr:          usr,
		mfa:          mfa,
		db:           st,
		requireEmail: config.GetBool("user.requireEmailVerificationForLogin"),
	}
}
