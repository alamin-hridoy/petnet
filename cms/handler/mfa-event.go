package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"brank.as/petnet/cms/storage"
	epb "brank.as/petnet/gunk/dsa/v1/email"
	upb "brank.as/petnet/gunk/dsa/v1/user"
	"brank.as/petnet/gunk/dsa/v2/partnerlist"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	sVcpb "brank.as/petnet/gunk/dsa/v2/service"
	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	pps "brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	rbipb "brank.as/rbac/gunk/v1/invite"
	mpb "brank.as/rbac/gunk/v1/mfa"
	prpb "brank.as/rbac/gunk/v1/permissions"
	rbupb "brank.as/rbac/gunk/v1/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OTPForm struct {
	Code string
}

type mfaEvent struct {
	eventID  string
	resource string
	action   tpb.ActionType
	data     []byte
}

func (s *Server) initMFAEvent(w http.ResponseWriter, r *http.Request, me *mfaEvent) error {
	log := logging.FromContext(r.Context())

	ctx := context.Background()
	if me.eventID == "" {
		ev, err := s.rbac.InitiateMFA(ctx, &mpb.InitiateMFARequest{
			UserID:      mw.GetUserID(r.Context()),
			Type:        mpb.MFA_EMAIL,
			Description: "status placeholder replace later",
		})
		if err != nil {
			st, ok := status.FromError(err)
			if !ok || st.Code() != codes.NotFound {
				logging.WithError(err, log).Error("initiating mfa")
				return err
			}
			return storage.MFANotFound
		}
		me.eventID = ev.GetEventID()
	}
	if _, err := s.pf.CreateEventData(ctx, &tpb.CreateEventDataRequest{
		EventData: &tpb.EventData{
			EventID:  me.eventID,
			Resource: me.resource,
			Action:   me.action,
			Data:     string(me.data),
		},
	}); err != nil {
		logging.WithError(err, log).Error("creating event data")
		return err
	}

	sess, err := s.sess.Get(r, userCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		return err
	}
	sess.Values[mfaEventID] = me.eventID
	if err := s.sess.Save(r, w, sess); err != nil {
		log.WithError(err).Error("saving session")
		return err
	}
	return nil
}

