package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainHandler(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"statusCode": 200, "body": "mock response", "isBase64Encoded": false}`)
	}))
	defer mockServer.Close()

	rieEndpoint = mockServer.URL

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lambdaEvent, err := buildLambdaEvent(r)
		if err != nil {
			http.Error(w, "Failed to build Lambda event: "+err.Error(), http.StatusInternalServerError)
			return
		}
		response, err := invokeLambda(lambdaEvent)
		if err != nil {
			http.Error(w, "Failed to invoke Lambda: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := mapLambdaResponseToHTTP(w, response); err != nil {
			http.Error(w, "Failed to process Lambda response: "+err.Error(), http.StatusInternalServerError)
		}
	})

	server := httptest.NewServer(handler)
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
