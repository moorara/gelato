package server

import (
	"context"

	"vertical/grpc-service/internal/idl/greetingpb"
)

type GreetMock struct {
	InCtx    context.Context
	InReq    *greetingpb.GreetRequest
	OutResp  *greetingpb.GreetResponse
	OutError error
}

// MockGreetingService is a mock implementation for greetingpb.GreetingServiceServer interface.
type MockGreetingService struct {
	GreetCounter int
	GreetMocks   []GreetMock
}

func (m *MockGreetingService) Greet(ctx context.Context, req *greetingpb.GreetRequest) (*greetingpb.GreetResponse, error) {
	i := m.GreetCounter
	m.GreetCounter++
	m.GreetMocks[i].InCtx = ctx
	m.GreetMocks[i].InReq = req
	return m.GreetMocks[i].OutResp, m.GreetMocks[i].OutError
}
