package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/c12s/metrics/internal/config"
	"github.com/c12s/metrics/internal/errors"
	"github.com/c12s/metrics/internal/models"
	"github.com/c12s/metrics/internal/utils"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type MetricsService struct {
	FileService        *LocalFileService
	QueryMetricsConfig *config.MetricsConfig
	AppConfig          *config.AppConfig
	NodeID             string
	UsageMetrics       *models.UsageMetrics
}

func NewMetricsService(fileService *LocalFileService, queryMetricsConfig *config.MetricsConfig, nodeID string, appConfig *config.AppConfig) *MetricsService {
	return &MetricsService{
		FileService:        fileService,
		QueryMetricsConfig: queryMetricsConfig,
		NodeID:             nodeID,
		AppConfig:          appConfig,
		UsageMetrics:       models.NewUsageMetrics(),
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

func (m *MetricsService) GetMetrics() *errors.ErrorStruct {
	byteQueryResultsFromCAdvisor, err := m.SendExternalGetRequestToMetricsEndpoint(m.AppConfig.GetCAdvisorAddress())
	if err != nil {
		log.Println("Byte query result from cAdvisor", err.GetErrorMessage())
		return err
	}
	actualMetricsValueFromCAdvisor, err := m.castResultsFromBytesToActualValue(byteQueryResultsFromCAdvisor, "cAdvisor")
	if err != nil {
		return err
	}
	byteQueryResultsFromNodeExporter, err := m.SendExternalGetRequestToMetricsEndpoint(m.AppConfig.GetNodeExporterAddress())
	if err != nil {
		log.Println("Byte query result from cAdvisor", err.GetErrorMessage())
		return err
	}
	actualMetricsValueFromNodeExporter, err := m.castResultsFromBytesToActualValue(byteQueryResultsFromNodeExporter, "node-exporter")
	if err != nil {
		return err
	}
	mergedSlicesForMetrics := append(*actualMetricsValueFromCAdvisor, *actualMetricsValueFromNodeExporter...)
	fileFormat := models.MetricFileFormat{
		NodeId:  m.NodeID,
		Metrics: mergedSlicesForMetrics,
	}
	byteFormatOfMetrics, err := m.formatMetricsIntoByteArray(&fileFormat)
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

func (ms *MetricsService) castResultsFromBytesToActualValue(readedBytes []byte, resultsScrapedFrom string) (*[]models.MetricData, *errors.ErrorStruct) {
	data := string(readedBytes)

	parser := expfmt.TextParser{}
	metrics, err := parser.TextToMetricFamilies(strings.NewReader(data))
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	var parsedMetrics []models.MetricData
	for _, mf := range metrics {
		for _, m := range mf.Metric {
			if _, exists := (*ms.QueryMetricsConfig.GetQueries())[*mf.Name]; exists {
				metric := ms.createMetricData(*mf.Name, m)
				parsedMetrics = append(parsedMetrics, metric)
				ms.UsageMetrics.UpdateUsageMetrics(metric)
			}
		}
	}
	if resultsScrapedFrom == "cAdvisor" {
		parsedMetrics = append(parsedMetrics, ms.UsageMetrics.GetCustomMetricDataFromCAdvisor()...)
	} else {
		parsedMetrics = append(parsedMetrics, ms.UsageMetrics.GetCustomMetricDataFromNodeExporter()...)

	}
	return &parsedMetrics, nil
}

func (ms *MetricsService) createMetricData(metricName string, m *dto.Metric) models.MetricData {
	metric := models.MetricData{
		MetricName: metricName,
		Labels:     make(map[string]string),
		Value:      ms.getValue(m),
	}
	metric.Timestamp = time.Now().Unix()
	for _, label := range m.Label {
		if *label.Value == "" {
			continue
		}
		metric.Labels[*label.Name] = *label.Value
	}
	return metric
}

func (m *MetricsService) getValue(passedValue *dto.Metric) float64 {
	if passedValue.Gauge != nil {
		return *passedValue.Gauge.Value
	} else if passedValue.Counter != nil {
		return *passedValue.Counter.Value
	} else if passedValue.Untyped != nil {
		return *passedValue.Untyped.Value
	}
	return 0
}

func (m *MetricsService) getMetricsToFileWriteFormat(data *[]models.MetricData) (*models.MetricFileFormat, *errors.ErrorStruct) {
	return &models.MetricFileFormat{
		NodeId:  m.NodeID,
		Metrics: *data,
	}, nil
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
	castedMetricsIntoProperMapStructure := utils.ConvertFromStringArrayToMapStringStruct(newMetrics)
	m.QueryMetricsConfig.AppendNewMetricsToDefaultMap(castedMetricsIntoProperMapStructure)
	log.Println(m.QueryMetricsConfig.GetQueries())
	err := m.GetMetrics()
	if err != nil {
		return err
	}
	return nil
}

func (ms MetricsService) SendExternalGetRequestToMetricsEndpoint(url string) ([]byte, *errors.ErrorStruct) {
	response, err := http.Get("http://" + url + "/metrics")
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	return body, nil
}
