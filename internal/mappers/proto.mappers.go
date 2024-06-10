package mappers

import (
	"log"

	"github.com/c12s/metrics/internal/errors"
	"github.com/c12s/metrics/internal/models"
	api "github.com/c12s/metrics/pkg/api"
	"github.com/c12s/metrics/pkg/external"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(err *errors.ErrorStruct) error {
	if err == nil {
		log.Println("Received nil error in MapError")
		return nil
	}
	log.Println("Error from gRPC mapper:", err)
	return status.Error(codes.Code(err.GetErrorStatus()), err.GetErrorMessage())
}

func MapFromModelFileFormatToGrpcPostResponseFormat(fileFormat models.MetricFileFormat) *api.PostNewMetricsResp {
	apiMetrics := &api.MetricsWithNodeID{
		Metrics: make([]*api.MetricData, len(fileFormat.Metrics)),
	}
	for i, metric := range fileFormat.Metrics {
		apiMetrics.Metrics[i] = &api.MetricData{
			MetricName: metric.MetricName,
			Labels:     metric.Labels,
			Value:      metric.Value,
			Timestamp:  metric.Timestamp,
		}
	}
	apiMetrics.NodeId = fileFormat.NodeId

	return &api.PostNewMetricsResp{
		Data: apiMetrics,
	}
}

func MapFromModelFileFormatToGrpcGetResponseFormat(fileFormat models.MetricFileFormat) *api.GetLatestMetricsResp {
	apiMetrics := &api.MetricsWithNodeID{
		Metrics: make([]*api.MetricData, len(fileFormat.Metrics)),
	}
	for i, metric := range fileFormat.Metrics {
		apiMetrics.Metrics[i] = &api.MetricData{
			MetricName: metric.MetricName,
			Labels:     metric.Labels,
			Value:      metric.Value,
			Timestamp:  metric.Timestamp,
		}
	}
	apiMetrics.NodeId = fileFormat.NodeId

	return &api.GetLatestMetricsResp{
		Data: apiMetrics,
	}
}

func MapFromApiExternalApplicationToModelExternalApplication(list []*api.ExternalApplication) []models.ExternalApplication {
	var result []models.ExternalApplication
	for _, apiApp := range list {
		modelApp := models.ExternalApplication{
			Address: apiApp.Address,
		}
		result = append(result, modelApp)
	}
	return result
}

func MapFromModelExternalApplicationToApiExternalApplication(list []models.ExternalApplication) []*api.ExternalApplication {
	var result []*api.ExternalApplication
	for _, modelApp := range list {
		apiApp := &api.ExternalApplication{
			Address: modelApp.Address,
		}
		result = append(result, apiApp)
	}
	return result
}

func MapFromExternalMetricDataToModelMetricData(list []*external.ExternalMetricData) []models.MetricData {
	var result []models.MetricData
	for _, apiApp := range list {
		modelApp := models.MetricData{
			MetricName: apiApp.MetricName,
			Labels:     apiApp.Labels,
			Value:      apiApp.Value,
			Timestamp:  apiApp.Timestamp,
		}
		result = append(result, modelApp)
	}
	return result
}
