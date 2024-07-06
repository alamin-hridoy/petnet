package handler

import (
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	cmsmw "brank.as/petnet/cms/mw"
	"brank.as/petnet/cms/paginator"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	pfsvc "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
)

type dsaServiceTempData struct {
	ServiceRequestList []ServiceList
	UserInfo           *User
	CompanyName        string
	PaginationData     paginator.Paginator
	SearchTerms        string
	HasLiveAccess      bool
	PresetPermission   map[string]map[string]bool
	ServiceRequest     bool
}

type ServiceList struct {
	ServiceName         string
	RequestDate         time.Time
	Applied             time.Time
	Status              string
	AllPartnersSelected bool
}

type By func(s1, s2 *ServiceList) bool

func (by By) Sort(ServiceLists []ServiceList) {
	ps := &ServiceListSorter{
		ServiceLists: ServiceLists,
		by:           by,
	}
	sort.Sort(ps)
}

type SBy func(s1, s2 *pfsvc.ServiceRequest) bool

func (by SBy) Sort(ServiceLists []*pfsvc.ServiceRequest) {
	ps := &SServiceListSorter{
		ServiceLists: ServiceLists,
		by:           by,
	}
	sort.Sort(ps)
}

type ServiceListSorter struct {
	ServiceLists []ServiceList
	by           func(s1, s2 *ServiceList) bool
}

type SServiceListSorter struct {
	ServiceLists []*pfsvc.ServiceRequest
	by           func(s1, s2 *pfsvc.ServiceRequest) bool
}

func (s *ServiceListSorter) Len() int {
	return len(s.ServiceLists)
}

func (s *ServiceListSorter) Swap(i, j int) {
	s.ServiceLists[i], s.ServiceLists[j] = s.ServiceLists[j], s.ServiceLists[i]
}

func (s *ServiceListSorter) Less(i, j int) bool {
	return s.by(&s.ServiceLists[i], &s.ServiceLists[j])
}

func (s *SServiceListSorter) Len() int {
	return len(s.ServiceLists)
}

func (s *SServiceListSorter) Swap(i, j int) {
	s.ServiceLists[i], s.ServiceLists[j] = s.ServiceLists[j], s.ServiceLists[i]
}

func (s *SServiceListSorter) Less(i, j int) bool {
	return s.by(s.ServiceLists[i], s.ServiceLists[j])
}

