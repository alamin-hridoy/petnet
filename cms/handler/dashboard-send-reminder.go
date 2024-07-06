package handler

import (
	"net/http"
	"net/url"

	epb "brank.as/petnet/gunk/dsa/v1/email"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type SendReminderForm struct {
	OrgID       string
	Email       string
	CurrentPage string
}

func (s *Server) postDashboardSendReminder(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var f SendReminderForm
	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := validation.ValidateStruct(&f,
		validation.Field(&f.OrgID, validation.Required, is.UUIDv4),
		validation.Field(&f.Email, validation.Required, is.Email),
		validation.Field(&f.CurrentPage, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	gp, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{
		OrgID: f.OrgID,
	})
	if err != nil || gp == nil || gp.Profile == nil {
		logging.WithError(err, log).Error("org profile info not found")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	isEmailSend := "1"
	if _, err := s.pf.SendOnboardingReminder(ctx, &epb.SendOnboardingReminderRequest{
		Email:  f.Email,
		OrgID:  f.OrgID,
		UserID: gp.Profile.UserID,
	}); err != nil {
		isEmailSend = "0"
		logging.WithError(err, log).Error("sending reminder")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if _, err := s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
		Profile: &ppb.OrgProfile{
			OrgID:        f.OrgID,
			ReminderSent: ppb.Boolean_True,
		},
	}); err != nil {
		logging.WithError(err, log).Error("sending reminder")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	u, _ := url.Parse(s.urls.Base)
	u.Path = dsaAppListPath
	q := u.Query()
	q.Add("page", f.CurrentPage)
	q.Add("emailsend", isEmailSend)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusSeeOther)
}
