package main

import (
	"flag"
	"os"

	"github.com/moorara/graceful"
	"github.com/moorara/health"
	"github.com/moorara/konfig"
	"github.com/moorara/observer"
	"github.com/moorara/observer/ohttp"
	"go.uber.org/zap"

	"horizontal/http-service/internal/controller"
	"horizontal/http-service/internal/gateway"
	"horizontal/http-service/internal/handler"
	"horizontal/http-service/internal/repository"
	"horizontal/http-service/internal/server"
	"horizontal/http-service/pkg/xhttp"
	"horizontal/http-service/version"
)

// Configurations
var config = struct {
	Name                 string
	HTTPPort             uint16
	Environment          string
	Region               string
	LogLevel             string
	OpenTelemetryAddress string
}{
	Name:        "http-service",
	HTTPPort:    4000,    // default
	Environment: "dev",   // default
	Region:      "local", // default
	LogLevel:    "debug", // default
}

func main() {
	// Get configurations
	_ = konfig.Pick(&config)
	flag.Parse()

	// CREATE AN OBSERVER

	observerOpts := []observer.Option{
		observer.WithMetadata(config.Name, version.Version, config.Environment, config.Region, map[string]string{}),
		observer.WithLogger(config.LogLevel),
	}

	if config.OpenTelemetryAddress != "" {
		observerOpts = append(observerOpts,
			observer.WithOpenTelemetry(config.OpenTelemetryAddress, nil),
		)
	}

	observer := observer.New(true, observerOpts...)
	observabilityMiddleware := ohttp.NewMiddleware(observer, ohttp.Options{})

	// CREATE GATEWAYS

	translateGateway, err := gateway.NewTranslateGateway()
	if err != nil {
		observer.Logger().Fatal("failed to create translate gateway", zap.Error(err))
	}

	// CREATE REPOSITORIES

	greetingRepository, err := repository.NewGreetingRepository()
	if err != nil {
		observer.Logger().Fatal("failed to create greeting repository", zap.Error(err))
	}

	// CREATE CONTROLLERS

	greetingController, err := controller.NewGreetingController(translateGateway, greetingRepository)
	if err != nil {
		observer.Logger().Fatal("failed to create greeting controller", zap.Error(err))
	}

	// CREATE HANDLERS

	greetingHandler, err := handler.NewGreetingHandler(greetingController)
	if err != nil {
		observer.Logger().Fatal("failed to create greetting handler", zap.Error(err))
	}

	// CREATE SERVERS

	// Create an HTTP health handler for health checking the service by external systems
	health.SetLogger(observer.Logger().Sugar())
	health.RegisterChecker(translateGateway, greetingRepository)
	healthHandler := health.HandlerFunc()

	httpServer, err := server.NewHTTPServer(healthHandler, greetingHandler, server.HTTPServerOptions{
		Port: config.HTTPPort,
		Middleware: []xhttp.Middleware{
			observabilityMiddleware,
		},
	})

	if err != nil {
		observer.Logger().Fatal("failed to create http server", zap.Error(err))
	}

	// Gracefully, connect the clients and start the servers
	// Gracefully, retry the lost connections
	// Gracefully, disconnect the clients and shutdown the servers on termination signals
	graceful.SetLogger(observer.Logger().Sugar())
	graceful.RegisterClient(translateGateway, greetingRepository)
	graceful.RegisterServer(httpServer)
	code := graceful.StartAndWait()

	os.Exit(code)
}
