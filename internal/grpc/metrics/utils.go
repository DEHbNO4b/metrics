package metrics

import (
	"github.com/DEHbNO4b/metrics/internal/data"
	pb "github.com/DEHbNO4b/metrics/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcMetricToDomain(m *pb.Metric) data.Metrics {
	ans := data.NewMetric()
	ans.ID = m.Id
	ans.MType = m.Type
	*ans.Delta = m.Delta
	*ans.Value = m.Value
	return ans
}
func domainMetricToGrps(m data.Metrics) *pb.Metric {
	ans := pb.Metric{}
	ans.Id = m.ID
	ans.Type = m.MType
	ans.Value = *m.Value
	ans.Delta = *m.Delta
	return &ans
}
func validate(m *pb.Metric) error {
	if m.Id == "" || m.Type == "" {
		return status.Error(codes.InvalidArgument, "wrong request data")
	}
	if m.Type != "gauge" && m.Type != "counter" {
		return status.Errorf(codes.InvalidArgument, "wrong metric type %s", m.Type)
	}
	return nil
}
