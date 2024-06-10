package servers

import (
	"context"

	"github.com/c12s/metrics/internal/config"
	"github.com/c12s/metrics/internal/mappers"
	"github.com/c12s/metrics/internal/service"

	api "github.com/c12s/metrics/pkg/api"
)

type MetricsGrpcServer struct {
	api.UnimplementedMetricsServer
	metricsService             *service.MetricsService
	externalApplicationsConfig *config.ExternalApplicationsConfig
}

func NewMetricsGrpcServer(metricsService *service.MetricsService, externalApplicationsConfig *config.ExternalApplicationsConfig) api.MetricsServer {
	return &MetricsGrpcServer{
		metricsService:             metricsService,
		externalApplicationsConfig: externalApplicationsConfig,
	}
}

func (m *MetricsGrpcServer) PostNewMetrics(ctx context.Context, req *api.NewMetricsArray) (*api.PostNewMetricsResp, error) {
	err := m.metricsService.ReloadQuery(req.Metrics)
	if err != nil {
		return nil, mappers.MapError(err)
	}
	data, err := m.metricsService.GetLatestMetrics()
	if err != nil {
		return nil, mappers.MapError(err)
	}
	apiMetrics := mappers.MapFromModelFileFormatToGrpcPostResponseFormat(*data)
	return apiMetrics, nil
}

func (m *MetricsGrpcServer) GetLatestMetrics(ctx context.Context, empty *api.GetLatestMetricsReq) (*api.GetLatestMetricsResp, error) {
	latestMetrics, err := m.metricsService.GetLatestMetrics()
	if err != nil {
		return nil, mappers.MapError(err)
	}
	apiMetrics := mappers.MapFromModelFileFormatToGrpcGetResponseFormat(*latestMetrics)
	return apiMetrics, nil
}

func (m *MetricsGrpcServer) PostNewExternalApplicationsList(ctx context.Context, applicationsList *api.ExternalApplicationsList) (*api.ExternalApplicationsList, error) {
	castedResults := mappers.MapFromApiExternalApplicationToModelExternalApplication(applicationsList.ExternalApplications)
	m.externalApplicationsConfig.LoadNewApplications(castedResults)
	castedResultsToApiExternalResult := mappers.MapFromModelExternalApplicationToApiExternalApplication(castedResults)
	return &api.ExternalApplicationsList{
		ExternalApplications: castedResultsToApiExternalResult,
	}, nil
}
