package handler

import (
	"net/http"

	"github.com/c12s/metrics/internal/service"
	"github.com/c12s/metrics/internal/utils"
)

type NewMetricsTest struct {
	Queries []string `json:"queries"`
}

type MetricsHandler struct {
	metricsService *service.MetricsService
}

func NewMetricsHandler(metricsService *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
	}
}

func (mh MetricsHandler) Test(rw http.ResponseWriter, h *http.Request) {
	mh.metricsService.GetMetricsFromPrometheus()
	utils.WriteResp("OK", 200, rw)
}

func (mh MetricsHandler) GetLatestMetrics(rw http.ResponseWriter, h *http.Request) {
	data, err := mh.metricsService.GetLatestMetrics()
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "/api/metrics", rw)
		return
	}
	utils.WriteResp(data, 200, rw)
}

func (mh MetricsHandler) PostNewMetrics(rw http.ResponseWriter, h *http.Request) {
	var metrics NewMetricsTest
	if !utils.DecodeJSONFromRequest(h, rw, &metrics) {
		utils.WriteErrorResp("Error while casting data into structure", 500, "/api/auth/register", rw)
		return
	}
	mh.metricsService.ReloadQuery(metrics.Queries)
	utils.WriteResp(map[string]string{
		"status": "OK",
	}, 200, rw)
}
