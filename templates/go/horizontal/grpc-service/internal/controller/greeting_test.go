package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"horizontal/grpc-service/internal/entity"
	"horizontal/grpc-service/internal/gateway"
	"horizontal/grpc-service/internal/repository"
)

func TestNewGreetingController(t *testing.T) {
	tests := []struct {
		name               string
		translateGateway   gateway.TranslateGateway
		greetingRepository repository.GreetingRepository
	}{
		{
			name:               "OK",
			translateGateway:   &MockTranslateGateway{},
			greetingRepository: &MockGreetingRepository{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			controller, err := NewGreetingController(tc.translateGateway, tc.greetingRepository)

			assert.NoError(t, err)
			assert.NotNil(t, controller)
		})
	}
}

func TestGreetingControllerGreet(t *testing.T) {
	tests := []struct {
		name                  string
		mockTranslateGateway  *MockTranslateGateway
		mockGeetingRepository *MockGreetingRepository
		ctx                   context.Context
		request               *entity.GreetRequest
		expectedResponse      *entity.GreetResponse
		expectedError         error
	}{
		{
			name: "GatewayFails",
			mockTranslateGateway: &MockTranslateGateway{
				GetValueMocks: []GetValueMock{
					{OutError: errors.New("error on calling translation service")},
				},
			},
			ctx: context.Background(),
			request: &entity.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: nil,
			expectedError:    errors.New("error on calling translation service"),
		},
		{
			name: "RepositoryFails",
			mockTranslateGateway: &MockTranslateGateway{
				GetValueMocks: []GetValueMock{
					{OutString: "Hello"},
				},
			},
			mockGeetingRepository: &MockGreetingRepository{
				CreateMocks: []CreateMock{
					{OutError: errors.New("error on storing data")},
				},
			},
			ctx: context.Background(),
			request: &entity.GreetRequest{
				Name: "Jane",
			},
			expectedError: errors.New("error on storing data"),
		},
		{
			name: "Success",
			mockTranslateGateway: &MockTranslateGateway{
				GetValueMocks: []GetValueMock{
					{OutString: "Hello"},
				},
			},
			mockGeetingRepository: &MockGreetingRepository{
				CreateMocks: []CreateMock{
					{OutTimestamp: time.Now()},
				},
			},
			ctx: context.Background(),
			request: &entity.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: &entity.GreetResponse{
				Greeting: "Hello, Jane!",
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			controller := &greetingController{
				translateGateway:   tc.mockTranslateGateway,
				greetingRepository: tc.mockGeetingRepository,
			}

			response, err := controller.Greet(tc.ctx, tc.request)

			assert.Equal(t, tc.expectedResponse, response)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
