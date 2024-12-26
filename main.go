package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

var port = getEnv("PORT", "8080")
var rieEndpoint = getEnv("RIE_ENDPOINT", "http://localhost:8080/2015-03-31/functions/function/invocations")

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("[Lambda URL Proxy] Listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
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

// See. https://docs.aws.amazon.com/lambda/latest/dg/urls-invocation.html#urls-payloads
func buildLambdaEvent(r *http.Request) (*events.APIGatewayV2HTTPRequest, error) {
	bodyBytes, _ := io.ReadAll(r.Body)
	body := string(bodyBytes)

	return &events.APIGatewayV2HTTPRequest{
		Version:               "2.0",
		RouteKey:              "$default",
		RawPath:               r.URL.Path,
		RawQueryString:        r.URL.RawQuery,
		Headers:               joinValuesWithComma(r.Header),
		QueryStringParameters: joinValuesWithComma(r.URL.Query()),
		Body:                  body,
		IsBase64Encoded:       false,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RouteKey:   "$default",
			Stage:      "$default",
			Time:       time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			TimeEpoch:  time.Now().Unix(),
			DomainName: r.Host,
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

func joinValuesWithComma(input map[string][]string) map[string]string {
	result := map[string]string{}
	for key, values := range input {
		result[key] = strings.Join(values, ",")
	}
	return result
}

func invokeLambda(event *events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(rieEndpoint, "application/json", strings.NewReader(string(eventData)))
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
