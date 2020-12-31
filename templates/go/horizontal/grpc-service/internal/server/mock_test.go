package server

import (
	"context"

	"horizontal/grpc-service/internal/idl/greetingpb"
)

type GreetMock struct {
	InCtx    context.Context
	InReq    *greetingpb.GreetRequest
	OutResp  *greetingpb.GreetResponse
	OutError error
}

// MockGreetingHandler is a mock implementation for handler.GreetingHandler interface.
type MockGreetingHandler struct {
	GreetCounter int
	GreetMocks   []GreetMock
}

func (m *MockGreetingHandler) Greet(ctx context.Context, req *greetingpb.GreetRequest) (*greetingpb.GreetResponse, error) {
	i := m.GreetCounter
	m.GreetCounter++
	m.GreetMocks[i].InCtx = ctx
	m.GreetMocks[i].InReq = req
	return m.GreetMocks[i].OutResp, m.GreetMocks[i].OutError
}
