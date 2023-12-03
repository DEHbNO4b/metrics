package agent

import (
	"context"
	"log"

	"github.com/DEHbNO4b/metrics/internal/data"
	pb "github.com/DEHbNO4b/metrics/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Client pb.MetricsClient
	Conn   *grpc.ClientConn
}

func NewGrpcClient(addr string) *GrpcClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := pb.NewMetricsClient(conn)
	return &GrpcClient{
		Client: c,
		Conn:   conn,
	}
}
func (g *GrpcClient) SendMetric(ctx context.Context, m data.Metrics, key string) {
	metric := domainMetricToGrpc(m)
	g.Client.AddSingle(ctx, &pb.AddSingleRequest{Metric: metric})

}
func (g *GrpcClient) SendMetrics(ctx context.Context, metrics []data.Metrics) {
	gm := make([]*pb.Metric, 0, 30)
	for _, el := range metrics {
		m := domainMetricToGrpc(el)
		gm = append(gm, m)
	}
	g.Client.AddMetrics(ctx, &pb.AddMetricsRequest{Metrics: gm})
}
