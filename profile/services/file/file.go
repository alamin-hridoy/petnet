package file

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	"brank.as/petnet/profile/storage"
)

type FileCore interface {
	UploadFile(ctx context.Context, fu storage.FileUpload) (*storage.FileUpload, error)
	ListFiles(ctx context.Context, oid string, f storage.FileUploadFilter) ([]storage.FileUpload, error)
	DeleteFileUpload(ctx context.Context, fu *fpb.DeleteFileUploadRequest) error
}

type Svc struct {
	fpb.UnimplementedFileServiceServer
	core FileCore
}

func New(core FileCore) *Svc {
	s := &Svc{
		core: core,
	}
	return s
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { fpb.RegisterFileServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return fpb.RegisterFileServiceHandlerFromEndpoint(ctx, mux, address, options)
}
