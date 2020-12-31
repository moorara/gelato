package mapper

import (
	"errors"

	"horizontal/http-service/internal/entity"
	"horizontal/http-service/internal/idl"
)

// GreetRequestIDLToDomain transforms the IDL-specific (wire or transport protocol) representation of GreetRequest to its domain-specific representation.
func GreetRequestIDLToDomain(req *idl.GreetRequest) (*entity.GreetRequest, error) {
	if req == nil {
		return nil, errors.New("greet request cannot be nil")
	}

	if req.Name == "" {
		return nil, errors.New("name cannot be empty")
	}

	return &entity.GreetRequest{
		Name: req.Name,
	}, nil
}

// GreetResponseDomainToIDL transforms the domain-specific representation of GreetResponse to its IDL-specific (wire or transport protocol) representation.
func GreetResponseDomainToIDL(resp *entity.GreetResponse) (*idl.GreetResponse, error) {
	if resp == nil {
		return nil, errors.New("greet response cannot be nil")
	}

	if resp.Greeting == "" {
		return nil, errors.New("greeting cannot be empty")
	}

	return &idl.GreetResponse{
		Greeting: resp.Greeting,
	}, nil
}
