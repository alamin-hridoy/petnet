// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package terminal

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

// TerminalServiceClient is the client API for TerminalService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TerminalServiceClient interface {
	// Create Remittance transaction.
	CreateRemit(ctx context.Context, in *CreateRemitRequest, opts ...grpc.CallOption) (*CreateRemitResponse, error)
	// Confirm and process Remittance transaction.
	ConfirmRemit(ctx context.Context, in *ConfirmRemitRequest, opts ...grpc.CallOption) (*ConfirmRemitResponse, error)
	// Get user profile by ID.
	ListRemit(ctx context.Context, in *ListRemitRequest, opts ...grpc.CallOption) (*ListRemitResponse, error)
	// Search remittance.
	LookupRemit(ctx context.Context, in *LookupRemitRequest, opts ...grpc.CallOption) (*LookupRemitResponse, error)
	// Disburse remittance.
	DisburseRemit(ctx context.Context, in *DisburseRemitRequest, opts ...grpc.CallOption) (*DisburseRemitResponse, error)
	// Search partner by transaction id.
	GetPartnerByTxnID(ctx context.Context, in *GetPartnerByTxnIDRequest, opts ...grpc.CallOption) (*GetPartnerByTxnIDResponse, error)
}

type terminalServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTerminalServiceClient(cc grpc.ClientConnInterface) TerminalServiceClient {
	return &terminalServiceClient{cc}
}

