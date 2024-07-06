package handler

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"

	"brank.as/rbac/idp/auth"
	"brank.as/rbac/serviceutil/logging"
)

const (
	cookieStoreName        = "default-cookie-store"
	openIDKey              = "session-open-id"
	cookieRedirectEmailKey = "redirect-email"
	authRetryKey           = "auth-retry"
)

type sessionData struct {
	MFAEvent  string            `json:"mfa_event"`
	MFAType   string            `json:"mfa_type"`
	LoginCtx  map[string]string `json:"login_ctx"`
	OpenIDCtx map[string]string `json:"open_id_ctx"`
}

type identity struct {
	Challenge   string
	UserID      string
	Email       string
	CodeError   string
	Remember    bool
	RememberFor time.Duration
	Sesssion    sessionData
}

type loginFormParams struct {
	emailPrefill string
	challenge    string
	loginErrors  map[string]string
	afterSignup  bool
	callbackErr  bool
}

func (s *server) getLoginHandler(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context()).WithField("method", "getLoginHandler")
	ctx := logging.WithLogger(r.Context(), log)
	chlg := r.URL.Query().Get("login_challenge")

	params := admin.NewGetLoginRequestParams().WithLoginChallenge(chlg)
	lr, err := s.hydra.GetLoginRequest(params)
	if err != nil {
		_, ok := err.(*admin.GetLoginRequestGone)
		if ok {
			log.WithField("challenge", chlg).Error("login challenge is already processed")
			http.Redirect(w, r, s.ErrRedirURL, http.StatusFound)
			return
		}
		log.WithError(err).Error("failed to get login request")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if login := s.checkSkip(w, r.WithContext(ctx), lr, chlg); !login {
		return
	}

	authBools, err := getAuthFields(*lr.Payload.RequestURL)
	if err != nil {
		log.Printf("failed to parse auth query: %s", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}

	callbackErr := authBools["callbackErr"]
	afterSignup := authBools["afterSignup"]
	login := authBools["login"]
	preRegister := authBools["pre_register"]
	register := authBools["register"]

	if login || preRegister {
		email, err := getEmailPrefill(*lr.Payload.RequestURL)
		if err != nil {
			// ignore, not important
			log.Printf("failed to get emailPrefill: %+v", err)
		}
		s.showLoginForm(
			&loginFormParams{emailPrefill: email, challenge: chlg, callbackErr: callbackErr},
			lr.Payload.Client, w, r)
		return
	}

	if register {
		// s.showRegisterForm(&registerFormParams{Challenge: challengeKey}, w, r)
		return
	}

	extURL, err := s.getExternalAuthURL(*lr.Payload.RequestURL, chlg)
	if err != nil {
		log.Printf("failed to get external auth URL params: %s", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
	if extURL != "" {
		http.Redirect(w, r, extURL, http.StatusFound)
		return
	}

	// hydra was unable to authenticate user, showing the login form
	s.showLoginForm(&loginFormParams{
		challenge:   chlg,
		afterSignup: afterSignup,
		callbackErr: callbackErr,
	}, lr.Payload.Client, w, r)
}

func (s *server) showLoginForm(p *loginFormParams, cl *models.OAuth2Client, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sess, err := s.cookieStore.Get(r, cookieStoreName)
	if err != nil {
		log.Printf("getting cookie store  %q", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}

	// if an email exists it means this is a redirect after user registration
	email, _ := sess.Values[cookieRedirectEmailKey].(string)
	sess.Values[cookieRedirectEmailKey] = ""
	if err := s.cookieStore.Save(r, w, sess); err != nil {
		log.Printf("saving session  %q", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
	tmpl := s.loginTpl
	if cfg, err := auth.ParseClientConfig(cl); err == nil && cfg.LoginTmpl != "" {
		tmpl = s.lookup(ctx, cfg.LoginTmpl, s.loginTpl)
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, struct {
		Challenge        string
		LoginErrors      map[string]string
		CSRFField        template.HTML
		IsProdEnv        bool
		AfterSignup      bool
		EmailPrefill     string
		RedirectionEmail string
		ProjectName      string
		urls
	}{
		Challenge:        p.challenge,
		LoginErrors:      p.loginErrors,
		CSRFField:        csrf.TemplateField(r),
		IsProdEnv:        s.environment == "production",
		AfterSignup:      p.afterSignup,
		EmailPrefill:     "",
		RedirectionEmail: email,
		ProjectName:      s.projectName,
		urls:             s.u,
	}); err != nil {
		log.Printf("failed to render login form: %s", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
}

type loginForm struct {
	Challenge string
	Email     string
	Password  string
	Remember  bool
	Admin     bool
}

func (s *server) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	fields := map[string]bool{
		"challenge": true,
		"email":     true,
		"password":  true,
		"remember":  true,
		"admin":     true,
	}
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "postLoginHandler")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var form loginForm
	if err := s.dcd.Decode(&form, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	extra := map[string]string{}
	for k, v := range r.PostForm {
		k = strings.ToLower(k)
		if fields[k] || len(v) == 0 {
			continue
		}
		for i := range v {
			v[i] = strings.TrimSpace(v[i])
		}
		extra[k] = strings.Join(v, ",")
	}
	log.WithField("extra", extra).Debug("form parsed")
	form.Email = strings.TrimSpace(form.Email)
	form.Password = strings.TrimSpace(form.Password)

	loginErrors := make(map[string]string)
	switch "" {
	case form.Email + form.Password:
		loginErrors["emailError"] = "Please enter email here"
		loginErrors["passwordError"] = "Please enter password here"
	case form.Email:
		loginErrors["username"] = "Please, provide a valid email address"
	case form.Password:
		loginErrors["email"] = form.Email
		loginErrors["passwordError"] = "Please enter password here"
	}
	log.WithField("errors", loginErrors).Debug("post handler errors")

	lr, err := s.hydra.GetLoginRequest(admin.NewGetLoginRequestParamsWithContext(ctx).
		WithLoginChallenge(form.Challenge))
	if err != nil {
		_, ok := err.(*admin.GetLoginRequestGone)
		if ok {
			log.WithField("challenge", form.Challenge).Error("login challenge is already processed")
			http.Redirect(w, r, s.ErrRedirURL, http.StatusFound)
			return
		}
		logging.WithError(err, log).Error("failed to get login request")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}

	if len(loginErrors) > 0 {
		s.showLoginForm(&loginFormParams{
			challenge:   form.Challenge,
			loginErrors: loginErrors,
		}, lr.Payload.Client, w, r)
		return
	}

	ctx = logging.WithLogger(ctx, log)
	ac := s.authCl(w, r, ctx, lr.Payload.Client)
	if ac == nil {
		log.WithField("client", *lr.Payload.Client).
			WithField("client_list", s.authClient).
			Error("missing auth client")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}

	ident, err := ac.Authenticate(ctx, auth.Challenge{
		Username:    form.Email,
		Password:    form.Password,
		HydraClient: lr.Payload.Client.ClientID,
		Extra:       parseURLParams(ctx, *lr.Payload.RequestURL, extra),
	}, nil)
	if err != nil {
		log.WithField("error", err).Error("error authenticate")
		switch e := auth.FromError(err); e.Code {
		case auth.NotFound, auth.PermissionDenied:
			if e.TrackAttempts && e.AttemptRemain == 0 {
				e.MergeDetails(map[string]string{"AttemptRemain": "You've reached the maximum allowed authentication input.Your account has been temporarily blocked."})
				s.showLoginForm(&loginFormParams{
					challenge:   form.Challenge,
					loginErrors: e.Errors,
				}, lr.Payload.Client, w, r)
				return
			}
			e.MergeDetails(map[string]string{
				"username": "Invalid email or password. Please try again.",
				"email":    form.Email,
			})
			s.showLoginForm(&loginFormParams{
				challenge:   form.Challenge,
				loginErrors: e.Errors,
			}, lr.Payload.Client, w, r)
		case auth.ExpiredPassword:
			e.MergeDetails(map[string]string{
				"username": "Password has expired. Please reset and try again.",
			})
			s.showLoginForm(&loginFormParams{
				challenge:   form.Challenge,
				loginErrors: e.Errors,
			}, lr.Payload.Client, w, r)
		case auth.ExistingSession:
			loginErrors["existingSession"] = "You are already signed in."
			s.showLoginForm(&loginFormParams{
				challenge:   form.Challenge,
				loginErrors: loginErrors,
			}, lr.Payload.Client, w, r)
		case auth.InvalidRecord:
			e.MergeDetails(map[string]string{"username": "Record invalid. Please contact support."})
			s.showLoginForm(&loginFormParams{
				challenge:   form.Challenge,
				loginErrors: e.Errors,
			}, lr.Payload.Client, w, r)
		case auth.Unknown:
			log.WithError(err).Error("unknown session error")
			http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		default:
			log.WithError(err).Error("login call failed")
			http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		}
		return
	}

	id := identity{
		Challenge:   form.Challenge,
		UserID:      ident.UserID,
		Remember:    form.Remember,
		Email:       form.Email,
		RememberFor: time.Duration(ac.Remember()),
		Sesssion: sessionData{
			MFAEvent:  ident.MFAEventID,
			MFAType:   ident.MFAType,
			LoginCtx:  ident.SessionContext(),
			OpenIDCtx: ident.LoginContext(),
		},
	}
	log.WithField("idt", ident).Debug("identity")
	if ident.MFAEventID != "" {
		s.serveOTPLogin(ctx, w, r, id, lr.Payload.Client)
		return
	}

	s.acceptLogin(logging.WithLogger(ctx, log), w, r, id)
}

func getAuthFields(authURL string) (map[string]bool, error) {
	parsed, err := url.Parse(authURL)
	if err != nil {
		return nil, err
	}
	m := map[string]bool{}

	q := parsed.Query()
	for _, t := range []string{
		"admin_login",
		"login",
		"register",
		"pre_register",
		"afterSignup",
		"autoLogin",
		"callbackErr",
	} {
		m[t] = q.Get(t) == "1"
	}

	return m, nil
}

func getEmailPrefill(authURL string) (string, error) {
	const prefillKey = "email_prefill"
	parsed, err := url.Parse(authURL)
	if err != nil {
		return "", err
	}
	return parsed.Query().Get(prefillKey), nil
}

func (s *server) getExternalAuthURL(authURL, challenge string) (string, error) {
	const authProviderKey = "auth_provider"
	parsed, err := url.Parse(authURL)
	if err != nil {
		return "", err
	}
	authProvider := parsed.Query().Get(authProviderKey)
	switch authProvider {
	default:
		return "", nil
	}
}

func (s *server) acceptLogin(ctx context.Context, w http.ResponseWriter, r *http.Request, idt identity) {
	log := logging.FromContext(ctx).WithField("method", "acceptLogin")

	pms := admin.NewGetLoginRequestParams().WithLoginChallenge(idt.Challenge)
	req, err := s.hydra.GetLoginRequest(pms)
	if err != nil {
		_, ok := err.(*admin.GetLoginRequestGone)
		if ok {
			log.WithField("challenge", idt.Challenge).Error("login challenge is already processed")
			http.Redirect(w, r, s.ErrRedirURL, http.StatusFound)
			return
		}
		logging.WithError(err, log).Error("failed to get login request")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}
	idt.Sesssion.LoginCtx["owner"] = req.Payload.Client.Owner

	log.WithField("auth client", *req.Payload.Client).Debug("accepting")
	params := admin.NewAcceptLoginRequestParams().
		WithLoginChallenge(idt.Challenge).
		WithBody(&models.AcceptLoginRequest{
			Context:     idt.Sesssion,
			Remember:    idt.Remember,
			RememberFor: int64(idt.RememberFor),
			Subject:     &idt.UserID,
		})
	acceptLogin, err := s.hydra.AcceptLoginRequest(params)
	if err != nil {
		logging.WithError(err, log).Error("failed to accept login request")
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
	log.WithField("accept payload", *acceptLogin.Payload).Debug("accepted")
	http.Redirect(w, r, *acceptLogin.Payload.RedirectTo, http.StatusFound)
}

func (s *server) checkSkip(w http.ResponseWriter, r *http.Request, lr *admin.GetLoginRequestOK, chlg string) bool {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("func", "checkSkip")
	// hydra was already able to authenticate the user, just accepts login and redirect
	if !*lr.Payload.Skip {
		log.Trace("no skip flag")
		return true
	}
	log = log.WithField("login", "skipped")
	ac := s.authCl(w, r, ctx, lr.Payload.Client)
	if ac == nil {
		return true
	}
	ident, err := ac.Lookup(ctx, auth.Challenge{
		ID:          loginSubj(lr),
		HydraClient: lr.Payload.Client.ClientID,
		Extra:       parseURLParams(ctx, *lr.Payload.RequestURL, nil),
	})
	if auth.FromError(err).Code == auth.NotFound {
		// Happen if the user is no longer found in database or account has been disabled.
		s.rejectLoginRedir(w, r, chlg, *lr.Payload.Subject)
		return false
	}
	if auth.FromError(err).Code == auth.InvalidRecord {
		return true
	}
	if err != nil {
		log.WithError(err).Error("failed to lookup identity")
		return true
	}
	if ident.ForceLogin {
		if err := s.rejectLogin(w, r, chlg, *lr.Payload.Subject); err != nil {
			logging.WithError(err, log).Error("reject login id")
		}
		log.WithField("ident", ident).Trace("force login")
		p, err := url.Parse(*lr.Payload.RequestURL)
		if err != nil {
			logging.WithError(err, log).Error("request URL")
			return true
		}
		q := p.Query()
		q.Add("login_hint", ident.UserID)
		p.RawQuery = q.Encode()
		http.Redirect(w, r, p.String(), http.StatusFound)
		return false
	}
	s.acceptLogin(logging.WithLogger(ctx, log), w, r, identity{
		Challenge:   chlg,
		UserID:      loginSubj(lr),
		Remember:    false,
		RememberFor: time.Duration(ac.Remember()),
		Sesssion: sessionData{
			LoginCtx:  ident.SessionContext(),
			OpenIDCtx: ident.LoginContext(),
		},
	})
	return false
}

// rejectLoginRedir fails the login attempt and redirects to configured error url.
func (s *server) rejectLoginRedir(w http.ResponseWriter, r *http.Request, chg, sub string) {
	if err := s.rejectLogin(w, r, chg, sub); err != nil {
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, s.ErrRedirURL, http.StatusFound)
}

func (s *server) rejectLogin(w http.ResponseWriter, r *http.Request, chg, sub string) error {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "rejectLogin")

	_, err := s.hydra.RejectLoginRequest(
		admin.NewRejectLoginRequestParams().
			WithLoginChallenge(chg).
			WithBody(&models.RejectRequest{
				ErrorHint: "failed to lookup user",
			}))
	if err != nil {
		log.WithError(err).Error("failed to reject login")
		return err
	}

	_, err = s.hydra.RevokeAuthenticationSession(
		admin.NewRevokeAuthenticationSessionParamsWithContext(ctx).WithSubject(sub))
	if err != nil {
		log.WithError(err).Error("failed to revoke authentication session")
		return err
	}
	return nil
}

func parseURLParams(ctx context.Context, reqURL string, extra map[string]string) map[string]string {
	u, err := url.Parse(reqURL)
	if err != nil {
		logging.WithError(err, logging.FromContext(ctx)).Error("request url invalid")
		return extra
	}
	if extra == nil {
		extra = map[string]string{}
	}
	for k, v := range u.Query() {
		key := "url." + k
		for i := range v {
			v[i] = strings.TrimSpace(v[i])
		}
		if len(v) == 0 {
			v = []string{"true"}
		}
		extra[key] = strings.Join(v, ",")
	}
	return extra
}

func loginSubj(lr *admin.GetLoginRequestOK) string {
	if lr == nil {
		return ""
	}
	if lr.Payload.OidcContext != nil &&
		lr.Payload.OidcContext.LoginHint != "" {
		return lr.Payload.OidcContext.LoginHint
	}
	return *lr.Payload.Subject
}
