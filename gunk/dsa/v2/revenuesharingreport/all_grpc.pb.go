// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package revenuesharingreport

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RevenueSharingReportServiceClient is the client API for RevenueSharingReportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RevenueSharingReportServiceClient interface {
	// Create Revenue Sharing Report
	CreateRevenueSharingReport(ctx context.Context, in *CreateRevenueSharingReportRequest, opts ...grpc.CallOption) (*CreateRevenueSharingReportResponse, error)
	// List Revenue Sharing Report
	GetRevenueSharingReportList(ctx context.Context, in *GetRevenueSharingReportListRequest, opts ...grpc.CallOption) (*GetRevenueSharingReportListResponse, error)
	// Update Revenue Sharing Report
	UpdateRevenueSharingReport(ctx context.Context, in *UpdateRevenueSharingReportRequest, opts ...grpc.CallOption) (*UpdateRevenueSharingReportResponse, error)
}

type revenueSharingReportServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRevenueSharingReportServiceClient(cc grpc.ClientConnInterface) RevenueSharingReportServiceClient {
	return &revenueSharingReportServiceClient{cc}
}

func (c *revenueSharingReportServiceClient) CreateRevenueSharingReport(ctx context.Context, in *CreateRevenueSharingReportRequest, opts ...grpc.CallOption) (*CreateRevenueSharingReportResponse, error) {
	out := new(CreateRevenueSharingReportResponse)
	err := c.cc.Invoke(ctx, "/petnet.v2.revenuesharingreport.RevenueSharingReportService/CreateRevenueSharingReport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *revenueSharingReportServiceClient) GetRevenueSharingReportList(ctx context.Context, in *GetRevenueSharingReportListRequest, opts ...grpc.CallOption) (*GetRevenueSharingReportListResponse, error) {
	out := new(GetRevenueSharingReportListResponse)
	err := c.cc.Invoke(ctx, "/petnet.v2.revenuesharingreport.RevenueSharingReportService/GetRevenueSharingReportList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *revenueSharingReportServiceClient) UpdateRevenueSharingReport(ctx context.Context, in *UpdateRevenueSharingReportRequest, opts ...grpc.CallOption) (*UpdateRevenueSharingReportResponse, error) {
	out := new(UpdateRevenueSharingReportResponse)
	err := c.cc.Invoke(ctx, "/petnet.v2.revenuesharingreport.RevenueSharingReportService/UpdateRevenueSharingReport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RevenueSharingReportServiceServer is the server API for RevenueSharingReportService service.
// All implementations must embed UnimplementedRevenueSharingReportServiceServer
// for forward compatibility
type RevenueSharingReportServiceServer interface {
	// Create Revenue Sharing Report
	CreateRevenueSharingReport(context.Context, *CreateRevenueSharingReportRequest) (*CreateRevenueSharingReportResponse, error)
	// List Revenue Sharing Report
	GetRevenueSharingReportList(context.Context, *GetRevenueSharingReportListRequest) (*GetRevenueSharingReportListResponse, error)
	// Update Revenue Sharing Report
	UpdateRevenueSharingReport(context.Context, *UpdateRevenueSharingReportRequest) (*UpdateRevenueSharingReportResponse, error)
	mustEmbedUnimplementedRevenueSharingReportServiceServer()
}

// UnimplementedRevenueSharingReportServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRevenueSharingReportServiceServer struct{}

func (UnimplementedRevenueSharingReportServiceServer) CreateRevenueSharingReport(context.Context, *CreateRevenueSharingReportRequest) (*CreateRevenueSharingReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRevenueSharingReport not implemented")
}

func (UnimplementedRevenueSharingReportServiceServer) GetRevenueSharingReportList(context.Context, *GetRevenueSharingReportListRequest) (*GetRevenueSharingReportListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRevenueSharingReportList not implemented")
}

func (UnimplementedRevenueSharingReportServiceServer) UpdateRevenueSharingReport(context.Context, *UpdateRevenueSharingReportRequest) (*UpdateRevenueSharingReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRevenueSharingReport not implemented")
}

func (UnimplementedRevenueSharingReportServiceServer) mustEmbedUnimplementedRevenueSharingReportServiceServer() {
}

// UnsafeRevenueSharingReportServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RevenueSharingReportServiceServer will
// result in compilation errors.
type UnsafeRevenueSharingReportServiceServer interface {
	mustEmbedUnimplementedRevenueSharingReportServiceServer()
}

func RegisterRevenueSharingReportServiceServer(s grpc.ServiceRegistrar, srv RevenueSharingReportServiceServer) {
	s.RegisterService(&RevenueSharingReportService_ServiceDesc, srv)
}

func _RevenueSharingReportService_CreateRevenueSharingReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRevenueSharingReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RevenueSharingReportServiceServer).CreateRevenueSharingReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.revenuesharingreport.RevenueSharingReportService/CreateRevenueSharingReport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RevenueSharingReportServiceServer).CreateRevenueSharingReport(ctx, req.(*CreateRevenueSharingReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RevenueSharingReportService_GetRevenueSharingReportList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRevenueSharingReportListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RevenueSharingReportServiceServer).GetRevenueSharingReportList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.revenuesharingreport.RevenueSharingReportService/GetRevenueSharingReportList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RevenueSharingReportServiceServer).GetRevenueSharingReportList(ctx, req.(*GetRevenueSharingReportListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RevenueSharingReportService_UpdateRevenueSharingReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRevenueSharingReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RevenueSharingReportServiceServer).UpdateRevenueSharingReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/petnet.v2.revenuesharingreport.RevenueSharingReportService/UpdateRevenueSharingReport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RevenueSharingReportServiceServer).UpdateRevenueSharingReport(ctx, req.(*UpdateRevenueSharingReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RevenueSharingReportService_ServiceDesc is the grpc.ServiceDesc for RevenueSharingReportService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RevenueSharingReportService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "petnet.v2.revenuesharingreport.RevenueSharingReportService",
	HandlerType: (*RevenueSharingReportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRevenueSharingReport",
			Handler:    _RevenueSharingReportService_CreateRevenueSharingReport_Handler,
		},
		{
			MethodName: "GetRevenueSharingReportList",
			Handler:    _RevenueSharingReportService_GetRevenueSharingReportList_Handler,
		},
		{
			MethodName: "UpdateRevenueSharingReport",
			Handler:    _RevenueSharingReportService_UpdateRevenueSharingReport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "brank.as/petnet/gunk/dsa/v2/revenuesharingreport/all.proto",
}
