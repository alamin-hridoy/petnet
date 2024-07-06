package user

import (
	"context"
	"time"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/integrations/email"
	"brank.as/rbac/usermgm/storage"
)

const resetDur = time.Hour * 24 * 3

type Svc struct {
	usr   UserStore
	org   OrgStore
	mfa   MFAStore
	mail  Mailer
	reset time.Duration

	orgInit       OrgInit
	autoOrg       bool
	autoApp       bool
	invReq        bool
	devEnv        bool
	usrExistErr   bool
	notifyDisable bool
	notifyEnable  bool
}

type Config struct {
	Env               string
	PublicSignup      bool
	AutoApprove       bool
	UserExistError    bool
	NotifyDisableUser bool
	NotifyEnableUser  bool
	ResetDuration     time.Duration
	MFATimeout        time.Duration
}

func New(c Config, us UserStore, org OrgStore, mail Mailer, orgInit OrgInit, mfa MFAStore) *Svc {
	if c.ResetDuration <= 0 {
		c.ResetDuration = resetDur
	}
	return &Svc{
		devEnv:        c.Env == "development",
		usr:           us,
		org:           org,
		mfa:           mfa,
		mail:          mail,
		reset:         c.ResetDuration,
		autoOrg:       c.PublicSignup,
		invReq:        !c.PublicSignup,
		autoApp:       c.AutoApprove,
		orgInit:       orgInit,
		usrExistErr:   c.UserExistError,
		notifyDisable: c.NotifyDisableUser,
		notifyEnable:  c.NotifyEnableUser,
	}
}

// Mailer for connecting to SMTP server
type Mailer interface {
	ForgotPassword(email, name, code string) error
	InviteUser(email, inviteCode string, inv email.Invite) error
	Approved(email string, ap email.Approved) error
	DisableUser(email string, dis email.User) error
	EnableUser(email string, en email.User) error
	// InviteExpired(email string) error
}

type UserStore interface {
	CreatePasswordReset(context.Context, string, time.Duration) (string, error)
	PasswordReset(ctx context.Context, code, pw string) error
	VerifyConfirmationCode(context.Context, string) (*storage.User, error)
	SetPasswordByID(ctx context.Context, id, pw string) error
	SetUsername(ctx context.Context, invite, username string) error
	CreateConfirmationCode(ctx context.Context, userID string) (string, error)
	GetConfirmationCode(ctx context.Context, uid string) (string, error)
	ValidateUserPass(ctx context.Context, id, password string) (string, string, error)
	CreateChangePassword(ctx context.Context, userID, eventID, newPass string) (string, error)

	CreateUser(context.Context, storage.User, storage.Credential) (*storage.User, error)
	ChangePassword(ctx context.Context, userID, eventID string) error
	GetUserByEmail(context.Context, string) (*storage.User, error)
	GetUserByInvite(context.Context, string) (*storage.User, error)
	GetUserByID(context.Context, string) (*storage.User, error)
	UpdateUserByID(context.Context, storage.User) (*storage.User, error)
	GetActiveMFAByType(ctx context.Context, user, mType string) ([]storage.MFA, error)
	DisableUser(ctx context.Context, id string) error
	EnableUser(ctx context.Context, id string) error
	UnlockUser(ctx context.Context, id string) error
}

type OrgStore interface {
	GetOrgByID(context.Context, string) (*storage.Organization, error)
	CreateOrg(context.Context, storage.Organization) (string, error)
}

type OrgInit interface {
	CreateOrg(context.Context, storage.Organization) (string, error)
	ActivateOrg(ctx context.Context, user string, org string) (string, error)
}

type MFAStore interface {
	InitiateMFA(ctx context.Context, c core.MFAChallenge) (*core.MFAChallenge, error)
	MFAuth(ctx context.Context, m core.MFAChallenge) (*core.MFAChallenge, error)
}
