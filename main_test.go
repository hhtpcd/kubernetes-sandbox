package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_sendOk(t *testing.T) {
	tests := []struct {
		name string
		exp  int
	}{
		{
			name: "returns HTTP 200",
			exp:  200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)

			w := httptest.NewRecorder()

			sendOk(w, req)
			resp := w.Result()

			if resp.StatusCode != tt.exp {
				t.Errorf("sendOk() expected %d, got %d", tt.exp, resp.StatusCode)
			}
		})
	}
}

func Test_healthz(t *testing.T) {
	tests := []struct {
		name string
		exp  int
	}{
		{
			name: "returns HTTP 200",
			exp:  200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/healthz", nil)

			w := httptest.NewRecorder()

			healthz(w, req)
			resp := w.Result()

			if resp.StatusCode != tt.exp {
				t.Errorf("healthz() expected %d, got %d", tt.exp, resp.StatusCode)
			}
		})
	}
}
