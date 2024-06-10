package config

import (
	"github.com/c12s/metrics/internal/client"
	"github.com/c12s/metrics/internal/models"
)

type ExternalApplicationsConfig struct {
	externalApplications []models.ExternalApplication
}

func NewExternalApplicationsConfig() *ExternalApplicationsConfig {
	return &ExternalApplicationsConfig{
		externalApplications: make([]models.ExternalApplication, 0),
	}
}

func (eac *ExternalApplicationsConfig) loadNewApplicationWithConfiguration(application models.ExternalApplication) {
	client, err := client.NewExternalMetricsClient(application.Address)
	if err != nil {
		return
	}
	application.ExternalClient = client
	eac.externalApplications = append(eac.externalApplications, application)
}

func (eac *ExternalApplicationsConfig) LoadNewApplications(applications []models.ExternalApplication) {
	for _, app := range applications {
		eac.loadNewApplicationWithConfiguration(app)
	}
}

func (eac *ExternalApplicationsConfig) GetExternalApplications() *[]models.ExternalApplication {
	return &eac.externalApplications
}