func (c *terminalServiceClient) CreateRemit(ctx context.Context, in *CreateRemitRequest, opts ...grpc.CallOption) (*CreateRemitResponse, error) {
	out := new(CreateRemitResponse)
	err := c.cc.Invoke(ctx, "/terminal.TerminalService/CreateRemit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terminalServiceClient) ConfirmRemit(ctx context.Context, in *ConfirmRemitRequest, opts ...grpc.CallOption) (*ConfirmRemitResponse, error) {
	out := new(ConfirmRemitResponse)
	err := c.cc.Invoke(ctx, "/terminal.TerminalService/ConfirmRemit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terminalServiceClient) ListRemit(ctx context.Context, in *ListRemitRequest, opts ...grpc.CallOption) (*ListRemitResponse, error) {
	out := new(ListRemitResponse)
	err := c.cc.Invoke(ctx, "/terminal.TerminalService/ListRemit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terminalServiceClient) LookupRemit(ctx context.Context, in *LookupRemitRequest, opts ...grpc.CallOption) (*LookupRemitResponse, error) {
	out := new(LookupRemitResponse)
	err := c.cc.Invoke(ctx, "/terminal.TerminalService/LookupRemit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terminalServiceClient) DisburseRemit(ctx context.Context, in *DisburseRemitRequest, opts ...grpc.CallOption) (*DisburseRemitResponse, error) {
	out := new(DisburseRemitResponse)
	err := c.cc.Invoke(ctx, "/terminal.TerminalService/DisburseRemit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terminalServiceClient) GetPartnerByTxnID(ctx context.Context, in *GetPartnerByTxnIDRequest, opts ...grpc.CallOption) (*GetPartnerByTxnIDResponse, error) {
	out := new(GetPartnerByTxnIDResponse)
	err := c.cc.Invoke(ctx, "/terminal.TerminalService/GetPartnerByTxnID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TerminalServiceServer is the server API for TerminalService service.
// All implementations must embed UnimplementedTerminalServiceServer
// for forward compatibility
type TerminalServiceServer interface {
	// Create Remittance transaction.
	CreateRemit(context.Context, *CreateRemitRequest) (*CreateRemitResponse, error)
	// Confirm and process Remittance transaction.
	ConfirmRemit(context.Context, *ConfirmRemitRequest) (*ConfirmRemitResponse, error)
	// Get user profile by ID.
	ListRemit(context.Context, *ListRemitRequest) (*ListRemitResponse, error)
	// Search remittance.
	LookupRemit(context.Context, *LookupRemitRequest) (*LookupRemitResponse, error)
	// Disburse remittance.
	DisburseRemit(context.Context, *DisburseRemitRequest) (*DisburseRemitResponse, error)
	// Search partner by transaction id.
	GetPartnerByTxnID(context.Context, *GetPartnerByTxnIDRequest) (*GetPartnerByTxnIDResponse, error)
	mustEmbedUnimplementedTerminalServiceServer()
}

// UnimplementedTerminalServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTerminalServiceServer struct{}

func (UnimplementedTerminalServiceServer) CreateRemit(context.Context, *CreateRemitRequest) (*CreateRemitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRemit not implemented")
}

func (UnimplementedTerminalServiceServer) ConfirmRemit(context.Context, *ConfirmRemitRequest) (*ConfirmRemitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmRemit not implemented")
}

func (UnimplementedTerminalServiceServer) ListRemit(context.Context, *ListRemitRequest) (*ListRemitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRemit not implemented")
}

func (UnimplementedTerminalServiceServer) LookupRemit(context.Context, *LookupRemitRequest) (*LookupRemitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LookupRemit not implemented")
}

func (UnimplementedTerminalServiceServer) DisburseRemit(context.Context, *DisburseRemitRequest) (*DisburseRemitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisburseRemit not implemented")
}

func (UnimplementedTerminalServiceServer) GetPartnerByTxnID(context.Context, *GetPartnerByTxnIDRequest) (*GetPartnerByTxnIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPartnerByTxnID not implemented")
}
func (UnimplementedTerminalServiceServer) mustEmbedUnimplementedTerminalServiceServer() {}

// UnsafeTerminalServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TerminalServiceServer will
// result in compilation errors.
type UnsafeTerminalServiceServer interface {
	mustEmbedUnimplementedTerminalServiceServer()
}

func RegisterTerminalServiceServer(s grpc.ServiceRegistrar, srv TerminalServiceServer) {
	s.RegisterService(&TerminalService_ServiceDesc, srv)
}

func _TerminalService_CreateRemit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRemitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerminalServiceServer).CreateRemit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terminal.TerminalService/CreateRemit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerminalServiceServer).CreateRemit(ctx, req.(*CreateRemitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerminalService_ConfirmRemit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfirmRemitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerminalServiceServer).ConfirmRemit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terminal.TerminalService/ConfirmRemit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerminalServiceServer).ConfirmRemit(ctx, req.(*ConfirmRemitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerminalService_ListRemit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRemitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerminalServiceServer).ListRemit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terminal.TerminalService/ListRemit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerminalServiceServer).ListRemit(ctx, req.(*ListRemitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerminalService_LookupRemit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LookupRemitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerminalServiceServer).LookupRemit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terminal.TerminalService/LookupRemit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerminalServiceServer).LookupRemit(ctx, req.(*LookupRemitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerminalService_DisburseRemit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisburseRemitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerminalServiceServer).DisburseRemit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terminal.TerminalService/DisburseRemit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerminalServiceServer).DisburseRemit(ctx, req.(*DisburseRemitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerminalService_GetPartnerByTxnID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPartnerByTxnIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerminalServiceServer).GetPartnerByTxnID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terminal.TerminalService/GetPartnerByTxnID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerminalServiceServer).GetPartnerByTxnID(ctx, req.(*GetPartnerByTxnIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TerminalService_ServiceDesc is the grpc.ServiceDesc for TerminalService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TerminalService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "terminal.TerminalService",
	HandlerType: (*TerminalServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRemit",
			Handler:    _TerminalService_CreateRemit_Handler,
		},
		{
			MethodName: "ConfirmRemit",
			Handler:    _TerminalService_ConfirmRemit_Handler,
		},
		{
			MethodName: "ListRemit",
			Handler:    _TerminalService_ListRemit_Handler,
		},
		{
			MethodName: "LookupRemit",
			Handler:    _TerminalService_LookupRemit_Handler,
		},
		{
			MethodName: "DisburseRemit",
			Handler:    _TerminalService_DisburseRemit_Handler,
		},
		{
			MethodName: "GetPartnerByTxnID",
			Handler:    _TerminalService_GetPartnerByTxnID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "brank.as/petnet/gunk/drp/v1/terminal/all.proto",
}