package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/c12s/metrics/internal/utils"
)

// var Queries = []string{
// 	"container_cpu_usage_seconds_total",
// 	"container_memory_usage_bytes",
// 	"container_network_receive_bytes_total",
// 	"container_network_transmit_bytes_total",
// 	"container_fs_usage_bytes",
// 	"container_fs_writes_bytes_total",
// 	"container_fs_reads_bytes_total",
// 	"container_start_time_seconds",
// 	"container_tasks_state",
// 	"node_cpu_seconds_total",
// 	"node_memory_MemAvailable_bytes",
// 	"node_disk_io_time_seconds_total",
// 	"node_disk_read_bytes_total",
// 	"node_disk_written_bytes_total",
// 	"node_network_receive_bytes_total",
// 	"node_network_transmit_bytes_total",
// 	"node_filesystem_avail_bytes",
// 	"node_filesystem_size_bytes",
// 	"node_load1",
// 	"node_load5",
// 	"node_load15",
// }

var Queries = map[string]struct{}{
	"container_cpu_usage_seconds_total":      {},
	"container_spec_cpu_quota":               {},
	"container_memory_usage_bytes":           {},
	"container_spec_memory_limit_bytes":      {},
	"container_fs_usage_bytes":               {},
	"container_spec_cpu_period":              {},
	"container_network_receive_bytes_total":  {},
	"container_network_transmit_bytes_total": {},
	"node_cpu_seconds_total":                 {},
	"node_memory_MemTotal_bytes":             {},
	"node_memory_MemAvailable_bytes":         {},
	"node_filesystem_size_bytes":             {},
	"node_filesystem_free_bytes":             {},
	"node_network_receive_bytes_total":       {},
	"node_network_transmit_bytes_total":      {},
	"machine_cpu_cores":                      {},
}
var basicCronTimerForScrapingMetrics = "30s"
var basicCronTimerForScrapingExternalMetrics = "30s"

type MetricsConfig struct {
	queries              map[string]struct{}
	cronTimer            string
	externalAppCronTimer string
}

func (mc *MetricsConfig) GetQueries() *map[string]struct{} {
	return &mc.queries
}

func (mc *MetricsConfig) GetCronTimer() string {
	return mc.cronTimer
}

func (mc *MetricsConfig) SetQueries(queries map[string]struct{}) {
	mc.queries = queries
}

func (mc *MetricsConfig) SetCronTimer(cronTimer string) {
	mc.cronTimer = cronTimer
}

func (mc *MetricsConfig) GetExternalCronTimer() string {
	return mc.externalAppCronTimer
}

func (mc *MetricsConfig) AppendNewMetricsToDefaultMap(newMetrics map[string]struct{}) {
	mc.queries = Queries
	for k, v := range newMetrics {
		mc.queries[k] = v
	}
}

func NewMetricsConfigWithPresetConfiguration() *MetricsConfig {
	return &MetricsConfig{
		queries:              Queries,
		cronTimer:            basicCronTimerForScrapingMetrics,
		externalAppCronTimer: basicCronTimerForScrapingExternalMetrics,
	}
}

func NewMetricsConfigLoadedFromEnv() (*MetricsConfig, error) {
	queriesFromENV := os.Getenv("APP_METRICS_CONFIG")
	cronTimer := os.Getenv("APP_METRICS_CRON_TIMER")
	externalAppCronTimer := os.Getenv("APP_METRICS_EXTERNAL_CRON_TIMER")
	if queriesFromENV == "" && cronTimer == "" && externalAppCronTimer == "" {
		fmt.Println("Queries or crone timers are not up")
		return nil, errors.ErrUnsupported
	}
	var queriesToSet map[string]struct{}
	if queriesFromENV != "" {
		queriesToSet = utils.ConvertFromCSVToMapStringStruct(os.Getenv("APP_METRICS_CONFIG"))
	} else {
		queriesToSet = Queries
	}
	if cronTimer == "" {
		cronTimer = basicCronTimerForScrapingMetrics
	}
	if externalAppCronTimer == "" {
		externalAppCronTimer = basicCronTimerForScrapingExternalMetrics
	}
	return &MetricsConfig{
		queries:              queriesToSet,
		cronTimer:            cronTimer,
		externalAppCronTimer: externalAppCronTimer,
	}, nil
}
