package startup

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/c12s/metrics/internal/config"
	"github.com/c12s/metrics/internal/handler"
	"github.com/c12s/metrics/internal/mappers"
	"github.com/c12s/metrics/internal/models"
	"github.com/c12s/metrics/internal/servers"
	"github.com/c12s/metrics/internal/service"

	pkgAPI "github.com/c12s/metrics/pkg/api"
	"github.com/c12s/metrics/pkg/external"
	"google.golang.org/grpc"
)

type App struct {
	appConfig                  *config.AppConfig
	metricsConfig              *config.MetricsConfig
	httpServer                 *servers.HttpServer
	grpcServer                 *grpc.Server
	externalApplicationsConfig *config.ExternalApplicationsConfig
	shutdownProcesses          []func()
}

func NewApp() (*App, error) {
	cfg := config.NewAppConfigFromEnv()
	if cfg == nil {
		return nil, errors.New("Basic server configuration is empty")
	}
	metricsCfg, err := config.NewMetricsConfigLoadedFromEnv()
	if err != nil {
		fmt.Println("Metrics from env are having some errors. Proceeding with preset configuration.")
		metricsCfg = config.NewMetricsConfigWithPresetConfiguration()
	}
	app := &App{
		appConfig:         cfg,
		metricsConfig:     metricsCfg,
		shutdownProcesses: make([]func(), 0),
	}
	app.init()
	return app, nil
}

func (app *App) init() {
	app.externalApplicationsConfig = config.NewExternalApplicationsConfig()
	fileService := service.NewLocalFileService()
	app.initializeNodeID(fileService)
	metricsService := service.NewMetricsService(fileService, app.metricsConfig, app.appConfig.GetNodeID(), app.appConfig)
	natsService, err := service.NewNatsService(app.appConfig.GetNatsAddress(), metricsService)
	if err != nil {
		log.Println(err)
	}
	app.shutdownProcesses = append(app.shutdownProcesses, func() {
		natsService.Disconnect()
	})
	cronService := service.NewCronService()
	app.shutdownProcesses = append(app.shutdownProcesses, func() {
		cronService.Stop()
	})
	natsService.InitializeMetricsSubscriber()
	metricsHandler := handler.NewMetricsHandler(metricsService)
	cronService.AddJob("@every "+app.metricsConfig.GetCronTimer(), func() {
		metricsService.GetMetrics()
		log.Println(app.metricsConfig.GetQueries())
	})
	cronService.AddJob("@every "+app.metricsConfig.GetExternalCronTimer(), func() {
		applicationsListtWithClients := app.externalApplicationsConfig.GetExternalApplications()
		if applicationsListtWithClients == nil || len(*applicationsListtWithClients) == 0 {
			return
		}
		var metrics []models.MetricData
		for _, application := range *applicationsListtWithClients {
			data, err := application.ExternalClient.ExternalLatestMetrics(context.Background(), &external.ExternalLatestMetricsReq{})
			if err != nil {
				continue
			}
			metrics = append(metrics, mappers.MapFromExternalMetricDataToModelMetricData(data.Metrics)...)
		}
		metricsService.WriteMetricsFromExternalApplication(metrics)
	})
	cronService.Start()
	metricsService.GetMetrics()

	// servers
	customHttpServer := servers.NewHttpServer(metricsHandler)
	configGrpcServer := servers.NewMetricsGrpcServer(metricsService, app.externalApplicationsConfig)
	s := grpc.NewServer()
	pkgAPI.RegisterMetricsServer(s, configGrpcServer)
	app.grpcServer = s
	err = app.startGrpcServer()
	if err != nil {
		log.Println(err)
	}
	app.httpServer = customHttpServer
	app.startHttpServer()
}

func (app *App) startGrpcServer() error {
	lis, err := net.Listen("tcp", ":"+app.appConfig.GetGRPCPort())
	if err != nil {
		return err
	}
	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := app.grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return nil
}

func (app *App) startHttpServer() {
	app.httpServer.InitServer(app.appConfig.GetServerPort())
	app.httpServer.Run()
}

func (app *App) GracefulStop() {
	app.grpcServer.GracefulStop()
	for i, shudownProcess := range app.shutdownProcesses {
		log.Println("Shutting down the process with index ", i)
		shudownProcess()
	}
}

func (app *App) initializeNodeID(fileService *service.LocalFileService) {
	readedNodeIDInBytes, errFromReading := fileService.ReadFromFile("/etc/c12s/nodeid")
	if errFromReading != nil {
		log.Fatalln(errFromReading.GetErrorMessage())
	}
	stringValueOfNodeID := string(readedNodeIDInBytes)
	app.appConfig.SetNodeID(stringValueOfNodeID)

}
