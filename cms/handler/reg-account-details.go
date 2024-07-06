package handler

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/csrf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	upb "brank.as/petnet/gunk/dsa/v1/user"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	mpb "brank.as/petnet/gunk/v1/mfa"
	"brank.as/petnet/serviceutil/logging"
)

type SignupForm struct {
	CSRFField            template.HTML
	Email                string
	InviteCode           string
	Password             string
	PasswordConfirmation string
	Errors               map[string]error
}

func (s *Server) getAccountDetails(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	form := SignupForm{CSRFField: csrf.TemplateField(r)}
	form.InviteCode = ""
	form.Email = ""
	template := s.templates.Lookup("account-details.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	queryParams := r.URL.Query()
	inviteCode, err := url.PathUnescape(queryParams.Get("invite_code"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	if inviteCode != "" {
		form.InviteCode = inviteCode
		ru, err := s.pf.RetrieveInvite(ctx, &upb.RetrieveInviteRequest{
			Code: inviteCode,
		})
		if err != nil {
			log.Error("Unable to get data by invite code")
		}
		if ru.Email != "" {
			form.Email = ru.Email
		}
	}

	if err := template.Execute(w, form); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postAccountDetails(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form SignupForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Email, is.Email, validation.Required),
		validation.Field(&form.Password, validation.Required, validation.Length(8, 64)),
		validation.Field(&form.PasswordConfirmation, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["Email"] != nil {
				formErr["Email"] = err["Email"]
			}
			if err["Password"] != nil {
				formErr["Password"] = err["Password"]
			}
			if err["PasswordConfirmation"] != nil {
				formErr["PasswordConfirmation"] = err["PasswordConfirmation"]
			}
		}
		logging.WithError(err, log).Error("invalid request")
	}

	if formErr["Password"] == nil && formErr["PasswordConfirmation"] == nil && !isPasswordValid(form.Password) {
		formErr["Password"] = errors.New("The password entered does not meet policy.")
	}

	if formErr["Password"] == nil && formErr["PasswordConfirmation"] == nil && form.Password != form.PasswordConfirmation {
		formErr["Password"] = errors.New("Password doesn't match")
		formErr["PasswordConfirmation"] = errors.New("Password doesn't match")
	}

	form.Errors = formErr
	form.CSRFField = csrf.TemplateField(r)
	if len(formErr) > 0 {
		template := s.templates.Lookup("account-details.html")
		if template == nil {
			log.Error("unable to load template")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if err := template.Execute(w, form); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		return
	}

	if form.Password != form.PasswordConfirmation {
		log.Error("password and password confirmation mismatch")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	fName := "placeholder"
	lName := "placeholder"
	if form.InviteCode != "" {
		ru, err := s.pf.RetrieveInvite(ctx, &upb.RetrieveInviteRequest{
			Code: form.InviteCode,
		})
		if err != nil {
			log.Error("Unable to get data by invite code")
		}
		if ru.GetFirstName() != "" {
			fName = ru.GetFirstName()
		}
		if ru.GetLastName() != "" {
			lName = ru.GetLastName()
		}
	}

	res, err := s.pf.Signup(ctx, &upb.SignupRequest{
		Username:   form.Email,
		FirstName:  fName,
		LastName:   lName,
		Email:      form.Email,
		Password:   form.Password,
		InviteCode: form.InviteCode,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.AlreadyExists {
			logging.WithError(err, log).Error("unable to signup")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	if err != nil {
		formErr["Email"] = errors.New("Email already exists")
		form.Errors = formErr
		form.CSRFField = csrf.TemplateField(r)
		if len(formErr) > 0 {
			template := s.templates.Lookup("account-details.html")
			if template == nil {
				log.Error("unable to load template")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			if err := template.Execute(w, form); err != nil {
				log.Infof("error with template execution: %+v", err)
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			return
		}
	}
	// start new one

	getPf, err := s.pf.GetUserProfileByEmail(ctx, &upb.GetUserProfileByEmailRequest{
		Email: form.Email,
	})
	if err != nil {
		log.Error("Unable to get user profile by email")
	}

	if getPf.GetProfile() != nil {
		comPf, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{
			OrgID: getPf.GetProfile().GetOrgID(),
		})
		if err != nil {
			log.Error("Unable to get company profile")
		}

		if comPf.GetProfile().GetStatus() == ppb.Status_Pending {
			s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
				Profile: &ppb.OrgProfile{
					OrgID:  getPf.GetProfile().GetOrgID(),
					Status: ppb.Status_Accepted,
				},
			})
		}
	}

	if getPf == nil {
		_, err = s.pf.GetProfile(ctx, &ppb.GetProfileRequest{
			OrgID: res.OrgID,
		})

		if status.Code(err) == codes.NotFound {
			_, err = s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
				Profile: &ppb.OrgProfile{
					UserID:  res.UserID,
					OrgID:   res.OrgID,
					OrgType: ppb.OrgType(ppb.OrgType_DSA),
				},
			})
			if err != nil {
				log.WithError(err).Error("creating org profile")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
		}
		_, err = s.pf.CreateUserProfile(ctx, &upb.CreateUserProfileRequest{
			Profile: &upb.Profile{
				UserID: res.UserID,
				OrgID:  res.OrgID,
				Email:  form.Email,
			},
		})
		if err != nil && status.Code(err) != codes.AlreadyExists {
			log.WithError(err).Error("creating user profile")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	sess.Values[sessionUserID] = res.UserID
	sess.Values[sessionEmail] = form.Email
	if err := s.sess.Save(r, w, sess); err != nil {
		log.WithError(err).Error("saving session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if form.InviteCode != "" {
		if !s.disableActionMFA {
			if _, err := s.pf.EnableMFA(ctx, &mpb.EnableMFARequest{
				UserID: res.UserID,
				Type:   mpb.MFA_EMAIL,
				Source: form.Email,
			}); err != nil {
				logging.WithError(err, log).Error("enable mfa")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
		}
		http.Redirect(w, r, loginPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, verifyAccountPath, http.StatusSeeOther)
}

func isPasswordValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
