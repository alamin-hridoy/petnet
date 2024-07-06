// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: brank.as/petnet/gunk/v1/mfa/all.proto

/*
Package mfa is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package mfa

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code

var (
	_ io.Reader
	_ status.Status
	_ = runtime.String
	_ = utilities.NewDoubleArray
	_ = metadata.Join
)

func request_MFAService_EnableMFA_0(ctx context.Context, marshaler runtime.Marshaler, client MFAServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq EnableMFARequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["UserID"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "UserID")
	}

	protoReq.UserID, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "UserID", err)
	}

	msg, err := client.EnableMFA(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_MFAService_EnableMFA_0(ctx context.Context, marshaler runtime.Marshaler, server MFAServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq EnableMFARequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["UserID"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "UserID")
	}

	protoReq.UserID, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "UserID", err)
	}

	msg, err := server.EnableMFA(ctx, &protoReq)
	return msg, metadata, err
}

// RegisterMFAServiceHandlerServer registers the http handlers for service MFAService to "mux".
// UnaryRPC     :call MFAServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterMFAServiceHandlerFromEndpoint instead.
func RegisterMFAServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server MFAServiceServer) error {
	mux.Handle("POST", pattern_MFAService_EnableMFA_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/petnet.v1.mfa.MFAService/EnableMFA")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_MFAService_EnableMFA_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_MFAService_EnableMFA_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})

	return nil
}

// RegisterMFAServiceHandlerFromEndpoint is same as RegisterMFAServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterMFAServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterMFAServiceHandler(ctx, mux, conn)
}

// RegisterMFAServiceHandler registers the http handlers for service MFAService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterMFAServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterMFAServiceHandlerClient(ctx, mux, NewMFAServiceClient(conn))
}

// RegisterMFAServiceHandlerClient registers the http handlers for service MFAService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "MFAServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "MFAServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "MFAServiceClient" to call the correct interceptors.
func RegisterMFAServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client MFAServiceClient) error {
	mux.Handle("POST", pattern_MFAService_EnableMFA_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/petnet.v1.mfa.MFAService/EnableMFA")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_MFAService_EnableMFA_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_MFAService_EnableMFA_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})

	return nil
}

var pattern_MFAService_EnableMFA_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 2}, []string{"v1", "mfa", "UserID"}, ""))

var forward_MFAService_EnableMFA_0 = runtime.ForwardResponseMessage
