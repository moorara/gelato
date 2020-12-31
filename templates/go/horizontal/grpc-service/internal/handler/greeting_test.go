package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"horizontal/grpc-service/internal/controller"
	"horizontal/grpc-service/internal/entity"
	"horizontal/grpc-service/internal/idl/greetingpb"
)

func TestNewGreetingHandler(t *testing.T) {
	tests := []struct {
		name               string
		greetingController controller.GreetingController
	}{
		{
			name:               "OK",
			greetingController: &MockGreetingController{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := NewGreetingHandler(tc.greetingController)

			assert.NoError(t, err)
			assert.NotNil(t, handler)
		})
	}
}

func TestGreetingHandlerGreet(t *testing.T) {
	tests := []struct {
		name                   string
		mockGreetingController *MockGreetingController
		ctx                    context.Context
		request                *greetingpb.GreetRequest
		expectedResponse       *greetingpb.GreetResponse
		expectedError          error
	}{
		{
			name:             "RequestMappingFails",
			ctx:              context.Background(),
			request:          nil,
			expectedResponse: nil,
			expectedError:    errors.New("greet request cannot be nil"),
		},
		{
			name: "ControllerFails",
			mockGreetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{OutError: errors.New("controller failed")},
				},
			},
			ctx: context.Background(),
			request: &greetingpb.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: nil,
			expectedError:    errors.New("controller failed"),
		},
		{
			name: "ResponseMappingFails",
			mockGreetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{OutResponse: nil},
				},
			},
			ctx: context.Background(),
			request: &greetingpb.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: nil,
			expectedError:    errors.New("greet response cannot be nil"),
		},
		{
			name: "Success",
			mockGreetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{
						OutResponse: &entity.GreetResponse{
							Greeting: "Hello, Jane!",
						},
					},
				},
			},
			ctx: context.Background(),
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
			handler := &greetingHandler{
				greetingController: tc.mockGreetingController,
			}

			response, err := handler.Greet(tc.ctx, tc.request)

			assert.Equal(t, tc.expectedResponse, response)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
