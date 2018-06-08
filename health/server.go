package health

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var grpcSrv *grpc.Server
var s *srv

// HealthCheckFunc ...
type HealthCheckFunc func() bool

func init() {
	s = &srv{
		probes: []HealthCheckFunc{},
	}
}

// RegisterHealthServer registers a health-check server and its implementation
// to the gRPC server. This must be called before invoking Serve.
func RegisterHealthServer(srv *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(srv, s)
	reflection.Register(srv)
}

func RegisterHealthCheckFunc(f HealthCheckFunc) {
	s.probes = append(s.probes, f)
}

func ShutdownFunc() func() {
	return func() {
		grpcSrv.GracefulStop()
	}
}

type srv struct {
	probes []HealthCheckFunc
}

func (s *srv) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	var status = grpc_health_v1.HealthCheckResponse_SERVING

	for _, p := range s.probes {
		if ok := p(); !ok {
			status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
			break
		}
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: status,
	}, nil
}
