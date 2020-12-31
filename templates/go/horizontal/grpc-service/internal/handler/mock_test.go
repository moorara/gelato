package handler

import (
	"context"

	"horizontal/grpc-service/internal/entity"
)

type GreetMock struct {
	InCtx       context.Context
	InRequest   *entity.GreetRequest
	OutResponse *entity.GreetResponse
	OutError    error
}

// MockGreetingController is a mock implementation for controller.GreetingController.
type MockGreetingController struct {
	GreetCounter int
	GreetMocks   []GreetMock
}

func (m *MockGreetingController) Greet(ctx context.Context, request *entity.GreetRequest) (*entity.GreetResponse, error) {
	i := m.GreetCounter
	m.GreetCounter++
	m.GreetMocks[i].InCtx = ctx
	m.GreetMocks[i].InRequest = request
	return m.GreetMocks[i].OutResponse, m.GreetMocks[i].OutError
}
