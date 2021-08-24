package lookup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	req := &Request{
		ID: "1234",
	}

	resp := &Response{
		Name: "Milad",
	}

	mock := MockService(t)
	mock.Expect().Lookup().WithArgs(req).Return(resp, nil)
	mock.Expect().Lookup().WithArgs(req).Return(resp, nil)

	p := &Provider{
		service: mock.Impl(),
	}

	name, err := p.Provide("1234")
	mock.Assert()
	assert.NoError(t, err)
	assert.Equal(t, "Milad", name)
}
