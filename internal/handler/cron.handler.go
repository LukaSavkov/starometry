package handler

import "github.com/c12s/metrics/internal/service"

type CronHandler struct {
	cronService *service.CronService
}

func NewCronHandler(cronService *service.CronService) *CronHandler {
	return &CronHandler{
		cronService: cronService,
	}
}
