// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package mfa

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

// MFAServiceClient is the client API for MFAService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MFAServiceClient interface {
	// Get MFA List by User ID.
	GetRegisteredMFA(ctx context.Context, in *GetRegisteredMFARequest, opts ...grpc.CallOption) (*GetRegisteredMFAResponse, error)
	// Enable MFA source.
	EnableMFA(ctx context.Context, in *EnableMFARequest, opts ...grpc.CallOption) (*EnableMFAResponse, error)
	// Disable MFA source.
	DisableMFA(ctx context.Context, in *DisableMFARequest, opts ...grpc.CallOption) (*DisableMFAResponse, error)
}

type mFAServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMFAServiceClient(cc grpc.ClientConnInterface) MFAServiceClient {
	return &mFAServiceClient{cc}
}

func (c *mFAServiceClient) GetRegisteredMFA(ctx context.Context, in *GetRegisteredMFARequest, opts ...grpc.CallOption) (*GetRegisteredMFAResponse, error) {
	out := new(GetRegisteredMFAResponse)
	err := c.cc.Invoke(ctx, "/mfa.MFAService/GetRegisteredMFA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mFAServiceClient) EnableMFA(ctx context.Context, in *EnableMFARequest, opts ...grpc.CallOption) (*EnableMFAResponse, error) {
	out := new(EnableMFAResponse)
	err := c.cc.Invoke(ctx, "/mfa.MFAService/EnableMFA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mFAServiceClient) DisableMFA(ctx context.Context, in *DisableMFARequest, opts ...grpc.CallOption) (*DisableMFAResponse, error) {
	out := new(DisableMFAResponse)
	err := c.cc.Invoke(ctx, "/mfa.MFAService/DisableMFA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MFAServiceServer is the server API for MFAService service.
// All implementations must embed UnimplementedMFAServiceServer
// for forward compatibility
type MFAServiceServer interface {
	// Get MFA List by User ID.
	GetRegisteredMFA(context.Context, *GetRegisteredMFARequest) (*GetRegisteredMFAResponse, error)
	// Enable MFA source.
	EnableMFA(context.Context, *EnableMFARequest) (*EnableMFAResponse, error)
	// Disable MFA source.
	DisableMFA(context.Context, *DisableMFARequest) (*DisableMFAResponse, error)
	mustEmbedUnimplementedMFAServiceServer()
}

// UnimplementedMFAServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMFAServiceServer struct{}

func (UnimplementedMFAServiceServer) GetRegisteredMFA(context.Context, *GetRegisteredMFARequest) (*GetRegisteredMFAResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRegisteredMFA not implemented")
}

func (UnimplementedMFAServiceServer) EnableMFA(context.Context, *EnableMFARequest) (*EnableMFAResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableMFA not implemented")
}

func (UnimplementedMFAServiceServer) DisableMFA(context.Context, *DisableMFARequest) (*DisableMFAResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableMFA not implemented")
}
func (UnimplementedMFAServiceServer) mustEmbedUnimplementedMFAServiceServer() {}

// UnsafeMFAServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MFAServiceServer will
// result in compilation errors.
type UnsafeMFAServiceServer interface {
	mustEmbedUnimplementedMFAServiceServer()
}

func RegisterMFAServiceServer(s grpc.ServiceRegistrar, srv MFAServiceServer) {
	s.RegisterService(&MFAService_ServiceDesc, srv)
}

func _MFAService_GetRegisteredMFA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRegisteredMFARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MFAServiceServer).GetRegisteredMFA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mfa.MFAService/GetRegisteredMFA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MFAServiceServer).GetRegisteredMFA(ctx, req.(*GetRegisteredMFARequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MFAService_EnableMFA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnableMFARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MFAServiceServer).EnableMFA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mfa.MFAService/EnableMFA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MFAServiceServer).EnableMFA(ctx, req.(*EnableMFARequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MFAService_DisableMFA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisableMFARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MFAServiceServer).DisableMFA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mfa.MFAService/DisableMFA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MFAServiceServer).DisableMFA(ctx, req.(*DisableMFARequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MFAService_ServiceDesc is the grpc.ServiceDesc for MFAService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MFAService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mfa.MFAService",
	HandlerType: (*MFAServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRegisteredMFA",
			Handler:    _MFAService_GetRegisteredMFA_Handler,
		},
		{
			MethodName: "EnableMFA",
			Handler:    _MFAService_EnableMFA_Handler,
		},
		{
			MethodName: "DisableMFA",
			Handler:    _MFAService_DisableMFA_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "brank.as/rbac/gunk/v1/mfa/all.proto",
}

// MFAAuthServiceClient is the client API for MFAAuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MFAAuthServiceClient interface {
	// Initiate MFA event.
	InitiateMFA(ctx context.Context, in *InitiateMFARequest, opts ...grpc.CallOption) (*InitiateMFAResponse, error)
	// Validate MFA value.
	ValidateMFA(ctx context.Context, in *ValidateMFARequest, opts ...grpc.CallOption) (*ValidateMFAResponse, error)
	// Retry MFA event.
	RetryMFA(ctx context.Context, in *RetryMFARequest, opts ...grpc.CallOption) (*RetryMFAResponse, error)
	// Record external MFA code.
	ExternalMFA(ctx context.Context, in *ExternalMFARequest, opts ...grpc.CallOption) (*ExternalMFAResponse, error)
}

type mFAAuthServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMFAAuthServiceClient(cc grpc.ClientConnInterface) MFAAuthServiceClient {
	return &mFAAuthServiceClient{cc}
}

func (c *mFAAuthServiceClient) InitiateMFA(ctx context.Context, in *InitiateMFARequest, opts ...grpc.CallOption) (*InitiateMFAResponse, error) {
	out := new(InitiateMFAResponse)
	err := c.cc.Invoke(ctx, "/mfa.MFAAuthService/InitiateMFA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mFAAuthServiceClient) ValidateMFA(ctx context.Context, in *ValidateMFARequest, opts ...grpc.CallOption) (*ValidateMFAResponse, error) {
	out := new(ValidateMFAResponse)
	err := c.cc.Invoke(ctx, "/mfa.MFAAuthService/ValidateMFA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mFAAuthServiceClient) RetryMFA(ctx context.Context, in *RetryMFARequest, opts ...grpc.CallOption) (*RetryMFAResponse, error) {
	out := new(RetryMFAResponse)
	err := c.cc.Invoke(ctx, "/mfa.MFAAuthService/RetryMFA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mFAAuthServiceClient) ExternalMFA(ctx context.Context, in *ExternalMFARequest, opts ...grpc.CallOption) (*ExternalMFAResponse, error) {
	out := new(ExternalMFAResponse)
	err := c.cc.Invoke(ctx, "/mfa.MFAAuthService/ExternalMFA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MFAAuthServiceServer is the server API for MFAAuthService service.
// All implementations must embed UnimplementedMFAAuthServiceServer
// for forward compatibility
type MFAAuthServiceServer interface {
	// Initiate MFA event.
	InitiateMFA(context.Context, *InitiateMFARequest) (*InitiateMFAResponse, error)
	// Validate MFA value.
	ValidateMFA(context.Context, *ValidateMFARequest) (*ValidateMFAResponse, error)
	// Retry MFA event.
	RetryMFA(context.Context, *RetryMFARequest) (*RetryMFAResponse, error)
	// Record external MFA code.
	ExternalMFA(context.Context, *ExternalMFARequest) (*ExternalMFAResponse, error)
	mustEmbedUnimplementedMFAAuthServiceServer()
}

// UnimplementedMFAAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMFAAuthServiceServer struct{}

func (UnimplementedMFAAuthServiceServer) InitiateMFA(context.Context, *InitiateMFARequest) (*InitiateMFAResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitiateMFA not implemented")
}

func (UnimplementedMFAAuthServiceServer) ValidateMFA(context.Context, *ValidateMFARequest) (*ValidateMFAResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateMFA not implemented")
}

func (UnimplementedMFAAuthServiceServer) RetryMFA(context.Context, *RetryMFARequest) (*RetryMFAResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetryMFA not implemented")
}

func (UnimplementedMFAAuthServiceServer) ExternalMFA(context.Context, *ExternalMFARequest) (*ExternalMFAResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExternalMFA not implemented")
}
func (UnimplementedMFAAuthServiceServer) mustEmbedUnimplementedMFAAuthServiceServer() {}

// UnsafeMFAAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MFAAuthServiceServer will
// result in compilation errors.
type UnsafeMFAAuthServiceServer interface {
	mustEmbedUnimplementedMFAAuthServiceServer()
}

func RegisterMFAAuthServiceServer(s grpc.ServiceRegistrar, srv MFAAuthServiceServer) {
	s.RegisterService(&MFAAuthService_ServiceDesc, srv)
}

func _MFAAuthService_InitiateMFA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitiateMFARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MFAAuthServiceServer).InitiateMFA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mfa.MFAAuthService/InitiateMFA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MFAAuthServiceServer).InitiateMFA(ctx, req.(*InitiateMFARequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MFAAuthService_ValidateMFA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateMFARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MFAAuthServiceServer).ValidateMFA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mfa.MFAAuthService/ValidateMFA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MFAAuthServiceServer).ValidateMFA(ctx, req.(*ValidateMFARequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MFAAuthService_RetryMFA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetryMFARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MFAAuthServiceServer).RetryMFA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mfa.MFAAuthService/RetryMFA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MFAAuthServiceServer).RetryMFA(ctx, req.(*RetryMFARequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MFAAuthService_ExternalMFA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExternalMFARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MFAAuthServiceServer).ExternalMFA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mfa.MFAAuthService/ExternalMFA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MFAAuthServiceServer).ExternalMFA(ctx, req.(*ExternalMFARequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MFAAuthService_ServiceDesc is the grpc.ServiceDesc for MFAAuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MFAAuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mfa.MFAAuthService",
	HandlerType: (*MFAAuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitiateMFA",
			Handler:    _MFAAuthService_InitiateMFA_Handler,
		},
		{
			MethodName: "ValidateMFA",
			Handler:    _MFAAuthService_ValidateMFA_Handler,
		},
		{
			MethodName: "RetryMFA",
			Handler:    _MFAAuthService_RetryMFA_Handler,
		},
		{
			MethodName: "ExternalMFA",
			Handler:    _MFAAuthService_ExternalMFA_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "brank.as/rbac/gunk/v1/mfa/all.proto",
}