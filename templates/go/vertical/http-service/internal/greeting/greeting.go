package greeting

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"vertical/http-service/pkg/client"
	"vertical/http-service/pkg/xhttp"
)

// GreetRequest is the HTTP (wire/transport protocol) model for a Greet request.
type GreetRequest struct {
	Name string `json:"name"`
}

// GreetResponse is the HTTP (wire/transport protocol) model for a Greet response.
type GreetResponse struct {
	Greeting string `json:"greeting"`
}

// Service implements the HTTP handlers for Greeting APIs.
type Service struct {
	dbClient        client.Client
	translateClient client.Client
}

// NewService creates a new instance of Greeting service.
func NewService(dbClient client.Client, translateClient client.Client) (*Service, error) {
	return &Service{
		dbClient:        dbClient,
		translateClient: translateClient,
	}, nil
}

// Greet is the handler for GreetingService::Greet endpoint.
func (s *Service) Greet(w http.ResponseWriter, r *http.Request) {
	req := new(GreetRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		xhttp.Error(w, err, http.StatusBadRequest)
		return
	}

	// Call an external service through its client
	// s.translateClient
	hello := "Hello"

	greeting := fmt.Sprintf("%s, %s!", hello, req.Name)

	// Interact with a data store through its client
	// s.dbClient

	resp := &GreetResponse{
		Greeting: greeting,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// RegisterRoutes registers the HTTP routes for greeting service.
// Middleware are applied from left to right (the first middleware is the most inner and the last middleware is the most outter).
func (s *Service) RegisterRoutes(router *mux.Router, middleware ...xhttp.Middleware) {
	greetHandler := s.Greet

	for _, mid := range middleware {
		greetHandler = mid.Wrap(greetHandler)
	}

	router.Name("Greet").Methods("POST").Path("/v1/greet").HandlerFunc(greetHandler)
}
