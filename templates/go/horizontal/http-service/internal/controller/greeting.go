package controller

import (
	"context"
	"fmt"

	"horizontal/http-service/internal/entity"
	"horizontal/http-service/internal/gateway"
	"horizontal/http-service/internal/repository"
)

// GreetingController is the interface for greeting business logic.
type GreetingController interface {
	Greet(context.Context, *entity.GreetRequest) (*entity.GreetResponse, error)
}

// greetingController implements GreetingController interface.
type greetingController struct {
	translateGateway   gateway.TranslateGateway
	greetingRepository repository.GreetingRepository
}

// NewGreetingController creates a new instance of GreetingController.
func NewGreetingController(translateGateway gateway.TranslateGateway, greetingRepository repository.GreetingRepository) (GreetingController, error) {
	return &greetingController{
		translateGateway:   translateGateway,
		greetingRepository: greetingRepository,
	}, nil
}

// Greet creates and returns a greeting for a given name!
func (c *greetingController) Greet(ctx context.Context, req *entity.GreetRequest) (*entity.GreetResponse, error) {
	// Call an external service through its gateway
	hello, err := c.translateGateway.GetValue(ctx, "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "en")
	if err != nil {
		return nil, err
	}

	greeting := fmt.Sprintf("%s, %s!", hello, req.Name)

	// Interact with a data store through its repository
	_, err = c.greetingRepository.Create(ctx, greeting)
	if err != nil {
		return nil, err
	}

	resp := &entity.GreetResponse{
		Greeting: greeting,
	}

	return resp, nil
}
