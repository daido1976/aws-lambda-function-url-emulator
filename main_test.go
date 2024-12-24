package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	mockRieServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		responseBody, _ := json.Marshal(map[string]interface{}{
			"message":       "Hello from Lambda!",
			"requestedBody": string(body),
		})

		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"statusCode": 200,
			"headers": map[string]string{
				"Content-Type": "application/json",
			},
			"body": string(responseBody),
		}); err != nil {
			http.Error(w, "Failed to encode mock response", http.StatusInternalServerError)
		}
	}))
	defer mockRieServer.Close()

	rieEndpoint = mockRieServer.URL

	server := httptest.NewServer(http.HandlerFunc(mainHandler))
	defer server.Close()

	req, _ := http.NewRequest("POST", server.URL+"/foo/bar?testkey=testvalue", strings.NewReader("test body"))
	req.Header.Set("Custom-Header", "HeaderValue")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Verify status code
	if resp.StatusCode != 200 {
		t.Errorf("expected statusCode to be 200, got %d", resp.StatusCode)
	}

	// Verify headers
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected Content-Type to be application/json, got %s", contentType)
	}

	// Verify response body
	var responseBody map[string]interface{}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// Verify message
	if responseBody["message"] != "Hello from Lambda!" {
		t.Errorf("expected message to be 'Hello from Lambda!', got %v", responseBody["message"])
	}

	// Verify request body
	requestedBody := responseBody["requestedBody"].(string)

	var requestedBodyJSON map[string]interface{}
	if err := json.Unmarshal([]byte(requestedBody), &requestedBodyJSON); err != nil {
		t.Fatalf("failed to unmarshal requestedBody: %v", err)
	}

	// Verify each field
	if requestedBodyJSON["version"] != "2.0" {
		t.Errorf("expected version to be '2.0', got %v", requestedBodyJSON["version"])
	}

	if requestedBodyJSON["rawPath"] != "/foo/bar" {
		t.Errorf("expected rawPath to be '/foo/bar', got %v", requestedBodyJSON["rawPath"])
	}

	if requestedBodyJSON["rawQueryString"] != "testkey=testvalue" {
		t.Errorf("expected rawQueryString to be 'testkey=testvalue', got %v", requestedBodyJSON["rawQueryString"])
	}

	if queryParams := requestedBodyJSON["queryStringParameters"].(map[string]interface{}); queryParams["testkey"] != "testvalue" {
		t.Errorf("expected query parameter 'testkey' to be 'testvalue', got %v", queryParams["testkey"])
	}

	headers := requestedBodyJSON["headers"].(map[string]interface{})
	if headers["Custom-Header"] != "HeaderValue" {
		t.Errorf("expected header 'Custom-Header' to be 'HeaderValue', got %v", headers["Custom-Header"])
	}

	if headers["User-Agent"] != "Go-http-client/1.1" {
		t.Errorf("expected header 'User-Agent' to be 'Go-http-client/1.1', got %v", headers["User-Agent"])
	}

	if requestedBodyJSON["body"] != "test body" {
		t.Errorf("expected body to be 'test body', got %v", requestedBodyJSON["body"])
	}

	if requestedBodyJSON["isBase64Encoded"] != false {
		t.Errorf("expected isBase64Encoded to be false, got %v", requestedBodyJSON["isBase64Encoded"])
	}

	requestContext := requestedBodyJSON["requestContext"].(map[string]interface{})
	httpContext := requestContext["http"].(map[string]interface{})

	if httpContext["method"] != "POST" {
		t.Errorf("expected method to be 'POST', got %v", httpContext["method"])
	}

	if httpContext["path"] != "/foo/bar" {
		t.Errorf("expected path to be '/foo/bar', got %v", httpContext["path"])
	}

	if httpContext["protocol"] != "HTTP/1.1" {
		t.Errorf("expected protocol to be 'HTTP/1.1', got %v", httpContext["protocol"])
	}

	if httpContext["userAgent"] != "Go-http-client/1.1" {
		t.Errorf("expected userAgent to be 'Go-http-client/1.1', got %v", httpContext["userAgent"])
	}
}
