package profile

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
)

func (h *Svc) ListProfiles(ctx context.Context, req *ppb.ListProfilesRequest) (*ppb.ListProfilesResponse, error) {
	log := logging.FromContext(ctx)

	sortBy := "ASC"
	if req.GetSortBy() == ppb.SortBy_DESC {
		sortBy = "DESC"
	}

	sortByColumn := "date_applied"
	if req.GetSortByColumn() == ppb.SortByColumn_CompanyName {
		sortByColumn = "bus_info_company_name"
	}

	var riskScores []int32
	for _, rs := range req.GetRiskScore() {
		r := ppb.RiskScore_value[rs.String()]
		riskScores = append(riskScores, r)
	}

	var sts []int32
	for _, st := range req.GetStatus() {
		s := ppb.Status_value[st.String()]
		sts = append(sts, s)
	}

	var subDoc string
	if req.GetSubmittedDocument() == ppb.SubmittedDocument_NotSubmitted {
		subDoc = "not-submitted"
	} else if req.GetSubmittedDocument() == ppb.SubmittedDocument_Submitted {
		subDoc = "submitted"
	}

	var OrgType string
	if req.GetOrgType() == ppb.OrgType_PetNet {
		OrgType = "1"
	} else if req.GetOrgType() == ppb.OrgType_DSA {
		OrgType = "2"
	}

	pfs, err := h.ps.GetOrgProfiles(ctx,
		storage.FilterList{
			Limit:             req.GetLimit(),
			Offset:            req.GetOffset(),
			CompanyName:       req.GetCompanyName(),
			SortBy:            sortBy,
			SortByColumn:      sortByColumn,
			RiskScore:         riskScores,
			Status:            sts,
			SubmittedDocument: subDoc,
			OrgType:           OrgType,
			IsProvider:        req.GetIsProvider(),
		})
	if err != nil {
		logging.WithError(err, log).Error("listing profiles")
		return nil, status.Error(codes.Internal, "failed to list profiles")
	}

	var ppfs []*ppb.OrgProfile
	for _, pf := range pfs {
		ppfs = append(ppfs, storageToProto(&pf))
	}

	res := &ppb.ListProfilesResponse{}
	if len(ppfs) != 0 {
		tot := pfs[0].Count
		next := int(req.GetOffset()) + tot + 1
		res.Total = int32(tot)
		res.Next = int32(next)
	}
	res.Profiles = ppfs
	return res, nil
}
