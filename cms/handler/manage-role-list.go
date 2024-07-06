package handler

import (
	"context"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"brank.as/petnet/cms/paginator"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	ppb "brank.as/rbac/gunk/v1/permissions"
	rbupb "brank.as/rbac/gunk/v1/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	RoleDetails struct {
		ID          string
		OrgID       string
		Name        string
		Description string
		CreatedBy   string
		UpdatedBy   string
		Updated     time.Time
	}

	RoleListTempData struct {
		CSRFField        template.HTML
		Roles            []RoleDetails
		ErrMsg           string
		Errors           map[string]error
		PresetPermission map[string]map[string]bool
		ServiceRequest   bool
		LoginUserInfo    *User
		PaginationData   paginator.Paginator
	}
)

func (s *Server) getManageRoleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	template := s.templates.Lookup("role-list.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	queryParams := r.URL.Query()
	sb, err := url.PathUnescape(queryParams.Get("sort"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	pageNumber, err := url.PathUnescape(queryParams.Get("page"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	var offset int32 = 0
	convertedPageNumber, _ := strconv.Atoi(pageNumber)
	if convertedPageNumber <= 0 {
		convertedPageNumber = 1
	} else {
		offset = limitPerPage*int32(convertedPageNumber) - limitPerPage
	}

	errorMsg, _ := url.PathUnescape(queryParams.Get("errorMsg"))
	errMsg := ""
	if errorMsg == "true" {
		errMsg = "Role is assigned to user(s) and can't be deleted"
	}

	sbv := ppb.SortBy_ASC
	if sb == "desc" {
		sbv = ppb.SortBy_DESC
	}

	oid := mw.GetOrgID(ctx)
	res, err := s.rbac.ListRole(ctx, &ppb.ListRoleRequest{
		OrgID:  oid,
		SortBy: sbv,
		Offset: offset,
		Limit:  limitPerPage,
	})
	if err != nil {
		log.Error("unable to connect api")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	userLists := []string{}
	for _, role := range res.GetRoles() {
		if role.Name == "Owner" {
			continue
		}
		if role.CreateUID != "" {
			userLists = append(userLists, role.CreateUID)
		}
		if role.UpdatedUID != "" {
			userLists = append(userLists, role.UpdatedUID)
		}
	}
	userLists = uniqueSlice(userLists)
	userInfoLists := map[string]*rbupb.User{}
	usrs, err := s.rbac.ListUsers(ctx, &rbupb.ListUsersRequest{
		ID:    userLists,
		OrgID: oid,
	})
	if err == nil {
		userInfoLists = usrs.GetUser()
	}
	formErr := validation.Errors{}
	var roleList []RoleDetails
	var SkipCount int32
	for _, role := range res.GetRoles() {
		if role.Name == "Owner" {
			SkipCount = SkipCount + 1
			continue
		}
		roleList = append(roleList, RoleDetails{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedBy:   userInfoLists[role.CreateUID].GetFirstName() + " " + userInfoLists[role.CreateUID].GetLastName(),
			UpdatedBy:   userInfoLists[role.UpdatedUID].GetFirstName() + " " + userInfoLists[role.UpdatedUID].GetLastName(),
			Updated:     role.GetUpdated().AsTime(),
		})
	}

	etd := s.getEnforceTemplateData(ctx)
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	tempData := RoleListTempData{
		CSRFField:        csrf.TemplateField(r),
		Roles:            roleList,
		Errors:           formErr,
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
		LoginUserInfo:    &usrInfo.UserInfo,
		ErrMsg:           errMsg,
	}
	if len(tempData.Roles) > 0 {
		tempData.PaginationData = paginator.NewPaginator(int32(convertedPageNumber), limitPerPage, (res.Total - SkipCount), r)
	}
	tempData.LoginUserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, tempData); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) GetRoleByName(ctx context.Context, name string) (*ppb.Role, error) {
	log := logging.FromContext(ctx)
	rls, err := s.rbac.ListRole(ctx, &ppb.ListRoleRequest{
		Name: name,
	})
	if err != nil {
		logging.WithError(err, log).Error("getting roles failed")
		return nil, status.Error(codes.NotFound, "role not found")
	}
	if rls == nil {
		return nil, status.Error(codes.NotFound, "role not found")
	}

	for _, v := range rls.GetRoles() {
		if v.Name == name {
			return v, nil
		}
	}

	return nil, status.Error(codes.NotFound, "role not found")
}
