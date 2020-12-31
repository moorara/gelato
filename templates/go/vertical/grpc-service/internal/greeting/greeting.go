package greeting

import (
	"context"
	"fmt"

	"vertical/grpc-service/internal/idl/greetingpb"
	"vertical/grpc-service/pkg/client"
)

// service implements greetingpb.GreetingServiceServer interface.
type service struct {
	dbClient        client.Client
	translateClient client.Client
}

// NewService creates a new instance of Greeting service.
func NewService(dbClient client.Client, translateClient client.Client) (greetingpb.GreetingServiceServer, error) {
	return &service{
		dbClient:        dbClient,
		translateClient: translateClient,
	}, nil
}

// Greet implements the GreetingService::Greet endpoint.
func (s *service) Greet(ctx context.Context, req *greetingpb.GreetRequest) (*greetingpb.GreetResponse, error) {
	// Call an external service through its client
	// s.translateClient
	hello := "Hello"

	greeting := fmt.Sprintf("%s, %s!", hello, req.Name)

	// Interact with a data store through its client
	// s.dbClient

	resp := &greetingpb.GreetResponse{
		Greeting: greeting,
	}

	return resp, nil
}
