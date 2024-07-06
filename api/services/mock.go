package services

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	plcl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	scl "brank.as/petnet/gunk/dsa/v2/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Mock struct {
	ListPartnersResp    *scl.ListServiceRequestResponse
	GetListPartnersResp *plcl.GetPartnerListResponse
}

func (m *Mock) SetGetPartnersResp(r *scl.ListServiceRequestResponse) {
	m.ListPartnersResp = r
}

func (m *Mock) SetGetPartnerListResponse(g *plcl.GetPartnerListResponse) {
	m.GetListPartnersResp = g
}

func (m *Mock) CreatePartners(ctx context.Context, in *spb.CreatePartnersRequest, opts ...grpc.CallOption) (*spb.CreatePartnersResponse, error) {
	return &spb.CreatePartnersResponse{}, nil
}

func (m *Mock) UpdatePartners(ctx context.Context, in *spb.UpdatePartnersRequest, opts ...grpc.CallOption) (*spb.UpdatePartnersResponse, error) {
	return &spb.UpdatePartnersResponse{}, nil
}

func (m *Mock) GetPartners(ctx context.Context, in *spb.GetPartnersRequest, opts ...grpc.CallOption) (*spb.GetPartnersResponse, error) {
	return &spb.GetPartnersResponse{}, nil
}

func (m *Mock) GetPartner(ctx context.Context, in *spb.GetPartnersRequest, opts ...grpc.CallOption) (*spb.GetPartnerResponse, error) {
	return &spb.GetPartnerResponse{
		Partner: &spb.Partner{},
	}, nil
}

func (m *Mock) DeletePartner(ctx context.Context, in *spb.DeletePartnerRequest, opts ...grpc.CallOption) (*spb.DeletePartnerResponse, error) {
	return &spb.DeletePartnerResponse{}, nil
}

func (m *Mock) ValidatePartnerAccess(ctx context.Context, in *spb.ValidatePartnerAccessRequest, opts ...grpc.CallOption) (*spb.ValidatePartnerAccessResponse, error) {
	return &spb.ValidatePartnerAccessResponse{}, nil
}

func (m *Mock) EnablePartner(ctx context.Context, in *spb.EnablePartnerRequest, opts ...grpc.CallOption) (*spb.EnablePartnerResponse, error) {
	return &spb.EnablePartnerResponse{}, nil
}

func (m *Mock) DisablePartner(ctx context.Context, in *spb.DisablePartnerRequest, opts ...grpc.CallOption) (*spb.DisablePartnerResponse, error) {
	return &spb.DisablePartnerResponse{}, nil
}

func (m *Mock) ListServiceRequest(ctx context.Context, in *scl.ListServiceRequestRequest, opts ...grpc.CallOption) (*scl.ListServiceRequestResponse, error) {
	if m.ListPartnersResp != nil {
		return m.ListPartnersResp, nil
	}
	return &scl.ListServiceRequestResponse{}, nil
}

func (m *Mock) AcceptServiceRequest(ctx context.Context, in *scl.ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) RejectServiceRequest(ctx context.Context, in *scl.ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) EnableServiceRequest(ctx context.Context, in *scl.ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) DisableServiceRequest(ctx context.Context, in *scl.ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) ValidateServiceAccess(ctx context.Context, in *scl.ValidateServiceAccessRequest, opts ...grpc.CallOption) (*scl.ValidateServiceAccessResponse, error) {
	return &scl.ValidateServiceAccessResponse{}, nil
}

func (m *Mock) AddServiceRequest(ctx context.Context, in *scl.AddServiceRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) ApplyServiceRequest(ctx context.Context, in *scl.ApplyServiceRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) AcceptUploadSvcRequest(ctx context.Context, in *scl.AcceptUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) RejectUploadSvcRequest(ctx context.Context, in *scl.RejectUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) ListUploadSvcRequest(ctx context.Context, in *scl.ListUploadSvcRequestRequest, opts ...grpc.CallOption) (*scl.ListUploadSvcRequestResponse, error) {
	return &scl.ListUploadSvcRequestResponse{}, nil
}

