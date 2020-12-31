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

	"vertical/http-service/internal/greeting"
	"vertical/http-service/internal/server"
	"vertical/http-service/pkg/client"
	"vertical/http-service/pkg/xhttp"
	"vertical/http-service/version"
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

	// CREATE CLIENTS

	dbClient, err := client.New("db-client")
	if err != nil {
		observer.Logger().Fatal("failed to create database client", zap.Error(err))
	}

	translateClient, err := client.New("translate-client")
	if err != nil {
		observer.Logger().Fatal("failed to create translate client", zap.Error(err))
	}

	// CREATE SERVICES

	greetingService, err := greeting.NewService(dbClient, translateClient)
	if err != nil {
		observer.Logger().Fatal("failed to create greetting service", zap.Error(err))
	}

	// CREATE SERVERS

	// Create an HTTP health handler for health checking the service by external systems
	health.SetLogger(observer.Logger().Sugar())
	health.RegisterChecker(dbClient, translateClient)
	healthHandler := health.HandlerFunc()

	httpServer, err := server.NewHTTPServer(healthHandler, greetingService, server.HTTPServerOptions{
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
	graceful.RegisterClient(dbClient, translateClient)
	graceful.RegisterServer(httpServer)
	code := graceful.StartAndWait()

	os.Exit(code)
}
