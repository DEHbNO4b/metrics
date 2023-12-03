package metrics

import (
	"context"
	"errors"

	pb "github.com/DEHbNO4b/metrics/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DEHbNO4b/metrics/internal/interfaces"
)

type MetricsAPI struct {
	pb.UnimplementedMetricsServer
	expert interfaces.MetricsStorage
}

func Register(gRPC *grpc.Server) {
	pb.RegisterMetricsServer(gRPC, &MetricsAPI{})
}
func (m *MetricsAPI) AddSingle(ctx context.Context, in *pb.AddSingleRequest) (*pb.AddSingleResponse, error) {
	res := pb.AddSingleResponse{}
	err := validate(in.Metric)
	if err != nil {
		return nil, err
	}
	metric := grpcMetricToDomain(in.Metric)
	m.expert.SetMetric(metric)
	return &res, nil
}

func (m *MetricsAPI) GetSingle(ctx context.Context, in *pb.GetSingleRequest) (*pb.GetSingleResponse, error) {
	res := pb.GetSingleResponse{}
	err := validate(in.Metric)
	if err != nil {
		return nil, err
	}
	metric, err := m.expert.GetMetric(grpcMetricToDomain(in.Metric))
	if err != nil {
		if errors.Is(err, interfaces.ErrNotContains) {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	dm := domainMetricToGrps(metric)
	res.Metric = dm
	return &res, nil

}

// func (m *MetricsAPI) AddMetrics(context.Context, *pb.AddMetricsRequest) (*pb.AddMetricsResponse, error) {
// 	res := pb.AddMetricsResponse{}

// 	return &res, nil
// }

func (m *MetricsAPI) GetMetrics(context.Context, *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	res := pb.GetMetricsResponse{}
	metrics := m.expert.GetMetrics()
	for _, el := range metrics {
		res.Metrics = append(res.Metrics, domainMetricToGrps(el))
	}

	return &res, nil
}
