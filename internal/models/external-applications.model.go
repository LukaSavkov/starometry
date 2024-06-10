package models

import "github.com/c12s/metrics/pkg/external"

type ExternalApplication struct {
	Address        string `json:"address"`
	ExternalClient external.ExternalMetricsClient
}
