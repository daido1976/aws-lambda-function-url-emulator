package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

var port = getEnv("PORT", "8080")
var rieEndpoint = getEnv("RIE_ENDPOINT", "http://localhost:9000/2015-03-31/functions/function/invocations")

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	http.HandleFunc("/", mainHandler)
	log.Printf("[Lambda URL Proxy] Listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request
	log.Printf("[Lambda URL Proxy] %s %s\n", r.Method, r.URL.String())

	// Build the Lambda event
	lambdaEvent, err := buildLambdaEvent(r)
	if err != nil {
		http.Error(w, "Failed to build Lambda event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Invoke the Lambda function
	response, err := invokeLambda(lambdaEvent)
	if err != nil {
		http.Error(w, "Failed to invoke Lambda: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Map Lambda response to HTTP response
	if err := mapLambdaResponseToHTTP(w, response); err != nil {
		http.Error(w, "Failed to process Lambda response: "+err.Error(), http.StatusInternalServerError)
	}
}

func buildLambdaEvent(r *http.Request) (*events.APIGatewayV2HTTPRequest, error) {
	bodyBytes, _ := io.ReadAll(r.Body)
	body := string(bodyBytes)

	headers := map[string]string{}
	for key, values := range r.Header {
		headers[key] = strings.Join(values, ",")
	}

	return &events.APIGatewayV2HTTPRequest{
		Version:               "2.0",
		RawPath:               r.URL.Path,
		RawQueryString:        r.URL.RawQuery,
		Headers:               headers,
		QueryStringParameters: convertQueryParams(r.URL.Query()),
		Body:                  body,
		IsBase64Encoded:       false,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method:    r.Method,
				Path:      r.URL.Path,
				Protocol:  r.Proto,
				SourceIP:  r.RemoteAddr,
				UserAgent: r.UserAgent(),
			},
		},
	}, nil
}

func convertQueryParams(values url.Values) map[string]string {
	params := make(map[string]string)
	for key, value := range values {
		params[key] = strings.Join(value, ",")
	}
	return params
}

func invokeLambda(event *events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", rieEndpoint, strings.NewReader(string(eventData)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var lambdaResponse events.APIGatewayV2HTTPResponse
	if err := json.NewDecoder(res.Body).Decode(&lambdaResponse); err != nil {
		return nil, err
	}
	return &lambdaResponse, nil
}

func mapLambdaResponseToHTTP(w http.ResponseWriter, lambdaResponse *events.APIGatewayV2HTTPResponse) error {
	// Set headers
	// NOTE: Headers must be set before calling WriteHeader, as changes to headers after this have no effect.
	// See. https://pkg.go.dev/net/http#ResponseWriter
	for key, value := range lambdaResponse.Headers {
		w.Header().Set(key, value)
	}

	// Set the status code
	w.WriteHeader(lambdaResponse.StatusCode)

	// Handle Base64 encoding
	if lambdaResponse.IsBase64Encoded {
		bodyBytes, err := base64.StdEncoding.DecodeString(lambdaResponse.Body)
		if err != nil {
			return err
		}
		w.Write(bodyBytes)
	} else {
		w.Write([]byte(lambdaResponse.Body))
	}
	return nil
}
