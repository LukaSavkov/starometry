package models

import (
	"time"
)

const (
	BytesInMB = 1024 * 1024
	BytesInGB = 1024 * 1024 * 1024
)

type ServiceUsageMetrics struct {
	CpuUsage        float64
	RamUsage        float64
	RamMax          float64
	DiskUsage       float64
	NetworkReceive  float64
	NetworkTransmit float64
}

type UsageMetrics struct {
	ActiveServices      map[string]ServiceUsageMetrics
	nodeCPUUsage        float64
	nodeCPUUser         float64
	nodeCPUSystem       float64
	nodeCPUIdle         float64
	nodeCPUTotal        float64
	nodeRAMAvailable    float64
	nodeRAMTotal        float64
	nodeDiskUsage       float64
	nodeDiskTotal       float64
	nodeNetworkReceive  float64
	nodeNetworkTransmit float64
}

// NewUsageMetrics creates a new UsageMetrics instance.
func NewUsageMetrics() *UsageMetrics {
	return &UsageMetrics{
		nodeCPUUsage:        0.0,
		nodeCPUUser:         0.0,
		nodeCPUSystem:       0.0,
		nodeCPUIdle:         0.0,
		nodeCPUTotal:        0.0,
		nodeRAMAvailable:    0.0,
		nodeRAMTotal:        0.0,
		nodeDiskUsage:       0.0,
		nodeDiskTotal:       0.0,
		nodeNetworkReceive:  0.0,
		nodeNetworkTransmit: 0.0,
		ActiveServices:      make(map[string]ServiceUsageMetrics),
	}
}

func (um *UsageMetrics) UpdateUsageMetrics(metric MetricData) {
	switch metric.MetricName {
	case "container_cpu_usage_seconds_total":
		if metric.Labels["name"] != "" {
			metricName := metric.Labels["name"]
			serviceMetrics := um.ActiveServices[metricName]
			serviceMetrics.CpuUsage = metric.Value
			um.ActiveServices[metricName] = serviceMetrics
		}
	case "container_memory_usage_bytes":
		if metric.Labels["name"] != "" {
			metricName := metric.Labels["name"]
			serviceMetrics := um.ActiveServices[metricName]
			serviceMetrics.RamUsage = metric.Value
			um.ActiveServices[metricName] = serviceMetrics
		}
	case "container_spec_memory_limit_bytes":
		if metric.Labels["name"] != "" {
			metricName := metric.Labels["name"]
			serviceMetrics := um.ActiveServices[metricName]
			serviceMetrics.RamMax = metric.Value
			um.ActiveServices[metricName] = serviceMetrics
		}
	case "container_fs_usage_bytes":
		if metric.Labels["name"] != "" {
			metricName := metric.Labels["name"]
			serviceMetrics := um.ActiveServices[metricName]
			serviceMetrics.DiskUsage = metric.Value
			um.ActiveServices[metricName] = serviceMetrics
		}
	case "container_network_receive_bytes_total":
		if metric.Labels["name"] != "" && metric.Labels["interface"] == "eth0" {
			metricName := metric.Labels["name"]
			serviceMetrics := um.ActiveServices[metricName]
			serviceMetrics.NetworkReceive = metric.Value
			um.ActiveServices[metricName] = serviceMetrics
		}
	case "container_network_transmit_bytes_total":
		if metric.Labels["name"] != "" && metric.Labels["interface"] == "eth0" {
			metricName := metric.Labels["name"]
			serviceMetrics := um.ActiveServices[metricName]
			serviceMetrics.NetworkTransmit = metric.Value
			um.ActiveServices[metricName] = serviceMetrics
		}
	case "node_cpu_seconds_total":
		um.nodeCPUTotal += metric.Value
		if metric.Labels["mode"] == "user" {
			um.nodeCPUUser += metric.Value
		}
		if metric.Labels["mode"] == "system" {
			um.nodeCPUSystem += metric.Value
		}
		if metric.Labels["mode"] == "idle" {
			um.nodeCPUIdle += metric.Value
		}
	case "node_memory_MemTotal_bytes":
		um.nodeRAMTotal = metric.Value
	case "node_memory_MemAvailable_bytes":
		um.nodeRAMAvailable = metric.Value
	case "node_filesystem_size_bytes":
		if metric.Labels["fstype"] != "rootfs" {
			um.nodeDiskTotal += metric.Value
		}
	case "node_filesystem_free_bytes":
		if metric.Labels["fstype"] != "rootfs" {
			um.nodeDiskUsage += metric.Value
		}

	case "node_network_receive_bytes_total":
		if metric.Labels["device"] == "eth0" {
			um.nodeNetworkReceive = metric.Value
		}
	case "node_network_transmit_bytes_total":
		if metric.Labels["device"] == "eth0" {
			um.nodeNetworkTransmit = metric.Value
		}
	}

}

