package handler

import (
	"encoding/json"
	"net/http"

	"brank.as/petnet/cms/storage"
	rbipb "brank.as/rbac/gunk/v1/invite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/serviceutil/logging"
	ppb "brank.as/rbac/gunk/v1/permissions"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	InviteForm struct {
		Email     string
		Role      string
		FirstName string
		LastName  string
	}
)

func (s *Server) postInviteMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var f InviteForm
	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := validation.ValidateStruct(&f,
		validation.Field(&f.Email, validation.Required, is.Email),
		validation.Field(&f.Role, is.UUIDv4),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["Email"] != nil {
				logging.WithError(err, log).Error("Invalid Email or Required Email")
			}
			if err["Role"] != nil {
				logging.WithError(err, log).Error("Role is Invalid")
			}
		}
		logging.WithError(err, log).Error("invalid request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	UserInfo := s.loadUserInfo(r)
	iur := &rbipb.InviteUserRequest{
		Email:           f.Email,
		Role:            f.Role,
		FirstName:       f.FirstName,
		LastName:        f.LastName,
		OrgID:           UserInfo.OrgID,
		OrgName:         UserInfo.OrgName,
		CustomEmailData: map[string]string{"firstName": f.FirstName},
	}

	d, err := json.Marshal(iur)
	if err != nil {
		logging.WithError(err, log).Error("marshal request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	me := &mfaEvent{
		resource: string(storage.InviteUser),
		action:   tpb.ActionType_Create,
		data:     d,
	}
	if err := s.initMFAEvent(w, r, me); err != nil {
		if err != storage.MFANotFound {
			logging.WithError(err, log).Error("initializing mfa event")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}

		res, err := s.rbac.InviteUser(ctx, iur)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok || st.Code() == codes.AlreadyExists {
				s.ManageUserListForm(w, r, "AlreadyExists")
				http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
				return
			}
		}

		au := &ppb.AddUserRequest{
			RoleID: f.Role,
			UserID: res.GetID(),
		}
		if _, err := s.rbac.AddUser(ctx, au); err != nil {
			logging.WithError(err, log).Error("adding user to role")
			return
		}
		http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
	}
	http.Redirect(w, r, manageUserListPath+"?show_otp=true", http.StatusSeeOther)
}
