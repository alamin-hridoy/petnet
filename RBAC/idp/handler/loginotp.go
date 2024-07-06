package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"

	"brank.as/rbac/idp/auth"
	"brank.as/rbac/serviceutil/logging"
)

const otpSession = "otp-validate"

func (s *server) serveOTPLogin(ctx context.Context, w http.ResponseWriter, r *http.Request, idt identity, cl *models.OAuth2Client) {
	log := logging.FromContext(ctx).WithField("method", "serveotplogin")

	sess, err := s.cookieStore.Get(r, cookieStoreName)
	if err != nil {
		logging.WithError(err, log).Error("session cookie")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}
	idt.Sesssion.LoginCtx = nil
	idt.Sesssion.OpenIDCtx = nil
	sess.Values[otpSession] = idt
	if err := sess.Save(r, w); err != nil {
		logging.WithError(err, log).Error("session cookie")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}

	tmpl := s.otpTpl
	if cfg, err := auth.ParseClientConfig(cl); err == nil && cfg.OTPTmpl != "" {
		tmpl = s.lookup(ctx, cfg.OTPTmpl, s.otpTpl)
	}
	// TODO(Chad): Handle returned MFA types
	if err := tmpl.Execute(w, otpTemplateData{
		OTPForm: OTPForm{
			CSRFField:       csrf.TemplateField(r),
			ProcessEndpoint: "/login/otp",
			CodeError:       idt.CodeError,
			Email:           idt.Email,
		},
		urls:        s.u,
		ProjectName: s.projectName,
	}); err != nil {
		log.WithError(err).Error("failed to render OTP form")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
	}
}

func (s *server) retryOTPLogin(ctx context.Context, w http.ResponseWriter, r *http.Request, form OTPForm) {
	log := logging.FromContext(ctx).WithField("method", "retryotplogin")
	log.WithField("form", form).Debug("retry otp login handler")

	sess, err := s.cookieStore.Get(r, cookieStoreName)
	if err != nil {
		logging.WithError(err, log).Error("session cookie")
		http.Redirect(w, r, s.u.LoginURL, http.StatusTemporaryRedirect)
		return
	}
	idt, ok := sess.Values[otpSession].(identity)
	if !ok {
		log.WithError(err).Error("session invalid")
		http.Redirect(w, r, s.u.LoginURL, http.StatusTemporaryRedirect)
		return
	}

	lr, err := s.hydra.GetLoginRequest(admin.NewGetLoginRequestParamsWithContext(ctx).
		WithLoginChallenge(idt.Challenge))
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

	ctx = logging.WithLogger(ctx, log)
	ac := s.authCl(w, r, ctx, lr.Payload.Client)
	if ac == nil {
		log.WithField("client", *lr.Payload.Client).
			WithField("client_list", s.authClient).
			Error("missing auth client")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}

	ident, err := ac.ResetMFA(ctx, auth.Challenge{
		ID:          idt.UserID,
		HydraClient: lr.Payload.Client.ClientID,
	}, auth.OTPChallenge{
		Code:  form.Code,
		Type:  idt.Sesssion.MFAType,
		Event: idt.Sesssion.MFAEvent,
	})
	if err != nil {
		switch e := auth.FromError(err); e.Code {
		case auth.NotFound:
			// Happen if the user is no longer found in database or account has been disabled.
			s.rejectLoginRedir(w, r, idt.Challenge, idt.UserID)
			return
		default:
			log.WithError(err).Error("failed to lookup identity")
			http.Error(w, genericErrMsg, http.StatusInternalServerError)
			return
		}
	}
	s.serveOTPLogin(ctx, w, r, identity{
		Challenge:   idt.Challenge,
		UserID:      ident.UserID,
		Remember:    idt.Remember,
		Email:       idt.Email,
		RememberFor: time.Duration(ac.Remember()),
		Sesssion: sessionData{
			MFAEvent: ident.MFAEventID,
			MFAType:  ident.MFAType,
		},
	}, lr.Payload.Client)
}

func (s *server) postOTPLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "postotplogin")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var form OTPForm
	if err := s.dcd.Decode(&form, r.PostForm); err != nil {
		log.WithError(err).Error("failed to parse OTP form")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}
	log.WithField("form", form).Debug("post otp login handler")

	if form.Retry != "" {
		s.retryOTPLogin(ctx, w, r, form)
		return
	}

	sess, err := s.cookieStore.Get(r, cookieStoreName)
	if err != nil {
		logging.WithError(err, log).Error("session cookie")
		http.Redirect(w, r, s.u.LoginURL, http.StatusTemporaryRedirect)
		return
	}
	idt, ok := sess.Values[otpSession].(identity)
	if !ok {
		log.WithError(err).Error("session invalid")
		http.Redirect(w, r, s.u.LoginURL, http.StatusTemporaryRedirect)
		return
	}

	lr, err := s.hydra.GetLoginRequest(admin.NewGetLoginRequestParamsWithContext(ctx).
		WithLoginChallenge(idt.Challenge))
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
		ID:          idt.UserID,
		HydraClient: lr.Payload.Client.ClientID,
	}, &auth.OTPChallenge{
		Code:  form.Code,
		Type:  idt.Sesssion.MFAType,
		Event: idt.Sesssion.MFAEvent,
	})
	if err != nil {
		switch e := auth.FromError(err); e.Code {
		case auth.NotFound:
			// Happen if the user is no longer found in database or account has been disabled.
			s.rejectLoginRedir(w, r, idt.Challenge, idt.UserID)
			return
		case auth.OTPInvalid:
			id := identity{
				Challenge:   idt.Challenge,
				UserID:      idt.UserID,
				Remember:    idt.Remember,
				Email:       idt.Email,
				CodeError:   err.Error(),
				RememberFor: time.Duration(ac.Remember()),
				Sesssion: sessionData{
					MFAEvent:  idt.Sesssion.MFAEvent,
					MFAType:   idt.Sesssion.MFAType,
					LoginCtx:  idt.Sesssion.LoginCtx,
					OpenIDCtx: idt.Sesssion.OpenIDCtx,
				},
			}
			s.serveOTPLogin(ctx, w, r, id, lr.Payload.Client)
			return
		default:
			log.WithError(err).Error("failed to lookup identity")
			http.Error(w, genericErrMsg, http.StatusInternalServerError)
			return
		}
	}

	s.acceptLogin(logging.WithLogger(ctx, log), w, r, identity{
		Challenge:   idt.Challenge,
		UserID:      ident.UserID,
		Remember:    idt.Remember,
		RememberFor: time.Duration(ac.Remember()),
		Sesssion: sessionData{
			LoginCtx:  ident.SessionContext(),
			OpenIDCtx: ident.LoginContext(),
		},
	})
}
