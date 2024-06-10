package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/c12s/metrics/internal/startup"
)

func main() {
	app, err := startup.NewApp()
	if err != nil {
		log.Fatalln(err)
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	<-shutdown

	app.GracefulStop()
}
