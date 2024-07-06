package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"

	"brank.as/rbac/idp/auth"
	"brank.as/rbac/serviceutil/logging"
)

func (s *server) getConsent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "getconsent")

	chg := r.URL.Query().Get("consent_challenge")
	csnt, err := s.hydra.GetConsentRequest(admin.NewGetConsentRequestParamsWithContext(ctx).
		WithConsentChallenge(chg))
	if err != nil {
		logging.WithError(err, log).Error("fetch consent request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	g := auth.Grant{
		Challenge: chg,
		UserID:    csnt.Payload.Subject,
		ClientID:  csnt.Payload.Client.ClientID,
		OwnerID:   csnt.Payload.Client.Owner,
		Granted:   csnt.Payload.RequestedScope,
		Remember:  true,
	}
	ct := s.authCl(w, r, ctx, csnt.Payload.Client).Consent()
	if ct == nil {
		log.Debug("no consent grantor")
		s.acceptConsent(w, r, nil, g, parseSession(csnt.Payload.Context))
		return
	}
	if len(csnt.Payload.RequestedScope) == 0 {
		log.Debug("no consent scopes requested")
		s.acceptConsent(w, r, ct, g, parseSession(csnt.Payload.Context))
		return
	}

	gr, err := ct.ServeGrant(ctx, g)
	if err != nil {
		logging.WithError(err, log).Error("serve grant")
		s.rejectConsent(w, r, auth.Grant{Challenge: chg})
		return
	}
	if gr.Skip {
		s.acceptConsent(w, r, ct, g, parseSession(csnt.Payload.Context))
		return
	}

	tmpl, cl := s.consentTpl, csnt.Payload.Client
	if cfg, err := auth.ParseClientConfig(cl); err == nil && cfg.ConsentTmpl != "" {
		tmpl = s.lookup(ctx, cfg.ConsentTmpl, s.consentTpl)
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, consentParams{
		ProjectName: s.projectName,
		Challenge:   chg,
		CSRFField:   csrf.TemplateField(r),
		IsProdEnv:   s.environment == "production",
		Scopes:      *gr,
		urls:        s.u,
	}); err != nil {
		log.WithError(err).Error("failed to render consent form")
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
}

type consentParams struct {
	ProjectName string
	Challenge   string
	CSRFField   template.HTML
	IsProdEnv   bool
	Scopes      auth.GrantDetail
	urls
}

type consentForm struct {
	Challenge string
	Cancel    string
	Scopes    []auth.Scope
	Groups    []auth.Group
}

func (s *server) postConsent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "postConsent")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var form consentForm
	if err := s.dcd.Decode(&form, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.WithField("form", form).WithField("postform", r.PostForm).Debug("post handler")
	if form.Challenge == "" {
		log.Error("missing challenge")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}

	csnt, err := s.hydra.GetConsentRequest(admin.NewGetConsentRequestParamsWithContext(ctx).
		WithConsentChallenge(form.Challenge))
	if err != nil {
		logging.WithError(err, log).Error("fetch consent request")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}
	ct := s.authCl(w, r, ctx, csnt.Payload.Client).Consent()
	if ct == nil {
		log.Debug("no consent grantor")
		g := auth.Grant{Challenge: form.Challenge, Granted: csnt.Payload.RequestedScope, Remember: true}
		s.acceptConsent(w, r, ct, g, parseSession(csnt.Payload.Context))
		return
	}
	gr := &auth.Grant{
		Challenge: form.Challenge,
		UserID:    csnt.Payload.Subject,
		ClientID:  csnt.Payload.Client.ClientID,
		OwnerID:   csnt.Payload.Client.Owner,
		Granted:   csnt.Payload.RequestedScope,
	}
	if form.Cancel != "" {
		s.rejectConsent(w, r, *gr)
		return
	}
	sess := parseSession(csnt.Payload.Context)

	// TODO: parse scopes from form to enable partial grants.
	gr, err = ct.Grant(ctx, *gr)
	if err != nil {
		logging.WithError(err, log).Error("grant consent")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}
	s.acceptConsent(w, r, ct, *gr, sess)
}

func (s *server) acceptConsent(w http.ResponseWriter, r *http.Request,
	ct auth.ConsentGrantor, gr auth.Grant, sess sessionData,
) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "idp.acceptconsent")

	if ct != nil {
		g, err := ct.Grant(ctx, gr)
		if err != nil {
			logging.WithError(err, log).Error("record consent grant")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if g.ID != "" {
			// record grant id
			sess.LoginCtx["grant_id"] = g.ID
		}
		gr.Granted = g.Granted
		gr.Remember = g.Remember
	}

	ret, err := s.hydra.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
		WithConsentChallenge(gr.Challenge).
		WithBody(&models.AcceptConsentRequest{
			GrantScope: gr.Granted,
			Remember:   gr.Remember,
			Session: &models.ConsentRequestSession{
				IDToken:     sess.OpenIDCtx,
				AccessToken: sess.LoginCtx,
			},
		}))
	if err != nil {
		logging.WithError(err, log).Error("hydra accept")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.WithField("subject", gr.UserID).Trace("granted consent")
	http.Redirect(w, r, *ret.Payload.RedirectTo, http.StatusFound)
}

func (s *server) rejectConsent(w http.ResponseWriter, r *http.Request, gr auth.Grant) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "rejectConsent")
	ret, err := s.hydra.RejectConsentRequest(admin.NewRejectConsentRequestParamsWithContext(ctx).
		WithConsentChallenge(gr.Challenge).WithBody(&models.RejectRequest{
		Error:      "consent not granted",
		StatusCode: http.StatusUnauthorized,
	}))
	if err != nil {
		logging.WithError(err, log).Error("reject consent")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}
	log.WithField("subject", gr.UserID).Trace("rejected consent")
	http.Redirect(w, r, *ret.Payload.RedirectTo, http.StatusTemporaryRedirect)
}

func parseSession(pl interface{}) sessionData {
	b, _ := json.Marshal(pl)
	s := sessionData{}
	json.Unmarshal(b, &s)
	return s
}
