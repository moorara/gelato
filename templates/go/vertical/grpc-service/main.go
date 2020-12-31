package main

import (
	"flag"
	"os"

	"github.com/moorara/graceful"
	"github.com/moorara/health"
	"github.com/moorara/konfig"
	"github.com/moorara/observer"
	"github.com/moorara/observer/ogrpc"
	"go.uber.org/zap"

	"vertical/grpc-service/internal/greeting"
	"vertical/grpc-service/internal/server"
	"vertical/grpc-service/pkg/client"
	"vertical/grpc-service/version"
)

// Configurations
var config = struct {
	Name                 string
	HTTPPort             uint16
	GRPCPort             uint16
	Environment          string
	Region               string
	LogLevel             string
	OpenTelemetryAddress string
}{
	Name:        "grpc-service",
	HTTPPort:    4000,    // default
	GRPCPort:    5000,    // default
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
	serverInterceptor := ogrpc.NewServerInterceptor(observer, ogrpc.Options{})

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

	grpcServer, err := server.NewGRPCServer(greetingService, server.GRPCServerOptions{
		Port:    config.GRPCPort,
		Options: serverInterceptor.ServerOptions(),
	})

	if err != nil {
		observer.Logger().Fatal("failed to create grpc server", zap.Error(err))
	}

	// Create an HTTP health handler for health checking the service by external systems
	health.SetLogger(observer.Logger().Sugar())
	health.RegisterChecker(dbClient, translateClient)
	healthHandler := health.HandlerFunc()

	httpServer, err := server.NewHTTPServer(healthHandler, server.HTTPServerOptions{
		Port: config.HTTPPort,
	})

	if err != nil {
		observer.Logger().Fatal("failed to create http server", zap.Error(err))
	}

	// Gracefully, connect the clients and start the servers
	// Gracefully, retry the lost connections
	// Gracefully, disconnect the clients and shutdown the servers on termination signals
	graceful.SetLogger(observer.Logger().Sugar())
	graceful.RegisterClient(dbClient, translateClient)
	graceful.RegisterServer(grpcServer, httpServer)
	code := graceful.StartAndWait()

	os.Exit(code)
}
