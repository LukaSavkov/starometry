package servers

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/c12s/metrics/internal/handler"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	server         *http.Server
	metricsHandler *handler.MetricsHandler
}

func NewHttpServer(metricsHandler *handler.MetricsHandler) *HttpServer {
	return &HttpServer{
		metricsHandler: metricsHandler,
	}
}

func (httpServer *HttpServer) ConfigureRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", httpServer.metricsHandler.Test).Methods("GET")
	router.HandleFunc("/latest", httpServer.metricsHandler.GetLatestMetrics).Methods("GET")
	router.HandleFunc("/place-new-config", httpServer.metricsHandler.PostNewMetrics).Methods("POST")
	return router
}

func (httpServer *HttpServer) InitServer(port string) {
	httpServer.server = &http.Server{
		Addr:         ":" + port,
		Handler:      httpServer.ConfigureRouter(),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
}

func (httpServer *HttpServer) GetHttpServer() *http.Server {
	return httpServer.server
}

func (httpServer *HttpServer) Run() {
	go func() {
		log.Println("HTTP Server running.")
		if err := httpServer.server.ListenAndServe(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, os.Kill)

	<-stopChan
	log.Println("Received terminate, graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.server.Shutdown(ctx); err != nil {
		log.Fatalf("Cannot gracefully shutdown: %v", err)
	}
	log.Println("Server stopped")
}
