// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package cashincashout

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

// CashInCashOutServiceClient is the client API for CashInCashOutService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CashInCashOutServiceClient interface {
	// Inquire
	CiCoInquire(ctx context.Context, in *CiCoInquireRequest, opts ...grpc.CallOption) (*CiCoInquireResponse, error)
	// Execute
	CiCoExecute(ctx context.Context, in *CiCoExecuteRequest, opts ...grpc.CallOption) (*CiCoExecuteResponse, error)
	// Retry
	CiCoRetry(ctx context.Context, in *CiCoRetryRequest, opts ...grpc.CallOption) (*CiCoRetryResponse, error)
	// OTP Confirm
	CiCoOTPConfirm(ctx context.Context, in *CiCoOTPConfirmRequest, opts ...grpc.CallOption) (*CiCoOTPConfirmResponse, error)
	// Validate
	CiCoValidate(ctx context.Context, in *CiCoValidateRequest, opts ...grpc.CallOption) (*CiCoValidateResponse, error)
	// Get Transaction List for CI/CO
	CICOTransactList(ctx context.Context, in *CICOTransactListRequest, opts ...grpc.CallOption) (*CICOTransactListResponse, error)
}

type cashInCashOutServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCashInCashOutServiceClient(cc grpc.ClientConnInterface) CashInCashOutServiceClient {
	return &cashInCashOutServiceClient{cc}
}