func (s *Server) getDsaServices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	template := s.templates.Lookup("dsa-services.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid := mw.GetOrgID(ctx)
	queryParams := r.URL.Query()
	searchTerms, err := url.PathUnescape(queryParams.Get("search-term"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}
	lSvc, err := s.pf.ListServiceRequest(ctx, &pfsvc.ListServiceRequestRequest{
		OrgIDs: []string{oid},
	})
	if err != nil {
		log.WithError(err).Error(err.Error())
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	etd := s.getEnforceTemplateData(ctx)
	svcLists := []string{
		pfsvc.ServiceType_BILLSPAYMENT.String(),
		pfsvc.ServiceType_REMITTANCE.String(),
		pfsvc.ServiceType_CASHINCASHOUT.String(),
		pfsvc.ServiceType_MICROINSURANCE.String(),
	}
	sListM := make(map[string]ServiceList)
	newReq := &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: pfsvc.ServiceType_REMITTANCE.String(),
	}
	var partnerListApplicantList []PartnerListApplicant
	var partnerListApplicantListCico []PartnerListApplicant
	var partnerListApplicantListBillsPayment []PartnerListApplicant
	var ptnrListsCico []string
	var ptnrListsBillsPayment []string
	PartnerListMaps := s.getPartnerListMaps(ctx, pfsvc.ServiceType_REMITTANCE.String())
	ap := PartnerListApplicant{
		Stype: AllPartners,
		Name:  AllPartners,
	}
	partnerListApplicantList = append(partnerListApplicantList, ap)
	partnerListApplicantListCico = append(partnerListApplicantListCico, ap)
	partnerListApplicantListBillsPayment = append(partnerListApplicantListBillsPayment, ap)
	svc, err := s.pf.GetPartnerList(ctx, newReq)
	if err != nil {
		log.Error("failed to Get Partner List")
	} else {
		for _, sv := range svc.GetPartnerList() {
			partnerListApplicantList = append(partnerListApplicantList, PartnerListApplicant{
				Stype: sv.GetStype(),
				Name:  sv.GetName(),
			})
		}
	}
	newReqq := &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: pfsvc.ServiceType_CASHINCASHOUT.String(),
	}
	svcCico, err := s.pf.GetPartnerList(ctx, newReqq)
	if err != nil {
		log.Error("failed to Get CiCo Partner List")
	} else {
		for _, sv := range svcCico.GetPartnerList() {
			ptnrListsCico = append(ptnrListsCico, sv.GetStype())
			partnerListApplicantListCico = append(partnerListApplicantListCico, PartnerListApplicant{
				Stype: sv.GetStype(),
				Name:  sv.GetName(),
			})
		}
	}
	newReqqBills := &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: pfsvc.ServiceType_BILLSPAYMENT.String(),
	}
	svcBillPayment, err := s.pf.GetPartnerList(ctx, newReqqBills)
	if err != nil {
		log.Error("failed to Get Bills Payment Partner List")
	} else {
		for _, sv := range svcBillPayment.GetPartnerList() {
			ptnrListsBillsPayment = append(ptnrListsBillsPayment, sv.GetStype())
			partnerListApplicantListBillsPayment = append(partnerListApplicantListBillsPayment, PartnerListApplicant{
				Stype: sv.GetStype(),
				Name:  sv.GetName(),
			})
		}
	}
	var selectedPartnerList []string
	var selectedPartnerListCico []string
	var selectedPartnerListBillsPayment []string
	allpartnersselected := false
	allpartnersselectedCico := false
	allpartnersselectedBillsPayment := false
	var acceptedPartners []string
	var acceptedPartnersCico []string
	var acceptedPartnersBillsPayment []string
	ress, _ := s.pf.ListServiceRequest(ctx, &pfsvc.ListServiceRequestRequest{
		OrgIDs:       []string{oid},
		Types:        []pfsvc.ServiceType{pfsvc.ServiceType_REMITTANCE},
		SortByColumn: pfsvc.ServiceSort_APPLIED,
	})
	if ress != nil {
		for _, vv := range ress.GetServiceRequst() {
			selectedPartnerList = append(selectedPartnerList, getPartnerFullName(PartnerListMaps, vv.Partner))
			if vv.Status == pfsvc.ServiceRequestStatus_ACCEPTED || vv.Status == pfsvc.ServiceRequestStatus_PENDING {
				acceptedPartners = append(acceptedPartners, vv.Partner)
			}
		}
	}
	CiCoPartnerListMaps := s.getPartnerListMaps(ctx, pfsvc.ServiceType_CASHINCASHOUT.String())
	ressCico, _ := s.pf.ListServiceRequest(ctx, &pfsvc.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []pfsvc.ServiceType{pfsvc.ServiceType_CASHINCASHOUT},
	})
	if ressCico != nil {
		for _, vv := range ressCico.GetServiceRequst() {
			selectedPartnerListCico = append(selectedPartnerListCico, getPartnerFullName(CiCoPartnerListMaps, vv.Partner))
			if vv.Status == pfsvc.ServiceRequestStatus_ACCEPTED || vv.Status == pfsvc.ServiceRequestStatus_PENDING {
				acceptedPartnersCico = append(acceptedPartnersCico, vv.Partner)
			}
		}
	}
	BillsPaymentPartnerListMaps := s.getPartnerListMaps(ctx, pfsvc.ServiceType_BILLSPAYMENT.String())
	ressBillsPayment, _ := s.pf.ListServiceRequest(ctx, &pfsvc.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []pfsvc.ServiceType{pfsvc.ServiceType_BILLSPAYMENT},
	})
	if ressBillsPayment != nil {
		for _, vv := range ressBillsPayment.GetServiceRequst() {
			selectedPartnerListBillsPayment = append(selectedPartnerListBillsPayment, getPartnerFullName(BillsPaymentPartnerListMaps, vv.Partner))
			if vv.Status == pfsvc.ServiceRequestStatus_ACCEPTED || vv.Status == pfsvc.ServiceRequestStatus_PENDING {
				acceptedPartnersBillsPayment = append(acceptedPartnersBillsPayment, vv.Partner)
			}
		}
	}
	hv, _ := cmsmw.InArray(ap, partnerListApplicantList)
	hvCico, _ := cmsmw.InArray(ap, partnerListApplicantListCico)
	hvBillsPayment, _ := cmsmw.InArray(ap, partnerListApplicantListBillsPayment)
	cpl := len(partnerListApplicantList)
	cplCico := len(partnerListApplicantListCico)
	cplBillsPayment := len(partnerListApplicantListBillsPayment)
	if hv {
		cpl = cpl - 1
	}
	if hvCico {
		cplCico = cplCico - 1
	}
	if hvBillsPayment {
		cplBillsPayment = cplBillsPayment - 1
	}
	if len(selectedPartnerList) >= cpl {
		allpartnersselected = true
	}
	if len(selectedPartnerListCico) >= cplCico {
		allpartnersselectedCico = true
	}
	if len(selectedPartnerListBillsPayment) >= cplBillsPayment {
		allpartnersselectedBillsPayment = true
	}
	pf, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
		return
	}
	transactionTypes, err := getTransactionTypes(pf.GetProfile())
	if err != nil {
		logging.WithError(err, log).Info("getting transactionTypes")
		return
	}
	dplRes, err := s.pf.GetDSAPartnerList(ctx, &spbl.DSAPartnerListRequest{TransactionType: transactionTypes})
	if err != nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	var ptnrLists []string
	if dplRes.DSAPartnerList != nil {
		for _, v := range dplRes.DSAPartnerList {
			ptnrLists = append(ptnrLists, v.Partner)
		}
	}
	if testEq(ptnrListsCico, acceptedPartnersCico) {
		allpartnersselectedCico = true
	}
	if testEq(ptnrListsBillsPayment, acceptedPartnersBillsPayment) {
		allpartnersselectedBillsPayment = true
	}
	if testEq(ptnrLists, acceptedPartners) {
		allpartnersselected = true
	}
	var sList []ServiceList
	for i, sn := range pfsvc.ServiceType_name {
		hV, _ := cmsmw.InArray(pfsvc.ServiceType_name[i], svcLists)
		if hV {
			allSelected := false
			switch pfsvc.ServiceType_name[i] {
			case pfsvc.ServiceType_REMITTANCE.String():
				allSelected = allpartnersselected
			case pfsvc.ServiceType_CASHINCASHOUT.String():
				allSelected = allpartnersselectedCico
			case pfsvc.ServiceType_BILLSPAYMENT.String():
				allSelected = allpartnersselectedBillsPayment
			}
			sListM[strings.ToLower(sn)] = ServiceList{
				ServiceName:         pfsvc.ServiceType_name[i],
				RequestDate:         time.Time{},
				Applied:             time.Time{},
				AllPartnersSelected: allSelected,
				Status:              pfsvc.ServiceRequestStatus_name[0],
			}
		}
	}
	stsSet := map[string]string{}
	stsSetCico := map[string]string{}
	stsSetBillsPayment := map[string]string{}
	for _, v := range lSvc.GetServiceRequst() {
		if v.Type.String() == pfsvc.ServiceType_REMITTANCE.String() {
			stsSet[v.GetOrgID()] = v.GetStatus().String()
		}
		if v.Type.String() == pfsvc.ServiceType_CASHINCASHOUT.String() {
			stsSetCico[v.GetOrgID()] = v.GetStatus().String()
		}
		if v.Type.String() == pfsvc.ServiceType_BILLSPAYMENT.String() {
			stsSetBillsPayment[v.GetOrgID()] = v.GetStatus().String()
		}
	}
	if len(lSvc.GetServiceRequst()) > 0 {

		SBy(func(s1, s2 *pfsvc.ServiceRequest) bool {
			return s1.Applied.AsTime().Before(s2.Applied.AsTime())
		}).Sort(lSvc.GetServiceRequst())

		for _, lv := range lSvc.GetServiceRequst() {
			allSelected := false
			var sts string
			if lv.GetStatus() == pfsvc.ServiceRequestStatus_PENDING && lv.Type.String() == pfsvc.ServiceType_REMITTANCE.String() {
				stsSet[lv.GetOrgID()] = pfsvc.ServiceRequestStatus_PENDING.String()
			}
			if lv.GetStatus() == pfsvc.ServiceRequestStatus_PENDING && lv.Type.String() == pfsvc.ServiceType_CASHINCASHOUT.String() {
				stsSetCico[lv.GetOrgID()] = pfsvc.ServiceRequestStatus_PENDING.String()
			}
			if lv.GetStatus() == pfsvc.ServiceRequestStatus_PENDING && lv.Type.String() == pfsvc.ServiceType_BILLSPAYMENT.String() {
				stsSetBillsPayment[lv.GetOrgID()] = pfsvc.ServiceRequestStatus_PENDING.String()
			}
			switch lv.Type.String() {
			case pfsvc.ServiceType_REMITTANCE.String():
				allSelected = allpartnersselected
				sts = stsSet[lv.GetOrgID()]
			case pfsvc.ServiceType_CASHINCASHOUT.String():
				allSelected = allpartnersselectedCico
				sts = stsSetCico[lv.GetOrgID()]
			case pfsvc.ServiceType_BILLSPAYMENT.String():
				allSelected = allpartnersselectedBillsPayment
				sts = stsSetBillsPayment[lv.GetOrgID()]
			default:
				sts = lv.GetStatus().String()
			}
			sListM[strings.ToLower(pfsvc.ServiceType_name[int32(lv.GetType())])] = ServiceList{
				ServiceName: pfsvc.ServiceType_name[int32(lv.GetType())],
				RequestDate: func() time.Time {
					if sts == pfsvc.ServiceRequestStatus_PENDING.String() ||
						sts == pfsvc.ServiceRequestStatus_ACCEPTED.String() ||
						sts == pfsvc.ServiceRequestStatus_REJECTED.String() {
						return lv.Applied.AsTime()
					}
					return lv.Created.AsTime()
				}(),
				Applied:             lv.Applied.AsTime(),
				Status:              sts,
				AllPartnersSelected: allSelected,
			}
		}
	}
	for i, v := range pfsvc.ServiceType_name {
		hV, _ := cmsmw.InArray(pfsvc.ServiceType_name[i], svcLists)
		if hV {
			sList = append(sList, sListM[strings.ToLower(v)])
		}
	}

	By(func(s1, s2 *ServiceList) bool {
		return s1.ServiceName < s2.ServiceName
	}).Sort(sList)

	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	tempData := dsaServiceTempData{
		ServiceRequestList: sList,
		UserInfo:           &usrInfo.UserInfo,
		CompanyName:        usrInfo.CompanyName,
		PaginationData:     paginator.Paginator{},
		SearchTerms:        searchTerms,
		HasLiveAccess:      s.hasLiveAccess(r.Context(), oid),
		PresetPermission:   etd.PresetPermission,
		ServiceRequest:     etd.ServiceRequests,
	}
	tempData.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, tempData); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func testEq(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
