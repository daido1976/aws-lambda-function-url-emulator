package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	mockRieServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"statusCode": 200, "body": "mock response", "isBase64Encoded": false}`)
	}))
	defer mockRieServer.Close()

	rieEndpoint = mockRieServer.URL

	server := httptest.NewServer(http.HandlerFunc(mainHandler))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/test?key=value", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected statusCode to be 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "mock response" {
		t.Errorf("expected Body to be 'mock response', got %s", string(body))
	}
}
