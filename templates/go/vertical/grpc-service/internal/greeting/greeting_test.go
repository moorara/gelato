package greeting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"vertical/grpc-service/internal/idl/greetingpb"
	"vertical/grpc-service/pkg/client"
)

func TestNewService(t *testing.T) {
	dbClient, _ := client.New("db-client")
	translateClient, _ := client.New("translate-client")

	tests := []struct {
		name            string
		dbClient        client.Client
		translateClient client.Client
	}{
		{
			name:            "OK",
			dbClient:        dbClient,
			translateClient: translateClient,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service, err := NewService(tc.dbClient, tc.translateClient)

			assert.NoError(t, err)
			assert.NotNil(t, service)
		})
	}
}

func TestGreetingHandlerGreet(t *testing.T) {
	dbClient, _ := client.New("db-client")
	translateClient, _ := client.New("translate-client")

	tests := []struct {
		name             string
		dbClient         client.Client
		translateClient  client.Client
		ctx              context.Context
		request          *greetingpb.GreetRequest
		expectedResponse *greetingpb.GreetResponse
		expectedError    error
	}{
		{
			name:            "Success",
			dbClient:        dbClient,
			translateClient: translateClient,
			ctx:             context.Background(),
			request: &greetingpb.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: &greetingpb.GreetResponse{
				Greeting: "Hello, Jane!",
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &service{
				dbClient:        tc.dbClient,
				translateClient: tc.translateClient,
			}

			response, err := service.Greet(tc.ctx, tc.request)

			assert.Equal(t, tc.expectedResponse, response)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