// todo make error handling return error to be shown on screen
func (s *Server) postConfirmMFAEvent(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form OTPForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Code, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	sess, err := s.sess.Get(r, userCookieName)
	if err != nil {
		logging.WithError(err, log).Error("fetching user cookie")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	eid := sess.Values[mfaEventID].(string)
	if eid == "" {
		logging.WithError(err, log).Error("event id is empty")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	resp, err := s.pf.GetEventData(ctx, &tpb.GetEventDataRequest{
		EventID: eid,
	})
	if err != nil {
		logging.WithError(err, log).Error("get event data")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var mfaFailed bool
	res := storage.ResourceType(resp.GetEventData().GetResource())
	if res != storage.ChangePassword {
		if _, err := s.rbac.ValidateMFA(ctx, &mpb.ValidateMFARequest{
			UserID:  mw.GetUserID(ctx),
			Type:    mpb.MFA_EMAIL,
			Token:   form.Code,
			EventID: eid,
		}); err != nil {
			logging.WithError(err, log).Error("validating mfa event")
			mfaFailed = true
		}
	}

	act := resp.GetEventData().GetAction()
	switch res {
	case storage.Status:
		switch act {
		case tpb.ActionType_Update:
			pf := &ppb.UpsertProfileRequest{}
			if err := json.Unmarshal([]byte(resp.GetEventData().GetData()), pf); err != nil {
				logging.WithError(err, log).Error("unmarshal event data")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			if mfaFailed {
				http.Redirect(w, r, "/dashboard/dsa-applicant-list/"+pf.Profile.OrgID+"?otp_failed", http.StatusSeeOther)
				return
			}
			if _, err := s.pf.UpsertProfile(ctx, pf); err != nil {
				logging.WithError(err, log).Error("updating status")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			ui, err := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: pf.Profile.UserID})
			if err != nil {
				logging.WithError(err, log).Info("getting user")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			if ui == nil || ui.User == nil {
				logging.WithError(err, log).Error("getting user")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			isEmailSend := "1"
			if _, err := s.pf.SendOnboardingReminder(ctx, &epb.SendOnboardingReminderRequest{
				Email:  ui.User.Email,
				OrgID:  pf.Profile.OrgID,
				UserID: pf.Profile.UserID,
			}); err != nil {
				isEmailSend = "0"
				logging.WithError(err, log).Error("sending reminder")
			}
			http.Redirect(w, r, fmt.Sprintf("/dashboard/dsa-applicant-list/%s?emailsend=%s", pf.Profile.OrgID, isEmailSend), http.StatusSeeOther)
		}
	case storage.EnableMFA:
		switch act {
		case tpb.ActionType_Update:
			if mfaFailed {
				http.Redirect(w, r, "/dashboard/mfa-config"+"?show_otp=true&otp_failed", http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/dashboard/mfa-config", http.StatusSeeOther)
		}
	case storage.RolePermission:
		switch act {
		case tpb.ActionType_Update:
			urp := &UpdateRolePermissions{}
			if err := json.Unmarshal([]byte(resp.GetEventData().GetData()), urp); err != nil {
				logging.WithError(err, log).Error("unmarshal event data")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			if mfaFailed {
				http.Redirect(w, r, "/dashboard/manage-role/edit/"+urp.ID+"?otp_failed", http.StatusSeeOther)
				return
			}
			if err := s.updateRolePermissions(ctx, log, *urp); err != nil {
				log.Error("updating role permission")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, manageRoleListPath, http.StatusSeeOther)
		}
	case storage.RoleDelete:
		switch act {
		case tpb.ActionType_Delete:
			urp := &prpb.DeleteRoleRequest{}
			if err := json.Unmarshal([]byte(resp.GetEventData().GetData()), urp); err != nil {
				logging.WithError(err, log).Error("unmarshal event data")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			if mfaFailed {
				http.Redirect(w, r, manageRoleListPath+"?otp_failed", http.StatusSeeOther)
				return
			}
			if _, err := s.rbac.DeleteRole(ctx, urp); err != nil {
				log.Error("deleting role error")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, manageRoleListPath, http.StatusSeeOther)
		}
	case storage.UserDelete:
		switch act {
		case tpb.ActionType_Delete:
			du := &rbupb.DisableUserRequest{}
			if err := json.Unmarshal([]byte(resp.GetEventData().GetData()), du); err != nil {
				logging.WithError(err, log).Error("unmarshal event data")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			if mfaFailed {
				http.Redirect(w, r, manageUserListPath+"?otp_failed", http.StatusSeeOther)
				return
			}
			ul, er := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{
				ID: du.UserID,
			})
			if er != nil {
				logging.WithError(err, log).Error("getting user failed.")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
			}
			if ul.User.GetInviteStatus() == pps.InviteSent {
				if _, err := s.rbac.CancelInvite(ctx, &rbipb.CancelInviteRequest{
					ID: du.UserID,
				}); err != nil {
					log.WithError(err).Error(err.Error())
					http.Redirect(w, r, errorPath, http.StatusSeeOther)
					return
				}
			}
			if _, err := s.pf.DeleteUserProfile(ctx, &upb.DeleteUserProfileRequest{
				UserID: du.UserID,
			}); err != nil {
				log.WithError(err).Error(err.Error())
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
		}
	case storage.InviteUser:
		switch act {
		case tpb.ActionType_Create:
			iur := &rbipb.InviteUserRequest{}
			if err := json.Unmarshal([]byte(resp.GetEventData().GetData()), iur); err != nil {
				logging.WithError(err, log).Error("unmarshal event data")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			if mfaFailed {
				http.Redirect(w, r, manageUserListPath+"?otp_failed", http.StatusSeeOther)
				return
			}
			res, err := s.rbac.InviteUser(ctx, iur)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok || st.Code() == codes.AlreadyExists {
					http.Redirect(w, r, manageUserListPath+"?errorMsg=AlreadyExists", http.StatusSeeOther)
					return
				}
			}
			if _, err := s.rbac.AddUser(ctx, &prpb.AddUserRequest{
				RoleID: iur.Role,
				UserID: res.GetID(),
			}); err != nil {
				logging.WithError(err, log).Error("adding user to role")
				return
			}
			http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
		}
	case storage.CreateRole:
		switch act {
		case tpb.ActionType_Create:
			cr := &prpb.CreateRoleRequest{}
			if err := json.Unmarshal([]byte(resp.GetEventData().GetData()), cr); err != nil {
				logging.WithError(err, log).Error("unmarshal event data")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			if mfaFailed {
				http.Redirect(w, r, manageRoleListPath+"?otp_failed", http.StatusSeeOther)
				return
			}
			if _, err := s.rbac.CreateRole(ctx, &prpb.CreateRoleRequest{
				Name:        cr.Name,
				Description: cr.Description,
			}); err != nil {
				logging.WithError(err, log).Error("Create Role")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, manageRoleListPath, http.StatusSeeOther)
		}
	case storage.ChangePassword:
		switch act {
		case tpb.ActionType_Update:
			cr := &rbupb.ConfirmUpdateRequest{}
			if err := json.Unmarshal([]byte(resp.GetEventData().GetData()), cr); err != nil {
				logging.WithError(err, log).Error("unmarshal event data")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			cr.MFAToken = form.Code
			if _, err := s.rbacUserAuth.ConfirmUpdate(ctx, cr); err != nil {
				logging.WithError(err, log).Error("Update Password")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, edPassPath+"?msgType=success", http.StatusSeeOther)
		}
	case storage.InviteProvider:
		switch act {
		case tpb.ActionType_Create:
			iur := &InviteProviderForm{}
			if err := json.Unmarshal([]byte(resp.GetEventData().GetData()), iur); err != nil {
				logging.WithError(err, log).Error("unmarshal event data")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			if mfaFailed {
				http.Redirect(w, r, providersListPath+"?otp_failed", http.StatusSeeOther)
				return
			}

			orgName := iur.FirstName + " " + iur.LastName

			rls, err := s.GetRoleByName(ctx, "provider")
			if err != nil {
				logging.WithError(err, log).Error("failed to get provider role.")
				http.Redirect(w, r, providersListPath+"?errorMsg=RoleNotFound", http.StatusSeeOther)
				return
			}

			iUser, err := s.rbac.InviteUser(ctx, &rbipb.InviteUserRequest{
				OrgName:         orgName,
				FirstName:       iur.FirstName,
				LastName:        iur.LastName,
				Email:           iur.Email,
				Role:            rls.ID,
				CustomEmailData: map[string]string{"firstName": iur.FirstName},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if !ok || st.Code() == codes.AlreadyExists {
					logging.WithError(err, log).Error("failed to send invite user for conflicting.")
					http.Redirect(w, r, providersListPath+"?errorMsg=AlreadyExists", http.StatusSeeOther)
					return
				}
				logging.WithError(err, log).Error("failed to send invite user.")
				http.Redirect(w, r, providersListPath+"?errorMsg=InviteUser", http.StatusSeeOther)
				return
			}

			if _, err := s.rbac.AddUser(ctx, &prpb.AddUserRequest{
				RoleID: rls.ID,
				UserID: iUser.GetID(),
			}); err != nil {
				logging.WithError(err, log).Error("adding user to role")
				http.Redirect(w, r, providersListPath+"?errorMsg=AddUser", http.StatusSeeOther)
				return
			}

			user, err := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{
				ID: iUser.GetID(),
			})
			if err != nil {
				logging.WithError(err, log).Error("failed to get user")
				http.Redirect(w, r, providersListPath+"?errorMsg=InviteUser", http.StatusSeeOther)
				return
			}

			var orgId string

			if user != nil && user.GetUser() != nil {
				orgId = user.GetUser().GetOrgID()
			}

			_, err = s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
				Profile: &ppb.OrgProfile{
					UserID:           iUser.GetID(),
					OrgID:            orgId,
					IsProvider:       true,
					OrgType:          ppb.OrgType(ppb.OrgType_DSA),
					Status:           ppb.Status_Pending,
					Partner:          iur.Provider,
					TransactionTypes: strings.Join(iur.TransactionTypes, ","),
					BusinessInfo: &ppb.BusinessInfo{
						CompanyName: orgName,
					},
				},
			})
			if err != nil {
				logging.WithError(err, log).Error("failed to upsert profile.")
				http.Redirect(w, r, providersListPath+"?errorMsg=CreateProfile", http.StatusSeeOther)
				return
			}

			_, err = s.pf.CreateUserProfile(ctx, &upb.CreateUserProfileRequest{
				Profile: &upb.Profile{
					UserID: iUser.ID,
					OrgID:  orgId,
					Email:  iur.Email,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if !ok || st.Code() == codes.AlreadyExists {
					logging.WithError(err, log).Error("failed to create user profile for conflicting.")
					http.Redirect(w, r, providersListPath+"?errorMsg=AlreadyExists", http.StatusSeeOther)
					return
				}
				logging.WithError(err, log).Error("failed to create user profile.")
				http.Redirect(w, r, providersListPath+"?errorMsg=CreateUserProfile", http.StatusSeeOther)
				return
			}

			gPs, err := s.pf.GetPartnerByStype(ctx, &partnerlist.GetPartnerByStypeRequest{
				Stype: iur.Provider,
			})
			if err != nil {
				logging.WithError(err, log).Error("failed to get provider service type.")
				http.Redirect(w, r, providersListPath+"?errorMsg=GetPartnerByStype", http.StatusSeeOther)
				return
			}
			if gPs == nil || gPs.PartnerList == nil {
				logging.WithError(err, log).Error("failed to get provider service type.")
				http.Redirect(w, r, providersListPath+"?errorMsg=GetPartnerByStype", http.StatusSeeOther)
				return
			}

			_, err = s.pf.AddServiceRequest(ctx, &sVcpb.AddServiceRequestRequest{
				OrgID:       orgId,
				Type:        sVcpb.ServiceType(sVcpb.ServiceType_value[gPs.PartnerList.GetServiceName()]),
				Partners:    []string{iur.Provider},
				AllPartners: false,
			})
			if err != nil {
				logging.WithError(err, log).Error("failed to create service request.")
				http.Redirect(w, r, providersListPath+"?errorMsg=ServiceReq", http.StatusSeeOther)
				return
			}

			_, err = s.pf.AcceptServiceRequest(ctx, &sVcpb.ServiceStatusRequestRequest{
				OrgID:     orgId,
				Partner:   iur.Provider,
				SvcName:   gPs.GetPartnerList().GetServiceName(),
				UpdatedBy: mw.GetUserID(ctx),
			})

			if err != nil {
				logging.WithError(err, log).Error("failed to accept service request.")
				http.Redirect(w, r, providersListPath+"?errorMsg=ServiceReq", http.StatusSeeOther)
				return
			}

			http.Redirect(w, r, providersListPath, http.StatusSeeOther)
			return
		}
	}
}
