// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package oauth2

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

// AuthClientServiceClient is the client API for AuthClientService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClientServiceClient interface {
	// Create auth client for authenticating users for an org site or product.
	CreateClient(ctx context.Context, in *CreateClientRequest, opts ...grpc.CallOption) (*CreateClientResponse, error)
	// Update auth client for authenticating users for an org site or product.
	UpdateClient(ctx context.Context, in *UpdateClientRequest, opts ...grpc.CallOption) (*UpdateClientResponse, error)
	// List all auth clients.
	ListClients(ctx context.Context, in *ListClientsRequest, opts ...grpc.CallOption) (*ListClientsResponse, error)
	// Disable an auth client.  Client is permanently disabled.
	DisableClient(ctx context.Context, in *DisableClientRequest, opts ...grpc.CallOption) (*DisableClientResponse, error)
}

type authClientServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClientServiceClient(cc grpc.ClientConnInterface) AuthClientServiceClient {
	return &authClientServiceClient{cc}
}

func (c *authClientServiceClient) CreateClient(ctx context.Context, in *CreateClientRequest, opts ...grpc.CallOption) (*CreateClientResponse, error) {
	out := new(CreateClientResponse)
	err := c.cc.Invoke(ctx, "/brankas.rbac.v1.oauth2.AuthClientService/CreateClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClientServiceClient) UpdateClient(ctx context.Context, in *UpdateClientRequest, opts ...grpc.CallOption) (*UpdateClientResponse, error) {
	out := new(UpdateClientResponse)
	err := c.cc.Invoke(ctx, "/brankas.rbac.v1.oauth2.AuthClientService/UpdateClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClientServiceClient) ListClients(ctx context.Context, in *ListClientsRequest, opts ...grpc.CallOption) (*ListClientsResponse, error) {
	out := new(ListClientsResponse)
	err := c.cc.Invoke(ctx, "/brankas.rbac.v1.oauth2.AuthClientService/ListClients", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClientServiceClient) DisableClient(ctx context.Context, in *DisableClientRequest, opts ...grpc.CallOption) (*DisableClientResponse, error) {
	out := new(DisableClientResponse)
	err := c.cc.Invoke(ctx, "/brankas.rbac.v1.oauth2.AuthClientService/DisableClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthClientServiceServer is the server API for AuthClientService service.
// All implementations must embed UnimplementedAuthClientServiceServer
// for forward compatibility
type AuthClientServiceServer interface {
	// Create auth client for authenticating users for an org site or product.
	CreateClient(context.Context, *CreateClientRequest) (*CreateClientResponse, error)
	// Update auth client for authenticating users for an org site or product.
	UpdateClient(context.Context, *UpdateClientRequest) (*UpdateClientResponse, error)
	// List all auth clients.
	ListClients(context.Context, *ListClientsRequest) (*ListClientsResponse, error)
	// Disable an auth client.  Client is permanently disabled.
	DisableClient(context.Context, *DisableClientRequest) (*DisableClientResponse, error)
	mustEmbedUnimplementedAuthClientServiceServer()
}

// UnimplementedAuthClientServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthClientServiceServer struct{}

func (UnimplementedAuthClientServiceServer) CreateClient(context.Context, *CreateClientRequest) (*CreateClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateClient not implemented")
}

func (UnimplementedAuthClientServiceServer) UpdateClient(context.Context, *UpdateClientRequest) (*UpdateClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateClient not implemented")
}

func (UnimplementedAuthClientServiceServer) ListClients(context.Context, *ListClientsRequest) (*ListClientsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListClients not implemented")
}

func (UnimplementedAuthClientServiceServer) DisableClient(context.Context, *DisableClientRequest) (*DisableClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableClient not implemented")
}
func (UnimplementedAuthClientServiceServer) mustEmbedUnimplementedAuthClientServiceServer() {}

// UnsafeAuthClientServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthClientServiceServer will
// result in compilation errors.
type UnsafeAuthClientServiceServer interface {
	mustEmbedUnimplementedAuthClientServiceServer()
}

func RegisterAuthClientServiceServer(s grpc.ServiceRegistrar, srv AuthClientServiceServer) {
	s.RegisterService(&AuthClientService_ServiceDesc, srv)
}

func _AuthClientService_CreateClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthClientServiceServer).CreateClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brankas.rbac.v1.oauth2.AuthClientService/CreateClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthClientServiceServer).CreateClient(ctx, req.(*CreateClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthClientService_UpdateClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthClientServiceServer).UpdateClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brankas.rbac.v1.oauth2.AuthClientService/UpdateClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthClientServiceServer).UpdateClient(ctx, req.(*UpdateClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthClientService_ListClients_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListClientsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthClientServiceServer).ListClients(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brankas.rbac.v1.oauth2.AuthClientService/ListClients",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthClientServiceServer).ListClients(ctx, req.(*ListClientsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthClientService_DisableClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisableClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthClientServiceServer).DisableClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brankas.rbac.v1.oauth2.AuthClientService/DisableClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthClientServiceServer).DisableClient(ctx, req.(*DisableClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthClientService_ServiceDesc is the grpc.ServiceDesc for AuthClientService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthClientService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "brankas.rbac.v1.oauth2.AuthClientService",
	HandlerType: (*AuthClientServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateClient",
			Handler:    _AuthClientService_CreateClient_Handler,
		},
		{
			MethodName: "UpdateClient",
			Handler:    _AuthClientService_UpdateClient_Handler,
		},
		{
			MethodName: "ListClients",
			Handler:    _AuthClientService_ListClients_Handler,
		},
		{
			MethodName: "DisableClient",
			Handler:    _AuthClientService_DisableClient_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "brank.as/rbac/gunk/v1/oauth2/all.proto",
}