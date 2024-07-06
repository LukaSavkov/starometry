package config

import "os"

type AppConfig struct {
	serverPort       string
	nodeExporterAddr string
	cAdvisorAddr     string
	natsAddr         string
	grpcPort         string
	nodeID           string
}

func NewAppConfigFromEnv() *AppConfig {
	return &AppConfig{
		serverPort:       os.Getenv("APP_PORT"),
		nodeExporterAddr: os.Getenv("NODE_EXPORTER_URL") + ":" + os.Getenv("NODE_EXPORTER_PORT"),
		cAdvisorAddr:     os.Getenv("CADVISOR_URL") + ":" + os.Getenv("CADVISOR_PORT"),
		natsAddr:         os.Getenv("NATS_URL") + ":" + os.Getenv("NATS_PORT"),
		grpcPort:         os.Getenv("GRPC_PORT"),
		nodeID:           "",
	}
}

func (ap *AppConfig) GetServerPort() string {
	return ap.serverPort
}

func (ap *AppConfig) GetNodeExporterAddress() string {
	return ap.nodeExporterAddr
}

func (ap *AppConfig) GetCAdvisorAddress() string {
	return ap.cAdvisorAddr
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
