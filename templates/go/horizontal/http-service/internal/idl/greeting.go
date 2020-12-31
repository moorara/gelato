package idl

/*
The files in this package should ideally be auto-generated
from the OpenAPI specifications in the top-level `idl` directory.
*/

import (
	"net/http"

	"github.com/gorilla/mux"

	"horizontal/http-service/pkg/xhttp"
)

// GreetRequest is the HTTP (wire/transport protocol) model for a Greet request.
type GreetRequest struct {
	Name string `json:"name"`
}

// GreetResponse is the HTTP (wire/transport protocol) model for a Greet response.
type GreetResponse struct {
	Greeting string `json:"greeting"`
}

// GreetingHandler is the interface for greeting handler functions.
type GreetingHandler interface {
	Greet(http.ResponseWriter, *http.Request)
}

// RegisterGreetingHandler registers the HTTP routes for greeting handler.
// Middleware are applied from left to right (the first middleware is the most inner and the last middleware is the most outter).
func RegisterGreetingHandler(router *mux.Router, handler GreetingHandler, middleware ...xhttp.Middleware) {
	greetHandler := handler.Greet

	for _, mid := range middleware {
		greetHandler = mid.Wrap(greetHandler)
	}

	router.Name("Greet").Methods("POST").Path("/v1/greet").HandlerFunc(greetHandler)
}
