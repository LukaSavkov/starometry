package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsService struct {
	natsConnection *nats.Conn
	metricsService *MetricsService
}

func NewNatsService(natsAddress string, metricsService *MetricsService) (*NatsService, error) {
	natsConnection, err := nats.Connect(fmt.Sprintf("nats://%s", natsAddress))
	if err != nil {
		return nil, err
	}
	log.Println("Nats connection succeed on address ", natsAddress)
	return &NatsService{
		natsConnection: natsConnection,
		metricsService: metricsService,
	}, nil
}

func (ns *NatsService) Disconnect() {
	if ns.natsConnection != nil {
		ns.natsConnection.Close()
		ns.natsConnection = nil
	}
}

func (ns NatsService) InitializeMetricsSubscriber() {
	ns.natsConnection.Subscribe("getMetrics", func(msg *nats.Msg) {
		writtenMetrics, err := ns.metricsService.GetLatestMetrics()
		if err != nil {
			log.Println(err)
		}
		jsonData, errFromCast := json.Marshal(writtenMetrics)
		if errFromCast != nil {
			log.Println(errFromCast)
		}
		msg.Respond([]byte(jsonData))
	})

}

func (ns NatsService) TestPublish() {
	ns.natsConnection.Publish("getMetrics", nil)
}
