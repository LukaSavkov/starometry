package config

import "os"

type AppConfig struct {
	serverPort     string
	prometheusAddr string
	natsAddr       string
	grpcPort       string
	nodeID         string
}

func NewAppConfigFromEnv() *AppConfig {
	return &AppConfig{
		serverPort:     os.Getenv("APP_PORT"),
		prometheusAddr: os.Getenv("PROMETHEUS_URL") + ":" + os.Getenv("PROMETHEUS_PORT"),
		natsAddr:       os.Getenv("NATS_URL") + ":" + os.Getenv("NATS_PORT"),
		grpcPort:       os.Getenv("GRPC_PORT"),
		nodeID:         "",
	}
}

func (ap *AppConfig) GetServerPort() string {
	return ap.serverPort
}

func (ap *AppConfig) GetPrometheusAddress() string {
	return ap.prometheusAddr
}

func (ap *AppConfig) GetNatsAddress() string {
	return ap.natsAddr
}

func (ap *AppConfig) GetGRPCPort() string {
	return ap.grpcPort
}

func (ap *AppConfig) GetNodeID() string {
	return ap.nodeID
}

func (ap *AppConfig) SetNodeID(nodeID string) {
	if nodeID == "" {
		return
	}
	ap.nodeID = nodeID
}
