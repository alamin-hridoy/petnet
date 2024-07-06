package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/benbjohnson/hashfs"
	"github.com/bojanz/currency"
	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/h2non/filetype"
	"github.com/kenshaw/goji"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	cmsmw "brank.as/petnet/cms/mw"
	"brank.as/petnet/cms/storage"
	"brank.as/petnet/gunk/dsa/v1/user"
	pfpb "brank.as/petnet/gunk/dsa/v2/profile"
	session "brank.as/petnet/profile/services/rbsession"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	pm "brank.as/petnet/svcutil/permission"
	rpmpb "brank.as/rbac/gunk/v1/permissions"
	rbupb "brank.as/rbac/gunk/v1/user"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	userCookieName    = "petnet-user"
	sessionCookieName = "petnet-session"
	sessionOrgID      = "org-id"
	sessionUserID     = "user-id"
	sessionEmail      = "email"
	mfaEventID        = "event-id"
	userInfoDetails   = "user-info-detail"
	userInfo          = "userinfo"
	profileInfo       = "profileinfo"
)

const (
	sessionCookieState = "state"
	authCodeURL        = "somerandomstring"
	timeZone           = "Asia/Manila"
	timeLayout         = "Jan. 2, 2006 15:04:05 MST"
)

const (
	rootPath   = "/root"
	errorPath  = "/error"
	loginPath  = "/login"
	logoutPath = "/logout"

	// profile
	profilePath = "/profile"
	edPassPath  = "/change-password"

	// dsa onboarding
	verifyAccountPath   = "/register/verify-account"
	businessInfoPath    = "/register/businessinfo"
	financialInfoPath   = "/register/financialinfo"
	accountInfoPath     = "/register/accountinfo"
	drpInfoPath         = "/register/drpinfo"
	registerSuccessPath = "/register/success"
	skipOnboardingPath  = "/register/skip-onboarding"

	// dashboard
	dashboardCheckoutList     = "/dashboard/checkout-list/:id"
	dsaAppListPath            = "/dashboard/dsa-applicant-list"
	dsaAppListDetailPath      = "/dashboard/dsa-applicant-list/:id"
	dashboardChangeStatusPath = "/dashboard/change-status"
	dashboardSendReminderPath = "/dashboard/send-reminder"
	dashboardCheckDSACodePath = "/dashboard/profile-dsa-code"

	// signup
	preliminaryScreenPath = "/registration/preliminary-screen"
	regAccountDetailsPath = "/registration/account-details"

	// Configuration route
	locationConfigurationPath                 = "/dashboard/location/:id"
	locationUpdatePath                        = "/dashboard/locationUpdate"
	locationConfigurationUpdatePath           = "/dashboard/location/:id/update/:lid"
	locationConfigurationDeletePath           = "/dashboard/location/:id/delete/:lid"
	feeManagementPath                         = "/dashboard/fee-management/:id"
	feeManagementUpdatePath                   = "/dashboard/fee-management/:id/update/:fid"
	feeManagementDeletePath                   = "/dashboard/fee-management/:id/delete/:fid"
	dsaServicesPath                           = "/dsa-services"
	dsaPrtSelPath                             = "/dsa-select-partners"
	dsaPrtCiCoSelPath                         = "/dsa-select-provider"
	dsaPrtBillsPaymentSelPath                 = "/dsa-req-docs-bills-payment"
	dsaAdiBillsPaymentSelPath                 = "/dsa-additional-docs-bills-payment"
	dsaReqDocsPath                            = "/dsa-req-docs"
	dsaReqDocsCicoPath                        = "/dsa-req-docs-cico"
	dsaReqDocsMicroInsurance                  = "/dsa-req-docs-micro-insurance"
	dsaAddiDocsPath                           = "/dsa-additional-docs"
	serviceReqPath                            = "/dashboard/service-req"
	postAddRemarkSvcReq                       = "/dashboard/add-service-req"
	changeStatusSvcReq                        = "/dashboard/status-service-req/:status/:id/:partner/:file"
	changeAllStatusSvcReq                     = "/dashboard/status-service-req/:status/:id/:partner"
	changeCiCoAllStatusSvcReq                 = "/dashboard/status-service-req-cico/:status/:id/:partner"
	changeMIAllStatusSvcReq                   = "/dashboard/status-service-req-mi/:status/:id/:partner"
	changeBillsPaymentAllStatusSvcReq         = "/dashboard/status-service-req-bills-payment/:status/:id/:partner"
	changeAllStatusSvcReqCicoPath             = "/dashboard/status-service-req-cico/:status/:id/:partner/:file"
	changeAllStatusSvcReqMIPath               = "/dashboard/status-service-req-mi/:status/:id/:partner/:file"
	changeAllStatusSvcReqBillsPaymentPath     = "/dashboard/status-service-req-bills-payment/:status/:id/:partner/:file"
	changeAjaxAllStatusSvcReqCicoPath         = "/dashboard/ajax-status-service-req-cico/:status/:id/:partner/:file"
	changeAjaxAllStatusSvcReqMIPath           = "/dashboard/ajax-status-service-req-mi/:status/:id/:partner/:file"
	changeAjaxAllStatusSvcReqBillsPaymentPath = "/dashboard/ajax-status-service-req-bills-payment/:status/:id/:partner/:file"
	servicesDetailsPath                       = "/dashboard/services-details/:id"
	servicesDetailsCicoPath                   = "/dashboard/services-details-cico/:id"
	servicesDetailsBillsPaymentPath           = "/dashboard/services-details-bills-payment/:id"
	servicesDetailsMIPath                     = "/dashboard/services-details-mi/:id"
	serviceCatalogPath                        = "/dashboard/service-catalog/:id"
	serviceCatalogPathEdit                    = "/dashboard/service-catalog-edit/:id"
	serviceCatalogDeletePath                  = "/dashboard/service-catalog-delete"

	currencyConfigurationPath = "/dashboard/currency/:id"

	// api key management
	apiKeyListPath        = "/api-key/:apienv"
	apiKeyGenerateGetPath = "/api-key/generate/:apienv"
	apiKeySuccessGetPath  = "/api-key/generate/:apienv/success"
	authorizationPath     = "/oauth2-authorization-code/:apienv"
	apiGuide              = "/api-guide"

	// image
	viewGCSFilePath = "/u/files/:id"

	// transaction
	transactionListGetPath           = "/dashboard/transactionslist/:id"
	transactionListBPGetPath         = "/dashboard/transactionslistbp/:id"
	transactionListMIGetPath         = "/dashboard/transactionslistmi/:id"
	transactionListCICOGetPath       = "/dashboard/transactionslistcico/:id"
	transactiondetailPath            = "/dashboard/transactions/:id/:oid"
	transactiondisbursedetailPath    = "/dashboard/transactions-disburse/:id/:oid"
	dsaTransactionListGetPath        = "/transactions"
	dsaTransactionDetailPath         = "/transactions/:id"
	dsaTransactionDisburseDetailPath = "/transactions-disburse/:id"
	transactionGetPath               = "/dashboard/transaction"

	// mfa
	confirmEventPath = "/mfa-confirm"
	// Manage Members
	manageUserListPath = "/dashboard/manage-users"
	// Manage Roles
	manageRoleListPath     = "/dashboard/manage-role"
	rolePermissionPath     = "/dashboard/manage-role/edit/:id"
	manageRoleCreatePath   = "/dashboard/manage-role-create"
	manageInviteMemberPath = "/dashboard/invite-member"
	manageDisableUser      = "/dashboard/disable-user/:id"
	manageEnableUser       = "/dashboard/enable-user/:id"
	manageDeleteRole       = "/dashboard/delete-role/:id"

	// Providers List
	providersListPath  = "/dashboard/providers-list"
	inviteProviderPath = "/dashboard/invite-provider"

	// Field view ajax call
	disclosurePath    = "/dashboard/disclosure/:id"
	partnerListPath   = "/dashboard/partner-list"
	commissionFeePath = "/dashboard/commission-mgt"

	commissionRemoveFeePath   = "/dashboard/commission-mgt-remove/:id"
	partnerListPathCreate     = "/dashboard/partner-list-create"
	revenueShareMgtPath       = "/dashboard/revenue-sharing-mgt/:id"
	revenueShareMgtRemovePath = "/dashboard/revenue-sharing-mgt-remove/:oid/:id"

	// Services
	partnerServicesPath = "/dashboard/partner-services"

	// providers
	providersPath            = "/dashboard/providers"
	manageProviderCreatePath = "/dashboard/provider-create"
	manageProviderDeletePath = "/dashboard/provider-delete/:stype"
	editProvidersPath        = "/dashboard/providers/:lid"
	providerDeletePath       = "/dashboard/providers/delete/:lid"
)

