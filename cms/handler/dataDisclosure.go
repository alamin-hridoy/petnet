package handler

import (
	"encoding/json"
	"net/http"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	rbupb "brank.as/rbac/gunk/v1/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kenshaw/goji"
)

type (
	DisclosureData struct {
		Field string
	}
)

func (s *Server) postSpecificFieldData(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing file ID query param")
		return
	}

	var form DisclosureData
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Field, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["Field"] != nil {
				log.Error("Field is Required")
			}
		}
		return
	}

	pf, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
		return
	}

	u, err := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: pf.Profile.UserID})
	if err != nil {
		logging.WithError(err, log).Info("getting user")
		return
	}
	resField := make(map[string]string)
	if form.Field == "email" {
		resField["email"] = u.User.Email
	} else if form.Field == "phn" {
		resField["phn"] = pf.GetProfile().GetBusinessInfo().PhoneNumber
	} else if form.Field == "comEmail" {
		resField["comEmail"] = pf.GetProfile().GetBusinessInfo().CompanyEmail
	} else if form.Field == "bankAcc" {
		resField["bankAcc"] = pf.GetProfile().GetAccountInfo().GetBankAccountNumber()
	} else if form.Field == "bankHold" {
		resField["bankHold"] = pf.GetProfile().GetAccountInfo().GetBankAccountHolder()
	}

	jsn, _ := json.Marshal(resField)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsn)
}
