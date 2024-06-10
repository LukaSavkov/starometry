package client

import (
	"github.com/c12s/metrics/pkg/external"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewExternalMetricsClient(address string) (external.ExternalMetricsClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return external.NewExternalMetricsClient(conn), nil
}
