package handler

import (
	"encoding/json"
	"net/http"

	"horizontal/http-service/internal/controller"
	"horizontal/http-service/internal/idl"
	"horizontal/http-service/internal/mapper"
	"horizontal/http-service/pkg/xhttp"
)

// GreetingHandler is an alias for the HTTP service interface.
type GreetingHandler = idl.GreetingHandler

// greetingHandler implements GreetingHandler (idl.GreetingHandler) interface.
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
func (h *greetingHandler) Greet(w http.ResponseWriter, r *http.Request) {
	req := new(idl.GreetRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		xhttp.Error(w, err, http.StatusBadRequest)
		return
	}

	domainReq, err := mapper.GreetRequestIDLToDomain(req)
	if err != nil {
		xhttp.Error(w, err, http.StatusBadRequest)
		return
	}

	domainResp, err := h.greetingController.Greet(r.Context(), domainReq)
	if err != nil {
		xhttp.Error(w, err, http.StatusInternalServerError)
		return
	}

	resp, err := mapper.GreetResponseDomainToIDL(domainResp)
	if err != nil {
		xhttp.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