// TODO(vitthal): Create error/error messages package for user facing errors
var errorDuplicateFile = errors.New("The file has already been uploaded. Please select another file.")

type Server struct {
	templates *template.Template

	*goji.Mux
	env          string
	assets       fs.FS
	assetFS      *hashfs.FS
	decoder      *schema.Decoder
	urls         URLMap
	rbac         identity
	rbacUserAuth partialIdentity
	pf           profile
	drpSB        drp
	drpLV        drp

	sess             *mw.Hydra
	disableActionMFA bool

	gcs     storage.GCS
	svcName string
	cnf     *viper.Viper

	// internal services
	remcoCommSvc iRemcoCommissionSvc
}

type URLMap struct {
	Base    string `json:"base"`
	SSO     string `json:"sso"`
	UserMgm string `json:"user_mgm"`
}

func NewServer(mux *goji.Mux,
	env string,
	logger *logrus.Entry,
	assets fs.FS,
	decoder *schema.Decoder,
	urls URLMap,
	hmw *mw.Hydra,
	cnf *viper.Viper,
	cl Cl,
	cs *Conns,
	store storage.GCS,
	svcName string,
	opts ...ServerOptions,
) (*Server, error) {
	s := &Server{
		Mux:              goji.New(),
		env:              env,
		assets:           assets,
		assetFS:          hashfs.NewFS(assets),
		decoder:          decoder,
		urls:             urls,
		sess:             hmw,
		disableActionMFA: cnf.GetBool("local.disableActionMFA"),
		rbac:             cl.rbac,
		rbacUserAuth:     cl.rbacUserAuth,
		pf:               cl.pf,
		drpSB:            cl.drpSB,
		drpLV:            cl.drpLV,
		gcs:              store,
		svcName:          svcName,
		cnf:              cnf,
	}
	if err := s.parseTemplates(); err != nil {
		return nil, err
	}

	// initialize optional dependencies
	for _, opt := range opts {
		opt(s)
	}

	csrfSecure := cnf.GetBool("csrf.secure")
	csrfSecret := cnf.GetString("csrf.secret")
	if csrfSecret == "" {
		return nil, errors.New("CSRF secret must not be empty")
	}
	mw.ChainHTTPMiddleware(s, logger,
		mw.CSRF([]byte(csrfSecret), csrf.Secure(csrfSecure), csrf.Path("/")),
	)

	if cnf.GetString("runtime.environment") == "localdev" {
		s.Mux.Use(mw.MockMiddleware)
	} else {
		s.Mux.Use(hmw.Middleware)
		s.Mux.Use(hmw.ForwardAuthMiddleware)
	}

	// dsa onboarding
	s.HandleFunc(goji.Get(verifyAccountPath), s.getVerifyAccount)
	s.HandleFunc(goji.Post(verifyAccountPath), s.postVerifyAccount)
	s.HandleFunc(goji.Get(businessInfoPath), s.getBusinessInfo)
	s.HandleFunc(goji.Post(businessInfoPath), s.postBusinessInfo)
	s.HandleFunc(goji.Get(accountInfoPath), s.getAccountInfo)
	s.HandleFunc(goji.Post(accountInfoPath), s.postAccountInfo)
	s.HandleFunc(goji.Get(financialInfoPath), s.getFinancialInfo)
	s.HandleFunc(goji.Post(financialInfoPath), s.postFinancialInfo)
	s.HandleFunc(goji.Get(drpInfoPath), s.getDRPInfo)
	s.HandleFunc(goji.Post(drpInfoPath), s.postDRPInfo)
	s.HandleFunc(goji.Get(registerSuccessPath), s.getRegisterSuccess)
	s.HandleFunc(goji.Get(skipOnboardingPath), s.skipOnboarding)
	// mfa
	s.HandleFunc(goji.Post(confirmEventPath), s.postConfirmMFAEvent)

	{ // `/dashboard/` router
		n := goji.NewSubMux()
		s.Handle(goji.NewPathSpec("/dashboard/*"), n)
		d := func(r string) string { return strings.TrimPrefix(r, "/dashboard") }

		// Admin enforcement middleware
		n.Use(cmsmw.PetnetAdmin(cnf))

		// Permission enforcement middleware
		v := cmsmw.ValidatePermission(cnf, rpmpb.NewValidationServiceClient(cs.GetIdInt()))

		// dashboard

		// dsa list and detail
		n.HandleFunc(goji.Get(d(dsaAppListPath)), s.getDSAApplicantList)
		n.Handle(goji.Get(d(dsaAppListDetailPath)), v(
			s.getDashboardDetails, pm.DSAListDetailRes, pm.ReadAct),
		)

		// checklist
		n.Handle(goji.Get(d(dashboardCheckoutList)), v(
			s.getDashboardCheckoutList, pm.DocChecklistRes, pm.ReadAct),
		)
		n.Handle(goji.Post(d(dashboardCheckoutList)), v(
			s.postDashboardDocument, pm.DocChecklistRes, pm.UpdateAct),
		)
		// fee
		n.Handle(goji.Get(d(feeManagementPath)), v(
			s.getFeeManagement, pm.FeesRes, pm.ReadAct),
		)
		n.Handle(goji.Post(d(feeManagementPath)), v(
			s.postFeeCommission, pm.FeesRes, pm.CreateAct),
		)
		n.Handle(goji.Post(d(feeManagementUpdatePath)), v(
			s.updateFeeCommission, pm.FeesRes, pm.UpdateAct),
		)
		n.Handle(goji.Post(d(feeManagementDeletePath)), v(
			s.deleteFeeCommission, pm.FeesRes, pm.DeleteAct),
		)
		// svc catalog
		n.Handle(goji.Get(d(serviceCatalogPath)), v(
			s.getPartnerCatalog, pm.SvcCatalogRes, pm.ReadAct),
		)
		n.Handle(goji.Post(d(serviceCatalogPath)), v(
			s.postPartnerCatalog, pm.SvcCatalogRes, pm.CreateAct),
		)
		n.Handle(goji.Post(d(serviceCatalogPathEdit)), v(
			s.postPartnerCatalog, pm.SvcCatalogRes, pm.UpdateAct),
		)
		n.Handle(goji.Post(d(serviceCatalogDeletePath)), v(
			s.deletePartnerCatalog, pm.SvcCatalogRes, pm.DeleteAct),
		)
		// todo we need individual handlers for create, update, delete so we can enforce
		// permissions
		// currency
		n.Handle(goji.Get(d(currencyConfigurationPath)), v(
			s.getCurrencyConfig, pm.CurrencyRes, pm.ReadAct),
		)
		n.Handle(goji.Post(d(currencyConfigurationPath)), v(
			s.changeCurrencyConfig, pm.CurrencyRes, pm.UpdateAct),
		)

		// location/branch
		n.Handle(goji.Get(d(locationConfigurationPath)), v(
			s.getLocationConfig, pm.BranchRes, pm.ReadAct),
		)
		n.Handle(goji.Post(d(locationConfigurationPath)), v(
			s.postLocationConfig, pm.BranchRes, pm.CreateAct),
		)

		n.Handle(goji.Post(d(locationUpdatePath)), v(
			s.postLocationUpdate, pm.BranchRes, pm.UpdateAct),
		)
		n.Handle(goji.Post(d(locationConfigurationUpdatePath)), v(
			s.updateLocationConfig, pm.BranchRes, pm.UpdateAct),
		)
		n.Handle(goji.Post(d(locationConfigurationDeletePath)), v(
			s.deleteLocationConfig, pm.BranchRes, pm.DeleteAct),
		)

		// transaction
		n.Handle(goji.Get(d(transactionListGetPath)), v(
			s.getTransactionListSandbox, pm.TransactionRes, pm.ReadAct),
		)
		n.Handle(goji.Get(d(transactionListBPGetPath)), v(
			s.getBPTransactionList, pm.TransactionRes, pm.ReadAct),
		)
		n.Handle(goji.Get(d(transactionListMIGetPath)), v(
			s.getMITransactionList, pm.TransactionRes, pm.ReadAct),
		)
		n.Handle(goji.Get(d(transactionListCICOGetPath)), v(
			s.getTransactionListCICOSandbox, pm.TransactionRes, pm.ReadAct),
		)
		n.Handle(goji.Get(d(transactiondetailPath)), v(
			s.getTransactionDetails, pm.TransactionRes, pm.ReadAct),
		)
		n.Handle(goji.Get(d(transactiondisbursedetailPath)), v(
			s.getTransactionDisburseDetails, pm.TransactionRes, pm.ReadAct),
		)

		n.Handle(goji.Get(d(transactionGetPath)), v(
			s.getTransactionCompanies, pm.TransactionRes, pm.ReadAct),
		)

		n.HandleFunc(goji.Post(d(dashboardSendReminderPath)), s.postDashboardSendReminder)
		n.HandleFunc(goji.Post(d(dashboardChangeStatusPath)), s.postDashboardChangeStatus)
		n.HandleFunc(goji.Post(d(dashboardCheckDSACodePath)), s.getProfileByDsaCode)

		// manage members
		n.HandleFunc(goji.Get(d(manageUserListPath)), s.getManageUserList)
		n.HandleFunc(goji.Post(d(manageUserListPath)), s.postResendEmailConfirm)
		n.HandleFunc(goji.Post(d(manageInviteMemberPath)), s.postInviteMember)
		n.HandleFunc(goji.Get(d(rolePermissionPath)), s.getRolePermission)
		n.HandleFunc(goji.Post(d(rolePermissionPath)), s.postRolePermission)

		// manage roles
		n.HandleFunc(goji.Get(d(manageRoleListPath)), s.getManageRoleList)
		n.HandleFunc(goji.Post(d(manageRoleCreatePath)), s.postManageRoleCreate)
		n.HandleFunc(goji.Get(d(manageDisableUser)), s.getManageDisableUser)
		n.HandleFunc(goji.Get(d(manageEnableUser)), s.getManageEnableUser)
		n.HandleFunc(goji.Get(d(manageDeleteRole)), s.getManageDeleteRole)

		// Providers List
		n.HandleFunc(goji.Get(d(providersListPath)), s.getProvidersList)
		n.HandleFunc(goji.Post(d(inviteProviderPath)), s.postInviteProvides)

		n.HandleFunc(goji.Post(d(disclosurePath)), s.postSpecificFieldData)
		n.HandleFunc(goji.Get(d(partnerListPath)), s.getPartnerLists)
		n.HandleFunc(goji.Post(d(partnerListPathCreate)), s.postPartnerLists)

		// Commission Fee
		n.HandleFunc(goji.Get(d(commissionFeePath)), s.doGetCommissionMgt)
		n.HandleFunc(goji.Delete(d(commissionRemoveFeePath)), s.doDeleteCommissionMgt)
		n.HandleFunc(goji.Post(d(commissionFeePath)), s.doPostCommissionMgt)

		// Services
		n.HandleFunc(goji.Get(d(partnerServicesPath)), s.getPartnerServices)
		n.HandleFunc(goji.Post(d(partnerServicesPath)), s.updatePartnerServicesStatus)

		// providersPath
		n.HandleFunc(goji.Get(d(providersPath)), s.doGetproviders)
		n.HandleFunc(goji.Post(d(manageProviderCreatePath)), s.createProvider)
		// n.HandleFunc(goji.Post(d(manageProviderDeletePath)), s.deleteProvider)
		n.HandleFunc(goji.Post(d(editProvidersPath)), s.postProviderUpdate)
		n.HandleFunc(goji.Post(d(providerDeletePath)), s.deleteProvider)

		if cnf.GetBool("feature.serviceRequest") {
			n.HandleFunc(goji.Get(d(serviceReqPath)), s.getServiceReq)
			n.HandleFunc(goji.Get(d(servicesDetailsPath)), s.getServicesDetails)
			n.HandleFunc(goji.Get(d(servicesDetailsCicoPath)), s.getServicesDetailsCico)
			n.HandleFunc(goji.Get(d(servicesDetailsBillsPaymentPath)), s.getServicesDetailsBillsPayment)
			n.HandleFunc(goji.Get(d(servicesDetailsMIPath)), s.getServicesDetailsMI)
			n.HandleFunc(goji.Post(d(postAddRemarkSvcReq)), s.postAddRemarkSvcReq)
			n.HandleFunc(goji.Get(d(changeStatusSvcReq)), s.getChangeStatusSvcReq)
			n.HandleFunc(goji.Get(d(changeAllStatusSvcReq)), s.getChangeAllStatusSvcReq)
			n.HandleFunc(goji.Get(d(changeCiCoAllStatusSvcReq)), s.getChangeCiCoAllStatusSvcReq)
			n.HandleFunc(goji.Get(d(changeBillsPaymentAllStatusSvcReq)), s.getChangeBillsPaymentAllStatusSvcReq)
			n.HandleFunc(goji.Get(d(changeMIAllStatusSvcReq)), s.getChangeMIAllStatusSvcReq)
			n.HandleFunc(goji.Get(d(changeAllStatusSvcReqCicoPath)), s.getChangeAllStatusSvcReqCico)
			n.HandleFunc(goji.Get(d(changeAllStatusSvcReqBillsPaymentPath)), s.getChangeAllStatusSvcReqBillsPayment)
			n.HandleFunc(goji.Get(d(changeAllStatusSvcReqMIPath)), s.getChangeAllStatusSvcReqMI)
			n.HandleFunc(goji.Get(d(changeAjaxAllStatusSvcReqCicoPath)), s.getAjaxChangeAllStatusSvcReqCico)
			n.HandleFunc(goji.Get(d(changeAjaxAllStatusSvcReqMIPath)), s.getAjaxChangeAllStatusSvcReqMI)
			n.HandleFunc(goji.Get(d(changeAjaxAllStatusSvcReqBillsPaymentPath)), s.getAjaxChangeAllStatusSvcReqBillsPayment)
		}
		// Revenue sharing
		n.HandleFunc(goji.Get(d(revenueShareMgtPath)), s.getRevenueSharingMgt)
		n.HandleFunc(goji.Post(d(revenueShareMgtPath)), s.postRevenueSharingMgt)
		n.HandleFunc(goji.Delete(d(revenueShareMgtRemovePath)), s.removeRevenueSharingMgt)
	}

	// profile
	s.HandleFunc(goji.Get(profilePath), s.getProfile)
	s.HandleFunc(goji.Post(profilePath), s.postProfile)
	s.HandleFunc(goji.Get(edPassPath), s.geteditPass)
	s.HandleFunc(goji.Post(edPassPath), s.posteditPass)

	if cnf.GetBool("feature.serviceRequest") {
		// dsa services
		s.HandleFunc(goji.Get(dsaServicesPath), s.getDsaServices)
		s.HandleFunc(goji.Get(dsaPrtSelPath), s.getDsaSelectPartners)
		s.HandleFunc(goji.Post(dsaPrtSelPath), s.postDsaSelectPartners)
		s.HandleFunc(goji.Get(dsaPrtCiCoSelPath), s.getDsaSelectCiCoPartners)
		s.HandleFunc(goji.Post(dsaPrtCiCoSelPath), s.postDsaSelectCicoPartners)
		s.HandleFunc(goji.Get(dsaPrtBillsPaymentSelPath), s.getDsaSelectBillsPaymentPartners)
		s.HandleFunc(goji.Post(dsaPrtBillsPaymentSelPath), s.postDsaSelectBillsPaymentPartners)
		s.HandleFunc(goji.Get(dsaReqDocsPath), s.getDsaReqDocs)
		s.HandleFunc(goji.Post(dsaReqDocsPath), s.postDsaReqDocs)
		s.HandleFunc(goji.Get(dsaReqDocsMicroInsurance), s.getDsaReqDocsMicroInsurance)
		s.HandleFunc(goji.Post(dsaReqDocsMicroInsurance), s.postDsaReqDocsMicroInsurance)
		s.HandleFunc(goji.Get(dsaAdiBillsPaymentSelPath), s.getdsaAdiBillsPaymentSelPath)
		s.HandleFunc(goji.Post(dsaAdiBillsPaymentSelPath), s.postDsaReqDocsBillsPayment)
		s.HandleFunc(goji.Get(dsaAddiDocsPath), s.getDsaAddiDocs)
		s.HandleFunc(goji.Post(dsaAddiDocsPath), s.postDsaAddiDocs)
		s.HandleFunc(goji.Get(dsaReqDocsCicoPath), s.getDsaReqDocsCico)
		s.HandleFunc(goji.Post(dsaReqDocsCicoPath), s.postDsaReqDocsCico)
	}

	// dsa transaction
	s.HandleFunc(goji.Get(dsaTransactionListGetPath), s.getDSATransactionListSandbox)
	s.HandleFunc(goji.Get(dsaTransactionDetailPath), s.getDSATransactionDetails)
	s.HandleFunc(goji.Get(dsaTransactionDisburseDetailPath), s.getDSATransactionDisburseDetails)

	// signup
	s.HandleFunc(goji.Get(preliminaryScreenPath), s.getPreliminaryScreen)
	s.HandleFunc(goji.Get(regAccountDetailsPath), s.getAccountDetails)
	s.HandleFunc(goji.Post(regAccountDetailsPath), s.postAccountDetails)

	s.HandleFunc(goji.Get(errorPath), s.handleError)
	s.HandleFunc(goji.Get(loginPath), s.handleLogin)
	s.HandleFunc(goji.Post(logoutPath), s.handleLogout)
	s.HandleFunc(goji.Get(logoutPath), s.handleLogout)
	s.HandleFunc(goji.NewPathSpec("/oauth2/callback"), s.handleCallback)

	s.HandleFunc(goji.Get(rootPath), s.getRoot)
	s.HandleFunc(goji.Post(rootPath), s.getRoot)

	// api key management
	s.HandleFunc(goji.Get(apiKeyListPath), s.getApiKeyList)
	s.HandleFunc(goji.Get(apiKeyGenerateGetPath), s.getApiKeyGenerate)
	s.HandleFunc(goji.Post(apiKeyGenerateGetPath), s.postApiKeyGenerate)
	s.HandleFunc(goji.Get(apiKeySuccessGetPath), s.getApiKeySuccess)
	s.HandleFunc(goji.Get(authorizationPath), s.getAuthorizationCode)
	s.HandleFunc(goji.Post(authorizationPath), s.postAuthorizationCode)

	s.HandleFunc(goji.Get(apiGuide), s.getApiGuide)

	// image
	s.HandleFunc(goji.Get(viewGCSFilePath), s.getViewGCSFile)
	s.HandleFunc(goji.Delete(viewGCSFilePath), s.removeFile)

	s.HandleFunc(goji.NewPathSpec("/*"), s.handleIndex)

	return s, nil
}

