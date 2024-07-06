// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package service

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ServiceServiceClient is the client API for ServiceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceServiceClient interface {
	// Add DSA service request partners.
	AddServiceRequest(ctx context.Context, in *AddServiceRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Apply for DSA service.
	ApplyServiceRequest(ctx context.Context, in *ApplyServiceRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// List DSA Services.
	ListServiceRequest(ctx context.Context, in *ListServiceRequestRequest, opts ...grpc.CallOption) (*ListServiceRequestResponse, error)
	// Accept for DSA service.
	AcceptServiceRequest(ctx context.Context, in *ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Reject for DSA service.
	RejectServiceRequest(ctx context.Context, in *ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Enable for DSA service.
	EnableServiceRequest(ctx context.Context, in *ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Disable for DSA service.
	DisableServiceRequest(ctx context.Context, in *ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Validate for DSA service.
	ValidateServiceAccess(ctx context.Context, in *ValidateServiceAccessRequest, opts ...grpc.CallOption) (*ValidateServiceAccessResponse, error)
	// Add DSA upload service request partners.
	AddUploadSvcRequest(ctx context.Context, in *AddUploadSvcRequestRequest, opts ...grpc.CallOption) (*AddUploadSvcRequestResponse, error)
	// update for DSA upload service request
	UpdateUploadSvcRequest(ctx context.Context, in *UpdateUploadSvcRequestRequest, opts ...grpc.CallOption) (*UpdateUploadSvcRequestResponse, error)
	// Accept for DSA upload service request
	AcceptUploadSvcRequest(ctx context.Context, in *AcceptUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Reject for DSA upload service request
	RejectUploadSvcRequest(ctx context.Context, in *RejectUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Remove for DSA upload service request
	RemoveUploadSvcRequest(ctx context.Context, in *RemoveUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// List DSA upload Services request.
	ListUploadSvcRequest(ctx context.Context, in *ListUploadSvcRequestRequest, opts ...grpc.CallOption) (*ListUploadSvcRequestResponse, error)
	// set status for DSA upload service request
	SetStatusUploadSvcRequest(ctx context.Context, in *SetStatusUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Add Remark for DSA upload service request
	AddRemarkSvcRequest(ctx context.Context, in *AddRemarkSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// get all DSA Services.
	GetAllServiceRequest(ctx context.Context, in *GetAllServiceRequestRequest, opts ...grpc.CallOption) (*GetAllServiceRequestResponse, error)
	// Update service request orgid.
	UpdateServiceRequestByOrgID(ctx context.Context, in *UpdateServiceRequestByOrgIDRequest, opts ...grpc.CallOption) (*UpdateServiceRequestByOrgIDResponse, error)
}

type serviceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceServiceClient(cc grpc.ClientConnInterface) ServiceServiceClient {
	return &serviceServiceClient{cc}
}

func (c *serviceServiceClient) AddServiceRequest(ctx context.Context, in *AddServiceRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/AddServiceRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) ApplyServiceRequest(ctx context.Context, in *ApplyServiceRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/ApplyServiceRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) ListServiceRequest(ctx context.Context, in *ListServiceRequestRequest, opts ...grpc.CallOption) (*ListServiceRequestResponse, error) {
	out := new(ListServiceRequestResponse)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/ListServiceRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) AcceptServiceRequest(ctx context.Context, in *ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/AcceptServiceRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) RejectServiceRequest(ctx context.Context, in *ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/RejectServiceRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) EnableServiceRequest(ctx context.Context, in *ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/EnableServiceRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) DisableServiceRequest(ctx context.Context, in *ServiceStatusRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/DisableServiceRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) ValidateServiceAccess(ctx context.Context, in *ValidateServiceAccessRequest, opts ...grpc.CallOption) (*ValidateServiceAccessResponse, error) {
	out := new(ValidateServiceAccessResponse)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/ValidateServiceAccess", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) AddUploadSvcRequest(ctx context.Context, in *AddUploadSvcRequestRequest, opts ...grpc.CallOption) (*AddUploadSvcRequestResponse, error) {
	out := new(AddUploadSvcRequestResponse)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/AddUploadSvcRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) UpdateUploadSvcRequest(ctx context.Context, in *UpdateUploadSvcRequestRequest, opts ...grpc.CallOption) (*UpdateUploadSvcRequestResponse, error) {
	out := new(UpdateUploadSvcRequestResponse)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/UpdateUploadSvcRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) AcceptUploadSvcRequest(ctx context.Context, in *AcceptUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/AcceptUploadSvcRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) RejectUploadSvcRequest(ctx context.Context, in *RejectUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/RejectUploadSvcRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) RemoveUploadSvcRequest(ctx context.Context, in *RemoveUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/RemoveUploadSvcRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) ListUploadSvcRequest(ctx context.Context, in *ListUploadSvcRequestRequest, opts ...grpc.CallOption) (*ListUploadSvcRequestResponse, error) {
	out := new(ListUploadSvcRequestResponse)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/ListUploadSvcRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) SetStatusUploadSvcRequest(ctx context.Context, in *SetStatusUploadSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/SetStatusUploadSvcRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) AddRemarkSvcRequest(ctx context.Context, in *AddRemarkSvcRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/AddRemarkSvcRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) GetAllServiceRequest(ctx context.Context, in *GetAllServiceRequestRequest, opts ...grpc.CallOption) (*GetAllServiceRequestResponse, error) {
	out := new(GetAllServiceRequestResponse)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/GetAllServiceRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceServiceClient) UpdateServiceRequestByOrgID(ctx context.Context, in *UpdateServiceRequestByOrgIDRequest, opts ...grpc.CallOption) (*UpdateServiceRequestByOrgIDResponse, error) {
	out := new(UpdateServiceRequestByOrgIDResponse)
	err := c.cc.Invoke(ctx, "/petnet.v2.service.ServiceService/UpdateServiceRequestByOrgID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServiceServer is the server API for ServiceService service.
// All implementations must embed UnimplementedServiceServiceServer
// for forward compatibility
type ServiceServiceServer interface {
	// Add DSA service request partners.
	AddServiceRequest(context.Context, *AddServiceRequestRequest) (*emptypb.Empty, error)
	// Apply for DSA service.
	ApplyServiceRequest(context.Context, *ApplyServiceRequestRequest) (*emptypb.Empty, error)
	// List DSA Services.
	ListServiceRequest(context.Context, *ListServiceRequestRequest) (*ListServiceRequestResponse, error)
	// Accept for DSA service.
	AcceptServiceRequest(context.Context, *ServiceStatusRequestRequest) (*emptypb.Empty, error)
	// Reject for DSA service.
	RejectServiceRequest(context.Context, *ServiceStatusRequestRequest) (*emptypb.Empty, error)
	// Enable for DSA service.
	EnableServiceRequest(context.Context, *ServiceStatusRequestRequest) (*emptypb.Empty, error)
	// Disable for DSA service.
	DisableServiceRequest(context.Context, *ServiceStatusRequestRequest) (*emptypb.Empty, error)
	// Validate for DSA service.
	ValidateServiceAccess(context.Context, *ValidateServiceAccessRequest) (*ValidateServiceAccessResponse, error)
	// Add DSA upload service request partners.
	AddUploadSvcRequest(context.Context, *AddUploadSvcRequestRequest) (*AddUploadSvcRequestResponse, error)
	// update for DSA upload service request
	UpdateUploadSvcRequest(context.Context, *UpdateUploadSvcRequestRequest) (*UpdateUploadSvcRequestResponse, error)
	// Accept for DSA upload service request
	AcceptUploadSvcRequest(context.Context, *AcceptUploadSvcRequestRequest) (*emptypb.Empty, error)
	// Reject for DSA upload service request
	RejectUploadSvcRequest(context.Context, *RejectUploadSvcRequestRequest) (*emptypb.Empty, error)
	// Remove for DSA upload service request
	RemoveUploadSvcRequest(context.Context, *RemoveUploadSvcRequestRequest) (*emptypb.Empty, error)
	// List DSA upload Services request.
	ListUploadSvcRequest(context.Context, *ListUploadSvcRequestRequest) (*ListUploadSvcRequestResponse, error)
	// set status for DSA upload service request
	SetStatusUploadSvcRequest(context.Context, *SetStatusUploadSvcRequestRequest) (*emptypb.Empty, error)
	// Add Remark for DSA upload service request
	AddRemarkSvcRequest(context.Context, *AddRemarkSvcRequestRequest) (*emptypb.Empty, error)
	// get all DSA Services.
	GetAllServiceRequest(context.Context, *GetAllServiceRequestRequest) (*GetAllServiceRequestResponse, error)
	// Update service request orgid.
	UpdateServiceRequestByOrgID(context.Context, *UpdateServiceRequestByOrgIDRequest) (*UpdateServiceRequestByOrgIDResponse, error)
	mustEmbedUnimplementedServiceServiceServer()
}

// UnimplementedServiceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServiceServer struct{}

func (UnimplementedServiceServiceServer) AddServiceRequest(context.Context, *AddServiceRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddServiceRequest not implemented")
}

func (UnimplementedServiceServiceServer) ApplyServiceRequest(context.Context, *ApplyServiceRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApplyServiceRequest not implemented")
}

func (UnimplementedServiceServiceServer) ListServiceRequest(context.Context, *ListServiceRequestRequest) (*ListServiceRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListServiceRequest not implemented")
}

func (UnimplementedServiceServiceServer) AcceptServiceRequest(context.Context, *ServiceStatusRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptServiceRequest not implemented")
}

func (UnimplementedServiceServiceServer) RejectServiceRequest(context.Context, *ServiceStatusRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RejectServiceRequest not implemented")
}

func (UnimplementedServiceServiceServer) EnableServiceRequest(context.Context, *ServiceStatusRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableServiceRequest not implemented")
}

func (UnimplementedServiceServiceServer) DisableServiceRequest(context.Context, *ServiceStatusRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableServiceRequest not implemented")
}

func (UnimplementedServiceServiceServer) ValidateServiceAccess(context.Context, *ValidateServiceAccessRequest) (*ValidateServiceAccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateServiceAccess not implemented")
}

func (UnimplementedServiceServiceServer) AddUploadSvcRequest(context.Context, *AddUploadSvcRequestRequest) (*AddUploadSvcRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUploadSvcRequest not implemented")
}

func (UnimplementedServiceServiceServer) UpdateUploadSvcRequest(context.Context, *UpdateUploadSvcRequestRequest) (*UpdateUploadSvcRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUploadSvcRequest not implemented")
}

func (UnimplementedServiceServiceServer) AcceptUploadSvcRequest(context.Context, *AcceptUploadSvcRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptUploadSvcRequest not implemented")
}

func (UnimplementedServiceServiceServer) RejectUploadSvcRequest(context.Context, *RejectUploadSvcRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RejectUploadSvcRequest not implemented")
}

func (UnimplementedServiceServiceServer) RemoveUploadSvcRequest(context.Context, *RemoveUploadSvcRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveUploadSvcRequest not implemented")
}

func (UnimplementedServiceServiceServer) ListUploadSvcRequest(context.Context, *ListUploadSvcRequestRequest) (*ListUploadSvcRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUploadSvcRequest not implemented")
}

func (UnimplementedServiceServiceServer) SetStatusUploadSvcRequest(context.Context, *SetStatusUploadSvcRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStatusUploadSvcRequest not implemented")
}

func (UnimplementedServiceServiceServer) AddRemarkSvcRequest(context.Context, *AddRemarkSvcRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRemarkSvcRequest not implemented")
}

func (UnimplementedServiceServiceServer) GetAllServiceRequest(context.Context, *GetAllServiceRequestRequest) (*GetAllServiceRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllServiceRequest not implemented")
}

func (UnimplementedServiceServiceServer) UpdateServiceRequestByOrgID(context.Context, *UpdateServiceRequestByOrgIDRequest) (*UpdateServiceRequestByOrgIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateServiceRequestByOrgID not implemented")
}
func (UnimplementedServiceServiceServer) mustEmbedUnimplementedServiceServiceServer() {}

// UnsafeServiceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServiceServer will
// result in compilation errors.
type UnsafeServiceServiceServer interface {
	mustEmbedUnimplementedServiceServiceServer()
}

func RegisterServiceServiceServer(s grpc.ServiceRegistrar, srv ServiceServiceServer) {
	s.RegisterService(&ServiceService_ServiceDesc, srv)
}

func _ServiceService_AddServiceRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddServiceRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).AddServiceRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/AddServiceRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).AddServiceRequest(ctx, req.(*AddServiceRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_ApplyServiceRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplyServiceRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).ApplyServiceRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/ApplyServiceRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).ApplyServiceRequest(ctx, req.(*ApplyServiceRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_ListServiceRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListServiceRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).ListServiceRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/ListServiceRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).ListServiceRequest(ctx, req.(*ListServiceRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_AcceptServiceRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceStatusRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).AcceptServiceRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/AcceptServiceRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).AcceptServiceRequest(ctx, req.(*ServiceStatusRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_RejectServiceRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceStatusRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).RejectServiceRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/RejectServiceRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).RejectServiceRequest(ctx, req.(*ServiceStatusRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_EnableServiceRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceStatusRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).EnableServiceRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/EnableServiceRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).EnableServiceRequest(ctx, req.(*ServiceStatusRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_DisableServiceRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceStatusRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).DisableServiceRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/DisableServiceRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).DisableServiceRequest(ctx, req.(*ServiceStatusRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_ValidateServiceAccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateServiceAccessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).ValidateServiceAccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/ValidateServiceAccess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).ValidateServiceAccess(ctx, req.(*ValidateServiceAccessRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_AddUploadSvcRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUploadSvcRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).AddUploadSvcRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/AddUploadSvcRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).AddUploadSvcRequest(ctx, req.(*AddUploadSvcRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_UpdateUploadSvcRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUploadSvcRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).UpdateUploadSvcRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/UpdateUploadSvcRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).UpdateUploadSvcRequest(ctx, req.(*UpdateUploadSvcRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_AcceptUploadSvcRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptUploadSvcRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).AcceptUploadSvcRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/AcceptUploadSvcRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).AcceptUploadSvcRequest(ctx, req.(*AcceptUploadSvcRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_RejectUploadSvcRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RejectUploadSvcRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).RejectUploadSvcRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/RejectUploadSvcRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).RejectUploadSvcRequest(ctx, req.(*RejectUploadSvcRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_RemoveUploadSvcRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveUploadSvcRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).RemoveUploadSvcRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/RemoveUploadSvcRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).RemoveUploadSvcRequest(ctx, req.(*RemoveUploadSvcRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_ListUploadSvcRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUploadSvcRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).ListUploadSvcRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/ListUploadSvcRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).ListUploadSvcRequest(ctx, req.(*ListUploadSvcRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_SetStatusUploadSvcRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetStatusUploadSvcRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).SetStatusUploadSvcRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/SetStatusUploadSvcRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).SetStatusUploadSvcRequest(ctx, req.(*SetStatusUploadSvcRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_AddRemarkSvcRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRemarkSvcRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).AddRemarkSvcRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/AddRemarkSvcRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).AddRemarkSvcRequest(ctx, req.(*AddRemarkSvcRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_GetAllServiceRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllServiceRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).GetAllServiceRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/GetAllServiceRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).GetAllServiceRequest(ctx, req.(*GetAllServiceRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceService_UpdateServiceRequestByOrgID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateServiceRequestByOrgIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServiceServer).UpdateServiceRequestByOrgID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.service.ServiceService/UpdateServiceRequestByOrgID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServiceServer).UpdateServiceRequestByOrgID(ctx, req.(*UpdateServiceRequestByOrgIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ServiceService_ServiceDesc is the grpc.ServiceDesc for ServiceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServiceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "petnet.v2.service.ServiceService",
	HandlerType: (*ServiceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddServiceRequest",
			Handler:    _ServiceService_AddServiceRequest_Handler,
		},
		{
			MethodName: "ApplyServiceRequest",
			Handler:    _ServiceService_ApplyServiceRequest_Handler,
		},
		{
			MethodName: "ListServiceRequest",
			Handler:    _ServiceService_ListServiceRequest_Handler,
		},
		{
			MethodName: "AcceptServiceRequest",
			Handler:    _ServiceService_AcceptServiceRequest_Handler,
		},
		{
			MethodName: "RejectServiceRequest",
			Handler:    _ServiceService_RejectServiceRequest_Handler,
		},
		{
			MethodName: "EnableServiceRequest",
			Handler:    _ServiceService_EnableServiceRequest_Handler,
		},
		{
			MethodName: "DisableServiceRequest",
			Handler:    _ServiceService_DisableServiceRequest_Handler,
		},
		{
			MethodName: "ValidateServiceAccess",
			Handler:    _ServiceService_ValidateServiceAccess_Handler,
		},
		{
			MethodName: "AddUploadSvcRequest",
			Handler:    _ServiceService_AddUploadSvcRequest_Handler,
		},
		{
			MethodName: "UpdateUploadSvcRequest",
			Handler:    _ServiceService_UpdateUploadSvcRequest_Handler,
		},
		{
			MethodName: "AcceptUploadSvcRequest",
			Handler:    _ServiceService_AcceptUploadSvcRequest_Handler,
		},
		{
			MethodName: "RejectUploadSvcRequest",
			Handler:    _ServiceService_RejectUploadSvcRequest_Handler,
		},
		{
			MethodName: "RemoveUploadSvcRequest",
			Handler:    _ServiceService_RemoveUploadSvcRequest_Handler,
		},
		{
			MethodName: "ListUploadSvcRequest",
			Handler:    _ServiceService_ListUploadSvcRequest_Handler,
		},
		{
			MethodName: "SetStatusUploadSvcRequest",
			Handler:    _ServiceService_SetStatusUploadSvcRequest_Handler,
		},
		{
			MethodName: "AddRemarkSvcRequest",
			Handler:    _ServiceService_AddRemarkSvcRequest_Handler,
		},
		{
			MethodName: "GetAllServiceRequest",
			Handler:    _ServiceService_GetAllServiceRequest_Handler,
		},
		{
			MethodName: "UpdateServiceRequestByOrgID",
			Handler:    _ServiceService_UpdateServiceRequestByOrgID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "brank.as/petnet/gunk/dsa/v2/service/all.proto",
}