func (c *cashInCashOutServiceClient) CiCoInquire(ctx context.Context, in *CiCoInquireRequest, opts ...grpc.CallOption) (*CiCoInquireResponse, error) {
	out := new(CiCoInquireResponse)
	err := c.cc.Invoke(ctx, "/cashincashout.CashInCashOutService/CiCoInquire", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cashInCashOutServiceClient) CiCoExecute(ctx context.Context, in *CiCoExecuteRequest, opts ...grpc.CallOption) (*CiCoExecuteResponse, error) {
	out := new(CiCoExecuteResponse)
	err := c.cc.Invoke(ctx, "/cashincashout.CashInCashOutService/CiCoExecute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cashInCashOutServiceClient) CiCoRetry(ctx context.Context, in *CiCoRetryRequest, opts ...grpc.CallOption) (*CiCoRetryResponse, error) {
	out := new(CiCoRetryResponse)
	err := c.cc.Invoke(ctx, "/cashincashout.CashInCashOutService/CiCoRetry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cashInCashOutServiceClient) CiCoOTPConfirm(ctx context.Context, in *CiCoOTPConfirmRequest, opts ...grpc.CallOption) (*CiCoOTPConfirmResponse, error) {
	out := new(CiCoOTPConfirmResponse)
	err := c.cc.Invoke(ctx, "/cashincashout.CashInCashOutService/CiCoOTPConfirm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cashInCashOutServiceClient) CiCoValidate(ctx context.Context, in *CiCoValidateRequest, opts ...grpc.CallOption) (*CiCoValidateResponse, error) {
	out := new(CiCoValidateResponse)
	err := c.cc.Invoke(ctx, "/cashincashout.CashInCashOutService/CiCoValidate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cashInCashOutServiceClient) CICOTransactList(ctx context.Context, in *CICOTransactListRequest, opts ...grpc.CallOption) (*CICOTransactListResponse, error) {
	out := new(CICOTransactListResponse)
	err := c.cc.Invoke(ctx, "/cashincashout.CashInCashOutService/CICOTransactList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CashInCashOutServiceServer is the server API for CashInCashOutService service.
// All implementations must embed UnimplementedCashInCashOutServiceServer
// for forward compatibility
type CashInCashOutServiceServer interface {
	// Inquire
	CiCoInquire(context.Context, *CiCoInquireRequest) (*CiCoInquireResponse, error)
	// Execute
	CiCoExecute(context.Context, *CiCoExecuteRequest) (*CiCoExecuteResponse, error)
	// Retry
	CiCoRetry(context.Context, *CiCoRetryRequest) (*CiCoRetryResponse, error)
	// OTP Confirm
	CiCoOTPConfirm(context.Context, *CiCoOTPConfirmRequest) (*CiCoOTPConfirmResponse, error)
	// Validate
	CiCoValidate(context.Context, *CiCoValidateRequest) (*CiCoValidateResponse, error)
	// Get Transaction List for CI/CO
	CICOTransactList(context.Context, *CICOTransactListRequest) (*CICOTransactListResponse, error)
	mustEmbedUnimplementedCashInCashOutServiceServer()
}

// UnimplementedCashInCashOutServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCashInCashOutServiceServer struct{}

func (UnimplementedCashInCashOutServiceServer) CiCoInquire(context.Context, *CiCoInquireRequest) (*CiCoInquireResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CiCoInquire not implemented")
}

func (UnimplementedCashInCashOutServiceServer) CiCoExecute(context.Context, *CiCoExecuteRequest) (*CiCoExecuteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CiCoExecute not implemented")
}

func (UnimplementedCashInCashOutServiceServer) CiCoRetry(context.Context, *CiCoRetryRequest) (*CiCoRetryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CiCoRetry not implemented")
}

func (UnimplementedCashInCashOutServiceServer) CiCoOTPConfirm(context.Context, *CiCoOTPConfirmRequest) (*CiCoOTPConfirmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CiCoOTPConfirm not implemented")
}

func (UnimplementedCashInCashOutServiceServer) CiCoValidate(context.Context, *CiCoValidateRequest) (*CiCoValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CiCoValidate not implemented")
}

func (UnimplementedCashInCashOutServiceServer) CICOTransactList(context.Context, *CICOTransactListRequest) (*CICOTransactListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CICOTransactList not implemented")
}
func (UnimplementedCashInCashOutServiceServer) mustEmbedUnimplementedCashInCashOutServiceServer() {}

// UnsafeCashInCashOutServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CashInCashOutServiceServer will
// result in compilation errors.
type UnsafeCashInCashOutServiceServer interface {
	mustEmbedUnimplementedCashInCashOutServiceServer()
}

func RegisterCashInCashOutServiceServer(s grpc.ServiceRegistrar, srv CashInCashOutServiceServer) {
	s.RegisterService(&CashInCashOutService_ServiceDesc, srv)
}

func _CashInCashOutService_CiCoInquire_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CiCoInquireRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CashInCashOutServiceServer).CiCoInquire(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cashincashout.CashInCashOutService/CiCoInquire",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CashInCashOutServiceServer).CiCoInquire(ctx, req.(*CiCoInquireRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CashInCashOutService_CiCoExecute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CiCoExecuteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CashInCashOutServiceServer).CiCoExecute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cashincashout.CashInCashOutService/CiCoExecute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CashInCashOutServiceServer).CiCoExecute(ctx, req.(*CiCoExecuteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CashInCashOutService_CiCoRetry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CiCoRetryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CashInCashOutServiceServer).CiCoRetry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cashincashout.CashInCashOutService/CiCoRetry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CashInCashOutServiceServer).CiCoRetry(ctx, req.(*CiCoRetryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CashInCashOutService_CiCoOTPConfirm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CiCoOTPConfirmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CashInCashOutServiceServer).CiCoOTPConfirm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cashincashout.CashInCashOutService/CiCoOTPConfirm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CashInCashOutServiceServer).CiCoOTPConfirm(ctx, req.(*CiCoOTPConfirmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CashInCashOutService_CiCoValidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CiCoValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CashInCashOutServiceServer).CiCoValidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cashincashout.CashInCashOutService/CiCoValidate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CashInCashOutServiceServer).CiCoValidate(ctx, req.(*CiCoValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CashInCashOutService_CICOTransactList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CICOTransactListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CashInCashOutServiceServer).CICOTransactList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cashincashout.CashInCashOutService/CICOTransactList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CashInCashOutServiceServer).CICOTransactList(ctx, req.(*CICOTransactListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CashInCashOutService_ServiceDesc is the grpc.ServiceDesc for CashInCashOutService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CashInCashOutService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cashincashout.CashInCashOutService",
	HandlerType: (*CashInCashOutServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CiCoInquire",
			Handler:    _CashInCashOutService_CiCoInquire_Handler,
		},
		{
			MethodName: "CiCoExecute",
			Handler:    _CashInCashOutService_CiCoExecute_Handler,
		},
		{
			MethodName: "CiCoRetry",
			Handler:    _CashInCashOutService_CiCoRetry_Handler,
		},
		{
			MethodName: "CiCoOTPConfirm",
			Handler:    _CashInCashOutService_CiCoOTPConfirm_Handler,
		},
		{
			MethodName: "CiCoValidate",
			Handler:    _CashInCashOutService_CiCoValidate_Handler,
		},
		{
			MethodName: "CICOTransactList",
			Handler:    _CashInCashOutService_CICOTransactList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "brank.as/petnet/gunk/drp/v1/cashincashout/all.proto",
}
