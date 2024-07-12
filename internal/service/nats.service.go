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
	subject := fmt.Sprintf("%s.metrics", ns.metricsService.NodeID)
	ns.natsConnection.Subscribe(subject, func(msg *nats.Msg) {
		log.Println("NATS REQUEST START")
		writtenMetrics, err := ns.metricsService.GetLatestMetrics()
		log.Println("NATS REQUEST FINISH GET METRICS")
		if err != nil {
			log.Println("ERR IZ METRIKA", err)
		}
		jsonData, errFromCast := json.Marshal(writtenMetrics)
		if errFromCast != nil {
			log.Println(errFromCast)
		}
		log.Println("NATS REQUEST FINISH MARSHAL")
		err2 := msg.Respond([]byte(jsonData))
		log.Println(err2)
		log.Println("NATS REQUEST FINISH")
	})

}

// func (ns NatsService) TestPublish() {
// 	ns.natsConnection.Publish("getMetrics", nil)
// }