func (s *Server) lookupTemplate(name string) *template.Template {
	if s.env == "localdev" {
		if err := s.parseTemplates(); err != nil {
			return nil
		}
	}
	return s.templates.Lookup(name)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context()).WithField("path", r.URL.Path)
	log.Trace("received")
	pth := strings.TrimPrefix(r.URL.Path, "/")
	switch {
	case r.URL.Path == "" || r.URL.Path == "/":
		ot := GetOrgType(r.Context())
		switch ot {
		case pfpb.OrgType_PetNet:
			http.Redirect(w, r, dsaAppListPath, http.StatusTemporaryRedirect)
			return
		case pfpb.OrgType_DSA:
			http.Redirect(w, r, "/api-guide", http.StatusTemporaryRedirect)
			return
		case pfpb.OrgType_UnknownOrgType:
			http.Redirect(w, r, rootPath, http.StatusTemporaryRedirect)
			return
		}
	case filepath.Ext(pth) == ".html":
		if err := s.doTemplate(w, r, strings.TrimPrefix(r.URL.Path, "/"), http.StatusOK); err != nil {
			log.Infof("unable to load template: %+v", err)
			http.Error(w, "unable to load template", http.StatusInternalServerError)
		}
		return
	case r.URL.Path == "/favicon.ico":
	}
	log.Trace("fs check")
	if _, err := fs.Stat(s.assetFS, strings.TrimPrefix(r.URL.Path, "/")); err == nil {
		w.Header().Set("Cache-Control", "max-age=86400")
		if _, h := hashfs.ParseName(r.URL.Path); h != "" {
			// if asset is hashed extend cache to 180 days
			w.Header().Set("Cache-Control", "max-age=15552000")
		}
		http.FileServer(http.FS(s.assetFS)).ServeHTTP(w, r)
		return
	} else {
		logging.WithError(err, log).Error("stat error")
	}
	log.Trace("error template")
	if err := s.doTemplate(w, r, "error.html", http.StatusNotFound); err != nil {
		log.WithError(err).Error("unable to load error page")
		http.Error(w, "unable to load error page", http.StatusInternalServerError)
	}
}

