package agent

import (
	"github.com/DEHbNO4b/metrics/internal/data"
	pb "github.com/DEHbNO4b/metrics/proto"
)

func domainMetricToGrpc(m data.Metrics) *pb.Metric {
	ans := pb.Metric{
		Type:  m.MType,
		Id:    m.ID,
		Value: *m.Value,
		Delta: *m.Delta,
	}
	return &ans
}

// func grpcMetricToDomain(m *pb.Metric) data.Metrics {
// 	ans := data.NewMetric()
// 	ans.MType = m.Type
// 	ans.ID = m.Id
// 	*ans.Delta = m.Delta
// 	*ans.Value = m.Value
// 	return ans
// }