// calculatePercentages calculates various percentages and returns them as MetricData.
func (um *UsageMetrics) GetCustomMetricDataFromCAdvisor() []MetricData {
	currentTime := time.Now().Unix()
	metrics := make([]MetricData, 0)
	for serviceName, serviceMetrics := range um.ActiveServices {
		metrics = append(metrics, MetricData{
			MetricName: "custom_service_cpu_usage",
			Value:      serviceMetrics.CpuUsage,
			Timestamp:  currentTime,
			Labels:     map[string]string{"service_name": serviceName},
		})
		metrics = append(metrics, MetricData{
			MetricName: "custom_service_ram_usage_mb",
			Value:      convertToMB(serviceMetrics.RamUsage),
			Timestamp:  currentTime,
			Labels:     map[string]string{"service_name": serviceName},
		})
		metrics = append(metrics, MetricData{
			MetricName: "custom_service_ram_max_mb",
			Value:      convertToMB(serviceMetrics.RamMax),
			Timestamp:  currentTime,
			Labels:     map[string]string{"service_name": serviceName},
		})
		metrics = append(metrics, MetricData{
			MetricName: "custom_service_disk_usage_mb",
			Value:      convertToMB(serviceMetrics.DiskUsage),
			Timestamp:  currentTime,
			Labels:     map[string]string{"service_name": serviceName},
		})
		metrics = append(metrics, MetricData{
			MetricName: "custom_service_network_receive_mb",
			Value:      convertToMB(serviceMetrics.NetworkReceive),
			Timestamp:  currentTime,
			Labels:     map[string]string{"service_name": serviceName},
		})
		metrics = append(metrics, MetricData{
			MetricName: "custom_service_network_transmit_mb",
			Value:      convertToMB(serviceMetrics.NetworkTransmit),
			Timestamp:  currentTime,
			Labels:     map[string]string{"service_name": serviceName},
		})
	}
	metricsCopy := &metrics
	um.ResetAllReadedMetrics()
	return *metricsCopy
}

func (um *UsageMetrics) GetCustomMetricDataFromNodeExporter() []MetricData {
	currentTime := time.Now().Unix()
	metrics := []MetricData{
		{
			MetricName: "custom_node_cpu_usage_percentage",
			Value:      calculateNodeCPUUsagePercentage(um.nodeCPUUser, um.nodeCPUSystem, um.nodeCPUTotal),
			Timestamp:  currentTime,
			Labels:     map[string]string{},
		},
		{
			MetricName: "custom_node_ram_available_mb",
			Value:      convertToMB(um.nodeRAMAvailable),
			Timestamp:  currentTime,
			Labels:     map[string]string{},
		},
		{
			MetricName: "custom_node_ram_total_mb",
			Value:      convertToMB(um.nodeRAMTotal),
			Timestamp:  currentTime,
			Labels:     map[string]string{},
		},
		{
			MetricName: "custom_node_disk_usage_gb",
			Value:      calculateDiskUsage(um.nodeDiskTotal, um.nodeDiskUsage),
			Timestamp:  currentTime,
			Labels:     map[string]string{},
		},
		{
			MetricName: "custom_node_disk_total_gb",
			Value:      convertToGB(um.nodeDiskTotal),
			Timestamp:  currentTime,
			Labels:     map[string]string{},
		},
		{
			MetricName: "custom_node_network_receive_mb",
			Value:      convertToMB(um.nodeNetworkReceive),
			Timestamp:  currentTime,
			Labels:     map[string]string{},
		},
		{
			MetricName: "custom_node_network_transmit_mb",
			Value:      convertToMB(um.nodeNetworkTransmit),
			Timestamp:  currentTime,
			Labels:     map[string]string{},
		},
	}
	metricsCopy := &metrics
	um.ResetAllReadedMetrics()
	return *metricsCopy
}
func (um *UsageMetrics) ResetAllReadedMetrics() {
	um.nodeCPUUsage = 0
	um.nodeCPUUser = 0
	um.nodeCPUSystem = 0
	um.nodeCPUIdle = 0
	um.nodeCPUTotal = 0
	um.nodeRAMAvailable = 0
	um.nodeRAMTotal = 0
	um.nodeDiskUsage = 0
	um.nodeDiskTotal = 0
	um.nodeNetworkReceive = 0
	um.nodeNetworkTransmit = 0
	um.ActiveServices = make(map[string]ServiceUsageMetrics)
}

func calculateNodeCPUUsagePercentage(nodeCPUUser, nodeCPUSystem, nodeCPUTotal float64) float64 {
	if nodeCPUTotal == 0 {
		return 0
	}
	return ((nodeCPUUser + nodeCPUSystem) / nodeCPUTotal) * 100
}

func convertToMB(bytes float64) float64 {
	return bytes / BytesInMB
}

func convertToGB(bytes float64) float64 {
	return bytes / BytesInGB
}

func calculateDiskUsage(total, usage float64) float64 {
	return (total - usage) / BytesInGB
}
