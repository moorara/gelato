package handler

import (
	"context"

	"horizontal/grpc-service/internal/controller"
	"horizontal/grpc-service/internal/idl/greetingpb"
	"horizontal/grpc-service/internal/mapper"
)

// GreetingHandler is an alias for the gRPC server interface.
type GreetingHandler = greetingpb.GreetingServiceServer

// greetingHandler implements GreetingHandler (greetingpb.GreetingServiceServer) interface.
type greetingHandler struct {
	greetingController controller.GreetingController
}

// NewGreetingHandler creates a new instance of GreetingHandler.
func NewGreetingHandler(greetingController controller.GreetingController) (GreetingHandler, error) {
	return &greetingHandler{
		greetingController: greetingController,
	}, nil
}

// Greet is the handler for GreetingService::Greet endpoint.
func (h *greetingHandler) Greet(ctx context.Context, req *greetingpb.GreetRequest) (*greetingpb.GreetResponse, error) {
	domainReq, err := mapper.GreetRequestIDLToDomain(req)
	if err != nil {
		return nil, err
	}

	domainResp, err := h.greetingController.Greet(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	resp, err := mapper.GreetResponseDomainToIDL(domainResp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
