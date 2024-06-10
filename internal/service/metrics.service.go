package service

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/c12s/metrics/internal/config"
	"github.com/c12s/metrics/internal/errors"
	"github.com/c12s/metrics/internal/models"
	"github.com/c12s/metrics/internal/utils"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type MetricsService struct {
	prometheusAPI      v1.API
	FileService        *LocalFileService
	Query              string
	QueryMetricsConfig *config.MetricsConfig
	NodeID             string
}

func NewMetricsService(prometheusAPI v1.API, fileService *LocalFileService, query string, queryMetricsConfig *config.MetricsConfig, nodeID string) *MetricsService {
	return &MetricsService{
		prometheusAPI:      prometheusAPI,
		FileService:        fileService,
		Query:              query,
		QueryMetricsConfig: queryMetricsConfig,
		NodeID:             nodeID,
	}
}

func (m *MetricsService) GetLatestMetrics() (*models.MetricFileFormat, *errors.ErrorStruct) {
	byteMetrics, err := m.FileService.ReadFromFile("./data/scraped-metrics.json")
	if err != nil {
		log.Println("Error while reading scraped metrics: ", err)
		return nil, err
	}
	byteMetricsExternal, errFromByteMetricsExternal := m.FileService.ReadFromFile("./data/scraped-metrics-external.json")
	if errFromByteMetricsExternal != nil {
		log.Println("Error reading external metrics, they're not scraped probably: ", errFromByteMetricsExternal)
	}
	metrics, err := m.formatMetricsFromByteArray(byteMetrics)
	if err != nil {
		log.Println("Error casting in metrics from bytes: ", err)
		return nil, err
	}
	if errFromByteMetricsExternal == nil {
		metricsExternal, err := m.formatMetricsFromByteArray(byteMetricsExternal)
		if err != nil {
			log.Println("Error while casting external metrics: ", err)
			return nil, err
		}
		metrics.Metrics = append(metrics.Metrics, metricsExternal.Metrics...)
	}
	metrics.NodeId = m.NodeID
	return metrics, nil
}

func (m *MetricsService) WriteMetricsFromExternalApplication(metrics []models.MetricData) *errors.ErrorStruct {
	fileFormat := &models.MetricFileFormat{
		Metrics: metrics,
		NodeId:  m.NodeID,
	}
	byteFormatOfMetrics, err := m.formatMetricsIntoByteArray(fileFormat)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.GetErrorMessage())
		return err
	}
	errorFromWrite := m.FileService.WriteToFile("data/scraped-metrics-external.json", byteFormatOfMetrics)
	if errorFromWrite != nil {
		log.Fatalf("Error occured during writing to file. Error %s", errorFromWrite)
		return errors.NewError(errorFromWrite.Error(), 500)
	}
	return nil
}

func (m *MetricsService) GetMetricsFromPrometheus() *errors.ErrorStruct {
	queryResults, err := m.queryPrometheus()
	if err != nil {
		log.Println("Query results error", err.GetErrorMessage())
		return err
	}
	fileFormat, err := m.getMetricsToFileWriteFormat(*queryResults)
	if err != nil {
		log.Println("File format error: ", err.GetErrorMessage())
		return err
	}
	byteFormatOfMetrics, err := m.formatMetricsIntoByteArray(fileFormat)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.GetErrorMessage())
		return err
	}
	errorFromWrite := m.FileService.WriteToFile("data/scraped-metrics.json", byteFormatOfMetrics)
	if errorFromWrite != nil {
		log.Fatalf("Error occured during writing to file. Error %s", errorFromWrite)
		return errors.NewError(errorFromWrite.Error(), 500)
	}
	return nil
}

func (m *MetricsService) formatMetricsIntoByteArray(fileFormat *models.MetricFileFormat) ([]byte, *errors.ErrorStruct) {
	jsonFileFormat, err := json.MarshalIndent(fileFormat, "", "    ")
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	return jsonFileFormat, nil
}

