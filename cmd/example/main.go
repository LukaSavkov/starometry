package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/c12s/metrics/pkg/external"
	"google.golang.org/grpc"
)

func main() {
	app := AppServer{}
	server := grpc.NewServer()
	external.RegisterExternalMetricsServer(server, app)
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type AppServer struct {
	external.UnimplementedExternalMetricsServer
}

func (a AppServer) ExternalLatestMetrics(context.Context, *external.ExternalLatestMetricsReq) (*external.ExternalMetricsArray, error) {
	return &external.ExternalMetricsArray{
		Metrics: []*external.ExternalMetricData{
			{
				MetricName: "demo_app_metric",
				Labels:     map[string]string{"label": "value"},
				Value:      5,
				Timestamp:  time.Now().Unix(),
			},
		},
	}, nil
}