func (s *Server) templateData(r *http.Request) TemplateData {
	return TemplateData{
		Env:       s.env,
		CSRFField: csrf.TemplateField(r),
	}
}

func (s *Server) doTemplate(w http.ResponseWriter, r *http.Request, name string, status int) error {
	template := s.lookupTemplate(name)
	if template == nil || isPartialTemplate(name) {
		template, status = s.templates.Lookup("error.html"), http.StatusNotFound
	}

	w.WriteHeader(status)
	return template.Execute(w, s.templateData(r))
}

type TemplateData struct {
	Env       string
	CSRFField template.HTML
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-forwarded-for")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func isPartialTemplate(name string) bool {
	return strings.HasSuffix(name, ".part.html")
}

func uniqueSlice(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (s *Server) parseTemplates() error {
	manilaTZ, err := time.LoadLocation(timeZone)
	if err != nil {
		return err
	}
	templates := template.New("cms-templates").Funcs(template.FuncMap{
		"minorToAmount": func(mamt, cur string) string {
			amt, err := currency.NewMinor(mamt, cur)
			if err != nil {
				return ""
			}
			return amt.Amount.Round().Number()
		},
		"formatStatus": func(status string) string {
			return strings.ToLower(strings.ReplaceAll(status, " ", "-"))
		},
		"formatTimestamp": func(ts *tspb.Timestamp, layout string) string {
			if layout == "" {
				layout = timeLayout
			}

			if !ts.IsValid() {
				return ""
			}
			return ts.AsTime().In(manilaTZ).Format(layout)
		},
		"formatDate": func(ts time.Time, layout string) string {
			if layout == "" {
				layout = timeLayout
			}
			if ts.IsZero() {
				return ""
			}
			return ts.In(manilaTZ).Format(layout)
		},
		"mkMap": func(val ...interface{}) (map[string]interface{}, error) {
			if len(val)%2 != 0 {
				return nil, errors.New("invalid map inputs")
			}
			mp := make(map[string]interface{}, len(val)/2)
			for i := 0; i < len(val); i++ {
				k, ok := val[i].(string)
				if !ok {
					return nil, errors.New("keys must be strings")
				}
				i++
				mp[k] = val[i]
			}
			return mp, nil
		},
		"assetHash": func(n string) string {
			return path.Join("/", s.assetFS.HashName(strings.TrimPrefix(path.Clean(n), "/")))
		},
		"riskScoreClass": func(n string) string {
			riskScores := map[string]string{"UnknownRiskScore": "gray-300", "NIL/NA": "gray-300", "High": "petnetpink", "Medium": "petnetorange", "Low": "petnetgreen"}
			for key, value := range riskScores {
				if key == n {
					return value
				}
			}
			return ""
		},
		"statusClass": func(n string) string {
			status := map[string]string{"Incomplete": "gray-300", "Accepted": "petnetlightblue text-white", "Completed": "petnetgreen text-white", "Rejected": "petnetpink text-white", "Pending": "yellow-400", "PendingDocuments": "blue-800 text-white", "Pending Documents": "blue-800 text-white"}
			for key, value := range status {
				if key == n {
					return value
				}
			}
			return ""
		},
		"serviceStatusClass": func(n string) string {
			status := map[string]string{"ACCEPTED": "petnetlighterblue", "PENDING": "petnetlightyellow", "PARTNERDRAFT": "petnetlightyellow", "REQDOCDRAFT": "petnetlightyellow", "REJECTED": "petnetstatuspink", "NOSTATUS": "petnetlightyellow", "ADDIDOCDRAFT": "petnetlightyellow"}
			for key, value := range status {
				if key == n {
					return value
				}
			}
			return ""
		},
		"userStatusClass": func(n string) string {
			status := map[string]string{"Expired": "gray-300", "Approved": "petnetlightblue text-white", "InProgress": "petnetgreen", "Revoked": "petnetpink text-white", "In Progress": "yellow-400", "Invited": "blue-800 text-white", "Invite Sent": "green-200 text-black"}
			for key, value := range status {
				if key == n {
					return value
				}
			}
			return ""
		},
		"servicesStatusClass": func(n string) string {
			status := map[string]string{"Accepted": "petnetlighterblue", "ACCEPTED": "petnetlighterblue", "Rejected": "petnetstatuspink", "REJECTED": "petnetstatuspink", "Pending": "petnetslightgray", "PENDING": "petnetslightgray", "DISABLED": "petnetlightyellow", "Disabled": "petnetlightyellow"}
			for key, value := range status {
				if key == n {
					return value
				}
			}
			return ""
		},
		"jsonStringify": func(data interface{}) string {
			mapB, _ := json.Marshal(data)
			return string(mapB)
		},
		"countPaginate": func(a, b int32) int32 {
			if a > 0 {
				c := a / b
				if a%b != 0 {
					c = c + 1
				}
				return c
			}
			return 0
		},
		"noescape": func(str string) template.HTML {
			return template.HTML(str)
		},
		"stringTitle": func(str string) string {
			return strings.Title(strings.ToLower(str))
		},
		"selectedCheck": func(ptnr string, SelectedPartners []string) bool {
			if len(SelectedPartners) == 0 {
				return false
			}
			hv, _ := cmsmw.InArray(ptnr, SelectedPartners)
			return hv
		},
		"stringContains": func(s string, substr string) bool {
			return strings.Contains(s, substr)
		},
		"getFileName": func(files map[string]string, file string) string {
			if f, ok := files[file]; ok {
				return f
			}
			return file
		},
	}).Funcs(sprig.FuncMap())

	tmpl, err := templates.ParseFS(s.assets, "templates/*/*.html")
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}

func (s *Server) loadUserInfo(r *http.Request) User {
	user := User{
		FirstName: "placeholder",
		LastName:  "placeholder",
	}
	uid := mw.GetUserID(r.Context())
	if uid == "" {
		return user
	}

	cu, err := s.rbac.GetUser(r.Context(), &rbupb.GetUserRequest{ID: uid})
	if err != nil {
		return user
	}

	user = userData(cu.GetUser())
	if user.FirstName == "" {
		user.FirstName = "placeholder"
	}

	if user.LastName == "" {
		user.LastName = "placeholder"
	}

	return user
}

func (s *Server) storeMultiToGCS(r *http.Request, name, oid string) ([]string, map[string]string, error) {
	ctx := r.Context()
	var us []string
	fn := map[string]string{}
	fhs := r.MultipartForm.File[name]
	for _, h := range fhs {
		if h.Filename == "" {
			return nil, nil, errors.New(fmt.Sprintf("%s: File Name is empty", name))
		}

		f, err := h.Open()
		if err != nil {
			return nil, nil, err
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, nil, err
		}
		sf := storage.NewFile(b, oid)
		u, err := s.gcs.Store(ctx, sf)
		if err != nil {
			return nil, nil, err
		}
		tus := strings.TrimSpace(u)
		us = append(us, tus)
		fn[tus] = h.Filename
	}
	return us, fn, nil
}

func (s *Server) storeSingleToGCS(r *http.Request, name, oid string) (string, map[string]string, error) {
	fn := map[string]string{}
	f, h, err := r.FormFile(name)
	if err != nil {
		return "", nil, err
	}
	if h.Filename == "" {
		return "", nil, errors.New(fmt.Sprintf("%s: File Name is empty", name))
	}
	defer f.Close()

	c, err := ioutil.ReadAll(f)
	if err != nil {
		return "", nil, err
	}
	ctx := r.Context()
	sf := storage.NewFile(c, oid)
	u, err := s.gcs.Store(ctx, sf)
	if err != nil {
		return "", nil, err
	}
	ut := strings.TrimSpace(u)
	fn[ut] = h.Filename
	return ut, fn, nil
}

// validateFileHeadersForDuplicateAndFileType validates multipart file headers for given types
func validateFileHeadersForDuplicateAndFileType(fileHeaders []*multipart.FileHeader, types map[string]bool) error {
	fileNamesMap := make(map[string]bool, len(fileHeaders))
	for _, h := range fileHeaders {
		// error if duplicate file present
		if _, ok := fileNamesMap[h.Filename]; ok {
			return errorDuplicateFile
		}

		f, err := h.Open()
		if err != nil {
			return err
		}

		defer f.Close()
		err = validateFileType(f, types)
		if err != nil {
			return err
		}

		fileNamesMap[h.Filename] = true
	}
	return nil
}

// validateFileType validates multipart file for given types
func validateFileType(f multipart.File, types map[string]bool) error {
	// Need to read all file for doc and docx file types
	buff, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	defer func() {
		// reset read pointer to start which was set by above read
		f.Seek(0, io.SeekStart)
	}()

	invalidErr := errors.New("Invalid File Format")

	kind, err := filetype.Match(buff)
	if err != nil {
		return invalidErr
	}

	if kind == filetype.Unknown {
		return invalidErr
	}

	if valid, ok := types[kind.MIME.Value]; ok && valid {
		return nil
	}

	return invalidErr
}

// deprecated
func (s *Server) validateMultiFileType(r *http.Request, name string, types []string) error {
	fhs := r.MultipartForm.File[name]
	for _, h := range fhs {
		f, err := h.Open()
		defer f.Close()
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		ft := http.DetectContentType(b)
		rtn := false
		for _, t := range types {
			if ft == t {
				rtn = true
			}
		}
		if !rtn {
			return errors.New("Invalid File Format")
		}
	}
	return nil
}

// deprecated
func (s *Server) validateSingleFileType(r *http.Request, name string, types []string) error {
	f, _, err := r.FormFile(name)
	if err != nil {
		return nil
	}
	defer f.Close()
	c, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	ft := http.DetectContentType(c)
	rtn := false
	for _, t := range types {
		if ft == t {
			rtn = true
		}
	}
	if !rtn {
		return errors.New("Invalid File Format")
	}
	return nil
}

func GetOrgType(ctx context.Context) pfpb.OrgType {
	ot := mw.GetOrgType(ctx)
	if ot == "" {
		return pfpb.OrgType_UnknownOrgType
	}
	oti, err := strconv.Atoi(ot)
	if err != nil {
		return pfpb.OrgType_UnknownOrgType
	}
	return pfpb.OrgType(oti)
}

type enforceTemplateData struct {
	PresetPermission map[string]map[string]bool
	ServiceRequests  bool
}

func (s *Server) getEnforceTemplateData(ctx context.Context) enforceTemplateData {
	res := enforceTemplateData{
		PresetPermission: s.getPermissionCheckingList(ctx),
		ServiceRequests:  s.cnf.GetBool("feature.serviceRequest"),
	}
	return res
}

func (s *Server) getPermissionCheckingList(ctx context.Context) map[string]map[string]bool {
	presetPermission := make(map[string]map[string]bool)
	IsPetnetAdmin := mw.IsPetnetOwner(ctx)
	isProvider := mw.IsProvider(ctx)
	presetPermission["ispetnetowner"] = make(map[string]bool) // for checking IsPetnetAdmin from template
	presetPermission["isprovider"] = make(map[string]bool)    // for checking IsProvider from template
	presetPermission["isprovider"]["permission"] = false      // for checking IsPetnetAdmin from template
	presetPermission["ispetnetowner"]["permission"] = false   // for checking IsPetnetAdmin from template
	for _, pv := range pm.Permissions {
		presetPermission[pv.Resource] = make(map[string]bool)
	}
	for _, pv := range pm.Permissions {
		for _, av := range pv.Actions {
			presetPermission[pv.Resource][av] = false
		}
	}
	for _, pv := range pm.Permissions {
		for _, av := range pv.Actions {
			if IsPetnetAdmin {
				presetPermission[pv.Resource][av] = true
			}
		}
	}
	if IsPetnetAdmin {
		presetPermission["ispetnetowner"]["permission"] = true // for checking IsPetnetAdmin from template
		return presetPermission
	}

	if isProvider {
		presetPermission["isprovider"]["permission"] = true
	}

	UserID := mw.GetUserID(ctx)
	Oid := mw.GetOrgID(ctx)
	for _, pv := range pm.Permissions {
		for _, av := range pv.Actions {
			_, err := s.rbac.ValidatePermission(ctx, &rpmpb.ValidatePermissionRequest{
				ID:       UserID,
				Resource: pv.Resource,
				Action:   av,
				OrgID:    Oid,
			})
			if err == nil {
				presetPermission[pv.Resource][av] = true
			}
		}
	}
	return presetPermission
}

func (s *Server) getUsrIdFromSessionUser(ctx context.Context, sess *sessions.Session) (string, error) {
	sesUid := sess.Values[sessionUserID]
	if sesUid != nil {
		return sesUid.(string), nil
	}

	uid := mw.GetUserID(ctx)
	if uid == "" {
		return "", errors.New("uid not forund")
	}

	return uid, nil
}

func (s *Server) getOrgIdFromSessionUser(ctx context.Context, sess *sessions.Session) (string, error) {
	sesOid := sess.Values[sessionOrgID]
	if sesOid != nil {
		return sesOid.(string), nil
	}

	sesMwOid := sess.Values[session.OrgID]
	if sesMwOid != nil {
		return sesMwOid.(string), nil
	}

	oid := mw.GetOrgID(ctx)
	if oid != "" {
		return oid, nil
	}

	uid, err := s.getUsrIdFromSessionUser(ctx, sess)
	if err != nil {
		return "", err
	}

	usrs, err := s.pf.GetUserProfile(ctx, &user.GetUserProfileRequest{
		UserID: uid,
	})
	if err != nil {
		return "", err
	}

	usrP := usrs.GetProfile()
	if usrP == nil {
		return "", errors.New("user info not forund")
	}

	return usrP.GetOrgID(), nil
}

func (s *Server) IsPetnetOwner(ctx context.Context, oid string, uid string) bool {
	res, err := s.pf.GetProfile(ctx, &pfpb.GetProfileRequest{
		OrgID: oid,
	})
	if err != nil {
		return false
	}
	pf := res.GetProfile()
	if pf == nil {
		return false
	}
	if pf.GetUserID() == "" {
		return false
	}
	if pf.GetUserID() == uid {
		return true
	}
	return false
}

type UserDetailsInfo struct {
	CompanyName  string
	UserInfo     User
	ProfileImage string
}

type UserInfo struct {
	ID           string
	OrgID        string
	OrgName      string
	FirstName    string
	LastName     string
	ProfileImage string
	Email        string
	Aemail       string
	InviteStatus string
	CountryCode  string
	Phone        string
}

type ProfileInfo struct {
	CompanyName  string
	ProfileImage string
}

func (s *Server) GetUserInfo(w http.ResponseWriter, r *http.Request) (res UserDetailsInfo) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := mw.GetOrgID(ctx)
	uid := mw.GetUserID(ctx)
	res = UserDetailsInfo{}
	if oid == "" && uid == "" {
		return
	}
	res.UserInfo = s.loadUserInfo(r)
	if oid != "" {
		pf, err := s.pf.GetProfile(r.Context(), &pfpb.GetProfileRequest{OrgID: oid})
		if err != nil {
			logging.WithError(err, log).Info("failed to get profile")
		}

		if pf != nil && pf.GetProfile() != nil && pf.GetProfile().GetBusinessInfo() != nil {
			res.CompanyName = pf.GetProfile().GetBusinessInfo().CompanyName
		}
	}

	if uid != "" {
		gp, err := s.pf.GetUserProfile(ctx, &user.GetUserProfileRequest{
			UserID: uid,
		})
		if err != nil {
			log.Error("failed to get user profile")
		}

		if gp != nil && gp.GetProfile() != nil {
			res.ProfileImage = gp.GetProfile().ProfilePicture
		}
	}

	return
}