func (m *MetricsService) queryPrometheus() (*model.Value, *errors.ErrorStruct) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := m.prometheusAPI.Query(ctx, m.Query, time.Now())
	if err != nil {
		log.Println("Err while querying the prometheus ", err)
		m.QueryMetricsConfig.SetQueries(config.Queries)
		return nil, errors.NewError(err.Error(), 500)
	}
	if len(warnings) > 0 {
		log.Printf("Prometheus query warnings: %v", warnings)
	}
	return &result, nil
}

func (m *MetricsService) getMetricsToFileWriteFormat(val model.Value) (*models.MetricFileFormat, *errors.ErrorStruct) {
	vector, ok := val.(model.Vector)
	if !ok {
		return nil, errors.NewError("Error: Expecting vector type from Prometheus response", 500)
	}
	var metricsFileFormat models.MetricFileFormat
	for _, sample := range vector {
		data := models.MetricData{
			MetricName: string(sample.Metric["__name__"]),
			Labels:     make(map[string]string),
			Value:      float64(sample.Value),
			Timestamp:  int64(sample.Timestamp),
		}
		for labelName, labelValue := range sample.Metric {
			if labelName == "__name__" || labelName == "job" || labelName == "instance" {
				continue
			}
			isDockerComposeLabel := strings.Contains(string(labelName), "container_label_com_docker_compose")
			isAllowedDockerComposeLabel := labelName == "container_label_com_docker_compose_service" ||
				labelName == "container_label_com_docker_compose_version" ||
				labelName == "container_label_com_docker_compose_project"

			if !isDockerComposeLabel || isAllowedDockerComposeLabel {
				data.Labels[string(labelName)] = string(labelValue)
			}
		}
		metricsFileFormat.Metrics = append(metricsFileFormat.Metrics, data)
	}
	return &metricsFileFormat, nil
}

func (m *MetricsService) formatMetricsFromByteArray(data []byte) (*models.MetricFileFormat, *errors.ErrorStruct) {
	var metrics models.MetricFileFormat
	err := json.Unmarshal(data, &metrics)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	return &metrics, nil
}

func (m *MetricsService) ReloadQuery(newMetrics []string) *errors.ErrorStruct {
	m.QueryMetricsConfig.SetQueries(newMetrics)
	m.Query = utils.ConvertFromStringArrayToPromQLQuery(*m.QueryMetricsConfig.GetQueries())
	err := m.GetMetricsFromPrometheus()
	if err != nil {
		return err
	}
	return nil
}

func (m *MetricsService) queryPrometheusWithRetry(retries int, delay time.Duration) (*model.Value, *errors.ErrorStruct) {
	var queryResults *model.Value
	var err *errors.ErrorStruct
	for i := 0; i < retries; i++ {
		queryResults, err = m.queryPrometheus()
		if err == nil {
			log.Printf("Attempt %d succeed.", i+1)
			return queryResults, nil
		}
		log.Printf("Attempt %d failed: %s. Retrying in %v...", i+1, err.GetErrorMessage(), delay)
		time.Sleep(delay)
	}
	log.Printf("All attempts failed. Last error: %s", err.GetErrorMessage())
	return nil, err
}

func (m *MetricsService) GetInitialMetricsFromPrometheusOnApplicationInit() *errors.ErrorStruct {
	retries := 3
	delay := 2 * time.Second

	queryResults, err := m.queryPrometheusWithRetry(retries, delay)
	if err != nil {
		log.Printf("Query results error: %s", err.GetErrorMessage())
		return err
	}

	fileFormat, err := m.getMetricsToFileWriteFormat(*queryResults)
	if err != nil {
		log.Printf("File format error: %s", err.GetErrorMessage())
		return err
	}

	byteFormatOfMetrics, err := m.formatMetricsIntoByteArray(fileFormat)
	if err != nil {
		log.Fatalf("Error occurred during marshaling: %s", err.GetErrorMessage())
		return err
	}

	errorFromWrite := m.FileService.WriteToFile("data/scraped-metrics.json", byteFormatOfMetrics)
	if errorFromWrite != nil {
		log.Fatalf("Error occurred during writing to file: %s", errorFromWrite.Error())
		return errors.NewError(errorFromWrite.Error(), 500)
	}

	return nil
}
