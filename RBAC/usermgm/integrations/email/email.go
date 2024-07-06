package email

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/mail"
	"os"
	"path/filepath"
	"time"

	"github.com/dschmidt/go-layerfs"

	"brank.as/rbac/svcutil/mailer"
	"brank.as/rbac/usermgm/integrations/email/assets"
)

type Mailer interface {
	ForgotPassword(email, name, code string) error
	ConfirmEmail(email, code string) error
	EmailMFA(email, code string) error
	InviteUser(email, inviteCode string, f Invite) error
	Approved(email string, f Approved) error
	DisableUser(email string, f User) error
	EnableUser(email string, f User) error
}

type MailSender struct {
	mailer *mailer.Email

	mainURL       string
	idpURL        string
	signupURL     string
	ecRedirectURL string
	siteURL       string

	subjects Subjects
	templates
}

type Subjects struct {
	ForgotPW       string
	ConfirmEmail   string
	EmailMFA       string
	UserInvite     string
	InviteApproved string
	UserDisable    string
	UserEnable     string
}

type templates struct {
	assetFS fs.FS

	forgotPW       *template.Template
	confirmEmail   *template.Template
	emailMFA       *template.Template
	userInvite     *template.Template
	inviteApproved *template.Template
	userDisable    *template.Template
	userEnable     *template.Template
}

type MailerConfig struct {
	Server                                             string
	Port                                               int
	Timeout                                            time.Duration
	Username, Password, FromMail, FromName             string
	MainURL, IDPURL, SignupURL, ECRedirectURL, SiteURL string
	CACert, CACommonName                               string
	Subjects                                           Subjects
}

type Invite struct {
	Username        string
	UserEmail       string
	Duration        string
	ExpiryDate      string
	CustomEmailData map[string]string
	RedirectURL     string
	Email           string
}

type User struct {
	CustomEmailData map[string]string
	Email           string
	FirstName       string
	LastName        string
	OrgName         string
	RedirectURL     string
}

type Approved struct {
	CompanyName string
	OrgID       string
	RedirectURL string
	Email       string
}

func New(conf MailerConfig, assetFS fs.FS) (Mailer, error) {
	tmpl, err := template.ParseFS(assets.DefaultFS, "templates/*.html")
	if err != nil {
		return nil, err
	}

	ttmpl, err := tmpl.ParseFS(assetFS, "templates/*.html")
	if err == nil {
		tmpl = ttmpl
	}

	t := templates{
		assetFS: layerfs.New(assetFS, assets.DefaultFS),

		forgotPW:       tmpl.Lookup("reset-password.html"),
		confirmEmail:   tmpl.Lookup("confirm-email.html"),
		emailMFA:       tmpl.Lookup("email-mfa.html"),
		userInvite:     tmpl.Lookup("user-invite.html"),
		inviteApproved: tmpl.Lookup("invite-approved.html"),
		userDisable:    tmpl.Lookup("user-disable.html"),
		userEnable:     tmpl.Lookup("user-enable.html"),
	}
	switch (*template.Template)(nil) {
	case t.forgotPW:
		return nil, fmt.Errorf("missing forgot password template")
	case t.confirmEmail:
		return nil, fmt.Errorf("missing confirm email template")
	}
	mc := mailer.EmailConfig{
		Server:   conf.Server,
		Port:     conf.Port,
		Timeout:  conf.Timeout,
		Username: conf.Username,
		Password: conf.Password,
		DefaultFrom: mail.Address{
			Name:    conf.FromName,
			Address: conf.FromMail,
		},
	}

	if conf.CACert != "" {
		f, err := os.Open(conf.CACert)
		if err != nil {
			return nil, err
		}
		mc.CACert = f
		defer f.Close()
		mc.CACommonName = conf.CACommonName
	}

	if conf.SignupURL == "" {
		conf.SignupURL = conf.IDPURL
	}

	mailr, err := mailer.NewEmail(mc)
	if err != nil {
		return nil, err
	}

	return &MailSender{
		mailer: mailr,

		mainURL:       conf.MainURL,
		idpURL:        conf.IDPURL,
		signupURL:     conf.SignupURL,
		ecRedirectURL: conf.ECRedirectURL,
		siteURL:       conf.SiteURL,

		templates: t,
		subjects:  conf.Subjects,
	}, nil
}

func (ms *MailSender) sendMail(
	email, toName, subject string,
	tmpl *template.Template,
	data interface{},
	imgDir string,
) error {
	e := mailer.EmailMessage{
		To:      mail.Address{Name: toName, Address: email},
		Subject: subject,
		Msg:     tmpl,
		Args:    data,
	}
	if err := e.AssetDir(ms.assetFS, filepath.Join("images", imgDir)); err != nil {
		return err
	}
	return ms.mailer.Send(e)
}

type forgotPasswordForm struct {
	ResetURL    string
	RedirectURL string
}

func (ms *MailSender) ForgotPassword(email, name, code string) error {
	const imgDir = "forgot-password"
	f := forgotPasswordForm{
		ResetURL: fmt.Sprintf("%s/user/set-password?reset_code=%s", ms.idpURL, code),
	}
	return ms.sendMail(email, "", ms.subjects.ForgotPW, ms.forgotPW, f, imgDir)
}

func (ms *MailSender) ConfirmEmail(email, code string) error {
	const imgDir = "confirm-email"
	f := struct{ ConfirmEmailURL string }{
		ConfirmEmailURL: fmt.Sprintf("%s?confirm_code=%s", ms.ecRedirectURL, code),
	}
	return ms.sendMail(email, "", ms.subjects.ConfirmEmail, ms.confirmEmail, f, imgDir)
}

func (ms *MailSender) EmailMFA(email, code string) error {
	const imgDir = "confirm-mfa"
	f := struct{ MFACode string }{MFACode: code}
	return ms.sendMail(email, "", ms.subjects.EmailMFA, ms.emailMFA, f, imgDir)
}

func (ms *MailSender) InviteUser(email, inviteCode string, f Invite) error {
	const imgDir = "user-invite"
	f.RedirectURL = fmt.Sprintf("%s?invite_code=%s", ms.signupURL, inviteCode)
	f.Email = email
	return ms.sendMail(email, "", ms.subjects.UserInvite, ms.userInvite, f, imgDir)
}

func (ms *MailSender) Approved(email string, f Approved) error {
	const imgDir = "invite-approved"
	f.RedirectURL = ms.siteURL + "/login"
	f.Email = email
	return ms.sendMail(email, "", ms.subjects.InviteApproved, ms.inviteApproved, f, imgDir)
}

func (ms *MailSender) DisableUser(email string, f User) error {
	const imgDir = "user-disable"
	f.RedirectURL = ms.siteURL
	f.Email = email
	return ms.sendMail(email, "", ms.subjects.UserDisable, ms.userDisable, f, imgDir)
}

func (ms *MailSender) EnableUser(email string, f User) error {
	const imgDir = "user-enable"
	f.RedirectURL = ms.siteURL + "/login"
	f.Email = email
	return ms.sendMail(email, "", ms.subjects.UserEnable, ms.userEnable, f, imgDir)
}