func (m *Mock) AddRemarkSvcRequest(ctx context.Context, in *scl.AddRemarkSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) AddUploadSvcRequest(ctx context.Context, in *scl.AddUploadSvcRequestRequest, opts ...grpc.CallOption) (*scl.AddUploadSvcRequestResponse, error) {
	return &scl.AddUploadSvcRequestResponse{}, nil
}

func (m *Mock) GetAllServiceRequest(ctx context.Context, in *scl.GetAllServiceRequestRequest, opts ...grpc.CallOption) (*scl.GetAllServiceRequestResponse, error) {
	return &scl.GetAllServiceRequestResponse{}, nil
}

func (m *Mock) RemoveUploadSvcRequest(ctx context.Context, in *scl.RemoveUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) SetStatusUploadSvcRequest(ctx context.Context, in *scl.SetStatusUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *Mock) UpdateUploadSvcRequest(ctx context.Context, in *scl.UpdateUploadSvcRequestRequest, opts ...grpc.CallOption) (*scl.UpdateUploadSvcRequestResponse, error) {
	return &scl.UpdateUploadSvcRequestResponse{}, nil
}

func (m *Mock) UpdateServiceRequestByOrgID(ctx context.Context, in *scl.UpdateServiceRequestByOrgIDRequest, opts ...grpc.CallOption) (*scl.UpdateServiceRequestByOrgIDResponse, error) {
	return &scl.UpdateServiceRequestByOrgIDResponse{}, nil
}

func (m *Mock) CreatePartnerList(ctx context.Context, in *plcl.CreatePartnerListRequest, opts ...grpc.CallOption) (*plcl.CreatePartnerListResponse, error) {
	return &plcl.CreatePartnerListResponse{}, nil
}

func (m *Mock) UpdatePartnerList(ctx context.Context, in *plcl.UpdatePartnerListRequest, opts ...grpc.CallOption) (*plcl.UpdatePartnerListResponse, error) {
	return &plcl.UpdatePartnerListResponse{}, nil
}

func (m *Mock) GetPartnerList(ctx context.Context, in *plcl.GetPartnerListRequest, opts ...grpc.CallOption) (*plcl.GetPartnerListResponse, error) {
	if m.GetListPartnersResp != nil {
		return m.GetListPartnersResp, nil
	}
	return &plcl.GetPartnerListResponse{}, nil
}

func (m *Mock) GetDSAPartnerList(ctx context.Context, in *plcl.DSAPartnerListRequest, opts ...grpc.CallOption) (*plcl.GetDSAPartnerListResponse, error) {
	return &plcl.GetDSAPartnerListResponse{}, nil
}

func (m *Mock) DeletePartnerList(ctx context.Context, in *plcl.DeletePartnerListRequest, opts ...grpc.CallOption) (*plcl.DeletePartnerListResponse, error) {
	return &plcl.DeletePartnerListResponse{}, nil
}

func (m *Mock) EnablePartnerList(ctx context.Context, in *plcl.EnablePartnerListRequest, opts ...grpc.CallOption) (*plcl.EnablePartnerListResponse, error) {
	return &plcl.EnablePartnerListResponse{}, nil
}

func (m *Mock) DisablePartnerList(ctx context.Context, in *plcl.DisablePartnerListRequest, opts ...grpc.CallOption) (*plcl.DisablePartnerListResponse, error) {
	return &plcl.DisablePartnerListResponse{}, nil
}

func (m *Mock) EnableMultiplePartnerList(ctx context.Context, in *plcl.EnableMultiplePartnerListRequest, opts ...grpc.CallOption) (*plcl.EnableMultiplePartnerListResponse, error) {
	return &plcl.EnableMultiplePartnerListResponse{}, nil
}

func (m *Mock) DisableMultiplePartnerList(ctx context.Context, in *plcl.DisableMultiplePartnerListRequest, opts ...grpc.CallOption) (*plcl.DisableMultiplePartnerListResponse, error) {
	return &plcl.DisableMultiplePartnerListResponse{}, nil
}

func (m *Mock) GetPartnerByStype(ctx context.Context, in *plcl.GetPartnerByStypeRequest, opts ...grpc.CallOption) (*plcl.GetPartnerByStypeResponse, error) {
	return &plcl.GetPartnerByStypeResponse{}, nil
}
