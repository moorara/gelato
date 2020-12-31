package greeting

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"vertical/http-service/pkg/client"
	"vertical/http-service/pkg/xhttp"
)

func TestNewService(t *testing.T) {
	dbClient, _ := client.New("db-client")
	translateClient, _ := client.New("translate-client")

	tests := []struct {
		name            string
		dbClient        client.Client
		translateClient client.Client
	}{
		{
			name:            "OK",
			dbClient:        dbClient,
			translateClient: translateClient,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service, err := NewService(tc.dbClient, tc.translateClient)

			assert.NoError(t, err)
			assert.NotNil(t, service)
		})
	}
}

func TestGreetingHandlerGreet(t *testing.T) {
	dbClient, _ := client.New("db-client")
	translateClient, _ := client.New("translate-client")

	tests := []struct {
		name               string
		dbClient           client.Client
		translateClient    client.Client
		req                *http.Request
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "RequestDecodingFails",
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{`)),
			expectedStatusCode: 400,
			expectedBody:       "unexpected EOF\n",
		},
		{
			name:               "Success",
			dbClient:           dbClient,
			translateClient:    translateClient,
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "name": "Jane" }`)),
			expectedStatusCode: 200,
			expectedBody:       "{\"greeting\":\"Hello, Jane!\"}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{
				dbClient:        tc.dbClient,
				translateClient: tc.translateClient,
			}

			rec := httptest.NewRecorder()
			service.Greet(rec, tc.req)

			res := rec.Result()
			b, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)
			body := string(b)

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
			assert.Equal(t, tc.expectedBody, body)
		})
	}
}

func TestRegisterRoutes(t *testing.T) {
	tests := []struct {
		name       string
		router     *mux.Router
		middleware []xhttp.Middleware
	}{
		{
			name:   "OK",
			router: mux.NewRouter(),
			middleware: []xhttp.Middleware{
				&MockMiddleware{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{}

			service.RegisterRoutes(tc.router, tc.middleware...)

			assert.NotEmpty(t, tc.router.Get("Greet"))
		})
	}
}
