package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_helloHandler(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		response string
		status   int
	}{
		{
			name:     "Test GET",
			method:   http.MethodGet,
			response: "Hello, World!",
			status:   http.StatusOK,
		},
		{
			name:     "Test POST",
			method:   http.MethodPost,
			response: "Method not allowed.",
			status:   http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			responseRecorder := httptest.NewRecorder()
			helloHandler(responseRecorder, req)

			if responseRecorder.Code != tt.status {
				t.Errorf("Want status '%d', got '%d'", tt.status, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tt.response {
				t.Errorf("Want '%s', got '%s'", tt.response, responseRecorder.Body)
			}
		})
	}
}
