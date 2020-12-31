package handler

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"horizontal/http-service/internal/controller"
	"horizontal/http-service/internal/entity"
	"horizontal/http-service/pkg/xhttp"
)

func TestNewGreetingHandler(t *testing.T) {
	tests := []struct {
		name               string
		greetingController controller.GreetingController
	}{
		{
			name:               "OK",
			greetingController: &MockGreetingController{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := NewGreetingHandler(tc.greetingController)

			assert.NoError(t, err)
			assert.NotNil(t, handler)
		})
	}
}

func TestGreetingHandlerGreet(t *testing.T) {
	tests := []struct {
		name                   string
		mockGreetingController *MockGreetingController
		req                    *http.Request
		expectedStatusCode     int
		expectedBody           string
	}{
		{
			name:               "RequestDecodingFails",
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{`)),
			expectedStatusCode: 400,
			expectedBody:       "unexpected EOF\n",
		},
		{
			name:               "RequestMappingFails",
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "name": "" }`)),
			expectedStatusCode: 400,
			expectedBody:       "name cannot be empty\n",
		},
		{
			name: "ControllerFails",
			mockGreetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{OutError: xhttp.NewServerError(errors.New("controller failed"), 500)},
				},
			},
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "name": "Jane" }`)),
			expectedStatusCode: 500,
			expectedBody:       "controller failed\n",
		},
		{
			name: "ResponseMappingFails",
			mockGreetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{OutResponse: &entity.GreetResponse{}},
				},
			},
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "name": "Jane" }`)),
			expectedStatusCode: 500,
			expectedBody:       "greeting cannot be empty\n",
		},
		{
			name: "Success",
			mockGreetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{
						OutResponse: &entity.GreetResponse{
							Greeting: "Hello, Jane!",
						},
					},
				},
			},
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "name": "Jane" }`)),
			expectedStatusCode: 200,
			expectedBody:       "{\"greeting\":\"Hello, Jane!\"}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := &greetingHandler{
				greetingController: tc.mockGreetingController,
			}

			rec := httptest.NewRecorder()
			handler.Greet(rec, tc.req)

			res := rec.Result()
			b, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)
			body := string(b)

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
			assert.Equal(t, tc.expectedBody, body)
		})
	}
}
