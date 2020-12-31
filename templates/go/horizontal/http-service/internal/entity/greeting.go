package entity

import "fmt"

// GreetRequest is the domain model for a Greet request.
type GreetRequest struct {
	Name string
}

// String implements fmt.Stringer interface.
func (r *GreetRequest) String() string {
	return fmt.Sprintf("GreetRequest{name=%s}", r.Name)
}

// GreetResponse is the domain model for a Greet response.
type GreetResponse struct {
	Greeting string
}

// String implements fmt.Stringer interface.
func (r *GreetResponse) String() string {
	return fmt.Sprintf("GreetResponse{greeting=%s}", r.Greeting)
}
