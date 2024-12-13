package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	rieEndpoint = "http://localhost:9000/2015-03-31/functions/function/invocations"
	port        = 8080
)

type LambdaEvent struct {
	Version               string            `json:"version"`
	RawPath               string            `json:"rawPath"`
	RawQueryString        string            `json:"rawQueryString"`
	Headers               map[string]string `json:"headers"`
	QueryStringParameters map[string]string `json:"queryStringParameters,omitempty"`
	RequestContext        struct {
		HTTP struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"http"`
	} `json:"requestContext"`
	Body            string `json:"body,omitempty"`
	IsBase64Encoded bool   `json:"isBase64Encoded"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("[Proxy] %s %s\n", r.Method, r.URL.String())

		// Convert the request into a LambdaEvent
		lambdaEvent := LambdaEvent{
			Version:        "2.0",
			RawPath:        r.URL.Path,
			RawQueryString: r.URL.RawQuery,
			Headers:        make(map[string]string),
			RequestContext: struct {
				HTTP struct {
					Method string `json:"method"`
					Path   string `json:"path"`
				} `json:"http"`
			}{
				HTTP: struct {
					Method string `json:"method"`
					Path   string `json:"path"`
				}{
					Method: r.Method,
					Path:   r.URL.Path,
				},
			},
			IsBase64Encoded: false,
		}

		// Populate headers
		for name, values := range r.Header {
			lambdaEvent.Headers[name] = values[0]
		}

		// Populate query string parameters
		queryParams := r.URL.Query()
		if len(queryParams) > 0 {
			lambdaEvent.QueryStringParameters = make(map[string]string)
			for key, values := range queryParams {
				lambdaEvent.QueryStringParameters[key] = values[0]
			}
		}

		// Read and set the body
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			lambdaEvent.Body = string(bodyBytes)
		}

		// Serialize the LambdaEvent to JSON
		lambdaEventJSON, err := json.Marshal(lambdaEvent)
		if err != nil {
			http.Error(w, "Failed to marshal event", http.StatusInternalServerError)
			return
		}

		// Send the request to RIE
		resp, err := http.Post(rieEndpoint, "application/json", bytes.NewBuffer(lambdaEventJSON))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to invoke Lambda: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Forward the Lambda response back to the client
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	log.Printf("[Proxy] Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
