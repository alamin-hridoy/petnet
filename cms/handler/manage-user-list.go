package handler

import (
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"brank.as/petnet/cms/paginator"
	ppu "brank.as/petnet/gunk/dsa/v1/user"
	pfl "brank.as/petnet/gunk/dsa/v2/profile"
	strg "brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	rbipb "brank.as/rbac/gunk/v1/invite"
	ppb "brank.as/rbac/gunk/v1/permissions"
	rbupb "brank.as/rbac/gunk/v1/user"
	"github.com/gorilla/csrf"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

type (
	UsersDetails struct {
		ID           string
		FirstName    string
		LastName     string
		InviteStatus string
		UserRole     string
		Resend       bool
		IsNotDeleted bool
		Created      time.Time
	}
	UsersListTempData struct {
		CSRFField        template.HTML
		Users            []UsersDetails
		RoleLists        []RoleDetails
		SearchTerm       string
		InviteUsrErr     string
		InviteStatusList []string
		PresetPermission map[string]map[string]bool
		ServiceRequest   bool
		LoginUserInfo    *User
		PaginationData   paginator.Paginator
	}

	VerifyEmailConfirm struct {
		UserID string
	}
)

func (s *Server) getManageUserList(w http.ResponseWriter, r *http.Request) {
	s.ManageUserListForm(w, r, "")
}

func (s *Server) ManageUserListForm(w http.ResponseWriter, r *http.Request, inviteErr string) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	template := s.templates.Lookup("manage-user.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	queryParams := r.URL.Query()
	pageNumber, err := url.PathUnescape(queryParams.Get("page"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	searchTerm, err := url.PathUnescape(queryParams.Get("search-term"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	newReq := &rbupb.ListUsersRequest{
		OrgID: mw.GetOrgID(ctx),
		Name:  searchTerm,
	}

	var offset int32 = 0
	convertedPageNumber, _ := strconv.Atoi(pageNumber)
	if convertedPageNumber <= 0 {
		convertedPageNumber = 1
	} else {
		offset = limitPerPage*int32(convertedPageNumber) - limitPerPage
	}
	newReq.Offset = offset
	newReq.Limit = limitPerPage

	sb, err := url.PathUnescape(queryParams.Get("sort"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	newReq.SortBy = rbupb.SortBy_ASC
	if sb == "desc" {
		newReq.SortBy = rbupb.SortBy_DESC
	}

	sbc, err := url.PathUnescape(queryParams.Get("sort_column"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	newReq.SortByColumn = rbupb.SortByColumn_UserName
	if sbc == "CreatedDate" {
		newReq.SortByColumn = rbupb.SortByColumn_CreatedDate
	}

	sts, err := url.PathUnescape(queryParams.Get("status"))
	if err != nil {
		log.Error("unable to decode url type param")
	}
	if sts != "" {
		sStrArr := strings.Split(sts, ",")
		sArr := []rbupb.Status{}
		for _, s := range sStrArr {
			switch s {
			case "Invite Sent":
				sArr = append(sArr, rbupb.Status_InviteSent)
			case "Revoked":
				sArr = append(sArr, rbupb.Status_Revoked)
			case "Approved":
				sArr = append(sArr, rbupb.Status_Approved)
			}
		}
		newReq.Status = sArr
	}

	isArr := []string{}
	for _, is := range rbupb.Status_name {
		if is == "InviteSent" || is == "Revoked" || is == "Approved" {
			if is == "InviteSent" {
				is = "Invite Sent"
			}
			isArr = append(isArr, is)
		}
	}
	oid := mw.GetOrgID(ctx)
	res, err := s.rbac.ListRole(ctx, &ppb.ListRoleRequest{
		OrgID: oid,
	})
	if err != nil {
		log.Error("listing roles")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	ul, err := s.rbac.ListUsers(ctx, newReq)
	if err != nil {
		log.Error("listing users")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	var userids []string
	for _, usr := range ul.Users {
		userids = append(userids, usr.GetID())
	}
	ur, err := s.rbac.ListUserRoles(ctx, &ppb.ListUserRolesRequest{
		UserID: userids,
	})
	if err != nil {
		log.Error("listing users roles")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	userRoles := make(map[string]string)
	for k, v := range ur.GetRoles() {
		roleNames := []string{}
		for _, rid := range v.UserRoles {
			for _, role := range res.GetRoles() {
				if rid == role.GetID() {
					roleNames = append(roleNames, role.GetName())
				}
			}
		}
		userRoles[k] = strings.Join(roleNames, ", ")
	}

	var userList []UsersDetails
	for _, user := range ul.Users {
		resend := false
		convertedStartDate := user.GetDeleted()
		opf, _ := s.pf.GetProfile(ctx, &pfl.GetProfileRequest{
			OrgID: user.GetOrgID(),
		})
		if opf.Profile.GetUserID() == user.GetID() {
			continue
		}
		pf, err := s.pf.GetUserProfile(ctx, &ppu.GetUserProfileRequest{
			UserID: user.GetID(),
		})
		if err != nil && user.GetInviteStatus() != strg.InviteSent {
			convertedStartDate = tspb.New(time.Now())
		}
		if err == nil {
			convertedStartDate = pf.GetProfile().GetDeleted()
		}
		if user.GetInviteStatus() == strg.Revoked {
			convertedStartDate = tspb.New(time.Now())
		}
		convertedEndDate := tspb.New(time.Time{})
		timeCheck := convertedStartDate.AsTime().Equal(convertedEndDate.AsTime())
		if user.GetInviteStatus() == strg.Revoked && !timeCheck {
			resend = true
			timeCheck = false
		}
		uRoleList := ""
		if val, ok := userRoles[user.GetID()]; ok {
			uRoleList = val
		}
		userList = append(userList, UsersDetails{
			ID:           user.GetID(),
			FirstName:    user.GetFirstName(),
			LastName:     user.GetLastName(),
			InviteStatus: user.GetInviteStatus(),
			Resend:       resend,
			UserRole:     uRoleList,
			IsNotDeleted: timeCheck,
			Created:      user.GetCreated().AsTime(),
		})
	}

	var roleLists []RoleDetails
	for _, role := range res.GetRoles() {
		if role.Name == "Owner" {
			continue
		}
		roleLists = append(roleLists, RoleDetails{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		})
	}
	InviteUsrErr := ""
	if inviteErr == "AlreadyExists" {
		InviteUsrErr = "Email address already exists"
	}

	errorMsg, _ := url.PathUnescape(queryParams.Get("errorMsg"))
	if errorMsg == "AlreadyExists" {
		InviteUsrErr = "Email address already exists"
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	etd := s.getEnforceTemplateData(ctx)
	tempData := UsersListTempData{
		CSRFField:        csrf.TemplateField(r),
		Users:            userList,
		RoleLists:        roleLists,
		SearchTerm:       searchTerm,
		InviteUsrErr:     InviteUsrErr,
		InviteStatusList: isArr,
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
		LoginUserInfo:    &usrInfo.UserInfo,
	}
	if len(tempData.Users) > 0 {
		tempData.PaginationData = paginator.NewPaginator(int32(convertedPageNumber), limitPerPage, ul.Total, r)
	}
	tempData.LoginUserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, tempData); err != nil {
		logging.WithError(err, log).Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postResendEmailConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	if err := r.ParseForm(); err != nil {
		log.Error("parsing form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var f VerifyEmailConfirm
	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		log.Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	gu, err := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: f.UserID})
	if err != nil {
		log.Error("unable to connect api to getUser")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	switch gu.User.InviteStatus {
	case "":
		log.Error("inviteStatus Empty is not supported for resend mail.")
		http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
	case "Approved":
		log.Error("inviteStatus Approved is not supported for resend mail.")
		http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
	case "Invite Sent":
		log.Error("InviteStatus Invite Sent is not supported for resend mail.")
		http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
	case "Expired":
		log.Error("inviteStatus Expired is not supported for resend mail.")
		http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
	case "Revoked":
		log.Error("inviteStatus Revoked is not supported for resend mail.")
		http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
	case "In Progress":
		log.Error("inviteStatus In Progress is not supported for resend mail.")
		http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
	default:
		if _, err := s.rbac.Resend(ctx, &rbipb.ResendRequest{
			ID: f.UserID,
		}); err != nil {
			log.Error("resending confirmation email")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
}