func (s *Server) GetUserInfoFromCookie(w http.ResponseWriter, r *http.Request, force bool) (res UserDetailsInfo) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	res = UserDetailsInfo{}
	usrsInf := UserDetailsInfo{}
	usr := User{}
	uData := UserInfo{}
	pData := ProfileInfo{}
	sess, err := s.sess.Get(r, userInfoDetails)
	if err != nil {
		log.WithError(err).Error("fetching session")
		return
	}
	usrInfo := sess.Values[userInfo]
	if usrInfo != nil {
		uDataS, _ := usrInfo.(string)
		json.Unmarshal([]byte(string(uDataS)), &uData)
	}

	profileInf := sess.Values[profileInfo]
	if profileInf != nil {
		pDataS, _ := profileInf.(string)
		json.Unmarshal([]byte(string(pDataS)), &pData)
	}

	if uData.FirstName == "" || pData.CompanyName == "" || force {
		usrsInf = s.GetUserInfo(w, r)
	}

	if uData.FirstName == "" && usrsInf.UserInfo != usr || force {
		uf, _ := json.Marshal(UserInfo{
			ID:           usrsInf.UserInfo.ID,
			OrgID:        usrsInf.UserInfo.OrgID,
			OrgName:      usrsInf.UserInfo.OrgName,
			FirstName:    usrsInf.UserInfo.FirstName,
			LastName:     usrsInf.UserInfo.LastName,
			ProfileImage: usrsInf.UserInfo.ProfileImage,
			Email:        usrsInf.UserInfo.Email,
			Aemail:       usrsInf.UserInfo.Aemail,
			InviteStatus: usrsInf.UserInfo.InviteStatus,
			CountryCode:  usrsInf.UserInfo.CountryCode,
			Phone:        usrsInf.UserInfo.Phone,
		})
		sess.Values[userInfo] = string(uf)
	}

	if pData.CompanyName == "" && usrsInf.CompanyName != "" || force {
		pf, _ := json.Marshal(ProfileInfo{
			CompanyName:  usrsInf.CompanyName,
			ProfileImage: usrsInf.ProfileImage,
		})
		sess.Values[profileInfo] = string(pf)
	}
	if uData.FirstName == "" || pData.CompanyName == "" || force {
		if err := s.sess.Save(r, w, sess); err != nil {
			log.WithError(err).Error("failed to save session")
			return
		}

		usrInfo = sess.Values[userInfo]
		if usrInfo != nil {
			uDataS, _ := usrInfo.(string)
			json.Unmarshal([]byte(string(uDataS)), &uData)
		}

		profileInf = sess.Values[profileInfo]
		if profileInf != nil {
			pDataS, _ := profileInf.(string)
			json.Unmarshal([]byte(string(pDataS)), &pData)
		}
	}
	res.UserInfo = User{
		ID:           uData.ID,
		OrgID:        uData.OrgID,
		OrgName:      uData.OrgName,
		FirstName:    uData.FirstName,
		LastName:     uData.LastName,
		ProfileImage: uData.ProfileImage,
		Email:        uData.Email,
		Aemail:       uData.Aemail,
		InviteStatus: uData.InviteStatus,
		CountryCode:  uData.CountryCode,
		Phone:        uData.Phone,
	}
	res.CompanyName = pData.CompanyName
	res.ProfileImage = pData.ProfileImage
	res.UserInfo.ProfileImage = pData.ProfileImage
	return
}
