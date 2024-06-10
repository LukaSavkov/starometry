package service

import (
	"log"

	"github.com/robfig/cron/v3"
)

type CronService struct {
	CronScheduler *cron.Cron
	CronIDs       []cron.EntryID
}

func NewCronService() *CronService {
	c := cron.New(cron.WithSeconds())
	return &CronService{
		CronScheduler: c,
	}
}

func (cs *CronService) Start() {
	cs.CronScheduler.Start()
}

func (cs *CronService) AddJob(spec string, job func()) {
	cronId, err := cs.CronScheduler.AddFunc(spec, job)
	if err != nil {
		log.Fatalln(err)
	}
	cs.CronIDs = append(cs.CronIDs, cronId)
}

func (cs *CronService) Stop() {
	cs.CronScheduler.Stop()
}

func (cs *CronService) TerminateJob(cronID cron.EntryID) {
	cs.CronScheduler.Remove(cronID)
}
