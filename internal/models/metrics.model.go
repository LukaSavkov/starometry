package models

type MetricData struct {
	MetricName string            `json:"metric_name"`
	Labels     map[string]string `json:"labels"`
	Value      float64           `json:"value"`
	Timestamp  int64             `json:"timestamp"`
}

type MetricFileFormat struct {
	NodeId    string       `json:"nodeId"`
	ClusterId string       `json:"clusterId"`
	Metrics   []MetricData `json:"metrics"`
}
