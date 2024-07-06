package riskassesment

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	rat "brank.as/petnet/gunk/dsa/v1/riskassesment"
)

type Svc struct {
	rat.UnimplementedRiskAssesmentServiceServer
	core RiskAssesmentCore
}

type RiskAssesmentCore interface {
	UpsertQuestion(ctx context.Context, in *rat.RiskAssesmentQuestionRequest) (*rat.RiskAssesmentQuestionResponse, error)
	UpsertMlTfQuestion(ctx context.Context, in *rat.RiskAssesmentQuestionRequest) (*rat.RiskAssesmentQuestionResponse, error)
	ListQuestion(ctx context.Context, in *rat.ListQuestionRequest) (*rat.ListQuestionResponse, error)
	ListMlTfQuestion(ctx context.Context, in *rat.ListQuestionRequest) (*rat.ListQuestionResponse, error)
}

func New(core RiskAssesmentCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { rat.RegisterRiskAssesmentServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return rat.RegisterRiskAssesmentServiceHandlerFromEndpoint(ctx, mux, address, options)
}
