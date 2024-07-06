package email

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/dschmidt/go-layerfs"

	"brank.as/rbac/usermgm/integrations/email/assets"
)

type MockSender struct {
	Email         string
	Code          string
	server        string
	port          int
	username      string
	password      string
	fromMail      string
	mainURL       string
	fromName      string
	idpURL        string
	signupURL     string
	ecRedirectURL string
	siteURL       string
	tls           *tls.Config

	templates
	subjects Subjects
}

func NewMock(conf MailerConfig, assetFS fs.FS) (Mailer, error) {
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
	}
	switch (*template.Template)(nil) {
	case t.forgotPW:
		return nil, fmt.Errorf("missing forgot password template")
	case t.confirmEmail:
		return nil, fmt.Errorf("missing confirm email template")
	}

	if conf.SignupURL == "" {
		conf.SignupURL = conf.IDPURL
	}

	return &MockSender{
		mainURL:       conf.MainURL,
		idpURL:        conf.IDPURL,
		signupURL:     conf.SignupURL,
		ecRedirectURL: conf.ECRedirectURL,
		siteURL:       conf.SiteURL,

		templates: t,
	}, nil
}

func (ms *MockSender) ForgotPassword(email, name, code string) error {
	subj := ms.subjects.ForgotPW
	fmt.Printf("email sent, to: %s, subject: %s, url: %s/user/set-password?reset_code=%s", email, subj, ms.idpURL, code)
	return nil
}

func (ms *MockSender) ConfirmEmail(email, code string) error {
	subj := ms.subjects.ConfirmEmail
	url := fmt.Sprintf("%s?confirm_code=%s", ms.ecRedirectURL, code)
	fmt.Printf("email sent, to: %s, subject: %s, url: %s\n", email, subj, url)
	i := 1
	for {
		p := "/tmp/code" + strconv.Itoa(i)
		if _, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				err := ioutil.WriteFile(p, []byte(code), 0o755)
				if err != nil {
					return err
				}
				break
			}
		}
		i++
	}
	return nil
}

func (ms *MockSender) EmailMFA(email, code string) error {
	ms.Email = email
	ms.Code = code
	subj := ms.subjects.EmailMFA
	fmt.Printf("email sent, to: %s, subject: %s, code: %s", email, subj, code)
	return nil
}

func (ms *MockSender) InviteUser(email, inviteCode string, f Invite) error {
	subj := ms.subjects.UserInvite
	fmt.Printf("email sent, to: %s, subject: %s, url: %s?invite_code=%s", email, subj, ms.signupURL, inviteCode)
	return nil
}

func (ms *MockSender) Approved(email string, f Approved) error {
	subj := ms.subjects.InviteApproved
	f.RedirectURL = ms.siteURL + "/login"
	f.Email = email
	fmt.Printf("email sent, to: %s, subject: %s, url: %s", email, subj, ms.siteURL+"/login")
	return nil
}

func (ms *MockSender) DisableUser(email string, f User) error {
	subj := ms.subjects.UserDisable
	fmt.Printf("email sent, to: %s, subject: %s, url: %s", email, subj, ms.siteURL)
	return nil
}

func (ms *MockSender) EnableUser(email string, f User) error {
	subj := ms.subjects.UserEnable
	fmt.Printf("email sent, to: %s, subject: %s, url: %s", email, subj, ms.siteURL)
	return nil
}

type NoopSender struct{}

func (NoopSender) ForgotPassword(_, _, _ string) error    { return nil }
func (NoopSender) ConfirmEmail(_, _ string) error         { return nil }
func (NoopSender) EmailMFA(_, _ string) error             { return nil }
func (NoopSender) InviteUser(_, _ string, _ Invite) error { return nil }
func (NoopSender) Approved(_ string, _ Approved) error    { return nil }
func (NoopSender) DisableUser(_ string, _ User) error     { return nil }
func (NoopSender) EnableUser(_ string, _ User) error      { return nil }
